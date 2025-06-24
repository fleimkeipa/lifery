package interfaces

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
	"golang.org/x/oauth2"
)

type OAuthRepository interface {
	GetAuthURL() string
	ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, code string) (*model.OAuthUserInfo, error)
}
