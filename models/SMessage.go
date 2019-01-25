package models

import "time"

type Message struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	Subject          string    `json:"subject"`
	Body             string    `json:"body"`
	Date             time.Time `json:"date"`
	ExpiryDate       uint      `json:"expiry_date"`
	IsReminder       bool      `json:"is_reminder"`
	NextReminderDate time.Time `json:"next_reminder_date"`

	CreatorID           uint `json:"creator_id"`
	ParentID     uint `json:"parent_id"`
	ReminderFrequencyID uint `json:"reminder_frequency_id"`
}
