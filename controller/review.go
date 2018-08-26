package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetReviews(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	if request.CurrentPage == 0 {
		request.CurrentPage = 1
	}
	offset := request.Limit*request.CurrentPage - request.Limit

	// Execute instructions
	var total uint
	reviews := make([]models.Review, 0)

	// Query in database
	if err := db.Where("module LIKE ?", "%"+request.Search+"%").
		Order("id asc").
		Offset(offset).Limit(request.Limit).Find(&reviews).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        reviews,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func CreateReview(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}
	review.UserID = currentUser.ID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validate
	rvw := make([]models.Review, 0)
	if db.Where("student_id = ? and module_id = ?", review.StudentID, review.ModuleId).
		Find(&rvw).RowsAffected >= 1 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: "Este alumno ya tiene una revision con este modulo",
		})
	}

	// Insert reviews in database
	if err := db.Create(&review).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    review.ID,
		Message: fmt.Sprintf("El revision del modulo %s se registro correctamente", review.Module),
	})
}

func UpdateReview(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update review in database
	rows := db.Model(&review).Update(review).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", review.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    review.ID,
		Message: fmt.Sprintf("Los datos del la revison con el modulo %s se actualizaron correctamente", review.Module),
	})
}

func DeleteReview(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation review exist
	if db.First(&review).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", review.ID),
		})
	}

	// Delete review in database
	if err := db.Delete(&review).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    review.ID,
		Message: fmt.Sprintf("El revision con el modulo %s se elimino correctamente", review.Module),
	})
}
