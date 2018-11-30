package coursemodel

import "time"

type Course struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	
	CourseStudents []CourseStudent `json:"course_students"`
}
