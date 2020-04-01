package models

import "time"

type UserProgram struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	UserID           uint `json:"user_id"`
	ProgramID        uint `json:"program_id"`
	License          bool `json:"license"`
	UserSubsidiaryID uint `json:"user_subsidiary_id"`
}
