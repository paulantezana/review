package models

import "time"

type MssMessageRecipient struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsRead       bool `json:"is_read"`
	RecipientID  uint `json:"recipient_id"`
	MssMessageID uint `json:"mss_message_id"`
}
