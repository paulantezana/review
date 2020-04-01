package models

import "time"

type MssGroupMessage struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Subject             string    `json:"subject"`
	Body                string    `json:"body"`
	BodyType            uint8     `json:"body_type"` // 0 = plain string || 1 == file
	FilePath            string    `json:"file_path"`
	ExpiryDate          uint      `json:"expiry_date"`
	IsReminder          bool      `json:"is_reminder"`
	NextReminderDate    time.Time `json:"next_reminder_date"`
	CreatorID           uint      `json:"creator_id"`
	//ReminderFrequencyID uint      `json:"reminder_frequency_id"`
}
