package monitoringcontroller

import (
	"fmt"
	"github.com/paulantezana/review/models"
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

// GetPollsPaginate
func GetPollsPaginate(c echo.Context) error {
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
	polls := make([]models.Poll, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", request.ProgramID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&polls).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        polls,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

// From student
func GetPollsPaginateStudent(c echo.Context) error {
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
	polls := make([]models.Poll, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND program_id = ? AND state = true", "%"+request.Search+"%", request.ProgramID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&polls).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        polls,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetPollByID(c echo.Context) error {
	// Get data request
	poll := models.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&poll, poll.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    poll,
	})
}

func CreatePoll(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	poll := models.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// set current programID
	//poll.ProgramID = currentUser.DefaultProgramID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert companies in database
	if err := db.Create(&poll).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    poll.ID,
		Message: fmt.Sprintf("La encuesta %s se registro correctamente", poll.Name),
	})
}

func UpdatePoll(c echo.Context) error {
	// Get data request
	poll := models.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update poll in database
	db.Model(&poll).Update(poll)

	// Update columns
	db.Model(&poll).UpdateColumns(map[string]interface{}{
		"start_date_enable": poll.StartDateEnable,
		"end_date_enable":   poll.EndDateEnable,
		"show_analyze":      poll.ShowAnalyze,
	})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    poll.ID,
		Message: fmt.Sprintf("Los datos del la encuesta %s se actualizaron correctamente", poll.Name),
	})
}

func UpdateStatePoll(c echo.Context) error {
	// Get data request
	poll := models.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update columns
	db.Model(&poll).UpdateColumn("state", poll.State)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    poll.ID,
		Message: fmt.Sprintf("Los datos del la encuesta %s se actualizaron correctamente", poll.Name),
	})
}

// DeletePoll delete poll by id
func DeletePoll(c echo.Context) error {
	// Get data request
	poll := models.Poll{}
	if err := c.Bind(&poll); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete poll in database
	if err := db.Delete(&poll).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    poll.ID,
		Message: fmt.Sprintf("La encuesta %s se elimino correctamente", poll.Name),
	})
}
