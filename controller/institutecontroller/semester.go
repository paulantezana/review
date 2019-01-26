package institutecontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetSemesters(c echo.Context) error {
	// Get data request
	semester := models.Semester{}
	if err := c.Bind(&semester); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	semesters := make([]models.Semester, 0)
	if err := db.Debug().Where("program_id = ?", semester.ProgramID).Order("sequence asc").Find(&semesters).
		Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    semesters,
	})
}

func CreateSemester(c echo.Context) error {
	// Get data request
	semester := models.Semester{}
	if err := c.Bind(&semester); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Create new semester
	if err := db.Create(&semester).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    semester.ID,
		Message: fmt.Sprintf("El semestera de estudios %s se registro exitosamente", semester.Name),
	})
}

func UpdateSemester(c echo.Context) error {
	// Get data request
	semester := models.Semester{}
	if err := c.Bind(&semester); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update semester in database
	rows := db.Model(&semester).Update(semester).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", semester.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    semester.ID,
		Message: fmt.Sprintf("Los datos del semestera de estudios %s se actualizaron correctamente", semester.Name),
	})
}

func DeleteSemester(c echo.Context) error {
	// Get data request
	semester := models.Semester{}
	if err := c.Bind(&semester); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&semester).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    semester.ID,
		Message: fmt.Sprintf("The semester %s was successfully deleted", semester.Name),
	})
}
