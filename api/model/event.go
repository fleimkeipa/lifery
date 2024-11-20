package model

import "time"

type Event struct {
	ID        int64       `json:"id" pg:",pk"`
	Name      string      `json:"name"`
	Items     []EventItem `json:"items"`
	Date      time.Time   `json:"date"`
	TimeStart time.Time   `json:"time_start"`
	TimeEnd   time.Time   `json:"time_end"`
}

type EventType int

const (
	EventTypeString EventType = iota + 1
	EventTypePhoto
	EventTypeVideo
	EventTypeVoiceRecord
)

type EventItem struct {
	Type EventType `json:"type"`
	Data string    `json:"data"`
}

type EventCreateRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	Date        time.Time   `json:"date"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
}

type EventUpdateRequest struct {
	ID          int64       `json:"id" pg:",pk"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	Date        time.Time   `json:"date"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
}

type EventList struct {
	Events []Event `json:"events"`
	Total  int     `json:"total"`
	PaginationOpts
}

type EventFindOpts struct {
	SuplierID Filter
	PaginationOpts
}
