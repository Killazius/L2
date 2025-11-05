package domain

type Event struct {
	ID     int64  `json:"id"`
	Date   string `json:"date"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}
