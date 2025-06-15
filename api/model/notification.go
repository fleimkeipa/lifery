package model

import "time"

type NotificationStatus int

const (
	NotificationStatusUnread NotificationStatus = 100 + iota
	NotificationStatusRead
)

type Notification struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Type      string             `json:"type"`
	Message   string             `json:"message"`
	Read      NotificationStatus `json:"read"`
	CreatedAt time.Time          `json:"created_at"`
}

type NotificationList struct {
	Notifications []Notification `json:"notifications"`
	Total         int            `json:"total"`
	PaginationOpts
}

type NotificationCreateInput struct {
	UserID  string `json:"user_id" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type NotificationUpdateInput struct {
	Read NotificationStatus `json:"read"`
}

type NotificationFindOpts struct {
	OrderByOpts
	UserID Filter
	Read   Filter
	FieldsOpts
	PaginationOpts
}
