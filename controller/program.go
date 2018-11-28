package controller

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetPrograms(c echo.Context) error {

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	programs := make([]models.Program, 0)
	if err := db.Find(&programs).
		Order("id desc").
		Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    programs,
	})
}

func GetProgramByID(c echo.Context) error {
	// Get data request
	program := models.Program{}
	if err := c.Bind(&program); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&program, program.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    program,
	})
}

type createProgramRequest struct {
	Name      string `json:"name"`
	DNI       string `json:"dni"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

func CreateProgram(c echo.Context) error {
	// Get data request
	request := createProgramRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// ------------------------------------
	// Starting transaction
	// ------------------------------------
	tr := db.Begin()
	program := models.Program{
		Name: request.Name,
	}
	// Create new program
	if err := tr.Create(&program).Error; err != nil {
		tr.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create user
	cc := sha256.Sum256([]byte(request.Password))
	pwd := fmt.Sprintf("%x", cc)

	user := models.User{
		UserName:  request.UserName,
		Password:  pwd,
		Email:     request.Email,
		ProgramID: program.ID,
		Profile:   "admin",
	}
	if err := tr.Create(&user).Error; err != nil {
		tr.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create teacher
	teacher := models.Teacher{
		DNI:       request.DNI,
		FirstName: request.FirstName,
		ProgramID: program.ID,
		UserID:    user.ID,
	}
	if err := tr.Create(&teacher).Error; err != nil {
		tr.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	tr.Commit()
	// ------------------------------------
	// End Transaction
	// ------------------------------------

	// Insert program in database

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    program.ID,
		Message: fmt.Sprintf("El programa de estudios %s se registro exitosamente", program.Name),
	})
}

func UpdateProgram(c echo.Context) error {
	// Get data request
	program := models.Program{}
	if err := c.Bind(&program); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update program in database
	rows := db.Model(&program).Update(program).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", program.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    program.ID,
		Message: fmt.Sprintf("Los datos del programa de estudios %s se actualizaron correctamente", program.Name),
	})
}
