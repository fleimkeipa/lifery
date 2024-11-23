package model

type Connect struct {
	ID       int64         `json:"id" pg:",pk"`
	Status   RequestStatus `json:"status"`
	UserID   string        `json:"user_id"`
	FriendID string        `json:"friend_id"`
}

type ConnectList struct {
	Connects []Connect `json:"conects"`
	Total    int       `json:"total"`
	PaginationOpts
}

type ConnectCreateRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FriendID string `json:"friend_id" binding:"required"`
}

type RequestStatus string

const (
	RequestStatusApproved RequestStatus = "approved"
	RequestStatusRejected RequestStatus = "rejected"
	RequestStatusPending  RequestStatus = "pending"
)

type ConnectUpdateRequest struct {
	Status RequestStatus `json:"status" binding:"required"`
}

type DisconnectRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
}

type ConnectFindOpts struct {
	Status   Filter
	FriendID Filter
	FieldsOpts
	PaginationOpts
}
