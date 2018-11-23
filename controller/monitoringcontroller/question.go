package monitoringcontroller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models/monitoring"
	"github.com/paulantezana/review/utilities"
)

// GetQuestions get all questions by poll
func GetQuestions(c echo.Context) error {
	// Get data request
	poll := monitoring.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	questions := make([]monitoring.Question, 0)

	// Query in database
	if err := db.Where("poll_id = ?", poll.ID).
		Order("id asc").Scan(&questions).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    questions,
	})
}

type createQuestionsRequest struct {
	Questions []monitoring.Question `json:"questions"`
}

// CreateQuestions create multiple questions
func CreateQuestions(c echo.Context) error {
	// Get data request
	request := createQuestionsRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert companies in database
	tr := db.Begin()

	for _, question := range request.Questions {
		if err := db.Create(&question).Error; err != nil {
			tr.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: fmt.Sprintf("%s", err),
			})
		}
	}

	// Commit transaction
	tr.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se inserto %d preguntas correctamente", len(request.Questions)),
	})
}

// UpdateQuestion update one question
func UpdateQuestion(c echo.Context) error {
	// Get data request
	question := monitoring.Question{}
	if err := c.Bind(&question); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update question in database
	rows := db.Model(&question).Update(question).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", question.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    question.ID,
		Message: fmt.Sprintf("Los datos del la pregunta %s se actualizaron correctamente", question.Name),
	})
}

// DeleteQuestion Delete one question
func DeleteQuestion(c echo.Context) error {
	// Get data request
	question := monitoring.Question{}
	if err := c.Bind(&question); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete question in database
	if err := db.Delete(&question).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    question.ID,
		Message: fmt.Sprintf("La pregunta %s se elimino correctamente", question.Name),
	})
}
