package uc

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/repositories"
	"github.com/fleimkeipa/lifery/util"
)

type OAuthUC struct {
	googleOAuthRepository   *repositories.GoogleOAuthRepository
	linkedinOAuthRepository *repositories.LinkedInOAuthRepository
	userUC                  *UserUC
}

func NewOAuthUC(googleOAuthRepository *repositories.GoogleOAuthRepository, linkedinOAuthRepository *repositories.LinkedInOAuthRepository, userUC *UserUC) *OAuthUC {
	return &OAuthUC{
		googleOAuthRepository:   googleOAuthRepository,
		linkedinOAuthRepository: linkedinOAuthRepository,
		userUC:                  userUC,
	}
}

func (o *OAuthUC) GetAuthURL(provider model.OAuthProvider) string {
	switch provider {
	case model.GoogleProvider:
		return o.googleOAuthRepository.GetAuthURL()
	case model.LinkedInProvider:
		return o.linkedinOAuthRepository.GetAuthURL()
	default:
		return ""
	}
}

func (o *OAuthUC) HandleCallback(ctx context.Context, provider model.OAuthProvider, code string) (*model.User, error) {
	var userInfo *model.OAuthUserInfo
	var err error

	switch provider {
	case model.GoogleProvider:
		userInfo, err = o.googleOAuthRepository.GetUserInfo(ctx, code)
		if err != nil {
			return nil, err
		}
	case model.LinkedInProvider:
		userInfo, err = o.linkedinOAuthRepository.GetUserInfo(ctx, code)
		if err != nil {
			return nil, err
		}
	}

	existingUser, err := o.userUC.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		return existingUser, nil
	}

	username := util.GenerateUsername(userInfo.GivenName, userInfo.FamilyName)
	password := util.GenerateRandomPassword()

	newUser := model.UserCreateInput{
		Username:        username,
		Email:           userInfo.Email,
		Password:        password,
		ConfirmPassword: password,
		AuthType:        string(provider),
	}

	user, err := o.userUC.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
