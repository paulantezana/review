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

func CreateQuizAnswerDetail(c echo.Context) error {
    // Get data request
    answerDetail := models.QuizAnswerDetail{}
    if err := c.Bind(&answerDetail); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Insert answers in database
    if err := DB.Create(&answerDetail).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    answerDetail.ID,
        Message: fmt.Sprintf("La empresa %d se registro correctamente", answerDetail.ID),
    })
}
