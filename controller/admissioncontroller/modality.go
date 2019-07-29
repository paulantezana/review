package admissioncontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetModalities(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	modalities := make([]models.AdmissionModality, 0)
	DB.Find(&modalities)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    modalities,
	})
}

func GetModalityById(c echo.Context) error {
	// Get data request
	admissionModality := models.AdmissionModality{}
	if err := c.Bind(&admissionModality); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	if err := DB.First(&admissionModality, models.AdmissionModality{ID: admissionModality.ID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admissionModality,
	})
}

func CreateModality(c echo.Context) error {
	// Get data request
	modality := models.AdmissionModality{}
	if err := c.Bind(&modality); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Create new modality
	if err := db.Create(&modality).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    modality.ID,
		Message: fmt.Sprintf("La modalidad %s se registro exitosamente", modality.Name),
	})
}

func UpdateModality(c echo.Context) error {
	// Get data request
	modality := models.AdmissionModality{}
	if err := c.Bind(&modality); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update modality in database
	db.Model(&modality).Update(modality)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    modality.ID,
		Message: fmt.Sprintf("La modalidad %s se actualizaron correctamente", modality.Name),
	})
}

func DeleteModality(c echo.Context) error {
	// Get data request
	modality := models.AdmissionModality{}
	if err := c.Bind(&modality); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&modality).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    modality.ID,
		Message: fmt.Sprintf("La modalidad %s se elimino correctamente", modality.Name),
	})
}
