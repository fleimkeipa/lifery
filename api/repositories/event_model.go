package repositories

type event struct {
	Date        string      `json:"date"`
	TimeStart   string      `json:"time_start"`
	TimeEnd     string      `json:"time_end"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       []eventItem `json:"items"`
	ID          int         `json:"id" pg:",pk"`
	OwnerID     int         `json:"owner_id" pg:",notnull"`
	Owner       *user       `json:"owner" pg:"rel:has-one"`
	Visibility  int         `json:"visibility"`
}

type eventItem struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}
