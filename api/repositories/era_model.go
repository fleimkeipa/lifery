package repositories

import "time"

type era struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Name      string    `json:"name"`
	OwnerID   int       `json:"owner_id"`
	ID        int       `json:"id" pg:",pk"`
}
