package main

import (
	"log"
	"os"

	"github.com/fleimkeipa/lifery/controller"
	_ "github.com/fleimkeipa/lifery/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories"
	"github.com/fleimkeipa/lifery/uc"
	"github.com/fleimkeipa/lifery/util"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type \"Bearer \" and then your API Token
func main() {
	// Start the application
	serveApplication()
}

func serveApplication() {
	// init config
	loadConfig()

	// Create a new Echo instance
	e := echo.New()

	// Configure Echo settings
	configureEcho(e)

	// Configure CORS middleware
	configureCORS(e)

	// Configure the logger
	sugar := configureLogger(e)
	defer sugar.Sync() // Clean up logger at the end

	// Initialize PostgreSQL client
	dbClient := initDB()
	defer dbClient.Close() // Clean up db connections at the end

	userUC := initUserUC(dbClient)
	userController := controller.NewUserHandlers(userUC)

	eraUC := initEraUC(dbClient)
	eraController := controller.NewEraController(eraUC)

	eventUC := initEventUC(dbClient)
	eventController := controller.NewEventController(eventUC)

	connectUC := initConnectUC(dbClient)
	connectController := controller.NewConnectHandlers(connectUC, userUC)

	notificationUC := initNotificationUC(dbClient)
	notificationController := controller.NewNotificationHandlers(notificationUC)

	emailUC := initEmailUC()
	googleOAuthUC := initGoogleOAuthUC(dbClient)
	authHandlers := controller.NewAuthHandlers(userUC, emailUC)
	oauthHandlers := controller.NewOAuthHandlers(googleOAuthUC)

	// Define authentication routes and handlers
	authRoutes := e.Group("/auth")
	authRoutes.POST("/login", authHandlers.Login)
	authRoutes.POST("/register", authHandlers.Register)
	authRoutes.POST("/forgot-password", authHandlers.ForgotPassword)
	authRoutes.POST("/reset-password", authHandlers.ResetPassword)

	oauthRoutes := e.Group("/oauth")
	oauthRoutes.Use(util.JWTAuthViewer)
	oauthRoutes.GET("/google/url", oauthHandlers.GoogleAuthURL)
	oauthRoutes.POST("/google/callback", oauthHandlers.GoogleCallback)

	// Add JWT authentication and authorization middleware
	adminRoutes := e.Group("")
	adminRoutes.Use(util.JWTAuthAdmin)

	// Define viewer routes
	viewerRoutes := e.Group("")
	viewerRoutes.Use(util.JWTAuthViewer)

	// Define user routes
	userRoutes := e.Group("")
	userRoutes.Use(util.JWTAuthEditor)

	// Define user routes
	userRoutes = userRoutes.Group("/user")
	userRoutes.PUT("/username", userController.UpdateUsername)
	userRoutes.PUT("/password", userController.UpdatePassword)

	// Define events routes
	eventsRoutes := userRoutes.Group("/events")
	eventsRoutes.POST("", eventController.Create)
	eventsRoutes.PATCH("/:id", eventController.Update)
	eventsRoutes.DELETE("/:id", eventController.Delete)
	eventsRoutes.GET("/:id", eventController.GetByID)

	// Define public events routes
	publicEventsRoutes := viewerRoutes.Group("/events")
	publicEventsRoutes.GET("", eventController.List)

	// Define eras routes
	erasRoutes := userRoutes.Group("/eras")
	erasRoutes.POST("", eraController.Create)
	erasRoutes.PATCH("/:id", eraController.Update)
	erasRoutes.DELETE("/:id", eraController.Delete)
	erasRoutes.GET("/:id", eraController.GetByID)

	// Define public eras routes
	publicErasRoutes := viewerRoutes.Group("/eras")
	publicErasRoutes.GET("", eraController.List)

	// Define connects routes
	connectsRoutes := userRoutes.Group("/connects")
	connectsRoutes.POST("", connectController.Create)
	connectsRoutes.PATCH("/:id", connectController.Update)
	connectsRoutes.DELETE("/:id", connectController.Delete)
	connectsRoutes.GET("", connectController.ConnectsRequests)

	// Define notifications routes
	notificationsRoutes := userRoutes.Group("/notifications")
	notificationsRoutes.GET("", notificationController.List)
	notificationsRoutes.PATCH("/:id", notificationController.Update)

	// Define public user search routes
	publicUsersSearchRoutes := viewerRoutes.Group("/users")
	publicUsersSearchRoutes.GET("/search", userController.Search)

	// Define user routes
	usersRoutes := adminRoutes.Group("/users")
	usersRoutes.GET("", userController.List)
	usersRoutes.GET("/:id", userController.GetByID)
	usersRoutes.POST("", userController.Create)
	usersRoutes.PATCH("/:id", userController.Update)
	usersRoutes.DELETE("/:id", userController.DeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}

func loadConfig() {
	if os.Getenv("STAGE") == "prod" {
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// Configures the Echo instance
func configureEcho(e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true

	e.Validator = pkg.NewValidator()

	// Add Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add Recover middleware
	e.Use(middleware.Recover())
}

// Configures CORS settings
func configureCORS(e *echo.Echo) {
	corsConfig := middleware.CORSWithConfig(middleware.CORSConfig{
		UnsafeWildcardOriginWithAllowCredentials: true,
		AllowCredentials:                         true,
		AllowOrigins:                             []string{"*"},
		AllowMethods:                             []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE},
		AllowHeaders:                             []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})

	e.Use(corsConfig)
}

// Configures the logger and adds it as middleware
func configureLogger(e *echo.Echo) *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(pkg.ZapLogger(logger))

	sugar := logger.Sugar()
	loggerHandler := controller.NewLogger(sugar)
	e.Use(loggerHandler.LoggerMiddleware)

	return sugar
}

// Initializes the PostgreSQL client
func initDB() *pg.DB {
	db := pkg.NewPSQLClient()
	if db == nil {
		log.Fatal("Failed to initialize PostgreSQL client")
	}

	log.Println("PostgreSQL client initialized successfully")
	return db
}

func initEraUC(db *pg.DB) *uc.EraUC {
	eraDBRepo := repositories.NewEraRepository(db)
	return uc.NewEraUC(eraDBRepo)
}

func initUserUC(db *pg.DB) *uc.UserUC {
	userDBRepo := repositories.NewUserRepository(db)
	return uc.NewUserUC(userDBRepo)
}

func initConnectUC(db *pg.DB) *uc.ConnectsUC {
	userDBRepo := repositories.NewUserRepository(db)
	connectDBRepo := repositories.NewConnectRepository(db)
	notificationDBRepo := repositories.NewNotificationRepository(db)

	userUC := uc.NewUserUC(userDBRepo)
	notificationUC := uc.NewNotificationUC(notificationDBRepo)
	return uc.NewConnectsUC(userUC, connectDBRepo, notificationUC)
}

func initEventUC(db *pg.DB) *uc.EventUC {
	userDBRepo := repositories.NewUserRepository(db)
	connectDBRepo := repositories.NewConnectRepository(db)
	eventDBRepo := repositories.NewEventRepository(db)
	notificationDBRepo := repositories.NewNotificationRepository(db)

	userUC := uc.NewUserUC(userDBRepo)
	notificationUC := uc.NewNotificationUC(notificationDBRepo)
	connectsUC := uc.NewConnectsUC(userUC, connectDBRepo, notificationUC)

	return uc.NewEventUC(eventDBRepo, connectsUC)
}

func initNotificationUC(db *pg.DB) *uc.NotificationUC {
	notificationDBRepo := repositories.NewNotificationRepository(db)
	return uc.NewNotificationUC(notificationDBRepo)
}

func initEmailUC() *uc.EmailUC {
	emailRepo := repositories.NewEmailRepository()
	return uc.NewEmailUC(emailRepo)
}

func initGoogleOAuthUC(db *pg.DB) *uc.GoogleOAuthUC {
	userDBRepo := repositories.NewUserRepository(db)
	userUC := uc.NewUserUC(userDBRepo)
	return uc.NewGoogleOAuthUC(userUC)
}
