package models

import "time"

type GroupMessageRecipient struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsRead    bool      `json:"is_read"`

	RecipientID      uint `json:"recipient_id"`
	RecipientGroupID uint `json:"recipient_group_id"`
	MessageID        uint `json:"message_id"`
}
