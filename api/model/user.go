package model

import "time"

const ZeroCreds = "zeroCreds"

type User struct {
	DeletedAt time.Time `json:"deleted_at" pg:",soft_delete"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username" binding:"required" pg:",unique"`
	Email     string    `json:"email" binding:"required" pg:",unique"`
	Password  string    `json:"password" binding:"required"`
	ID        int64     `json:"id" pg:",pk"`
	RoleID    uint      `json:"role_id"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	PaginationOpts
}

type Owner struct {
	Username string `json:"username" pg:",unique"`
	Email    string `json:"email" pg:",unique"`
	ID       int64  `json:"id" pg:",pk"`
	RoleID   uint   `json:"role_id"`
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
