package controller

import (
	"fmt"
	"net/http"
	"time"

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

// moduleResponse struct
type moduleResponse struct {
	ID              uint   `json:"id"`
	Sequence        uint   `json:"sequence"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Points          uint   `json:"points"`
	Hours           uint   `json:"hours"`
	Semester        string `json:"semester"`
	StudentID       uint   `json:"student_id"`
	StudentDNI      string `json:"student_dni"`
	StudentFullName string `json:"student_full_name"`
}

// detailResponse struct
type detailResponse struct {
	ID                 uint      `json:"id" gorm:"primary_key"`
	Hours              uint      `json:"hours"`
	Note               uint      `json:"note"`
	NoteAppreciation   uint      `json:"note_appreciation"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	RUC                string    `json:"ruc"`
	NombreORazonSocial string    `json:"nombre_o_razon_social"`
    Direccion              string `json:"direccion"`
}

// consResponse struct
type consResponse struct {
	Module  moduleResponse   `json:"module"`
	Success bool             `json:"success"`
	Detail  []detailResponse `json:"detail"`
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

	// Find reviews
	moduleResponses := make([]moduleResponse, 0)
	if err := db.Table("reviews").
		Select("modules.id, modules.name, modules.sequence, modules.points, modules.hours, modules.semester, students.id as student_id, students.dni as student_dni, students.full_name as student_full_name").
		Joins("INNER JOIN modules on reviews.module_id = modules.id").
		Joins("INNER JOIN students on reviews.student_id = students.id").
		Where("reviews.id = ?", review.ID).
		Scan(&moduleResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Find detailResponse
	detailResponses := make([]detailResponse, 0)
	if err := db.Table("review_details").
		Select("review_details.hours, review_details.note, review_details.note_appreciation,review_details.start_date, review_details.end_date, companies.ruc, companies.nombre_o_razon_social, companies.direccion").
		Joins("INNER JOIN companies on review_details.company_id = companies.id").
		Where("review_details.review_id = ?", review.ID).
		Scan(&detailResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, consResponse{
		Success: true,
		Module:  moduleResponses[0],
		Detail:  detailResponses,
	})
}
