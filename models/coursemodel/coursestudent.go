package coursemodel

type CourseStudent struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	DNI      string  `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName string  `json:"full_name" gorm:"type:varchar(250)"`
	Phone    string  `json:"phone" gorm:"type:varchar(32)"`
	State    bool    `json:"state" gorm:"default:'true'"`
	Gender   string  `json:"gender"`
	Year     uint    `json:"year"`
	Payment  float32 `json:"payment"`

	StudentID uint `json:"student_id"`

	CourseID  uint `json:"course_id"`
	ProgramID uint `json:"program_id"`

	CourseExams []CourseExam `json:"course_exams"`
}
