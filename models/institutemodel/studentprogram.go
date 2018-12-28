package institutemodel

type StudentProgram struct {
	StudentID uint `json:"student_id" gorm:"primary_key"`
	ProgramID uint `json:"program_id" gorm:"primary_key"`
}
