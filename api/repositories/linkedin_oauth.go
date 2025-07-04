package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"

	"golang.org/x/oauth2"
)

type LinkedInOAuthRepository struct {
	config *oauth2.Config
}

func NewLinkedInOAuthRepository() *LinkedInOAuthRepository {
	config := &oauth2.Config{
		ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
		ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("LINKEDIN_REDIRECT_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
			TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
		},
	}
	return &LinkedInOAuthRepository{config: config}
}

func (l *LinkedInOAuthRepository) GetAuthURL() string {
	return l.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (l *LinkedInOAuthRepository) GetUserInfo(ctx context.Context, code string) (*model.OAuthUserInfo, error) {
	token, err := l.config.Exchange(ctx, code)
	if err != nil {
		return nil, pkg.NewError(err, "failed to exchange code for token", http.StatusInternalServerError)
	}

	client := l.config.Client(ctx, token)

	// Get user info from /v2/userinfo endpoint
	userInfoReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.linkedin.com/v2/userinfo", nil)
	userInfoReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	userInfoReq.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	userInfoResp, err := client.Do(userInfoReq)
	if err != nil {
		return nil, pkg.NewError(err, "userinfo request failed", http.StatusInternalServerError)
	}
	defer userInfoResp.Body.Close()

	// Check response status
	if userInfoResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(userInfoResp.Body)
		return nil, pkg.NewError(nil, fmt.Sprintf("userinfo API returned status %d: %s, body: %s", userInfoResp.StatusCode, userInfoResp.Status, string(bodyBytes)), http.StatusInternalServerError)
	}

	var userInfoData struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfoData); err != nil {
		return nil, pkg.NewError(err, "failed to decode userinfo response", http.StatusInternalServerError)
	}

	return &model.OAuthUserInfo{
		ID:         userInfoData.Sub,
		Email:      userInfoData.Email,
		Name:       userInfoData.Name,
		GivenName:  userInfoData.GivenName,
		FamilyName: userInfoData.FamilyName,
	}, nil
}
