package admissioncontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

type admissionsPaginateResponse struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`

	StudentID uint `json:"student_id"`
	ProgramID uint `json:"program_id"`
	UserID    uint `json:"user_id"`

	State bool `json:"state"`

	DNI      string `json:"dni"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type admissionsPaginateRequest struct {
	Search             string `json:"search"`
	CurrentPage        uint   `json:"current_page"`
	Limit              uint   `json:"limit"`
	AdmissionSettingID uint   `json:"admission_setting_id"`
	ProgramID          uint   `json:"program_id"`
}

func (r *admissionsPaginateRequest) validate() uint {
	con := config.GetConfig()
	if r.Limit == 0 {
		r.Limit = con.Global.Paginate
	}
	if r.CurrentPage == 0 {
		r.CurrentPage = 1
	}
	offset := r.Limit*r.CurrentPage - r.Limit
	return offset
}

func GetAdmissionsPaginate(c echo.Context) error {
	// Get data request
	request := admissionsPaginateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Pagination calculate
	offset := request.validate()

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	var total uint
	admissionsPaginateResponses := make([]admissionsPaginateResponse, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.student_id, admissions.program_id, admissions.state, students.dni , students.full_name, users.id as user_id, users.email, users.avatar").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN users on students.user_id = users.id").
		Where("students.dni LIKE ? AND admissions.admission_setting_id = ? AND admissions.program_id = ?", "%"+request.Search+"%", request.AdmissionSettingID, request.ProgramID).
		Or("lower(students.full_name) LIKE lower(?) AND admissions.admission_setting_id = ? AND admissions.program_id = ?", "%"+request.Search+"%", request.AdmissionSettingID, request.ProgramID).
		Order("admissions.id desc").
		Offset(offset).Limit(request.Limit).Scan(&admissionsPaginateResponses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.ResponsePaginate{
		Success:     true,
		Data:        admissionsPaginateResponses,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

type admissionsPaginateExamResponse struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Observation string `json:"observation"`
	Exonerated  bool   `json:"exonerated"`

	ExamNote float32   `json:"exam_note"`
	ExamDate time.Time `json:"exam_date"`

	StudentID uint `json:"student_id"`
	ProgramID uint `json:"program_id"`
	UserID    uint `json:"user_id"`

	State bool `json:"state"`

	DNI      string `json:"dni"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type admissionModelRequest struct {
	Admission models.Admission `json:"admission"`
	Student   models.Student   `json:"student"`
	User      models.User      `json:"user"`
}

func GetAdmissionsByID(c echo.Context) error {
	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query admission
	DB.First(&admission)

	// Query Student
	student := models.Student{}
	DB.First(&student, models.Student{ID: admission.StudentID})

	// Query user
	user := models.User{}
	DB.First(&user, models.User{ID: student.UserID})

	// Reset response
	user.Password = ""
	user.Key = ""

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: admissionModelRequest{
			Student:   student,
			Admission: admission,
			User:      user,
		},
	})
}

func GetAdmissionsPaginateExam(c echo.Context) error {
	// Get data request
	request := admissionsPaginateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Pagination calculate
	offset := request.validate()

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	var total uint
	admissionsPaginateExamResponses := make([]admissionsPaginateExamResponse, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exam_note, admissions.exam_date, admissions.exonerated, admissions.admission_date, admissions.year, admissions.student_id, admissions.program_id, admissions.state, students.dni , students.full_name, users.id as user_id, users.email, users.avatar").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN users on students.user_id = users.id").
		Where("students.dni LIKE ? AND admissions.admission_setting_id = ? AND admissions.program_id = ? AND admissions.state = true", "%"+request.Search+"%", request.AdmissionSettingID, request.ProgramID).
		Or("lower(students.full_name) LIKE lower(?) AND admissions.admission_setting_id = ? AND admissions.program_id = ? AND admissions.state = true", "%"+request.Search+"%", request.AdmissionSettingID, request.ProgramID).
		Order("admissions.id desc").
		Offset(offset).Limit(request.Limit).Scan(&admissionsPaginateExamResponses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.ResponsePaginate{
		Success:     true,
		Data:        admissionsPaginateExamResponses,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

type updateStudentAdmissionRequest struct {
	Student models.Student `json:"student"`
	User    models.User    `json:"user"`
}

func UpdateStudentAdmission(c echo.Context) error {
	// Get data request
	request := updateStudentAdmissionRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	TX := DB.Begin()

	// Validate if exist student
	if request.Student.ID == 0 {
		// has password new user account
		cc := sha256.Sum256([]byte(request.Student.DNI + "ST"))
		pwd := fmt.Sprintf("%x", cc)

		// Insert user in database
		request.User.UserName = request.Student.DNI + "ST"
		request.User.Password = pwd
		request.User.RoleID = 5

		if err := TX.Create(&request.User).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Insert student in database
		request.Student.UserID = request.User.ID
		request.Student.StudentStatusID = 2
		if err := TX.Create(&request.Student).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		// Update data
		rows := TX.Model(&request.Student).Update(request.Student).RowsAffected
		if rows == 0 {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", "No se pudo actualizar los datos del es")})
		}
		rows = TX.Model(&request.User).Update(request.User).RowsAffected
		if rows == 0 {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", "No se pudo actualizar los datos del es")})
		}

		// Query data user
		DB.First(&request.User, models.User{ID: request.User.ID})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    request,
		Message: fmt.Sprintf("El estudiante %s se registro correctamente", request.Student.FullName),
	})
}

func CreateAdmission(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	student := models.Student{}
	DB.First(&student, models.Student{ID: admission.StudentID})

	// Validation admission
	countV := utilities.Counter{}
	if err := DB.Raw("SELECT count(*) as count FROM admissions WHERE student_id = ? AND admission_setting_id = ? AND state = true", admission.StudentID, admission.AdmissionSettingID).
		Scan(&countV).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if countV.Count >= 1 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El estudiante %s ya esta registrado en el proceso de admision del año %d", student.FullName, admission.Year),
		})
	}

	// start transaction
	TX := DB.Begin()

	// Insert admission
	admission.AdmissionDate = time.Now()
	admission.UserID = currentUser.ID
	if err := TX.Create(&admission).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	//  Update relations   StudentProgram by default
	TX.Exec("UPDATE student_programs SET by_default = false WHERE student_id = ?", admission.StudentID)

	// Insert new Relation program and student
	studentProgram := models.StudentProgram{
		StudentID:     admission.StudentID,
		ProgramID:     admission.ProgramID,
		ByDefault:     true,
		YearAdmission: admission.Year,
	}
	if err := TX.Create(&studentProgram).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert student history
	studentHistory := models.StudentHistory{
		StudentID:   admission.StudentID,
		UserID:      currentUser.ID,
		Description: fmt.Sprintf("Admision"),
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
		Data:    admission,
		Message: fmt.Sprintf("El estudiante %s se registro correctamente", student.FullName),
	})
}

func UpdateAdmission(c echo.Context) error {
	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update admission
	rows := db.Model(&admission).Update(admission).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", admission.ID),
		})
	}

	// Query student
	db.First(&admission, models.Admission{ID: admission.ID})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admission,
		Message: fmt.Sprintf("Los datos del admision %d se actualizaron correctamente", admission.ID),
	})
}

func CancelAdmission(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	TX := DB.Begin()

	// Execute query
	if err := TX.Model(admission).UpdateColumn("state", false).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find admission details
	TX.First(&admission)

	// Insert new state student
	studentHistory := models.StudentHistory{
		Description: "Admisión anulada",
		StudentID:   admission.StudentID,
		UserID:      currentUser.ID,
		Date:        time.Now(),
		Type:        2,
	}
	if err := TX.Create(&studentHistory).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admission.ID,
		Message: fmt.Sprintf("Se anuló la admision con el id %d", admission.ID),
	})
}

type getNextClassroom struct {
	Classroom uint `json:"classroom"`
	Seat      uint `json:"seat"`
}

func GetNextClassroomAdmission(c echo.Context) error {
	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// get admission setting
	admissionSetting := models.AdmissionSetting{}
	DB.First(&admissionSetting, models.AdmissionSetting{ID: admission.AdmissionSettingID})

	// query las admission
	DB.Last(&admission, models.Admission{AdmissionSettingID: admissionSetting.ID})

	// increment
	admission.Seat++

	// Calculate validations
	if admission.Classroom == 0 {
		admission.Classroom = 1
	}
	if admission.Seat > admissionSetting.Seats {
		admission.Seat = 1
		admission.Classroom++
	}

	// return query
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: getNextClassroom{
			Classroom: admission.Classroom,
			Seat:      admission.Seat,
		},
	})
}

func UpdateExamAdmission(c echo.Context) error {
	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Set current time now
	admission.ExamDate = time.Now()

	// Update module in database
	rows := db.Model(&admission).Update(admission).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", admission.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admission.ID,
		Message: fmt.Sprintf("Los datos del modulo %d se actualizaron correctamente", admission.ID),
	})
}

type aStudentDFResponse struct {
	ID          uint      `json:"id"`
	DNI         string    `json:"dni"`
	FullName    string    `json:"full_name"`
	Phone       string    `json:"phone"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	BirthDate   time.Time `json:"birth_date"`
	BirthPlace  string    `json:"birth_place"`
	Country     string    `json:"country"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	Region      string    `json:"region"`
	MarketStall string    `json:"market_stall"`
	CivilStatus string    `json:"civil_status"`
	IsWork      string    `json:"is_work"` // y si || n = no
}

type aDFResponse struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`

	StudentID uint `json:"-"`
	ProgramID uint `json:"-"`

	Student aStudentDFResponse `json:"student"`
	Program models.Program     `json:"program"`
}

type fileADFResponse struct {
	Subsidiary models.Subsidiary `json:"subsidiary"`
	Admissions []aDFResponse     `json:"admissions"`
}

func FileAdmissionDF(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Response slice
	aDFResponses := make([]aDFResponse, 0)

	// Query all students
	for _, ID := range request.IDs {
		// Query admission
		aDFResponse := aDFResponse{}
		DB.Raw("SELECT * FROM admissions WHERE id = ? LIMIT 1", ID).Scan(&aDFResponse)

		// Query student
		aStudentDFResponse := aStudentDFResponse{}
		DB.Raw("SELECT * FROM students WHERE id = ? LIMIT 1", aDFResponse.StudentID).Scan(&aStudentDFResponse)
		aDFResponse.Student = aStudentDFResponse

		// Query program
		program := models.Program{}
		DB.First(&program, models.Program{ID: aDFResponse.ProgramID})
		aDFResponse.Program = program

		// Append array
		aDFResponses = append(aDFResponses, aDFResponse)
	}

	// Query program
	subsidiary := models.Subsidiary{}
	DB.First(&subsidiary, models.Subsidiary{ID: request.SubsidiaryID})

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: fileADFResponse{
			Admissions: aDFResponses,
			Subsidiary: subsidiary,
		},
	})
}

