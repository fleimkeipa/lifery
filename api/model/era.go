package model

import "time"

type Era struct {
	TimeRange time.Time `json:"time_range"`
	Name      string    `json:"name"`
	ID        int64     `json:"id" pg:",pk"`
}

type EraCreateRequest struct {
	TimeRange time.Time `json:"time_range"`
	Name      string    `json:"name"`
}

type EraUpdateRequest struct {
	TimeRange time.Time `json:"time_range"`
	Name      string    `json:"name"`
	ID        int64     `json:"id" pg:",pk"`
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
