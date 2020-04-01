package models

import "time"

type UserSubsidiary struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	UserID       uint          `json:"user_id"`
	SubsidiaryID uint          `json:"subsidiary_id"`
	License      bool          `json:"license"`
	UserPrograms []UserProgram `json:"user_programs"`
}
