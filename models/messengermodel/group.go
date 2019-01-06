package messengermodel

import "time"

type Group struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	IsActive bool      `json:"is_active"`
}
