package models

import "time"

type QuizQuestion struct {
    ID        uint      `json:"id" gorm:"primary_key"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Name      string    `json:"name"`
    Position  uint      `json:"position"`
    State     bool      `json:"state" gorm:"default:'true'"`

    TypeQuestionID uint `json:"type_question_id"`
    QuizID         uint `json:"quiz_id"`

    MultipleQuizQuestions []MultipleQuizQuestion `json:"multiple_quiz_questions"`
    //AnswerDetails     []AnswerDetail     `json:"answer_details"`
}