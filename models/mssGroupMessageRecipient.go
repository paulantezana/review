package models

import "time"

type MssGroupMessageRecipient struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsRead           bool `json:"is_read"`
	RecipientID      uint `json:"recipient_id"`
	RecipientGroupID uint `json:"recipient_group_id"`
    MssGroupMessageID        uint `json:"mss_group_message_id"`
}
