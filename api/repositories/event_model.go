package repositories

import "time"

type event struct {
	Date        string      `json:"date"`
	TimeStart   string      `json:"time_start"`
	TimeEnd     string      `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []eventItem `json:"items"`
	ID          int         `json:"id" pg:",pk"`
	UserID      int         `json:"user_id" pg:",notnull"`
	User        *user       `json:"user" pg:"rel:has-one"`
	Visibility  int         `json:"visibility"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at,omitempty" pg:",soft_delete"`
}

type eventItem struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}
