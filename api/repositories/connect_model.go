package repositories

type connect struct {
	User     *user `json:"user" pg:"rel:has-one,on_delete:CASCADE"`
	Friend   *user `json:"friend" pg:"rel:has-one,on_delete:CASCADE"`
	ID       int   `json:"id" pg:",pk"`
	Status   int   `json:"status"`
	UserID   int   `json:"user_id" pg:",notnull"`
	FriendID int   `json:"friend_id" pg:",notnull"`
}
