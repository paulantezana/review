package models

import "time"

type ProgramUser struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	ProgramID uint      `json:"program_id"`
	License   bool      `json:"license"`

	SubsidiaryUserID uint `json:"subsidiary_user_id"`
}
