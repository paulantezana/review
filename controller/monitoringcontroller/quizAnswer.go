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
	DB.Last(&answer, models.QuizAnswer{QuizID: answer.QuizID, StudentID: answer.StudentID})

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    answer,
		Message: fmt.Sprintf("La empresa %d se registro correctamente", answer.ID),
	})
}

// Multiple question
type multipleQuestionR struct {
	ID      uint   `json:"id"`
	Label   string `json:"label"`
	Correct bool   `json:"correct"`
}

// Question response
type questionR struct {
	ID                uint                `json:"id"`
	Name              string              `json:"name"`
	Answer            string              `json:"answer"`
	MultipleQuestions []multipleQuestionR `json:"multiple_questions"`
}

// Analyze Attempts
type attemptR struct {
	Note      uint        `json:"note"`
	Questions []questionR `json:"questions"`
}

// response struct
type analyzeQuizAnswerResponse struct {
	Attempts []attemptR `json:"attempts"`
}

//GetAnalyzeQuizAnswer
func GetAnalyzeQuizAnswerByStudent(c echo.Context) error {
	// Get data request
	answer := models.QuizAnswer{}
	if err := c.Bind(&answer); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// analyzeQuizAnswerResponse
	analyze := analyzeQuizAnswerResponse{}

	// Get questions
	questions := make([]answerSummary, 0)
	if err := DB.Table("quiz_questions").Select("id, name, type_question_id").
		Where("quiz_id = ?", answer.QuizID).Scan(&questions).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query answer
	quizAnswers := make([]models.QuizAnswer, 0)
	DB.Find(&quizAnswers, models.QuizAnswer{QuizID: answer.QuizID, StudentID: answer.StudentID})

	for _, quizAnswer := range quizAnswers {
		attemptR := attemptR{}
		for _, question := range questions {
			// Prepare struct question
			questionR := questionR{
				ID:   question.ID,
				Name: question.Name,
			}

			// Query answers
			answerDetail := answerDetailSummary{}
			if err := DB.Table("quiz_answer_details").Select("id, answer").
				Where("quiz_question_id = ? AND quiz_answer_id = ?", question.ID, quizAnswer.ID).
				Limit(1).
				Scan(&answerDetail).Error; err != nil {
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
            questionR.Answer = answerDetail.Answer

			// Query multiple questions
			multipleQuestionRs := make([]multipleQuestionR, 0)
			if err := DB.Table("multiple_quiz_questions").Select("id, label, correct").
				Where("quiz_question_id = ?", question.ID).
				Scan(&multipleQuestionRs).Error; err != nil {
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
			questionR.MultipleQuestions = multipleQuestionRs

			// Set Questions
			attemptR.Questions = append(attemptR.Questions, questionR)
		}
		// Set Attempts
		analyze.Attempts = append(analyze.Attempts, attemptR)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    analyze,
	})
}

// CreateQuizAnswer
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
	DB.First(&quizAnswer, models.QuizAnswer{QuizID: answer.QuizID, StudentID: answer.StudentID})

	// Validate
	if quizAnswer.ID >= 1 {
		answer.Attempts = quizAnswer.Attempts + 1
	} else {
		answer.Attempts = 1
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
	QuizQuestionID uint   `json:"quiz_question_id"`
	QuizAnswerID   uint   `json:"quiz_answer_id"`
	Answer         string `json:"answer"`
	Current        uint   `json:"current"`
	Total          uint   `json:"total"`
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
		QuizAnswerID:   request.QuizAnswerID,
		Answer:         request.Answer,
	}
	if err := DB.Create(&answerDetail).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query
	quizAnswer := models.QuizAnswer{}
	if err := DB.First(&quizAnswer, models.QuizAnswer{ID: answerDetail.QuizAnswerID}).Error; err != nil {
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

func TimeFinishQuizAnswer(c echo.Context) error {
	// Get data request
	quizAnswer := models.QuizAnswer{}
	if err := c.Bind(&quizAnswer); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// update quiz
	quizAnswer.Step = 2
	DB.Model(&quizAnswer).Update(quizAnswer)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quizAnswer,
		Message: fmt.Sprintf("El tiempo ha finalizado."),
	})
}
