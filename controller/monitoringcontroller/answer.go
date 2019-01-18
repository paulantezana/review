package monitoringcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models/monitoringmodel"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func CreateAnswer(c echo.Context) error {
    // Get data request
    answer := monitoringmodel.Answer{}
    if err := c.Bind(&answer); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Insert companies in database
    if err := DB.Create(&answer).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    answer.ID,
        Message: fmt.Sprintf("La empresa %d se registro correctamente", answer.ID),
    })
}

type answerDetailSummary struct {
    ID     uint   `json:"id" gorm:"primary_key"`
    Answer string `json:"answer"`
}
type multipleQuestionSummary struct {
    ID    uint   `json:"id"`
    Label string `json:"label"`
}
type answerSummary struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    TypeQuestionID uint `json:"type_question_id"`

    MultipleQuestions []multipleQuestionSummary `json:"multiple_questions"`
    AnswerDetails []answerDetailSummary `json:"answer_details"`
} 

func GetAnswerAll(c echo.Context) error {
    // Get data request
    poll := monitoringmodel.Poll{}
    if err := c.Bind(&poll); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Get questions
    questions := make([]answerSummary,0)
    if err := DB.Table("questions").Select("id, name, type_question_id").Where("poll_id = ?", poll.ID).Scan(&questions).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Get query answers
    for k, question := range questions {
        answerDetails := make([]answerDetailSummary,0)
        if err := DB.Table("answer_details").Select("id, answer").Where("question_id = ?", question.ID).Scan(&answerDetails).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }

        multipleQuestions := make([]multipleQuestionSummary,0)
        if err := DB.Table("multiple_questions").Select("id, label").Where("question_id = ?", question.ID).Scan(&multipleQuestions).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }

        questions[k].AnswerDetails =answerDetails
        questions[k].MultipleQuestions = multipleQuestions
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    questions,
    })
}