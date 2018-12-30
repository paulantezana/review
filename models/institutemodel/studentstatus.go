package institutemodel

type StudentStatus struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

// Set User's table name to be `profiles`
func (StudentStatus) TableName() string {
	return "student_status"
}
