package models

import "time"

type Course struct {
	ID                      uint      `json:"id" gorm:"primary_key"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Price                   float32   `json:"price"`
	ResolutionAuthorization string    `json:"resolution_authorization"`
	DateExam                time.Time `json:"date_exam"`

	CourseStudents []CourseStudent `json:"course_students"`
}
