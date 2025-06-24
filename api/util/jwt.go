package util

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// retrieve JWT key from .env file
var privateKey = []byte(os.Getenv("JWT_KEY"))

// GenerateJWT generate JWT token
func GenerateJWT(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.RoleID,
		"iat":      time.Now().Unix(),
		"eat":      time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	if claims["eat"].(float64) < float64(time.Now().Unix()) {
		return errors.New("token expired")
	}

	if claims == nil {
		return errors.New("invalid token claims. claims is nil")
	}

	return nil
}

// validate Admin role
func ValidateAdminRoleJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userRole := model.UserRole(claims["role"].(float64))
	if userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid admin token provided")
}

// validate Editor role
func ValidateEditorRoleJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userRole := model.UserRole(claims["role"].(float64))
	if userRole == model.EditorRole || userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid editor or admin token provided")
}

// validate Viewer role
func ValidateViewerRoleJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userRole := model.UserRole(claims["role"].(float64))
	if userRole == model.EditorRole || userRole == model.ViewerRole || userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid editor, viewer or admin token provided")
}

// GetUserIDOnToken return user id
func GetUserIDOnToken(c echo.Context) (string, error) {
	token, err := getToken(c)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("invalid id claims")
	}

	return id, nil
}

// GetOwnerFromToken returns the owner details from the JWT token
func GetOwnerFromToken(c echo.Context) (model.TokenOwner, error) {
	token, err := getToken(c)
	if err != nil {
		return model.TokenOwner{}, err
	}

	if !token.Valid {
		return model.TokenOwner{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid token claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid id claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid username claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid email claims")
	}

	role, ok := claims["role"].(float64)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid role claims")
	}

	return model.TokenOwner{
		ID:       id,
		Username: username,
		Email:    email,
		RoleID:   model.UserRole(role),
	}, nil
}

// GetOwnerFromCtx returns the owner details from the context
func GetOwnerFromCtx(ctx context.Context) model.TokenOwner {
	owner, ok := ctx.Value("user").(model.TokenOwner)
	if ok {
		return owner
	}

	return model.TokenOwner{}
}

// GetOwnerIDFromCtx returns the owner id from the context string type
func GetOwnerIDFromCtx(ctx context.Context) string {
	owner, ok := ctx.Value("user").(model.TokenOwner)
	if ok {
		return owner.ID
	}

	return ""
}

// check token validity
func getToken(context echo.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})

	return token, err
}

// extract token from request Authorization header
func getTokenFromRequest(c echo.Context) string {
	bearerToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func IsUserPublic(c echo.Context) bool {
	token, err := getToken(c)
	if err != nil {
		return true
	}

	if !token.Valid {
		return true
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}

	userRole := model.UserRole(claims["role"].(float64))
	if userRole == model.EditorRole || userRole == model.ViewerRole || userRole == model.AdminRole {
		return false
	}

	return true
}

// GenerateResetToken generates a JWT token for password reset
func GenerateResetToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"type":     "password_reset",
		"iat":      time.Now().Unix(),
		"eat":      time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiry
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", pkg.NewError(err, "failed to generate reset token", http.StatusInternalServerError)
	}

	return tokenString, nil
}

// ValidateResetToken validates a password reset token and returns the user info
func ValidateResetToken(tokenString string) (*model.User, error) {
	jwtParser := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkg.NewError(nil, "unexpected signing method: "+token.Header["alg"].(string), http.StatusBadRequest)
		}
		return privateKey, nil
	}

	token, err := jwt.Parse(tokenString, jwtParser)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse reset token", http.StatusInternalServerError)
	}

	if !token.Valid {
		return nil, pkg.NewError(nil, "invalid token", http.StatusBadRequest)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, pkg.NewError(nil, "invalid token claims", http.StatusBadRequest)
	}

	// Check if token is expired
	if claims["eat"].(float64) < float64(time.Now().Unix()) {
		return nil, pkg.NewError(nil, "token expired", http.StatusBadRequest)
	}

	// Check if token is for password reset
	if claims["type"] != "password_reset" {
		return nil, pkg.NewError(nil, "invalid token type", http.StatusBadRequest)
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, pkg.NewError(nil, "invalid id claims", http.StatusBadRequest)
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, pkg.NewError(nil, "invalid username claims", http.StatusBadRequest)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, pkg.NewError(nil, "invalid email claims", http.StatusBadRequest)
	}

	return &model.User{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

func GenerateRandomPassword() string {
	// Generate a random 32-character password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, 32)

	// Generate random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to a simple random generation if crypto/rand fails
		for i := range b {
			b[i] = charset[i%len(charset)]
		}
		return string(b)
	}

	// Map random bytes to charset characters
	for i := range b {
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(b)
}
