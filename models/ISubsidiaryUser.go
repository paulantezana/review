package models

import "time"

type SubsidiaryUser struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       uint      `json:"user_id"`
	SubsidiaryID uint      `json:"subsidiary_id"`
	License      bool      `json:"license"`

	ProgramUsers []ProgramUser `json:"program_users, omitempty"`
}
