package coursescontroller

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetCourseStudentsPaginate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	courseStudents := make([]models.CourseStudent, 0)

	// Query in database
	if err := db.Where("course_id = ?", request.CourseID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&courseStudents).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        courseStudents,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func CreateCourseStudent(c echo.Context) error {
	// Get data request
	courseStudent := models.CourseStudent{}
	if err := c.Bind(&courseStudent); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validation
	VCStudent := models.CourseStudent{}
	DB.First(&VCStudent, models.CourseStudent{DNI: courseStudent.DNI, CourseID: courseStudent.CourseID})
	if VCStudent.ID != 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El estudiante %s ya está matriculado en este curso", courseStudent.FullName),
		})
	}

	// Insert courseStudents in database
	if err := DB.Create(&courseStudent).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    courseStudent.ID,
		Message: fmt.Sprintf("La participante %s se registro correctamente", courseStudent.FullName),
	})
}

func UpdateCourseStudent(c echo.Context) error {
	// Get data request
	courseStudent := models.CourseStudent{}
	if err := c.Bind(&courseStudent); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update course in database
	rows := db.Model(&courseStudent).Update(courseStudent).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", courseStudent.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    courseStudent.ID,
		Message: fmt.Sprintf("Los datos del la participante %s se actualizaron correctamente", courseStudent.FullName),
	})
}

func DeleteCourseStudent(c echo.Context) error {
	// Get data request
	courseStudent := models.CourseStudent{}
	if err := c.Bind(&courseStudent); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Delete course in database
	if err := db.Delete(&courseStudent).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    courseStudent.ID,
		Message: fmt.Sprintf("La participante %s se elimino correctamente", courseStudent.FullName),
	})
}

type actCourseStudentDetail struct {
	Student models.CourseStudent `json:"student"`
	Program models.Program       `json:"program"`
}

type actCourseStudentResponse struct {
	Course   models.Course            `json:"course"`
	Students []actCourseStudentDetail `json:"students"`
}

func ActCourseStudent(c echo.Context) error {
	// Get data request
	courseStudents := make([]models.CourseStudent, 0)
	if err := c.Bind(&courseStudents); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Prepare struct response
	actCourseStudentDetails := make([]actCourseStudentDetail, 0)

	// Query
	for _, cStudent := range courseStudents {
		// Query student
		student := models.CourseStudent{}
		DB.First(&student, models.CourseStudent{ID: cStudent.ID})

		// Query program
		program := models.Program{}
		DB.First(&program, models.Program{ID: student.ProgramID})

		// Set current student
		actCourseStudentDetail := actCourseStudentDetail{
			Student: student,
			Program: program,
		}
		actCourseStudentDetails = append(actCourseStudentDetails, actCourseStudentDetail)
	}

	course := models.Course{}
	if actCourseStudentDetails[0].Student.ID >= 1 {
		DB.First(&course, models.Course{ID: actCourseStudentDetails[0].Student.CourseID})
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: actCourseStudentResponse{
			Course:   course,
			Students: actCourseStudentDetails,
		},
	})
}

// GetTempUploadCourseStudentBySubsidiary download template
func GetTempUploadCourseStudentBySubsidiary(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	programs := make([]models.Program, 0)
	if err := DB.Find(&programs, models.Program{SubsidiaryID: request.SubsidiaryID}).Order("id desc").Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get excel file
	fileDir := "templates/templateCourseStudent.xlsx"
	excel, err := excelize.OpenFile(fileDir)
	if err != nil {
		fmt.Println(err)
	}
	excel.DeleteSheet("ProgramIDS") // Delete sheet
	excel.NewSheet("ProgramIDS")    // Create new sheet

	excel.SetCellValue("ProgramIDS", "A1", "ID")
	excel.SetCellValue("ProgramIDS", "B1", "Programa De Estudios")

	// Set styles
	excel.SetColWidth("ProgramIDS", "B", "B", 35)
	excel.SetCellStyle("ProgramIDS", "A1", "B1", 2)

	// Set data
	for i := 0; i < len(programs); i++ {
		excel.SetCellValue("ProgramIDS", fmt.Sprintf("A%d", i+2), programs[i].ID)
		excel.SetCellValue("ProgramIDS", fmt.Sprintf("B%d", i+2), programs[i].Name)
	}
	excel.SetActiveSheet(1)

	// Save excel file by the given path.
	err = excel.SaveAs(fileDir)
	if err != nil {
		fmt.Println(err)
	}

	// Return file excel
	return c.File(fileDir)
}

// SetTempUploadStudent set upload student
func SetTempUploadStudentBySubsidiary(c echo.Context) error {
	// Source
	ID := c.FormValue("course_id")
	IDu, _ := strconv.ParseUint(ID, 0, 32)
	CourseID := uint(IDu)

	// FromFile
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	auxDir := "temp/" + file.Filename
	dst, err := os.Create(auxDir)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// ---------------------
	// Read File whit Excel
	// ---------------------
	excel, err := excelize.OpenFile(auxDir)
	if err != nil {
		return err
	}

	// GET CONNECTION DATABASE
	DB := provider.GetConnection()
	defer DB.Close()

	// Prepare
	ignoreCols := 1
	counter := 0
	TX := DB.Begin()

	// Get all the rows in the student.
	rows := excel.GetRows("CourseStudents")
	for k, row := range rows {

		if k >= ignoreCols {
			// Validate required fields
			if row[0] == "" || row[1] == "" {
				break
			}

			// program id
			pri, _ := strconv.ParseUint(strings.TrimSpace(row[0]), 0, 32)
			currentProgram := uint(pri)

			yea, _ := strconv.ParseUint(strings.TrimSpace(row[5]), 0, 32)
			rowYear := uint(yea)

			note, _ := strconv.ParseFloat(strings.TrimSpace(row[6]), 32)
			rowNote := float32(note)

			// DATABASE MODELS
			// Create model student
			student := models.CourseStudent{
				DNI:       strings.TrimSpace(row[1]),
				FullName:  strings.TrimSpace(row[2]),
				Phone:     strings.TrimSpace(row[3]),
				Gender:    strings.TrimSpace(row[4]),
				Year:      rowYear,
				Note:      rowNote,
				ProgramID: currentProgram,
				CourseID:  CourseID,
			}

			// Create Student
			if err := TX.Create(&student).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{
					Message: fmt.Sprintf("Ocurrió un error al insertar el alumno %s con "+
						"DNI: %s es posible que este alumno ya este en la base de datos o los datos son incorrectos, "+
						"Error: %s, no se realizo ninguna cambio en la base de datos", student.FullName, student.DNI, err),
				})
			}

			// Counter total operations success
			counter++
		}
	}
	TX.Commit()

	// Response success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se guardo %d registros en la base de datos", counter),
	})
}
