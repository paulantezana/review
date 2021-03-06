package monitoringcontroller

import (
	"fmt"
	"github.com/paulantezana/review/models"
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
)

// GetTypeQuestions get all type questions
func GetTypeQuestions(c echo.Context) error {

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	typeQuestions := make([]models.TypeQuestion, 0)
	if err := db.Find(&typeQuestions).
		Order("id desc").
		Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    typeQuestions,
	})
}
