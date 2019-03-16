package monitoringcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func CreateQuizAnswer(c echo.Context) error {
	// Get data request
	answer := models.QuizAnswer{}
	if err := c.Bind(&answer); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Validate
	//validateAnswer := models.QuizAnswer{}
	//if rows := DB.First(&validateAnswer, models.Answer{
	//	StudentID: answer.StudentID,
	//	PollID:    answer.PollID,
	//}).RowsAffected; rows >= 1 {
	//	return c.JSON(http.StatusOK, utilities.Response{
	//		Message: fmt.Sprintf("Esta encuesta ya fue resulto por este estudiante"),
	//	})
	//}

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
