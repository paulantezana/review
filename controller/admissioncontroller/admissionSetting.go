package admissioncontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetAdmissionSettings(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	admSettings := make([]models.AdmissionSetting, 0)
	DB.Order("year asc").Find(&admSettings, models.AdmissionSetting{SubsidiaryID: request.SubsidiaryID})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admSettings,
	})
}

// GetStudentDetailByID get student detail
func GetAdmissionSettingByID(c echo.Context) error {
	// Get data request
	admSetting := models.AdmissionSetting{}
	if err := c.Bind(&admSetting); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&admSetting, admSetting.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    admSetting,
	})
}

func CreateAdmissionSetting(c echo.Context) error {
	// Get data request
	admSetting := models.AdmissionSetting{}
	if err := c.Bind(&admSetting); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Create new payment
	if err := db.Create(&admSetting).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    admSetting.ID,
		Message: fmt.Sprintf("La admisión con el id %d se registro exitosamente", admSetting.ID),
	})
}

func UpdateAdmissionSetting(c echo.Context) error {
	// Get data request
	admSetting := models.AdmissionSetting{}
	if err := c.Bind(&admSetting); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update payment in database
	rows := db.Model(&admSetting).Update(admSetting).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", admSetting.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    admSetting.ID,
		Message: fmt.Sprintf("La admisión con el id %d se actualizaron correctamente", admSetting.ID),
	})
}

func DeleteAdmissionSetting(c echo.Context) error {
	// Get data request
	admSetting := models.AdmissionSetting{}
	if err := c.Bind(&admSetting); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&admSetting).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admSetting.ID,
		Message: fmt.Sprintf("La admisión con el id %d se elimino correctamente", admSetting.ID),
	})
}