package model

import "time"

type Event struct {
	Date        time.Time   `json:"date"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Visibility  Visibility  `json:"visibility"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
}

type Visibility int

const (
	EventVisibilityPublic Visibility = iota + 1
	EventVisibilityPrivate
	EventVisibilityJustMe
)

type EventType int

const (
	EventTypeString EventType = iota + 10
	EventTypePhoto
	EventTypeVideo
	EventTypeVoiceRecord
)

type EventItem struct {
	Data string    `json:"data"`
	Type EventType `json:"type"`
}

type EventCreateRequest struct {
	Date        string      `json:"date"`
	TimeStart   string      `json:"time_start"`
	TimeEnd     string      `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	Visibility  Visibility  `json:"visibility"`
}

type EventUpdateRequest struct {
	Date        string      `json:"date"`
	TimeStart   string      `json:"time_start"`
	TimeEnd     string      `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	ID          int64       `json:"id"`
	Visibility  Visibility  `json:"visibility"`
}

type EventList struct {
	Events []Event `json:"events"`
	Total  int     `json:"total"`
	PaginationOpts
}

type EventFindOpts struct {
	UserID     Filter
	Visibility Filter
	PaginationOpts
	OrderByOpts
}
