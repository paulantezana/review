package models

import "time"

type StudentHistory struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	StudentID   uint      `json:"student_id"`
	UserID      uint      `json:"user_id"`
	Type        uint      `json:"type"` // 1 = create || update  // 2 = delete || null // 3 = print || view
	Date        time.Time `json:"date"`
}
