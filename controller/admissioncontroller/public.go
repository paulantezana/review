package admissioncontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

// Show all results
func GetAdmissionExamAllResults(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Admission setting
	admissionSettings := make([]models.AdmissionSetting, 0)
	DB.Find(&admissionSettings)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admissionSettings,
	})
}

type resultsResponse struct {
	ID      uint        `json:"id"`
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

// Result by id
func GetAdmissionExamResultsById(c echo.Context) error {
	// Get data request
	admissionSetting := models.AdmissionSetting{}
	if err := c.Bind(&admissionSetting); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query Admission Setting
	DB.First(&admissionSetting, models.AdmissionSetting{ID: admissionSetting.ID})

	// Query programs
	programs := make([]models.Program, 0)
	DB.Find(&programs)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: resultsResponse{
			ID:      admissionSetting.ID,
			Name:    admissionSetting.Name,
			Content: programs,
		},
	})
}

type aERByProgramIdResponse struct {
	ID          uint    `json:"id"`
	Observation bool    `json:"observation"`
	ExamNote    float32 `json:"exam_note"`
	DNI         string  `json:"dni"`
	FullName    string  `json:"full_name"`
}
type aERByProgramIdRequest struct {
	SettingID uint `json:"setting_id"`
	ProgramID uint `json:"program_id"`
}

// Result by program
func GetAdmissionExamResultsByProgramId(c echo.Context) error {
	// Get data request
	request := aERByProgramIdRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query program
	program := models.Program{}
	DB.Find(&program, models.Program{ID: request.ProgramID})

	// Execute instructions
	aERByProgramIdResponses := make([]aERByProgramIdResponse, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.exonerated, admissions.exam_note,  students.dni, students.full_name").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Where("admissions.program_id = ? AND admissions.admission_setting_id = ?", request.ProgramID, request.SettingID).
		Order("admissions.exam_note desc").Scan(&aERByProgramIdResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: resultsResponse{
			ID:      program.ID,
			Name:    program.Name,
			Content: aERByProgramIdResponses,
		},
	})
}

// Subsidiary institute detail
func GetSubsidiariesDetail(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	subsidiaries := make([]models.Subsidiary, 0)
	if err := DB.Order("id desc").Find(&subsidiaries).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find programs
	for i := range subsidiaries {
		DB.Model(&subsidiaries[i]).Related(&subsidiaries[i].Programs)
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    subsidiaries,
	})
}
