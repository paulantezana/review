package institutecontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

type programUserResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProgramID uint   `json:"program_id"`
	License   bool   `json:"license"`
	Name      string `json:"name"`
}

type programUserRequest struct {
	UserID       uint `json:"user_id"`
	SubsidiaryID uint `json:"subsidiary_id"`
}

// Update
func ProgramsUserUpdate(c echo.Context) error {
    // Get data request
    programUser := models.ProgramUser{}
    if err := c.Bind(&programUser); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Update module in database
    DB.Model(&programUser).Where("id = ?", programUser.ID).UpdateColumn("license", programUser.License)

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    programUser.ID,
        Message: fmt.Sprintf("Los datos del se actualizaron correctamente"),
    })
}

// Get all programs licenses by user
func GetProgramsUserByUserID(c echo.Context) error {
	// Get data request
	request := programUserRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query Programs
	programs := make([]models.Subsidiary, 0)
	if err := DB.Raw("SELECT * FROM programs WHERE id NOT IN (SELECT program_id  FROM program_users WHERE user_id = ?) AND subsidiary_id = ?", request.UserID, request.SubsidiaryID).
		Scan(&programs).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Start Transaction
	TR := DB.Begin()

	// Insert SubsidiaryUsers
	for _, program := range programs {
		programUser := models.ProgramUser{
			UserID:    request.UserID,
			ProgramID: program.ID,
		}
		if err := TR.Create(&programUser).Error; err != nil {
			TR.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// End Transaction
	TR.Commit()

	// Query SubsidiaryUsers
	programUsers := make([]programUserResponse, 0)
	if err := DB.Table("program_users").
		Select("program_users.id, program_users.user_id, program_users.program_id, program_users.license, programs.name").
		Joins("INNER JOIN programs ON programs.id = program_users.program_id").
		Where("program_users.user_id = ? AND programs.subsidiary_id = ?", request.UserID, request.SubsidiaryID).
		Scan(&programUsers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    programUsers,
	})
}

func GetProgramsUserByStudentIDLicense(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query student
	student := models.Student{}
	if err := DB.First(&student, models.Student{UserID: currentUser.ID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query programs
	programs := make([]models.Program, 0)
	if err := DB.Debug().Raw("SELECT id, name FROM programs WHERE id "+
		"IN (SELECT program_id FROM student_programs WHERE student_id = ?)", student.ID).
		Scan(&programs).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    programs,
	})
}