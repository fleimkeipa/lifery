package model

import "time"

type Era struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    string    `json:"user_id"`
	ID        string    `json:"id"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EraCreateRequest struct {
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Color     string `json:"color"`
	Name      string `json:"name"`
}

type EraUpdateRequest struct {
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	ID        int64  `json:"id"`
}

type EraList struct {
	Eras  []Era `json:"eras"`
	Total int   `json:"total"`
	PaginationOpts
}

type EraFindOpts struct {
	Name   Filter
	UserID Filter
	OrderByOpts
	PaginationOpts
}
