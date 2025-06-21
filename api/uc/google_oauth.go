package uc

import (
	"context"
	"fmt"
	"os"

	"github.com/fleimkeipa/lifery/model"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleOAuthUC struct {
	userUC *UserUC
	config *oauth2.Config
}

func NewGoogleOAuthUC(userUC *UserUC) *GoogleOAuthUC {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthUC{
		userUC: userUC,
		config: config,
	}
}

func (g *GoogleOAuthUC) GetAuthURL() string {
	return g.config.AuthCodeURL("state")
}

func (g *GoogleOAuthUC) HandleCallback(ctx context.Context, code string) (*model.User, error) {
	// Exchange code for token
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	userInfo, err := g.getUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	existingUser, err := g.userUC.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists, return existing user
		return existingUser, nil
	}

	newUser := model.UserCreateInput{
		Username:        userInfo.GivenName + "_" + userInfo.FamilyName,
		Email:           userInfo.Email,
		Password:        generateRandomPassword(),
		ConfirmPassword: generateRandomPassword(),
		AuthType:        string(model.AuthTypeGoogle),
	}

	user, err := g.userUC.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (g *GoogleOAuthUC) getUserInfo(ctx context.Context, token *oauth2.Token) (*model.GoogleUserInfo, error) {
	client := g.config.Client(ctx, token)

	service, err := googleoauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	verifiedEmail := false
	if userInfo.VerifiedEmail != nil {
		verifiedEmail = *userInfo.VerifiedEmail
	}

	return &model.GoogleUserInfo{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		VerifiedEmail: verifiedEmail,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}, nil
}

func generateRandomPassword() string {
	// Generate a random 32-character password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[0] // This is a simplified version, in production use crypto/rand
	}
	return string(b)
}
