package admissioncontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetPayments(c echo.Context) error {
	// Get data request
	payment := models.Payment{}
	if err := c.Bind(&payment); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	payments := make([]models.Payment, 0)
	DB.Find(&payments)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    payments,
	})
}

func CreatePayment(c echo.Context) error {
	// Get data request
	payment := models.Payment{}
	if err := c.Bind(&payment); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Create new payment
	if err := db.Create(&payment).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    payment.ID,
		Message: fmt.Sprintf("El paymenta de estudios %s se registro exitosamente", payment.Reason),
	})
}

func UpdatePayment(c echo.Context) error {
	// Get data request
	payment := models.Payment{}
	if err := c.Bind(&payment); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update payment in database
	rows := db.Model(&payment).Update(payment).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", payment.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    payment.ID,
		Message: fmt.Sprintf("Los datos del paymenta de estudios %s se actualizaron correctamente", payment.Reason),
	})
}

func DeletePayment(c echo.Context) error {
	// Get data request
	payment := models.Payment{}
	if err := c.Bind(&payment); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&payment).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    payment.ID,
		Message: fmt.Sprintf("The payment %s was successfully deleted", payment.Reason),
	})
}
