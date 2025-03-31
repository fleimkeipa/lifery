package model

type Connect struct {
	ID       string        `json:"id"`
	UserID   string        `json:"user_id"`
	FriendID string        `json:"friend_id"`
	User     User          `json:"user"`
	Friend   User          `json:"friend"`
	Status   RequestStatus `json:"status"`
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
	OrderByOpts
	Status Filter
	UserID Filter
	FieldsOpts
	PaginationOpts
}
