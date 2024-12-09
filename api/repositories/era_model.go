package repositories

type era struct {
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	OwnerID   int    `json:"owner_id" pg:",notnull"`
	Owner     *user  `json:"owner" pg:"rel:has-one"`
	ID        int    `json:"id" pg:",pk"`
}

type eraGetResponse struct {
	TimeStart string    `json:"time_start"`
	TimeEnd   string    `json:"time_end"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	OwnerID   int       `json:"owner_id"`
	Owner     *miniUser `json:"owner"`
	ID        int       `json:"id"`
	tableName struct{}  `pg:",discard_unknown_columns"`
}

type miniUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
	RoleID   uint   `json:"role_id"`
}
