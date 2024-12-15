package repositories

import "time"

type user struct {
	DeletedAt time.Time `json:"deleted_at" pg:",soft_delete"`
	CreatedAt time.Time `json:"created_at"`
	Connects  []int     `json:"connects"`
	Username  string    `json:"username" pg:",unique"`
	Email     string    `json:"email" pg:",unique"`
	Password  string    `json:"password"`
	ID        int       `json:"id" pg:",pk"`
	RoleID    uint      `json:"role_id"`
}

type userConnects struct {
	Connects []user `json:"connects"`
}
