package repositories

type era struct {
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	UserID    int    `json:"user_id" pg:",notnull"`
	User      *user  `json:"user" pg:"rel:has-one"`
	ID        int    `json:"id" pg:",pk"`
}

type eraGetResponse struct {
	TimeStart string    `json:"time_start"`
	TimeEnd   string    `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    int       `json:"user_id"`
	User      *miniUser `json:"user"`
	ID        int       `json:"id"`
	tableName struct{}  `pg:",discard_unknown_columns"`
}

type miniUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
	RoleID   uint   `json:"role_id"`
}
