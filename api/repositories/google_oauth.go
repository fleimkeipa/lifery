package repositories

import (
	"context"
	"net/http"
	"os"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
)

type GoogleOAuthRepository struct {
	config *oauth2.Config
}

func NewGoogleOAuthRepository() *GoogleOAuthRepository {
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
	return &GoogleOAuthRepository{config: config}
}

func (l *GoogleOAuthRepository) GetAuthURL() string {
	return l.config.AuthCodeURL("state")
}

func (l *GoogleOAuthRepository) GetUserInfo(ctx context.Context, code string) (*model.OAuthUserInfo, error) {
	token, err := l.config.Exchange(ctx, code)
	if err != nil {
		return nil, pkg.NewError(err, "failed to exchange code for token", http.StatusInternalServerError)
	}

	client := l.config.Client(ctx, token)

	service, err := googleoauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, pkg.NewError(err, "failed to create oauth2 service", http.StatusInternalServerError)
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, pkg.NewError(err, "failed to get user info", http.StatusInternalServerError)
	}

	return &model.OAuthUserInfo{
		ID:         userInfo.Id,
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		GivenName:  userInfo.GivenName,
		FamilyName: userInfo.FamilyName,
	}, nil
}
