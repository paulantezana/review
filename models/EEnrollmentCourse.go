package models

import "time"

type EnrollmentCourse struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EnrollmentID uint `json:"enrollment_id"`
	CourseID     uint `json:"course_id"`
}
