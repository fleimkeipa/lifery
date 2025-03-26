package repositories

import "time"

type user struct {
	CreatedAt time.Time  `json:"created_at"`
	Connects  []*connect `json:"connects" pg:"rel:has_many,on_delete:CASCADE"`
	Eras      []*era     `json:"eras" pg:"rel:has_many,on_delete:CASCADE"`
	Username  string     `json:"username" pg:",unique"`
	Email     string     `json:"email" pg:",unique"`
	Password  string     `json:"password"`
	ID        int        `json:"id" pg:",pk"`
	RoleID    uint       `json:"role_id"`
}

type userConnects struct {
	Connects []user `json:"connects"`
}
