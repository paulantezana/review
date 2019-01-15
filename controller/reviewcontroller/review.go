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

type reviewsMainResponse struct {
	Validates reviewEnablesResponse `json:"validates"`
	Reviews   []reviewsResponse     `json:"reviews"`
}

type getReviewsRequest struct {
	StudentID uint `json:"student_id"`
	ProgramID uint `json:"program_id"`
}

// GetReviews functions get all reviews
func GetReviews(c echo.Context) error {
	// Get data request
	request := getReviewsRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Find StudentProgramID
	studentProgram := institutemodel.StudentProgram{}
	DB.First(&studentProgram, institutemodel.StudentProgram{StudentID: request.StudentID, ProgramID: request.ProgramID})

	// Query in database
	reviewsResponses := make([]reviewsResponse, 0)
	if err := DB.Table("reviews").
		Select("reviews.id, reviews.approbation_date, reviews.module_id, modules.name, modules.sequence, reviews.teacher_id, teachers.first_name as teacher_first_name, teachers.last_name as teacher_last_name").
		Joins("INNER JOIN modules ON reviews.module_id = modules.id").
		Joins("INNER JOIN teachers ON reviews.teacher_id = teachers.id").
		Order("reviews.id desc").
		Where("reviews.student_program_id = ?", studentProgram.ID).
		Scan(&reviewsResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// validation
	countReviews := len(reviewsResponses) // all review count
	var allModules uint                   // all modules count
	if err := DB.Model(&institutemodel.Module{}).Where("program_id = ?", request.ProgramID).Count(&allModules).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Calculate validations
	reviewEnablesResponse := reviewEnablesResponse{}
	if allModules == uint(countReviews) && allModules != 0 {
		reviewEnablesResponse.Consolidate = true
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data: reviewsMainResponse{
			Reviews:   reviewsResponses,
			Validates: reviewEnablesResponse,
		},
	})
}

type reviewRequest struct {
	ProgramID uint               `json:"program_id"`
	StudentID uint               `json:"student_id"`
	Review    reviewmodel.Review `json:"review"`
}

