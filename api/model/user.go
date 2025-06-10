package model

import "time"

const ZeroCreds = "zeroCreds"

type User struct {
	CreatedAt time.Time  `json:"created_at"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	ID        string     `json:"id"`
	Connects  []*Connect `json:"connects"`
	RoleID    UserRole   `json:"role_id"`
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
