package repositories

import "time"

type era struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *user     `json:"user" pg:"rel:has-one,fk:user_id"`
	TimeStart string    `json:"time_start"`
	TimeEnd   string    `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    int       `json:"user_id" pg:",notnull,on_delete:CASCADE"`
	ID        int       `json:"id" pg:",pk"`
}
