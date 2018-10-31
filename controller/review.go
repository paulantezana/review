package controller

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
)

// GetReviews funstions get reviews by student_id
func GetReviews(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	reviews := make([]models.Review, 0)

	// Query in database
	if err := db.Where("student_id = ?", student.ID).
		Order("id asc").Find(&reviews).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    reviews,
	})
}

// CreateReview function create new review
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
		Message: fmt.Sprintf("El revision del modulo se registro correctamente"),
	})
}

// UpdateReview function update review
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
		Message: fmt.Sprintf("Los datos del la revison con el modulo se actualizaron correctamente"),
	})
}

// DeleteReview function delete review
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
		Message: fmt.Sprintf("El revision con el modulo se elimino correctamente"),
	})
}

// GetActaReview function get data acta
func GetActaReview(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Find quotations in database by RequirementID  ========== Quotations, Providers, Users
	modules := make([]models.Module, 0)
	if err := db.Table("reviews").
		Select("modules.id, modules.name, modules.sequence, modules.points, modules.hours, modules.semester").
		Joins("INNER JOIN modules on reviews.module_id = modules.id").
		Order("modules.sequence asc").
		Where("WHERE reviews.id = ?", review.ID).
		Scan(&modules).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    "Hola",
	})
}

// consResponse struct
type consResponse struct {
	Module  models.Module `json:"module"`
	Success bool          `json:"success"`
}

// GetConstReview function get data constancy
func GetConstReview(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Find quotations in database by RequirementID  ========== Quotations, Providers, Users
	modules := make([]models.Module, 0)
	if err := db.Table("reviews").
		Select("modules.id, modules.name, modules.sequence, modules.points, modules.hours, modules.semester").
		Joins("INNER JOIN modules on reviews.module_id = modules.id").
		Where("reviews.id = ?", review.ID).
		Scan(&modules).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, consResponse{
		Success: true,
		Module:  modules[0],
	})
}
