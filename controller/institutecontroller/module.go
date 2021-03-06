package institutecontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetModules(c echo.Context) error {
	// Get data request
	module := models.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	modules := make([]models.Module, 0)

	// Query in database
	if err := db.Find(&modules, &module).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query semester by module
	for k, module := range modules {
		semesters := make([]models.ModuleSemester, 0)
		if err := db.Find(&semesters, models.ModuleSemester{ModuleID: module.ID}).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	modules := make([]models.Module, 0)
	if err := db.Where("name LIKE ? AND program_id = ?", "%"+request.Search+"%", request.ProgramID).
		Limit(5).Find(&modules).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    modules,
	})
}

func CreateModule(c echo.Context) error {
	// Get data request
	module := models.Module{}
	if err := c.Bind(&module); err != nil {
		return err
	}
	//module.ProgramID = currentUser.ProgramID

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Insert modules in database
	if err := db.Create(&module).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Delete all relations by semesters current module
	if err := DB.Delete(models.ModuleSemester{}, models.ModuleSemester{ModuleID: module.ID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Update module in database
	rows := DB.Model(&module).Update(module).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Delete module in database
	if err := DB.Delete(&module).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    module.ID,
		Message: fmt.Sprintf("El modulo %s se elimino correctamente", module.Name),
	})
}