type licenceADF struct {
	ID            uint      `json:"id"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`
	Classroom     uint      `json:"classroom"`
	Seat          uint      `json:"seat"`

	DNI      string `json:"dni"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Program  string `json:"program"`
}

func LicenseAdmissionDF(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query all students
	licenceADFs := make([]licenceADF, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.classroom, admissions.seat, "+
			"students.dni, students.full_name, programs.name as program, "+
			"users.avatar").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN users ON students.user_id = users.id").
		Joins("INNER JOIN programs ON admissions.program_id = programs.id").
		Where("admissions.id IN (?)", request.IDs).
		Scan(&licenceADFs).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    licenceADFs,
	})
}

type listADF struct {
	ID            uint      `json:"id"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`
	Classroom     uint      `json:"classroom"`
	Seat          uint      `json:"seat"`

	DNI      string `json:"dni"`
	FullName string `json:"full_name"`
	Program  string `json:"program"`
}

func ListAdmissionDF(c echo.Context) error {
	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query all students
	listADFs := make([]listADF, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.classroom, admissions.seat, "+
			"students.dni, students.full_name, programs.name as program ").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN programs ON admissions.program_id = programs.id").
		Where("admissions.admission_setting_id = ? AND admissions.state = true", admission.AdmissionSettingID).
		Scan(&listADFs).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    listADFs,
	})
}

func ExportAdmission(c echo.Context) error {
	// Get data request
	request := models.Admission{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// get admissions
	admissions := make([]models.Admission, 0)
	DB.Where("state = ? AND admission_setting_id = ?", request.State, request.AdmissionSettingID).Find(&admissions)

	// Create file excel
	file := exportExcel(admissions)

	// Return object
	return c.File(file)
}

func ExportAdmissionByIds(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// get admissions
	admissions := make([]models.Admission, 0)
	DB.Where("id IN (?)", request.IDs).Find(&admissions)

	// Create file excel
	file := exportExcel(admissions)

	// Return object
	return c.File(file)
}

func ReportAdmissionGeneral(c echo.Context) error {
	// Get data request
	request := models.Admission{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	reports := make([]listADF, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.year, admissions.admission_date, admissions.exonerated, admissions.state, "+
			"students.full_name, students.dni, programs.name as program ").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN programs ON admissions.program_id = programs.id").
		Where("admissions.admission_setting_id = ?", request.AdmissionSettingID).
		Scan(&reports).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    reports,
	})
}

func ExportAdmissionExamResults(c echo.Context) error {
    // Get data request
    request := models.AdmissionSetting{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Details admission settings
    if err := DB.First(&request,models.AdmissionSetting{ID: request.ID}).Error; err != nil {
        return err
    }

    // Query programs
    programs := make([]models.Program,0)
    if err := DB.Raw("SELECT * FROM programs WHERE subsidiary_id = ?", request.SubsidiaryID).Scan(&programs).Error; err != nil {
        return err
    }

    // CREATE EXCEL FILE
    excel := excelize.NewFile()

    // Create sheets
    for _, program := range programs {
        // Create sheet name
        sheetName := program.Name

        // Create new sheet
        excel.NewSheet(sheetName)

        // Query all admission by admission setting
        admissions := make([]models.Admission, 0)
        if err := DB.Where("program_id = ? AND admission_setting_id = ?", program.ID, request.ID).Find(&admissions).Error; err != nil {
           return err
        }

        // Set header values
        excel.SetCellValue(sheetName, "A1", "ID")
        excel.SetCellValue(sheetName, "B1", "DNI")
        excel.SetCellValue(sheetName, "C1", "Apellidos y Nombres")
        excel.SetCellValue(sheetName, "D1", "Nota")

        // Format style sheets
        excel.SetColWidth(sheetName,"B","B",10)
        excel.SetColWidth(sheetName,"C","C",35)
        excel.SetColWidth(sheetName,"D","D",8)

        for key, admission := range admissions {
            // Query get student all data
            student := models.Student{}
            DB.First(&student, models.Student{ID: admission.StudentID})

            // Fills data
            excel.SetCellValue(sheetName, fmt.Sprintf("A%d", key+2), admission.ID)
            excel.SetCellValue(sheetName, fmt.Sprintf("B%d", key+2), student.DNI)
            excel.SetCellValue(sheetName, fmt.Sprintf("C%d", key+2), student.FullName)
            excel.SetCellValue(sheetName, fmt.Sprintf("D%d", key+2), admission.ExamNote)
        }
    }

    // Default sheet active
    excel.SetActiveSheet(1)

    // save file
    err := excel.SaveAs("temp/admission.xlsx")
    if err != nil {
        fmt.Println(err)
    }

    // Return string directory
    return c.File("temp/admission.xlsx")
}

func exportExcel(admissions []models.Admission) string {
	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// CREATE EXCEL FILE
	excel := excelize.NewFile()

	// Create new sheet
	sheet1 := excel.NewSheet("Sheet1")

	// Set header values
	excel.SetCellValue("Sheet1", "A1", "ID")
	excel.SetCellValue("Sheet1", "B1", "Programa de estudios")
	excel.SetCellValue("Sheet1", "C1", "DNI")
	excel.SetCellValue("Sheet1", "D1", "Apellidos y Nombres")
	excel.SetCellValue("Sheet1", "E1", "Celular")
	excel.SetCellValue("Sheet1", "F1", "Email")
	excel.SetCellValue("Sheet1", "G1", "Sexo")
	excel.SetCellValue("Sheet1", "H1", "Fecha Nacimiento")
	excel.SetCellValue("Sheet1", "I1", "Lugar de nacimiento")
	excel.SetCellValue("Sheet1", "J1", "Distrito")
	excel.SetCellValue("Sheet1", "K1", "Provincia")
	excel.SetCellValue("Sheet1", "L1", "Región")
	excel.SetCellValue("Sheet1", "M1", "Pais")
	excel.SetCellValue("Sheet1", "N1", "Direccion")
	excel.SetCellValue("Sheet1", "O1", "Estado civil")
	excel.SetCellValue("Sheet1", "P1", "Trabaja")
	excel.SetCellValue("Sheet1", "Q1", "Puesto")

	// query
	for key, admission := range admissions {
		// Query get student all data
		student := models.Student{}
		DB.First(&student, models.Student{ID: admission.StudentID})

		// Query user
		user := models.User{}
		DB.First(&user, models.User{ID: student.UserID})

		// Query user
		program := models.Program{}
		DB.First(&program, models.Program{ID: admission.ProgramID})

		// Set values in excel file
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", key+2), admission.ID)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", key+2), program.Name)
		excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", key+2), student.DNI)
		excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", key+2), student.FullName)
		excel.SetCellValue("Sheet1", fmt.Sprintf("E%d", key+2), student.Phone)
		excel.SetCellValue("Sheet1", fmt.Sprintf("F%d", key+2), user.Email)
		excel.SetCellValue("Sheet1", fmt.Sprintf("G%d", key+2), student.Gender)
		excel.SetCellValue("Sheet1", fmt.Sprintf("H%d", key+2), student.BirthDate)
		excel.SetCellValue("Sheet1", fmt.Sprintf("I%d", key+2), student.BirthPlace)
		excel.SetCellValue("Sheet1", fmt.Sprintf("J%d", key+2), student.District)
		excel.SetCellValue("Sheet1", fmt.Sprintf("K%d", key+2), student.Province)
		excel.SetCellValue("Sheet1", fmt.Sprintf("L%d", key+2), student.Region)
		excel.SetCellValue("Sheet1", fmt.Sprintf("M%d", key+2), student.Country)
		excel.SetCellValue("Sheet1", fmt.Sprintf("N%d", key+2), student.Address)
		excel.SetCellValue("Sheet1", fmt.Sprintf("O%d", key+2), student.CivilStatus)
		excel.SetCellValue("Sheet1", fmt.Sprintf("P%d", key+2), student.IsWork)
		excel.SetCellValue("Sheet1", fmt.Sprintf("Q%d", key+2), student.MarketStall)
	}

	// Default active sheet
	excel.SetActiveSheet(sheet1)

	// save file
	err := excel.SaveAs("temp/admission.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	// Return string directory
	return "temp/admission.xlsx"
}
