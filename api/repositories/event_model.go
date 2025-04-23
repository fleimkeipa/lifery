package repositories

import "time"

type event struct {
	CreatedAt   time.Time   `json:"created_at"`
	DeletedAt   time.Time   `json:"deleted_at,omitempty" pg:",soft_delete"`
	UpdatedAt   time.Time   `json:"updated_at"`
	User        *user       `json:"user" pg:"rel:has-one"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Date        time.Time   `json:"date"`
	Items       []eventItem `json:"items"`
	ID          int         `json:"id" pg:",pk"`
	Visibility  int         `json:"visibility"`
	UserID      int         `json:"user_id" pg:",notnull"`
}

type eventItem struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}