// CreateReview function create new review
func CreateReview(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := reviewRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	request.Review.CreatorID = currentUser.ID

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Validate
	//rvw := make([]reviewmodel.Review, 0)
	//if DB.Where("student_id = ? and module_id = ?", review.StudentID, review.ModuleId).
	//	Find(&rvw).RowsAffected >= 1 {
	//	return c.JSON(http.StatusOK, utilities.Response{
	//		Message: "Este alumno ya tiene una revision con este modulo",
	//	})
	//}

	// start transaction
	TX := DB.Begin()

	// Find StudentProgramID
	studentProgram := institutemodel.StudentProgram{}
	TX.First(&studentProgram, institutemodel.StudentProgram{StudentID: request.StudentID, ProgramID: request.ProgramID})
	request.Review.StudentProgramID = studentProgram.ID

	// Insert reviews in database
	if err := TX.Create(&request.Review).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert History student
	studentHistory := institutemodel.StudentHistory{
		StudentID:   request.StudentID,
		UserID:      currentUser.ID,
		Description: fmt.Sprintf("Revisón de prácticas del modulo %d", request.Review.ModuleId),
		Date:        time.Now(),
		Type:        1,
	}
	if err := TX.Create(&studentHistory).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    request.Review.ID,
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
type actResponse struct {
	Student institutemodel.Student `json:"student"`
	Module  institutemodel.Module  `json:"module"`
	Details []detailResponse       `json:"details"`
	Review  reviewResponse         `json:"review"`
}

// GetActaReview function get data acta
func GetActaReview(c echo.Context) error {
	// Get data request
	review := reviewmodel.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query current review
	if err := DB.First(&review).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query student program
	studentProgram := institutemodel.StudentProgram{}
	if err := DB.First(&studentProgram, institutemodel.StudentProgram{ID: review.StudentProgramID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query student
	student := institutemodel.Student{}
	if err := DB.First(&student, institutemodel.Student{ID: studentProgram.StudentID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query module
	module := institutemodel.Module{}
	if err := DB.First(&module, institutemodel.Module{ID: review.ModuleId}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find detailResponse
	detailResponses := make([]detailResponse, 0)
	if err := DB.Table("review_details").
		Select("review_details.hours, review_details.note,review_details.start_date, review_details.end_date, companies.ruc, companies.name_social_reason, companies.address, companies.phone").
		Joins("INNER JOIN companies on review_details.company_id = companies.id").
		Where("review_details.review_id = ?", review.ID).
		Scan(&detailResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find Review
	reviewResponses := make([]reviewResponse, 0)
	if err := DB.Table("reviews").
		Select("reviews.id, reviews.approbation_date, teachers.first_name as teacher_first_name, teachers.last_name as teacher_last_name").
		Joins("INNER JOIN teachers on reviews.teacher_id = teachers.id").
		Where("reviews.id = ?", review.ID).
		Scan(&reviewResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: actResponse{
			Module:  module,
			Student: student,
			Details: detailResponses,
			Review:  reviewResponses[0],
		},
	})
}

// consResponse struct
type consResponse struct {
	Student institutemodel.Student `json:"student"`
	Details []detailResponse       `json:"details"`
	Review  reviewmodel.Review     `json:"review"`
	Module  institutemodel.Module  `json:"module"`
}

// GetConstReview function get data constancy
func GetConstReview(c echo.Context) error {
	// Get data request
	review := reviewmodel.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query current review
	if err := DB.First(&review).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query student program
	studentProgram := institutemodel.StudentProgram{}
	if err := DB.First(&studentProgram, institutemodel.StudentProgram{ID: review.StudentProgramID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query student
	student := institutemodel.Student{}
	if err := DB.First(&student, institutemodel.Student{ID: studentProgram.StudentID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query module
	module := institutemodel.Module{}
	if err := DB.First(&module, institutemodel.Module{ID: review.ModuleId}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find detailResponse
	detailResponses := make([]detailResponse, 0)
	if err := DB.Table("review_details").
		Select("review_details.hours, review_details.note, review_details.start_date, review_details.end_date, companies.ruc, companies.name_social_reason, companies.address, companies.phone").
		Joins("INNER JOIN companies on review_details.company_id = companies.id").
		Where("review_details.review_id = ?", review.ID).
		Scan(&detailResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: consResponse{
			Module:  module,
			Details: detailResponses,
			Review:  review,
			Student: student,
		},
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
	Student institutemodel.Student `json:"student"`
	Reviews []reviewModuleResponse `json:"reviews"`
}

// GetConsolidateReview function get data constancy
func GetConsolidateReview(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query student program
	studentProgram := institutemodel.StudentProgram{}
	if err := DB.First(&studentProgram, institutemodel.StudentProgram{
		StudentID: request.StudentID,
		ProgramID: request.ProgramID,
	}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query student
	student := institutemodel.Student{}
	if err := DB.First(&student, institutemodel.Student{ID: studentProgram.StudentID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find reviews
	reviewModuleResponses := make([]reviewModuleResponse, 0)
	if err := DB.Table("reviews").
		Select("reviews.id, reviews.approbation_date, modules.id as module_id, modules.sequence as module_sequence, modules.name as module_name, modules.description as module_description, modules.points as module_points, modules.hours as module_hours").
		Joins("INNER JOIN modules on reviews.module_id = modules.id").
		Where("reviews.student_program_id  = ?", studentProgram.ID).
		Scan(&reviewModuleResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// consult review detail
	for key, review := range reviewModuleResponses {
		redR := make([]reviewDetailResponse, 0)
		if err := DB.Table("review_details").
			Select("review_details.id, review_details.hours, review_details.note, review_details.start_date, review_details.end_date, companies.name_social_reason as company_name_social_reason, companies.address as company_address").
			Joins("INNER JOIN companies on review_details.company_id = companies.id").
			Where("review_details.review_id  = ?", review.ID).
			Scan(&redR).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		reviewModuleResponses[key].ReviewDetails = redR
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: consolidateResponse{
			Reviews: reviewModuleResponses,
			Student: student,
		},
	})
}
