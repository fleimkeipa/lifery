package repositories

import (
	"database/sql"
	"time"
)

type event struct {
	Date      time.Time    `json:"date"`
	TimeStart time.Time    `json:"time_start"`
	TimeEnd   time.Time    `json:"time_end"`
	Name      string       `json:"name"`
	Items     []eventItem  `json:"items"`
	ID        int          `json:"id" pg:",pk"`
	OwnerID   int          `json:"owner_id"`
	Private   sql.NullBool `json:"private"`
}

type eventType int

const (
	eventTypeString eventType = iota + 1
	eventTypePhoto
	eventTypeVideo
	eventTypeVoiceRecord
)

type eventItem struct {
	Data string    `json:"data"`
	Type eventType `json:"type"`
}
