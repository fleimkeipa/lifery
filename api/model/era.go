package model

import "time"

type Era struct {
	ID        int64     `json:"id" pg:",pk"`
	Name      string    `json:"name"`
	TimeRange time.Time `json:"time_range"`
}

type EraCreateRequest struct {
	Name      string    `json:"name"`
	TimeRange time.Time `json:"time_range"`
}

type EraUpdateRequest struct {
	ID        int64     `json:"id" pg:",pk"`
	Name      string    `json:"name"`
	TimeRange time.Time `json:"time_range"`
}

type EraList struct {
	Eras  []Era `json:"eras"`
	Total int   `json:"total"`
	PaginationOpts
}

type EraFindOpts struct {
	Name Filter
	PaginationOpts
}
