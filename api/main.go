package main

import (
	"log"

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
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

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

	// Initialize Redis client
	redisClient := initCache()

	eraDBRepo := repositories.NewEraRepository(dbClient)
	eraDBUC := uc.NewEraUC(eraDBRepo)
	eraController := controller.NewEraController(eraDBUC)

	eventCacheRepo := repositories.NewCacheRepository(redisClient)
	eventCacheUC := uc.NewEventCacheUC(eventCacheRepo)

	eventDBRepo := repositories.NewEventRepository(dbClient)
	eventDBUC := uc.NewEventUC(eventDBRepo, eventCacheUC)

	eventController := controller.NewEventController(eventDBUC)

	userRepo := repositories.NewUserRepository(dbClient)
	userUC := uc.NewUserUC(userRepo)
	userHandlers := controller.NewUserHandlers(userUC)

	authHandlers := controller.NewAuthHandlers(userUC)

	// Define authentication routes and handlers
	authRoutes := e.Group("/auth")
	authRoutes.POST("/login", authHandlers.Login)
	authRoutes.POST("/register", authHandlers.Register)

	// Add JWT authentication and authorization middleware
	adminRoutes := e.Group("")
	adminRoutes.Use(util.JWTAuth)

	// Define viewer routes
	viewerRoutes := e.Group("")
	viewerRoutes.Use(util.JWTAuthViewer)

	// Define events routes
	eventsRoutes := viewerRoutes.Group("/events")
	eventsRoutes.POST("", eventController.Create)
	eventsRoutes.PATCH("/:id", eventController.Update)
	eventsRoutes.DELETE("/:id", eventController.Delete)
	eventsRoutes.GET("", eventController.List)
	eventsRoutes.GET("/:id", eventController.GetByID)

	// Define eras routes
	erasRoutes := viewerRoutes.Group("/eras")
	erasRoutes.POST("", eraController.Create)
	erasRoutes.PATCH("/:id", eraController.Update)
	erasRoutes.DELETE("/:id", eraController.Delete)
	erasRoutes.GET("", eraController.List)
	erasRoutes.GET("/:id", eraController.GetByID)

	// Define user routes
	usersRoutes := adminRoutes.Group("/users")
	usersRoutes.GET("", userHandlers.List)
	usersRoutes.GET("/:id", userHandlers.GetByID)
	usersRoutes.POST("", userHandlers.Create)
	usersRoutes.PUT("/:id", userHandlers.UpdateUser)
	usersRoutes.DELETE("/:id", userHandlers.DeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// Configures the Echo instance
func configureEcho(e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true

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
		AllowMethods:                             []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
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

// Initializes the Cache client
func initCache() *redis.Client {
	cache := pkg.NewRedisClient()
	if cache == nil {
		log.Fatal("Failed to initialize Cache client")
	}

	log.Println("Cache client initialized successfully")
	return cache
}
