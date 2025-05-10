package repositories

import "time"

type user struct {
	CreatedAt time.Time  `json:"created_at"`
	Username  string     `json:"username" pg:",unique"`
	Email     string     `json:"email" pg:",unique"`
	Password  string     `json:"password"`
	Connects  []*connect `json:"connects" pg:"rel:has_many,on_delete:CASCADE"`
	Eras      []*era     `json:"eras" pg:"rel:has_many,on_delete:CASCADE"`
	ID        int        `json:"id" pg:",pk"`
	RoleID    UserRole   `json:"role_id"`
}

type userConnects struct {
	Connects []user `json:"connects"`
}
