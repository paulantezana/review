package monitoring

// Answer struct
type Answer struct {
	ID    uint `json:"id" gorm:"primary_key"`
	State bool `json:"state" gorm:"default:'true'"`

	StudentID uint `json:"student_id"`

	AnswerDetails []AnswerDetail `json:"answer_details, omitempty"`
}
