package model

type Connect struct {
	ID       string        `json:"id"`
	Status   RequestStatus `json:"status"`
	UserID   string        `json:"user_id"`
	User     User          `json:"user"`
	FriendID string        `json:"friend_id"`
	Friend   User          `json:"friend"`
}

type ConnectList struct {
	Connects []Connect `json:"connects"`
	Total    int       `json:"total"`
	PaginationOpts
}

type ConnectCreateRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
}

type RequestStatus int

const (
	RequestStatusPending RequestStatus = 100 + iota
	RequestStatusApproved
	RequestStatusRejected
)

type ConnectUpdateRequest struct {
	Status RequestStatus `json:"status" binding:"required"`
}

type DisconnectRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
}

type ConnectFindOpts struct {
	Status Filter
	UserID Filter
	FieldsOpts
	OrderByOpts
	PaginationOpts
}
