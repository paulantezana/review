package models

import "time"

type LanguageCourseExam struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Note      float32   `json:"note"`
	Date      time.Time `json:"date"`

	CourseStudentID uint `json:"course_student_id"`
}
