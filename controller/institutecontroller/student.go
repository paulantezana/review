package institutecontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/institutemodel"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

// GetStudents function get all students
func GetStudents(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

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
	students := make([]institutemodel.Student, 0)

	if currentUser.RoleID == 1 {
		// Query in database
		if err := db.Where("lower(full_name) LIKE lower(?)", "%"+request.Search+"%").
			Or("dni LIKE ?", "%"+request.Search+"%").
			Order("id asc").
			Offset(offset).Limit(request.Limit).Find(&students).
			Offset(-1).Limit(-1).Count(&total).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		// Query in database
		if err := db.Where("lower(full_name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", currentUser.DefaultProgramID).
			Or("dni LIKE ? AND program_id = ?", "%"+request.Search+"%", currentUser.DefaultProgramID).
			Order("id asc").
			Offset(offset).Limit(request.Limit).Find(&students).
			Offset(-1).Limit(-1).Count(&total).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customStudent := make([]institutemodel.Student, 0)
		for _, student := range students {
			customStudent = append(customStudent, institutemodel.Student{
				ID:       student.ID,
				FullName: student.FullName,
			})
		}
		return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
			Success:     true,
			Data:        customStudent,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
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

// StudentDetailResponse response struct
type StudentDetailResponse struct {
	StudentDetail []StudentDetail        `json:"student_detail"`
	Student       institutemodel.Student `json:"student"`
}

// GetStudentDetailByID get student detail
func GetStudentDetailByID(c echo.Context) error {
	// Get data request
	student := institutemodel.Student{}
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

	// Find quotations in database by RequirementID  ========== Quotations, Providers, Users
	//StudentDetails := make([]StudentDetail, 0)
	//if err := db.Table("reviews").
	//	Select("companies.nombre_o_razon_social as company_name, modules.name as module_name, modules.sequence as module_sequence, review_details.start_date, review_details.end_date, review_details.note, review_details.hours").
	//	Joins("INNER JOIN review_details on reviews.id = review_details.review_id").
	//	Joins("INNER JOIN companies on review_details.company_id = companies.id").
	//	Joins("INNER JOIN modules on reviews.module_id = modules.id").
	//	Order("modules.sequence asc").
	//	Where("reviews.student_id = ?", student.ID).
	//	Scan(&StudentDetails).Error; err != nil {
	//		return c.NoContent(http.StatusInternalServerError)
	//}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    student,
	})
}

// GetStudentDetailByID get student detail
func GetStudentDetailByDNI(c echo.Context) error {
	// Get data request
	student := institutemodel.Student{}
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
	students := make([]institutemodel.Student, 0)
	if err := db.Where("lower(full_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("dni LIKE ?", "%"+request.Search+"%").
		Limit(10).Find(&students).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	customStudents := make([]institutemodel.Student, 0)
	for _, student := range students {
		customStudents = append(customStudents, institutemodel.Student{
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

func CreateStudent(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	student := institutemodel.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}

	// Set program ID
	if student.DefaultProgramID == 0 {
		student.DefaultProgramID = currentUser.DefaultProgramID
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	tx := DB.Begin()

	// has password new user account
	cc := sha256.Sum256([]byte(student.DNI + "ST"))
	pwd := fmt.Sprintf("%x", cc)

	// Insert user in database
	userAccount := models.User{
		UserName: student.DNI + "ST",
		Password: pwd,
		RoleID:   5,
	}
	if err := tx.Create(&userAccount).Error; err != nil {
		tx.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert student in database
	student.UserID = userAccount.ID
	student.StudentStatusID = 1
	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	tx.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    student.ID,
		Message: fmt.Sprintf("El estudiante %s se registro correctamente", student.FullName),
	})
}

func UpdateStudent(c echo.Context) error {
	// Get data request
	student := institutemodel.Student{}
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
	student := institutemodel.Student{}
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

// GetTempUploadStudent dowloand template
func GetTempUploadStudent(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Return file sa
	if currentUser.RoleID == 1 {
		fileDir := "templates/templateStudentSA.xlsx"
		xlsx, err := excelize.OpenFile(fileDir)
		if err != nil {
			fmt.Println(err)
		}
		xlsx.NewSheet("ProgramIDS")

		// get connection
		db := config.GetConnection()
		defer db.Close()

		// Execute instructions
		programs := make([]institutemodel.Program, 0)
		if err := db.Find(&programs).Order("id desc").Error; err != nil {
			return err
		}

		xlsx.SetCellValue("ProgramIDS", "A1", "ID")
		xlsx.SetCellValue("ProgramIDS", "B1", "Programa De Estudios")

		for i := 0; i < len(programs); i++ {
			xlsx.SetCellValue("ProgramIDS", fmt.Sprintf("A%d", i+2), programs[i].ID)
			xlsx.SetCellValue("ProgramIDS", fmt.Sprintf("B%d", i+2), programs[i].Name)
		}
		xlsx.SetActiveSheet(1)

		// Save xlsx file by the given path.
		err = xlsx.SaveAs(fileDir)
		if err != nil {
			fmt.Println(err)
		}

		return c.File(fileDir)
	}

	// Return file admin
	return c.File("templates/templateStudent.xlsx")
}

// SetTempUploadStudent set upload student
func SetTempUploadStudent(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

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
	xlsx, err := excelize.OpenFile(auxDir)
	if err != nil {
		return err
	}

	// Prepare
	students := make([]institutemodel.Student, 0)
	ignoreCols := 1

	// Get all the rows in the student.
	rows := xlsx.GetRows("student")
	for k, row := range rows {

		if k >= ignoreCols {
			// Validate required fields
			if row[0] == "" {
				break
			}

			// program id
			var currentProgram uint
			currentProgram = currentUser.DefaultProgramID

			if currentProgram == 0 {
				u, _ := strconv.ParseUint(strings.TrimSpace(row[5]), 0, 32)
				currentProgram = uint(u)
			}

			ay, _ := strconv.ParseUint(strings.TrimSpace(row[3]), 0, 32)
			py, _ := strconv.ParseUint(strings.TrimSpace(row[4]), 0, 32)

			admissionYear := uint(ay)
			promotionYear := uint(py)

			students = append(students, institutemodel.Student{
				DNI:              strings.TrimSpace(row[0]),
				FullName:         strings.TrimSpace(row[1]),
				Phone:            strings.TrimSpace(row[2]),
				AdmissionYear:    admissionYear,
				PromotionYear:    promotionYear,
				DefaultProgramID: currentProgram,
				StudentStatusID: 1,
			})
		}
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert students in database
	tr := db.Begin()
	for _, student := range students {

		// has password new user account
		cc := sha256.Sum256([]byte(student.DNI + "ST"))
		pwd := fmt.Sprintf("%x", cc)

		// New Account
		userAccount := models.User{
			UserName: student.DNI + "ST",
			Password: pwd,
			RoleID:   5,
		}

		// Insert user in database
		if err := tr.Create(&userAccount).Error; err != nil {
			tr.Rollback()
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		student.UserID = userAccount.ID
		if err := tr.Create(&student).Error; err != nil {
			tr.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{
				Message: fmt.Sprintf("Ocurrió un error al insertar el alumno %s con "+
					"DNI: %s es posible que este alumno ya este en la base de datos o los datos son incorrectos, "+
					"Error: %s, no se realizo ninguna cambio en la base de datos", student.FullName, student.DNI, err),
			})
		}
	}
	tr.Commit()

	// Response success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se guardo %d registros den la base de datos", len(students)),
	})
}
