package institutemodel

type TeacherAction struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Action      string `json:"action"`
	Description string `json:"description"`

	TeacherID uint `json:"teacher_id"`
}
