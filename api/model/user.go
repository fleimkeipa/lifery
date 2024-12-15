package model

import "time"

const ZeroCreds = "zeroCreds"

type User struct {
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	Connects  []int     `json:"connects"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	ID        string    `json:"id"`
	RoleID    uint      `json:"role_id"`
}

type UserConnects struct {
	Connects []User `json:"connects"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	PaginationOpts
}

type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type UserFindOpts struct {
	Username Filter
	Email    Filter
	RoleID   Filter
	FieldsOpts
	PaginationOpts
}
