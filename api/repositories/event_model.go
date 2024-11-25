package repositories

import (
	"database/sql"
	"time"
)

type Event struct {
	Date      time.Time    `json:"date"`
	TimeStart time.Time    `json:"time_start"`
	TimeEnd   time.Time    `json:"time_end"`
	Name      string       `json:"name"`
	Items     []EventItem  `json:"items"`
	ID        int          `json:"id" pg:",pk"`
	OwnerID   int          `json:"owner_id"`
	Private   sql.NullBool `json:"private"`
}

type EventType int

const (
	EventTypeString EventType = iota + 1
	EventTypePhoto
	EventTypeVideo
	EventTypeVoiceRecord
)

type EventItem struct {
	Data string    `json:"data"`
	Type EventType `json:"type"`
}
