package models

import "time"

// Category struct
type PostCategory struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	Name      string `json:"name"`
	ParentID  uint   `json:"parent_id"`
	State     bool   `json:"state" gorm:"default:'true'"`
	ProgramId uint   `json:"program_id"`

	Posts []Post `json:"posts"`
}
