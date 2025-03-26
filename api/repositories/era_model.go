package repositories

import (
	"time"
)

type era struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    int       `json:"user_id" pg:",notnull,on_delete:CASCADE"`
	User      *user     `json:"user" pg:"rel:has-one,fk:user_id"`
	ID        int       `json:"id" pg:",pk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
