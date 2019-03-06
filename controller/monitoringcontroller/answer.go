package monitoringcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func CreateAnswer(c echo.Context) error {
	// Get data request
	answer := models.Answer{}
	if err := c.Bind(&answer); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Validate
	validateAnswer := models.Answer{}
	if rows := DB.First(&validateAnswer, models.Answer{
		StudentID: answer.StudentID,
		PollID:    answer.PollID,
	}).RowsAffected; rows >= 1 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("Esta encuesta ya fue resulto por este estudiante"),
		})
	}

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
	ID     uint   `json:"id"`
	Answer string `json:"answer"`
}
type multipleQuestionSummary struct {
	ID    uint   `json:"id"`
	Label string `json:"label"`
}
type answerSummary struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	TypeQuestionID uint   `json:"type_question_id"`

	MultipleQuestions []multipleQuestionSummary `json:"multiple_questions"`
	AnswerDetails     []answerDetailSummary     `json:"answer_details"`
}

func GetAnswerSummary(c echo.Context) error {
	// Get data request
	poll := models.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Get questions
	questions := make([]answerSummary, 0)
	if err := DB.Table("questions").Select("id, name, type_question_id").
		Where("poll_id = ?", poll.ID).Scan(&questions).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get query answers
	for k, question := range questions {
		answerDetails := make([]answerDetailSummary, 0)
		if err := DB.Table("answer_details").Select("id, answer").
			Where("question_id = ?", question.ID).Scan(&answerDetails).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		multipleQuestions := make([]multipleQuestionSummary, 0)
		if err := DB.Table("multiple_questions").Select("id, label").
			Where("question_id = ?", question.ID).Scan(&multipleQuestions).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		questions[k].AnswerDetails = answerDetails
		questions[k].MultipleQuestions = multipleQuestions
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    questions,
	})
}

type multipleQuestionOne struct {
	ID    uint   `json:"id"`
	Label string `json:"label"`
}

type navigateRequest struct {
	ID      uint `json:"id"`
	Current uint `json:"current"`
}

type answerNavigateResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Student   struct {
		ID       uint   `json:"id"`
		DNI      string `json:"dni"`
		FullName string `json:"full_name"`
		Gender   string `json:"gender"`
	} `json:"student"`
	Questions []struct {
		ID                uint                  `json:"id"`
		Name              string                `json:"name"`
		TypeQuestionID    uint                  `json:"type_question_id"`
		MultipleQuestions []multipleQuestionOne `json:"multiple_questions"`
		Answer            string                `json:"answer"`
	} `json:"questions"`
}

// Navigate answers
func GetAnswerNavigate(c echo.Context) error {
	// Get data request
	request := navigateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Prepare
	answerNavigate := answerNavigateResponse{}

	// Query answers
	var total uint
	DB.Model(&models.Answer{}).Where("poll_id = ?", request.ID).Count(&total)

	// Query current answer
	answer := models.Answer{}
	DB.Raw("SELECT * FROM answers "+
		"WHERE poll_id = ? "+
		"ORDER BY created_at DESC "+
		"LIMIT 1 "+
		"OFFSET ? ", request.ID, request.Current).Scan(&answer)
	answerNavigate.ID = answer.ID
	answerNavigate.CreatedAt = answer.CreatedAt

	// Query student
	student := models.Student{}
	if err := DB.First(&student, models.Student{ID: answer.StudentID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	answerNavigate.Student.ID = student.ID
	answerNavigate.Student.FullName = student.FullName
	answerNavigate.Student.DNI = student.DNI
	answerNavigate.Student.Gender = student.Gender

	// Get all questions by pollID
	if err := DB.Table("questions").Select("id, name, type_question_id").
		Where("poll_id = ?", request.ID).Scan(&answerNavigate.Questions).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get query answers
	for k, question := range answerNavigate.Questions {
		answerDetails := make([]answerDetailSummary, 0)
		if err := DB.Table("answer_details").Select("id, answer").
			Where("question_id = ? AND answer_id = ?", question.ID, answer.ID).
			Scan(&answerDetails).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}



		//multipleQuestions := make([]multipleQuestionOne, 0)
		//if err := DB.Table("multiple_questions").Select("id, label").
		//	Where("question_id = ?", question.ID).Scan(&multipleQuestions).Error; err != nil {
		//	return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		//}

		// check
		if len(answerDetails) >= 1 {
			answerNavigate.Questions[k].Answer = answerDetails[0].Answer
		}

        // Query multiple questions by multiple options
        if question.TypeQuestionID == 3 {
            // convert string to UINT == Convert answer to to id multiple question
            u, err := strconv.ParseUint(answerNavigate.Questions[k].Answer, 0, 32)

            // Check error
            if err == nil {
                mul := models.MultipleQuestion{}
                DB.First(&mul,models.MultipleQuestion{
                    QuestionID: question.ID,
                    ID: uint(u),
                })
                answerNavigate.Questions[k].Answer = mul.Label
            }
        }

        // Query multiple questions by multiple checks
        if question.TypeQuestionID == 4 {
            // Split string by ","
            ans := strings.Split(answerNavigate.Questions[k].Answer,",")

            answ := ""
            for i, an := range ans {
                // convert string to UINT == Convert answer to to id multiple question
                u, err := strconv.ParseUint(strings.TrimSpace(an), 0, 32)

                // Check error
                if err != nil {
                    break
                }

                // Query multiple question label
                mul := models.MultipleQuestion{}
                DB.First(&mul,models.MultipleQuestion{
                    QuestionID: question.ID,
                    ID: uint(u),
                })
                if i>=1 {
                    answ += "," + mul.Label
                }else {
                    answ = mul.Label
                }
            }
            answerNavigate.Questions[k].Answer = answ
        }
	}

	// return request
	return c.JSON(http.StatusOK, utilities.ResponsePaginate{
		Success:     true,
		Data:        answerNavigate,
		CurrentPage: request.Current,
		Total:       total,
	})
}
