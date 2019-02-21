package publiccontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

type admissionsPaginateExamResultResponse struct {
	ID         uint    `json:"id"`
	Exonerated bool    `json:"exonerated"`
	ExamNote   float32 `json:"exam_note"`
	Program    string  `json:"program"`
	DNI        string  `json:"dni"`
	FullName   string  `json:"full_name"`
	Year       uint    `json:"year"`
}

func GetAdmissionExamResults(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Required params
	//if request.SubsidiaryID == 0 {
	//    c.JSON(http.StatusOK,utilities.Response{Message: "EL parametro subsidiary_id es obligatorio"})
	//}

	// Required params
	//if request.SubsidiaryID == 0 {
	//    c.JSON(http.StatusOK,utilities.Response{Message: "EL parametro subsidiary_id es obligatorio"})
	//}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	admissionsPaginateExamResponses := make([]admissionsPaginateExamResultResponse, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.exonerated, admissions.exam_note, admissions.exam_date, programs.name as program, students.dni, students.full_name, admissions.year").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN programs ON admissions.program_id = programs.id").
		Order("admissions.exam_note desc").Scan(&admissionsPaginateExamResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admissionsPaginateExamResponses,
	})
}
