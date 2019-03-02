package models

import "time"

type StudentStatus struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

// Set User's table name to be `profiles`
func (StudentStatus) TableName() string {
	return "student_status"
}
