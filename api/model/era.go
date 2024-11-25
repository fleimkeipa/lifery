package model

import "time"

type Era struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	ID        string    `json:"id"`
}

type EraCreateRequest struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
}

type EraUpdateRequest struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
}

type EraList struct {
	Eras  []Era `json:"eras"`
	Total int   `json:"total"`
	PaginationOpts
}

type EraFindOpts struct {
	Name   Filter
	UserID Filter
	PaginationOpts
}
