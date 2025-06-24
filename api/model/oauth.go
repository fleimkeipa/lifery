package model

type OAuthProvider string

const (
	GoogleProvider   OAuthProvider = "google"
	LinkedInProvider OAuthProvider = "linkedin"
)

type OAuthUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type GoogleAuthRequest struct {
	Code string `json:"code" validate:"required"`
}

type LinkedInAuthRequest struct {
	Code string `json:"code" validate:"required"`
}
