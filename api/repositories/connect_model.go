package repositories

type connect struct {
	ID       int   `json:"id" pg:",pk"`
	Status   int   `json:"status"`
	UserID   int   `json:"user_id" pg:",notnull"`
	User     *user `json:"user" pg:"rel:has-one,on_delete:CASCADE"`
	FriendID int   `json:"friend_id" pg:",notnull"`
	Friend   *user `json:"friend" pg:"rel:has-one,on_delete:CASCADE"`
}
