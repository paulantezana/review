package institutemodel

type StudentStatus struct {
	ID   uint `json:"id" gorm:"primary_key"`
	Name uint `json:"name"`
}

// Set User's table name to be `profiles`
func (StudentStatus) TableName() string {
	return "student_status"
}
