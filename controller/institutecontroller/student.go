package institutecontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/paulantezana/review/models"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

// GetStudents function get all students
func GetStudentsPaginate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	students := make([]models.Student, 0)

	// Query in database
	if err := db.Where("lower(full_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("dni LIKE ?", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&students).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        students,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetStudentsPaginateByProgram(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	total := utilities.Counter{}
	students := make([]models.Student, 0)

	// Query in database
	DB.Raw("SELECT * FROM students "+
		"WHERE id IN (SELECT student_id FROM student_programs WHERE program_id = ?) "+
		"AND (lower(full_name) LIKE lower(?) OR dni LIKE ?) ORDER BY id desc "+
		"OFFSET ? LIMIT ?", request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%", offset, request.Limit).Scan(&students)

	// Query students count total
	DB.Raw("SELECT count(*) FROM students "+
		"WHERE id IN (SELECT student_id FROM student_programs WHERE program_id = ?) "+
		"AND (lower(full_name) LIKE lower(?) OR dni LIKE ?)", request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%").Scan(&total)

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        students,
		Total:       total.Count,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

// StudentDetail struct
type StudentDetail struct {
	CompanyName    string    `json:"company_name"`
	ModuleName     string    `json:"module_name"`
	ModuleSequence uint      `json:"module_sequence"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Hours          uint      `json:"hours"`
	Note           uint      `json:"note"`
}

// GetStudentDetailByID get student detail
func GetStudentByID(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&student, student.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    student,
	})
}

// GetStudentDetailByID get student detail
func GetStudentByDNI(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.Where("dni = ?", student.DNI).First(&student).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: "No se encontro ningun registro"})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    student,
	})
}

func GetStudentSearch(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	students := make([]models.Student, 0)
	if err := db.Where("lower(full_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("dni LIKE ?", "%"+request.Search+"%").
		Limit(10).Find(&students).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	customStudents := make([]models.Student, 0)
	for _, student := range students {
		customStudents = append(customStudents, models.Student{
			ID:       student.ID,
			FullName: student.FullName,
			DNI:      student.DNI,
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customStudents,
	})
}

type customStudentRequest struct {
	Student   models.Student `json:"student"`
	User      models.User    `json:"user"`
	ProgramID uint           `json:"program_id"`
}

func CreateStudent(c echo.Context) error {
	// Get data request
	request := customStudentRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Set program ID
	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	TX := DB.Begin()

	// has password new user account
	cc := sha256.Sum256([]byte(request.Student.DNI + "ST"))
	pwd := fmt.Sprintf("%x", cc)

	// Insert user in database
	userAccount := models.User{
		UserName: request.Student.DNI + "ST",
		Password: pwd,
		RoleID:   5,
	}
	if err := TX.Create(&userAccount).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert student in database
	request.Student.UserID = userAccount.ID
	request.Student.StudentStatusID = 1
	if err := TX.Create(&request.Student).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create relations
	if request.ProgramID >= 1 {
		studentProgram := models.StudentProgram{
			StudentID: request.Student.ID,
			ProgramID: request.ProgramID,
			ByDefault: true,
		}
		if err := TX.Create(&studentProgram).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    request.Student.ID,
		Message: fmt.Sprintf("El estudiante %s se registro correctamente", request.Student.FullName),
	})
}

func UpdateStudent(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update student in database
	rows := db.Model(&student).Update(student).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", "No se pudo actualizar")})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    student.ID,
		Message: fmt.Sprintf("Los datos del estudiante %s se actualizaron correctamente", student.FullName),
	})
}

func DeleteStudent(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete student in database
	if err := db.Delete(&student).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    student.ID,
		Message: fmt.Sprintf("El estudiante %s se elimino correctamente", student.FullName),
	})
}

type getStudentProgramsResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Level         string `json:"level"`
	SubsidiaryID  uint   `json:"subsidiary_id"`
	ByDefault     bool   `json:"by_default"`
	YearAdmission uint   `json:"year_admission"`
	YearPromotion uint   `json:"year_promotion"`
}

func GetStudentPrograms(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query
	studentPrograms := make([]getStudentProgramsResponse, 0)
	if err := DB.Table("programs").
		Select("programs.id, programs.name, programs.level, programs.subsidiary_id, student_programs.by_default, student_programs.year_admission, student_programs.year_promotion").
		Joins("INNER JOIN student_programs ON programs.id = student_programs.program_id").
		Order("programs.id desc").
		Where("student_programs.student_id = ?", student.ID).
		Scan(&studentPrograms).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    studentPrograms,
	})
}

func GetStudentHistory(c echo.Context) error {
	// Get data request
	student := models.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Find history
	studentHistory := make([]models.StudentHistory, 0)
	db.Find(&studentHistory, models.StudentHistory{StudentID: student.ID})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    studentHistory,
	})
}

// GetTempUploadStudentBySubsidiary download template
func GetTempUploadStudentBySubsidiary(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	programs := make([]models.Program, 0)
	if err := DB.Find(&programs, models.Program{SubsidiaryID: request.SubsidiaryID}).Order("id desc").Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get excel file
	fileDir := "templates/templateStudentSA.xlsx"
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

func GetTempUploadStudentByProgram(c echo.Context) error {
	// Return file excel
	return c.File("templates/templateStudent.xlsx")
}

// SetTempUploadStudent set upload student
func SetTempUploadStudentBySubsidiary(c echo.Context) error {
	// Source
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
	DB := config.GetConnection()
	defer DB.Close()

	// Prepare
	ignoreCols := 5
	counter := 0
	TX := DB.Begin()

	// Get all the rows in the student.
	rows := excel.GetRows("Student")
	for k, row := range rows {

		if k >= ignoreCols {
			// Validate required fields
			if row[0] == "" || row[1] == "" {
				break
			}

			// program id
			u, _ := strconv.ParseUint(strings.TrimSpace(row[0]), 0, 32)
			currentProgram := uint(u)

			// Create model student
			student := models.Student{
				DNI:      strings.TrimSpace(row[1]),
				FullName: strings.TrimSpace(row[2]),
				Phone:    strings.TrimSpace(row[3]),
				Gender:   strings.TrimSpace(row[5]),
				//BirthDate:          strings.TrimSpace(row[5]),
				BirthPlace:      strings.TrimSpace(row[7]),
				District:        strings.TrimSpace(row[8]),
				Province:        strings.TrimSpace(row[9]),
				Region:          strings.TrimSpace(row[10]),
				Country:         strings.TrimSpace(row[11]),
				Address:         strings.TrimSpace(row[12]),
				CivilStatus:     strings.TrimSpace(row[13]),
				IsWork:          strings.TrimSpace(row[14]),
				MarketStall:     strings.TrimSpace(row[15]),
				StudentStatusID: 1,
			}

			// has password new user account
			cc := sha256.Sum256([]byte(student.DNI + "ST"))
			pwd := fmt.Sprintf("%x", cc)

			// New Account
			userAccount := models.User{
				UserName: student.DNI + "ST",
				Email:    strings.TrimSpace(row[4]),
				Password: pwd,
				RoleID:   5,
			}

			// Insert user in database
			if err := TX.Create(&userAccount).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
			student.UserID = userAccount.ID // Set new user id

			if err := TX.Create(&student).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{
					Message: fmt.Sprintf("Ocurrió un error al insertar el alumno %s con "+
						"DNI: %s es posible que este alumno ya este en la base de datos o los datos son incorrectos, "+
						"Error: %s, no se realizo ninguna cambio en la base de datos", student.FullName, student.DNI, err),
				})
			}

			// Relation student
			studentProgram := models.StudentProgram{
				ProgramID: currentProgram,
				StudentID: student.ID,
			}
			if err := TX.Create(&studentProgram).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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

// SetTempUploadStudent set upload student
func SetTempUploadStudentByProgram(c echo.Context) error {
	// Get program ID
	idp := c.FormValue("id")
	u, _ := strconv.ParseUint(idp, 0, 32)
	currentProgramID := uint(u)

	// Source
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
	DB := config.GetConnection()
	defer DB.Close()

	// Prepare
	ignoreCols := 5
	counter := 0
	TX := DB.Begin()

	// Get all the rows in the student.
	rows := excel.GetRows("Student")
	for k, row := range rows {

		if k >= ignoreCols {
			// Validate required fields
			if row[0] == "" || row[1] == "" {
				break
			}

			// Create model student
			student := models.Student{
				DNI:      strings.TrimSpace(row[0]),
				FullName: strings.TrimSpace(row[1]),
				Phone:    strings.TrimSpace(row[2]),
				Gender:   strings.TrimSpace(row[4]),
				//BirthDate:          strings.TrimSpace(row[5]),
				BirthPlace:      strings.TrimSpace(row[6]),
				District:        strings.TrimSpace(row[7]),
				Province:        strings.TrimSpace(row[8]),
				Region:          strings.TrimSpace(row[9]),
				Country:         strings.TrimSpace(row[10]),
				Address:         strings.TrimSpace(row[11]),
				CivilStatus:     strings.TrimSpace(row[12]),
				IsWork:          strings.TrimSpace(row[13]),
				MarketStall:     strings.TrimSpace(row[14]),
				StudentStatusID: 1,
			}

			// has password new user account
			cc := sha256.Sum256([]byte(student.DNI + "ST"))
			pwd := fmt.Sprintf("%x", cc)

			// New Account
			userAccount := models.User{
				UserName: student.DNI + "ST",
				Password: pwd,
				Email:    strings.TrimSpace(row[3]),
				RoleID:   5,
			}

			// Insert user in database
			if err := TX.Create(&userAccount).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
			student.UserID = userAccount.ID // Set new user id

			if err := TX.Create(&student).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{
					Message: fmt.Sprintf("Ocurrió un error al insertar el alumno %s con "+
						"DNI: %s es posible que este alumno ya este en la base de datos o los datos son incorrectos, "+
						"Error: %s, no se realizo ninguna cambio en la base de datos", student.FullName, student.DNI, err),
				})
			}

			// Relation student
			studentProgram := models.StudentProgram{
				ProgramID: currentProgramID,
				StudentID: student.ID,
			}
			if err := TX.Create(&studentProgram).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
