package monitoringcontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models/monitoringmodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

// DeleteMultipleQuestion Delete one question
func DeleteMultipleQuestion(c echo.Context) error {
	// Get data request
	multipleQuestion := monitoringmodel.MultipleQuestion{}
	if err := c.Bind(&multipleQuestion); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete question in database
	if err := db.Delete(&multipleQuestion).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    multipleQuestion.ID,
		Message: fmt.Sprintf("La pregunta multiple %s se elimino correctamente", multipleQuestion.Label),
	})
}
