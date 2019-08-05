package models

import "time"

type LanguageCourseStudent struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DNI       string    `json:"dni" gorm:" type:varchar(15)"`
	FullName  string    `json:"full_name" gorm:"type:varchar(250)"`
	Phone     string    `json:"phone" gorm:"type:varchar(32)"`
	State     bool      `json:"state" gorm:"default:'true'"`
	Gender    string    `json:"gender"`
	Year      uint      `json:"year"`
	Payment   float32   `json:"payment"`
	Note      float32   `json:"note"`

	StudentID uint `json:"student_id"`

	CourseID  uint `json:"course_id"`
	ProgramID uint `json:"program_id"`

	CourseExams []LanguageCourseExam `json:"course_exams"`
}
