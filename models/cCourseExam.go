package models

import "time"

type CourseExam struct {
	ID   uint      `json:"id" gorm:"primary_key"`
	Note float32   `json:"note"`
	Date time.Time `json:"date"`

	CourseStudentID uint `json:"course_student_id"`
}
