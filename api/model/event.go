package model

type Event struct {
	Date        string      `json:"date"`
	TimeStart   string      `json:"time_start"`
	TimeEnd     string      `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []EventItem `json:"items"`
	ID          string      `json:"id"`
	OwnerID     string      `json:"owner_id"`
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
}
