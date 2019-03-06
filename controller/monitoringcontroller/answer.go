package monitoringcontroller

import (
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
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
		answerDetail := answerDetailSummary{}
		if err := DB.Table("answer_details").Select("id, answer").
			Where("question_id = ? AND answer_id = ?", question.ID, answer.ID).
		    Limit(1).
			Scan(&answerDetail).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Set answer
        answerNavigate.Questions[k].Answer = answerDetail.Answer

        // Query multiple questions by multiple options
        if question.TypeQuestionID == 3 {
            answerNavigate.Questions[k].Answer = getQuestionAnswerMultiple(answerNavigate.Questions[k].Answer,question.ID)
        }

        // Query multiple questions by multiple checks
        if question.TypeQuestionID == 4 {
            answerNavigate.Questions[k].Answer = getQuestionAnswerCheck(answerNavigate.Questions[k].Answer,question.ID)
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

func ExportExcelAnswers(c echo.Context) error {
    // Get data request
    poll := models.Poll{}
    if err := c.Bind(&poll); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query answers
    answers := make([]models.Answer,0)
    if err := DB.Find(&answers, models.Answer{PollID: poll.ID}).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // =========================================
    // CREATE NEW FILE EXCEL
    excel := excelize.NewFile()

    // Create a new sheet.
    answerSheet := "Answers"
    answerSheetIndex := excel.NewSheet(answerSheet)


    // Get all questions by pollID
    questions := make([]models.Question,0)
    if err := DB.Find(&questions,models.Question{PollID: poll.ID}).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Excel headers
    excel.SetCellValue(answerSheet,"A4","DNI")
    excel.SetCellValue(answerSheet,"B4","Apellidos y Nombres")
    excel.SetCellValue(answerSheet,"C4","Fecha")
    for k, qus := range questions {
        excel.SetCellValue(answerSheet, utilities.ConvertNumberToCharScheme(uint(k + 4)) + "4",qus.Name)
    }

    // Query answer details
    // START
    for k, answer := range answers {

        // Query student
        student := models.Student{}
        if err := DB.First(&student, models.Student{ID: answer.StudentID}).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }

        // Set data student to excel sheet
        excel.SetCellValue(answerSheet,fmt.Sprintf("A%d",k + 5),student.DNI)
        excel.SetCellValue(answerSheet,fmt.Sprintf("B%d",k + 5),student.FullName)
        excel.SetCellValue(answerSheet,fmt.Sprintf("C%d",k + 5),answer.CreatedAt)

        // Query answer detail
        for i, question := range questions {
            answerDetail := answerDetailSummary{}
            if err := DB.Table("answer_details").Select("id, answer").
                Where("question_id = ? AND answer_id = ?", question.ID, answer.ID).
                Limit(1).
                Scan(&answerDetail).Error; err != nil {
                return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
            }

            // Query multiple questions by multiple options
            if question.TypeQuestionID == 3 {
                answerDetail.Answer = getQuestionAnswerMultiple(answerDetail.Answer,question.ID)
            }

            // Query multiple questions by multiple checks
            if question.TypeQuestionID == 4 {
                answerDetail.Answer = getQuestionAnswerCheck(answerDetail.Answer,question.ID)
            }

            // Set answerDetail answers in sheet
            excel.SetCellValue(answerSheet,fmt.Sprintf("%s%d",utilities.ConvertNumberToCharScheme(uint(i + 4)),k + 5),answerDetail.Answer)
        }
    }
    // Query answer details
    // END

    // Set active sheet of the workbook.
    excel.SetActiveSheet(answerSheetIndex)

    // Save excel file by the given path.
    err := excel.SaveAs("./temp/answers.xlsx")
    if err != nil {
        fmt.Println(err)
    }

    // Return file excel
    return c.File("./temp/answers.xlsx")
}

// getQuestionAnswerCheck
func getQuestionAnswerCheck(an string, questionID uint) string  {
    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Split string by ","
    ans := strings.Split(strings.TrimSpace(an),",")

    anr := ""
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
            QuestionID: questionID,
            ID: uint(u),
        })

        //
        if i>=1 {
            anr += "," + mul.Label
        }else {
            anr = mul.Label
        }
    }

    // Validate
    if anr == "" {
        return an
    }

    // return data
    return anr
}

// getQuestionAnswerMultiple
func getQuestionAnswerMultiple(an string, questionID uint) string  {
    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // convert string to UINT == Convert answer to to id multiple question
    u, err := strconv.ParseUint(strings.TrimSpace(an), 0, 32)

    // Check error
    anr := an
    if err == nil {
        mul := models.MultipleQuestion{}
        DB.First(&mul,models.MultipleQuestion{
            QuestionID: questionID,
            ID: uint(u),
        })
        anr = mul.Label
    }

    // Return data
    return anr
}