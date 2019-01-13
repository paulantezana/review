package institutemodel

type StudentProgram struct {
	StudentID uint `json:"student_id"`
	ProgramID uint `json:"program_id"`
	ByDefault bool `json:"by_default"`
}
