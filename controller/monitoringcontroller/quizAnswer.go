package monitoringcontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetLastQuizAnswer(c echo.Context) error {
    // Get data request
    answer := models.QuizAnswer{}
    if err := c.Bind(&answer); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query answer
    DB.Last(&answer,models.QuizAnswer{QuizID: answer.QuizID, StudentID: answer.StudentID})

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    answer,
        Message: fmt.Sprintf("La empresa %d se registro correctamente", answer.ID),
    })
}

func CreateQuizAnswer(c echo.Context) error {
	// Get data request
	answer := models.QuizAnswer{}
	if err := c.Bind(&answer); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query answer
    quizAnswer := models.QuizAnswer{}
	DB.First(&quizAnswer,models.QuizAnswer{QuizID: answer.QuizID, StudentID: answer.StudentID})

    // Validate
    if quizAnswer.ID >= 1 {
        answer.Attempts =  quizAnswer.Attempts + 1
    }else {
        answer.Attempts =  1
        answer.Step = 1
    }

	// Insert answers in database
	if err := DB.Create(&answer).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    answer,
		Message: fmt.Sprintf("La empresa %d se registro correctamente", answer.ID),
	})
}

type answerDetailRequest struct {
    QuizQuestionID uint `json:"quiz_question_id"`
    QuizAnswerID   uint `json:"quiz_answer_id"`
    Current uint `json:"current"`
    Total uint `json:"total"`
}
func CreateQuizAnswerDetail(c echo.Context) error {
    // Get data request
    request := answerDetailRequest{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Insert answers in database
    answerDetail := models.QuizAnswerDetail{
        QuizQuestionID: request.QuizQuestionID,
        QuizAnswerID: request.QuizAnswerID,
    }
    if err := DB.Create(&answerDetail).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Query
    quizAnswer := models.QuizAnswer{}
    if err := DB.First(&quizAnswer,models.QuizAnswer{ID: answerDetail.QuizAnswerID}).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Update quiz in database
    if request.Current == request.Total {
        quizAnswer.Step = quizAnswer.Step + 1
    }
    quizAnswer.CurrentQuestion = request.Current + 1
    DB.Model(&quizAnswer).Update(quizAnswer)

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    quizAnswer,
        Message: fmt.Sprintf("La empresa se registro correctamente"),
    })
}
