package models

import "time"

type TeacherAction struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Action      string    `json:"action"`
	Description string    `json:"description"`

	TeacherID uint `json:"teacher_id"`
}
