package institutemodel

type StudentHistory struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Description string `json:"description"`
	StudentID   uint   `json:"student_id"`
	UserID      uint   `json:"user_id"`
}
