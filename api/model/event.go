package model

import "time"

type Event struct {
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at"`
	Date        time.Time   `json:"date"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Items       []EventItem `json:"items"`
	Visibility  Visibility  `json:"visibility"`
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

type EventCreateInput struct {
	Date        time.Time   `json:"date"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	Visibility  Visibility  `json:"visibility"`
}

type EventUpdateInput struct {
	Date        time.Time   `json:"date"`
	TimeStart   time.Time   `json:"time_start"`
	TimeEnd     time.Time   `json:"time_end"`
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	Visibility  Visibility  `json:"visibility"`
}

type EventList struct {
	Events []Event `json:"events"`
	Total  int     `json:"total"`
	PaginationOpts
}

type EventFindOpts struct {
	OrderByOpts
	UserID     Filter
	Visibility Filter
	PaginationOpts
}
