package institutemodel

type TeacherProgram struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	TeacherID uint   `json:"teacher_id"`
	ProgramID uint   `json:"program_id"`
	ByDefault bool   `json:"by_default"`
	Type      string `json:"type"` // cross // career
}
