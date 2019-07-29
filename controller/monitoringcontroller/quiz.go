package monitoringcontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetQuizzesPaginate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	quizzes := make([]models.Quiz, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", request.ProgramID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&quizzes).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        quizzes,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

// Get all quizzes
func GetQuizzesAllByDiplomat(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	quizzes := make([]models.Quiz, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND quiz_diplomat_id = ?", "%"+request.Search+"%", request.QuizDiplomatID).
		Order("id asc").Find(&quizzes).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quizzes,
	})
}

// Query By Diplomat student
func GetQuizzesAllByDiplomatStudent(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	quizzes := make([]models.Quiz, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND quiz_diplomat_id = ? AND state = true", "%"+request.Search+"%", request.QuizDiplomatID).
		Order("id asc").Find(&quizzes).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quizzes,
	})
}

// From student
func GetQuizzesPaginateStudent(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	quizzes := make([]models.Quiz, 0)

	// Validations
	filterIDS, err := validateRestrictions("quizzes", currentUser)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND program_id = ? AND state = true AND id IN (?)", "%"+request.Search+"%", request.ProgramID, filterIDS).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&quizzes).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        quizzes,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetQuizByID(c echo.Context) error {
	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&quiz, quiz.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quiz,
	})
}

func CreateQuiz(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// set current programID
	//poll.ProgramID = currentUser.DefaultProgramID

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Insert companies in database
	if err := db.Create(&quiz).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quiz.ID,
		Message: fmt.Sprintf("El cuestionario %s se registro correctamente", quiz.Name),
	})
}

func UpdateQuiz(c echo.Context) error {
	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update poll in database
	db.Model(&quiz).Update(quiz)

	// Update columns
	db.Model(&quiz).UpdateColumns(map[string]interface{}{
		"start_date_enable": quiz.StartDateEnable,
		"end_date_enable":   quiz.EndDateEnable,
		"limit_time_enable": quiz.LimitTimeEnable,
		"show_analyze":      quiz.ShowAnalyze,
	})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quiz.ID,
		Message: fmt.Sprintf("Los datos del cuestionario %s se actualizaron correctamente", quiz.Name),
	})
}

func UpdateStateQuiz(c echo.Context) error {
	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update columns
	db.Model(&quiz).UpdateColumn("state", quiz.State)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quiz.ID,
		Message: fmt.Sprintf("Los datos dell cuestionario %s se actualizaron correctamente", quiz.Name),
	})
}

// DeletePoll delete quiz by id
func DeleteQuiz(c echo.Context) error {
	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Delete poll in database
	if err := db.Delete(&quiz).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quiz.ID,
		Message: fmt.Sprintf("El cuestionario %s se elimino correctamente", quiz.Name),
	})
}
