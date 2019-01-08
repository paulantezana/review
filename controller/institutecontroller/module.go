package institutecontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models/institutemodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetModules(c echo.Context) error {
	// Get data request
	module := institutemodel.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	modules := make([]institutemodel.Module, 0)

	// Query in database
	if err := db.Find(&modules, &module).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Query semester by module
	for k, module := range modules {
		semesters := make([]institutemodel.ModuleSemester, 0)
		if err := db.Find(&semesters, institutemodel.ModuleSemester{ModuleID: module.ID}).Error; err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		modules[k].Semesters = semesters
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    modules,
	})
}

func GetModuleSearch(c echo.Context) error {

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	modules := make([]institutemodel.Module, 0)
	if err := db.Where("name LIKE ? AND program_id = ?", "%"+request.Search+"%", request.ProgramID).
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
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	module := institutemodel.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}
	//module.ProgramID = currentUser.ProgramID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert modules in database
	if err := db.Create(&module).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
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
	module := institutemodel.Module{}
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
	module := institutemodel.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

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
