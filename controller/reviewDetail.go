package controller

import (
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetReviewsDetail(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	reviews := make([]models.Review, 0)

	// Query in database
	if err := db.Where("review_id = ?", review.ID).
		Order("id asc").Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    reviews,
	})
}
