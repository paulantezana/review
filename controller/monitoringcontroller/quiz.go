package monitoringcontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetQuizzesPaginate(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
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

func GetQuizByID(c echo.Context) error {
	// Get data request
	quiz := models.Quiz{}
	if err := c.Bind(&quiz); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
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
	db := config.GetConnection()
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
	db := config.GetConnection()
	defer db.Close()

	// Update poll in database
	rows := db.Model(&quiz).Update(quiz).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", quiz.ID),
		})
	}

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
	db := config.GetConnection()
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
