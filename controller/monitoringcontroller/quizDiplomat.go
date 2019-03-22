package monitoringcontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetQuizDiplomatPaginate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	quizDiplomats := make([]models.QuizDiplomat, 0)

	// Query in database
	if err := DB.Where("lower(name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", request.ProgramID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&quizDiplomats).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        quizDiplomats,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

// From student
func GetQuizDiplomatPaginateStudent(c echo.Context) error {
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
	DB := config.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	quizDiplomats := make([]models.QuizDiplomat, 0)

	// Validations
	filterIDS, err := validateRestrictions("quiz_diplomats", currentUser)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query in database
	if err := DB.Where("lower(name) LIKE lower(?) AND id IN (?)", "%"+request.Search+"%", filterIDS).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&quizDiplomats).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        quizDiplomats,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetQuizDiplomatByID(c echo.Context) error {
	// Get data request
	quizDiplomat := models.QuizDiplomat{}
	if err := c.Bind(&quizDiplomat); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&quizDiplomat, quizDiplomat.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quizDiplomat,
	})
}

func CreateQuizDiplomat(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	quizDiplomat := models.QuizDiplomat{}
	if err := c.Bind(&quizDiplomat); err != nil {
		return err
	}

	// set current programID
	//poll.ProgramID = currentUser.DefaultProgramID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert companies in database
	if err := db.Create(&quizDiplomat).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    quizDiplomat.ID,
		Message: fmt.Sprintf("El cuestionario %s se registro correctamente", quizDiplomat.Name),
	})
}

func UpdateQuizDiplomat(c echo.Context) error {
	// Get data request
	quizDiplomat := models.QuizDiplomat{}
	if err := c.Bind(&quizDiplomat); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update poll in database
	db.Model(&quizDiplomat).Update(quizDiplomat)

	// Update columns
	//db.Model(&quizDiplomat).UpdateColumns(map[string]interface{}{
	//    "start_date_enable": quizDiplomat.StartDateEnable,
	//    "end_date_enable":   quiz.EndDateEnable,
	//    "limit_time_enable": quiz.LimitTimeEnable,
	//    "show_analyze":      quiz.ShowAnalyze,
	//})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quizDiplomat.ID,
		Message: fmt.Sprintf("Los datos del cuestionario %s se actualizaron correctamente", quizDiplomat.Name),
	})
}

func UpdateStateQuizDiplomat(c echo.Context) error {
	// Get data request
	quizDiplomat := models.QuizDiplomat{}
	if err := c.Bind(&quizDiplomat); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update columns
	db.Model(&quizDiplomat).UpdateColumn("state", quizDiplomat.State)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    quizDiplomat.ID,
		Message: fmt.Sprintf("Los datos dell cuestionario %s se actualizaron correctamente", quizDiplomat.Name),
	})
}

// DeletePoll delete quiz by id
func DeleteQuizDiplomat(c echo.Context) error {
	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
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
