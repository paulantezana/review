package monitoringcontroller

import (
	"net/http"

	"github.com/paulantezana/review/models/monitoring"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

// GetTypeQuestions get all type questions
func GetTypeQuestions(c echo.Context) error {

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	typeQuestions := make([]monitoring.TypeQuestion, 0)
	if err := db.Find(&typeQuestions).
		Order("id desc").
		Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    typeQuestions,
	})
}
