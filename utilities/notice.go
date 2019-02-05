package utilities

import "time"

type Notice struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
	RecipientID uint      `json:"recipient_id"`
	Avatar      string    `json:"avatar"`
	Type        string    `json:"type"`
}
