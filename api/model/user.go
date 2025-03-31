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
	RoleID    uint       `json:"role_id"`
}

type UserConnects struct {
	Connects []Connect `json:"connects"`
	Total    int       `json:"total"`
	PaginationOpts
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	PaginationOpts
}

type UserCreateRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
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
