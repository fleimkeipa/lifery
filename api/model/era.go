package model

type Era struct {
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	OwnerID   string `json:"owner_id"`
	ID        string `json:"id"`
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
	PaginationOpts
}
