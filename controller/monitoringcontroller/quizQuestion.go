package monitoringcontroller

import (
	"fmt"
	"github.com/paulantezana/review/models"
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

// GetQuestions get all questions by poll
func GetQuizQuestions(c echo.Context) error {
	// Get data request
    quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
    quizQuestions := make([]models.QuizQuestion, 0)

	// Query in database
	if err := db.Where("poll_id = ?", quiz.ID).
		Order("position asc").Find(&quizQuestions).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

    // Quiz Questions
	for k, question := range quizQuestions {
        multipleQuizQuestion := make([]models.MultipleQuizQuestion, 0)
		if err := db.Where("question_id = ?", question.ID).
			Order("id asc").Find(&multipleQuizQuestion).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
        quizQuestions[k].MultipleQuizQuestions = multipleQuizQuestion
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quizQuestions,
	})
}

type createQuizQuestionsRequest struct {
	Questions []models.Question `json:"questions"`
}

// CreateQuestions create multiple questions
func CreateQuizQuestions(c echo.Context) error {
	// Get data request
	request := createQuizQuestionsRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert companies in database
	tr := db.Begin()

	insetCount := 0
	updateCount := 0
	for _, question := range request.Questions {
		if question.ID == 0 {
			if err := tr.Create(&question).Error; err != nil {
				tr.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
			insetCount++
		} else {
			if err := tr.Model(&question).Update(question).Error; err != nil {
				tr.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
			updateCount++
		}
	}

	// Commit transaction
	tr.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se inserto %d preguntas y se actualizo %d preguntas de manera exitosa", insetCount, updateCount),
	})
}

// UpdateQuestion update one question
func UpdateQuizQuestion(c echo.Context) error {
	// Get data request
    quizQuestion := models.QuizQuestion{}
	if err := c.Bind(&quizQuestion); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update question in database
	rows := db.Model(&quizQuestion).Update(quizQuestion).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", quizQuestion.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quizQuestion.ID,
		Message: fmt.Sprintf("Los datos del la pregunta %s se actualizaron correctamente", quizQuestion.Name),
	})
}

// DeleteQuestion Delete one question
func DeleteQuizQuestion(c echo.Context) error {
	// Get data request
    quizQuestion := models.QuizQuestion{}
	if err := c.Bind(&quizQuestion); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete question in database
	if err := db.Delete(&quizQuestion).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quizQuestion.ID,
		Message: fmt.Sprintf("La pregunta %s se elimino correctamente", quizQuestion.Name),
	})
}
