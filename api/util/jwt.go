package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fleimkeipa/lifery/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// retrieve JWT key from .env file
var privateKey = []byte(os.Getenv("JWT_KEY"))

// GenerateJWT generate JWT token
func GenerateJWT(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       strconv.Itoa(int(user.ID)),
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

	userRole := uint(claims["role"].(float64))
	if userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid admin token provided")
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

	userRole := uint(claims["role"].(float64))
	if userRole == model.ViewerRole || userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid viewer or admin token provided")
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
		RoleID:   uint(role),
	}, nil
}

// GetOwnerFromCtx returns the owner details from the context
func GetOwnerFromCtx(ctx context.Context) *model.TokenOwner {
	owner, ok := ctx.Value("user").(model.TokenOwner)
	if ok {
		return &owner
	}

	return nil
}

// GetOwnerIDFromCtx returns the owner id from the context int type
func GetOwnerIDFromCtx(ctx context.Context) int64 {
	owner, ok := ctx.Value("user").(model.TokenOwner)
	if ok {
		intID, _ := strconv.Atoi(owner.ID)
		return int64(intID)
	}

	return 0
}

// GetStrOwnerIDFromCtx returns the owner id from the context string type
func GetStrOwnerIDFromCtx(ctx context.Context) string {
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
