package model

import "time"

const ZeroCreds = "zeroCreds"

type AuthType string

const (
	AuthTypeEmail    AuthType = "email"
	AuthTypeGoogle   AuthType = "google"
	AuthTypeLinkedIn AuthType = "linkedin"
)

type User struct {
	CreatedAt time.Time  `json:"created_at"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	ID        string     `json:"id"`
	Connects  []*Connect `json:"connects"`
	RoleID    UserRole   `json:"role_id"`
	AuthType  AuthType   `json:"auth_type"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	PaginationOpts
}

type UserCreateInput struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required, email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	AuthType        string `json:"auth_type"`
}

type UserFindOpts struct {
	OrderByOpts
	Username Filter
	Email    Filter
	RoleID   Filter
	FieldsOpts
	PaginationOpts
}

type UserConnectsFindOpts struct {
	OrderByOpts
	UserID Filter
	FieldsOpts
	PaginationOpts
}

type UpdateUsernameRequest struct {
	Username string `json:"username" validate:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}
