package repositories

type notification struct {
	ID        int    `json:"id" pg:",pk"`
	UserID    int    `json:"user_id" pg:",notnull"`
	Type      string `json:"type" pg:",notnull"`
	Message   string `json:"message" pg:",notnull"`
	Read      int    `json:"read" pg:",notnull"`
	CreatedAt string `json:"created_at" pg:",notnull"`
}
