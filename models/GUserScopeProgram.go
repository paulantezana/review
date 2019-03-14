package models

import "time"

type UserScopeProgram struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	License   bool      `json:"license"`

	AppModuleID   uint `json:"app_module_id"`
	ProgramUserID uint `json:"program_user_id"`
}
