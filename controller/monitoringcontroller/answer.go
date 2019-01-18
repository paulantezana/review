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
    questions := make([]monitoringmodel.Question,0)
    if err := DB.Order("position asc").Find(&questions, monitoringmodel.Question{PollID: poll.ID}).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Get query answers
    for k, question := range questions {
        answerDetails := make([]monitoringmodel.AnswerDetail,0)
        if err := DB.Find(&answerDetails, monitoringmodel.AnswerDetail{QuestionID: question.ID}).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
        multipleQuestions := make([]monitoringmodel.MultipleQuestion,0)
        if err := DB.Find(&multipleQuestions, monitoringmodel.MultipleQuestion{QuestionID: question.ID}).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }

        questions[k].AnswerDetails =answerDetails
        questions[k].MultipleQuestions = multipleQuestions
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    questions,
    })
}