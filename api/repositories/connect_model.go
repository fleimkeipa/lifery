package repositories

type connect struct {
	ID       int           `json:"id" pg:",pk"`
	Status   requestStatus `json:"status"`
	UserID   int           `json:"user_id"`
	FriendID int           `json:"friend_id"`
}

type requestStatus string

const (
	requestStatusApproved requestStatus = "approved"
	requestStatusRejected requestStatus = "rejected"
	requestStatusPending  requestStatus = "pending"
)
