package interfaces

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
)

type OAuthRepository interface {
	GetAuthURL() string
	GetUserInfo(ctx context.Context, code string) (*model.OAuthUserInfo, error)
}
