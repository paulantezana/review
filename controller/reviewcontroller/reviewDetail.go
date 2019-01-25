package reviewcontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

// reviewDetailByReviewResponse struct
type reviewDetailByReviewResponse struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Hours     uint      `json:"hours"`
	Note      uint      `json:"note"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	ReviewID    uint   `json:"review_id"`
	CompanyID   uint   `json:"company_id"`
	CompanyName string `json:"company_name"`
}

// GetReviewsDetailByReview function
func GetReviewsDetailByReview(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	reviewDetailByReviewResponse := make([]reviewDetailByReviewResponse, 0)

	// Query in database
	if err := db.Table("review_details").
		Select("review_details.id, review_details.hours, review_details.note, review_details.start_date, review_details.end_date, review_details.company_id, companies.name_social_reason as company_name").
		Joins("INNER JOIN companies on review_details.company_id = companies.id").
		Order("review_details.id asc").
		Where("review_details.review_id = ?", review.ID).
		Scan(&reviewDetailByReviewResponse).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    reviewDetailByReviewResponse,
	})
}

// DeleteReview function delete review
func DeleteReviewDetail(c echo.Context) error {
	// Get data request
	reviewDetail := models.ReviewDetail{}
	if err := c.Bind(&reviewDetail); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete review in database
	if err := db.Delete(&reviewDetail).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    reviewDetail.ID,
		Message: fmt.Sprintf("El el detalle de la revision se elimino correctamente"),
	})
}
