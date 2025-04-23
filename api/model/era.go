package model

import "time"

type Era struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    string    `json:"user_id"`
	ID        string    `json:"id"`
}

type EraCreateInput struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Color     string    `json:"color" validate:"required,iscolor"`
	Name      string    `json:"name"`
}

type EraUpdateInput struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Color     string    `json:"color" validate:"required,iscolor"`
	Name      string    `json:"name"`
}

type EraList struct {
	Eras  []Era `json:"eras"`
	Total int   `json:"total"`
	PaginationOpts
}

type EraFindOpts struct {
	OrderByOpts
	Name   Filter
	UserID Filter
	PaginationOpts
}
