package models

import "time"

type Payment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Amount    float32   `json:"amount"`
	Reason    string    `json:"reason"`
	Category  string    `json:"category"` // admission || enrollment
}
