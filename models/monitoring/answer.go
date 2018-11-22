package monitoring

type Answer struct {
    ID        uint `json:"id" gorm:"primary_key"`
    StudentID uint `json:"student_id"`
    State     bool `json:"state" gorm:"default:'true'"`

    AnswerDetails []AnswerDetail `json:"answer_details"`
}
