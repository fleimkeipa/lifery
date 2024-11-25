package repositories

type Connect struct {
	ID       int           `json:"id" pg:",pk"`
	Status   RequestStatus `json:"status"`
	UserID   string        `json:"user_id"`
	FriendID string        `json:"friend_id"`
}

type RequestStatus string

const (
	RequestStatusApproved RequestStatus = "approved"
	RequestStatusRejected RequestStatus = "rejected"
	RequestStatusPending  RequestStatus = "pending"
)
