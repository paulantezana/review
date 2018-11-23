package monitoring

type AnswerDetail struct {
    ID             uint   `json:"id" gorm:"primary_key"`
    Answer         string `json:"answer"`
    State          bool   `json:"state" gorm:"default:'true'"`

    QuestionID     uint   `json:"question_id"`
    TypeQuestionID uint   `json:"type_question_id"`
    AnswerID       uint   `json:"answer_id"`
}