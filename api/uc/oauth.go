package uc

import (
	"context"
	"fmt"

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

	// Check if user exists
	existingUser, err := o.userUC.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists, return existing user
		return existingUser, nil
	}

	// User does not exist, create new user
	newUser := model.UserCreateInput{
		Username: userInfo.Name,
		Email:    userInfo.Email,
		Password: util.GenerateRandomPassword(),
	}

	user, err := o.userUC.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
