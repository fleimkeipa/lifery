package repositories

import (
	"time"
)

type event struct {
	Date       time.Time   `json:"date"`
	TimeStart  time.Time   `json:"time_start"`
	TimeEnd    time.Time   `json:"time_end"`
	Name       string      `json:"name"`
	Items      []eventItem `json:"items"`
	ID         int         `json:"id" pg:",pk"`
	OwnerID    int         `json:"owner_id" pg:",notnull"`
	Owner      *user       `json:"owner" pg:"rel:has-one"`
	Visibility int         `json:"visibility"`
}

type eventItem struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}
