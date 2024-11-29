package repositories

type connect struct {
	ID       int           `json:"id" pg:",pk"`
	Status   requestStatus `json:"status"`
	UserID   int           `json:"user_id" pg:",notnull"`
	User     *user         `json:"user" pg:"rel:has-one"`
	FriendID int           `json:"friend_id" pg:",notnull"`
	Friend   *user         `json:"friend" pg:"rel:has-one"`
}

type requestStatus string

const (
	requestStatusApproved requestStatus = "approved"
	requestStatusRejected requestStatus = "rejected"
	requestStatusPending  requestStatus = "pending"
)
