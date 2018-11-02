package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetModules(c echo.Context) error {
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
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	modules := make([]models.Module, 0)

	// Query in database
	if err := db.Where("program_id = ?", currentUser.ProgramID).Order("sequence asc").Find(&modules).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    modules,
	})
}

func GetModuleSearch(c echo.Context) error {
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
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	modules := make([]models.Module, 0)
	if err := db.Where("name LIKE ? AND program_id = ?", "%"+request.Search+"%", currentUser.ProgramID).
		Limit(5).Find(&modules).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    modules,
	})
}

func CreateModule(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	module := models.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}
	module.ProgramID = currentUser.ProgramID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert modules in database
	if err := db.Create(&module).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    module.ID,
		Message: fmt.Sprintf("El modulo %s se registro correctamente", module.Name),
	})
}

func UpdateModule(c echo.Context) error {
	// Get data request
	module := models.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update module in database
	rows := db.Model(&module).Update(module).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", module.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    module.ID,
		Message: fmt.Sprintf("Los datos del modulo %s se actualizaron correctamente", module.Name),
	})
}

func DeleteModule(c echo.Context) error {
	// Get data request
	module := models.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation module exist
	if db.First(&module).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", module.ID),
		})
	}

	// Delete module in database
	if err := db.Delete(&module).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    module.ID,
		Message: fmt.Sprintf("El modulo %s se elimino correctamente", module.Name),
	})
}
