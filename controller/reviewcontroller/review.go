package reviewcontroller

import (
	"fmt"
	"github.com/paulantezana/review/models/institutemodel"
	"github.com/paulantezana/review/models/reviewmodel"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

// GetReviews functions get reviews by student_id
type reviewsResponse struct {
	ID               uint      `json:"id"`
	ApprobationDate  time.Time `json:"approbation_date"`
	ModuleId         uint      `json:"module_id"`
	Name             string    `json:"name"`
	Sequence         uint      `json:"sequence"`
	Semester         string    `json:"semester"`
	TeacherID        uint      `json:"teacher_id"`
	TeacherFirstName string    `json:"teacher_first_name"`
	TeacherLastName  string    `json:"teacher_last_name"`
}

type reviewEnablesResponse struct {
	Consolidate bool `json:"consolidate"`
}

type getReviewsResponse struct {
	Message   string                `json:"message"`
	Success   bool                  `json:"success"`
	Data      interface{}           `json:"data"`
	Validates reviewEnablesResponse `json:"validates"`
}

// GetReviews functions get all reviews
func GetReviews(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	student := institutemodel.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query in database
	reviewsResponses := make([]reviewsResponse, 0)
	if err := db.Table("reviews").
		Select("reviews.id, reviews.approbation_date, modules.id as module_id, modules.name, modules.semester, modules.sequence, teachers.id as teacher_id, teachers.first_name as teacher_first_name, teachers.last_name as teacher_last_name").
		Joins("INNER JOIN modules on reviews.module_id = modules.id").
		Joins("INNER JOIN teachers on reviews.teacher_id = teachers.id").
		Order("reviews.id asc").
		Where("reviews.student_id = ?", student.ID).
		Scan(&reviewsResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// validation
	allReviews := len(reviewsResponses) // all review count
	var allModules uint                 // all modules count
	if err := db.Model(&institutemodel.Module{}).Where("program_id = ?", currentUser.DefaultProgramID).Count(&allModules).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	reviewEnablesResponse := reviewEnablesResponse{}

	if allModules == uint(allReviews) && allModules != 0 {
		reviewEnablesResponse.Consolidate = true
	}

	// Return response
	return c.JSON(http.StatusCreated, getReviewsResponse{
		Success:   true,
		Data:      reviewsResponses,
		Validates: reviewEnablesResponse,
	})
}

// CreateReview function create new review
func CreateReview(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	review := reviewmodel.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}
	review.UserID = currentUser.ID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validate
	rvw := make([]reviewmodel.Review, 0)
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
	review := reviewmodel.Review{}
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
	review := reviewmodel.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

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
	ID               uint      `json:"id" gorm:"primary_key"`
	Hours            uint      `json:"hours"`
	Note             uint      `json:"note"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	RUC              string    `json:"ruc"`
	NameSocialReason string    `json:"name_social_reason"`
	Address          string    `json:"address"`
	Phone            string    `json:"phone"`
}

type reviewResponse struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	TeacherLastName  string    `json:"teacher_last_name"`
	TeacherFirstName string    `json:"teacher_first_name"`
	ApprobationDate  time.Time `json:"approbation_date"`
	TeacherID        uint      `json:"teacher_id"`
}

// consResponse struct
type actaResponse struct {
	Success bool             `json:"success"`
	Module  moduleResponse   `json:"module"`
	Detail  []detailResponse `json:"detail"`
	Review  reviewResponse   `json:"review"`
}

// GetActaReview function get data acta
func GetActaReview(c echo.Context) error {
	// Get data request
	review := reviewmodel.Review{}
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
		Select("review_details.hours, review_details.note,review_details.start_date, review_details.end_date, companies.ruc, companies.name_social_reason, companies.address, companies.phone").
		Joins("INNER JOIN companies on review_details.company_id = companies.id").
		Where("review_details.review_id = ?", review.ID).
		Scan(&detailResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Find Review
	reviewResponses := make([]reviewResponse, 0)
	if err := db.Table("reviews").
		Select("reviews.id, reviews.approbation_date, teachers.first_name as teacher_first_name, teachers.last_name as teacher_last_name").
		Joins("INNER JOIN teachers on reviews.teacher_id = teachers.id").
		Where("reviews.id = ?", review.ID).
		Scan(&reviewResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Response data
	return c.JSON(http.StatusOK, actaResponse{
		Success: true,
		Module:  moduleResponses[0],
		Detail:  detailResponses,
		Review:  reviewResponses[0],
	})
}

// consResponse struct
type consResponse struct {
	Success bool               `json:"success"`
	Module  moduleResponse     `json:"module"`
	Detail  []detailResponse   `json:"detail"`
	Review  reviewmodel.Review `json:"review"`
}

// GetConstReview function get data constancy
func GetConstReview(c echo.Context) error {
	// Get data request
	review := reviewmodel.Review{}
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
		Select("review_details.hours, review_details.note, review_details.start_date, review_details.end_date, companies.ruc, companies.name_social_reason, companies.address, companies.phone").
		Joins("INNER JOIN companies on review_details.company_id = companies.id").
		Where("review_details.review_id = ?", review.ID).
		Scan(&detailResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// find review
	if err := db.First(&review, review.ID).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, consResponse{
		Success: true,
		Module:  moduleResponses[0],
		Detail:  detailResponses,
		Review:  review,
	})
}

type reviewDetailResponse struct {
	ID                      uint      `json:"id" gorm:"primary_key"`
	Hours                   uint      `json:"hours"`
	Note                    uint      `json:"note"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	CompanyNameSocialReason string    `json:"company_name_social_reason"`
	CompanyAddress          string    `json:"company_address"`
}

type reviewModuleResponse struct {
	ID                uint                   `json:"id"`
	ApprobationDate   time.Time              `json:"approbation_date"`
	ModuleID          uint                   `json:"module_id"`
	ModuleSequence    uint                   `json:"module_sequence"`
	ModuleName        string                 `json:"module_name"`
	ModuleDescription string                 `json:"module_description"`
	ModulePoints      uint                   `json:"module_points"`
	ModuleHours       uint                   `json:"module_hours"`
	ModuleSemester    string                 `json:"module_semester"`
	ReviewDetails     []reviewDetailResponse `json:"review_details"`
}

type consolidateResponse struct {
	Success bool                   `json:"success"`
	Student institutemodel.Student `json:"student"`
	Reviews []reviewModuleResponse `json:"reviews"`
}

// GetConsolidateReview function get data constancy
func GetConsolidateReview(c echo.Context) error {
	// Get data request
	student := institutemodel.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Find reviews
	reviewModuleResponses := make([]reviewModuleResponse, 0)
	if err := db.Table("reviews").
		Select("reviews.id, reviews.approbation_date, modules.id as module_id, modules.sequence as module_sequence, modules.name as module_name, modules.description as module_description, modules.points as module_points, modules.hours as module_hours, modules.semester as module_semester").
		Joins("INNER JOIN modules on reviews.module_id = modules.id").
		Where("reviews.student_id  = ?", student.ID).
		Scan(&reviewModuleResponses).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Find current student
	db.First(&student, student.ID)

	// consult review detail
	for key, review := range reviewModuleResponses {
		redR := make([]reviewDetailResponse, 0)
		if err := db.Table("review_details").
			Select("review_details.id, review_details.hours, review_details.note, review_details.start_date, review_details.end_date, companies.name_social_reason as company_name_social_reason, companies.address as company_address").
			Joins("INNER JOIN companies on review_details.company_id = companies.id").
			Where("review_details.review_id  = ?", review.ID).
			Scan(&redR).Error; err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		reviewModuleResponses[key].ReviewDetails = redR
	}

	return c.JSON(http.StatusOK, consolidateResponse{
		Success: true,
		Reviews: reviewModuleResponses,
		Student: student,
	})
}