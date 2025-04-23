package repositories

import "time"

type era struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *user     `json:"user" pg:"rel:has-one,fk:user_id"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    int       `json:"user_id" pg:",notnull,on_delete:CASCADE"`
	ID        int       `json:"id" pg:",pk"`
}
