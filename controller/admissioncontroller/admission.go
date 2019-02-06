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
	Search      string `json:"search"`
	CurrentPage uint   `json:"current_page"`
	Limit       uint   `json:"limit"`
    AdmissionSettingID uint `json:"admission_setting_id"`
	ProgramID   uint   `json:"program_id"`
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

	// find if exist student
	//st := models.Student{}
	//if request.Student.ID >= 1 {
	//    DB.First(&st, models.Student{ID: request.Student.ID})
	//}else {
	//    DB.First(&st, models.Student{DNI: request.Student.DNI})
	//}

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
	if err := DB.Raw("SELECT count(*) as count FROM admissions WHERE student_id = ? AND year = ? AND state = true", admission.StudentID, admission.Year).
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
    Classroom     uint      `json:"classroom"`
    Seat          uint      `json:"seat"`
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
    DB.Last(&admission,models.Admission{AdmissionSettingID: admissionSetting.ID})

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
    return c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Data: getNextClassroom{
            Classroom: admission.Classroom,
            Seat: admission.Seat,
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

type fileAdmissionResponse struct {
	Students   []models.Student  `json:"students"`
	Subsidiary models.Subsidiary `json:"subsidiary"`
	Program    models.Program    `json:"program"`
}

func FileAdmission(c echo.Context) error {
	// Get data request
	admissions := make([]models.Admission, 0)
	if err := c.Bind(&admissions); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	students := make([]models.Student, 0)

	// Query all students
	for _, admission := range admissions {
		// Query get admission all data --- current admission
		DB.First(&admission, models.Admission{ID: admission.ID})

		// Query get student all data
		student := models.Student{}
		DB.First(&student, models.Student{ID: admission.StudentID})

		// Append array
		students = append(students, student)
	}

	// Query program
	program := models.Program{}
	subsidiary := models.Subsidiary{}
	if len(admissions) >= 1 {
		DB.First(&program, models.Program{ID: admissions[0].ProgramID})
		DB.First(&subsidiary, models.Subsidiary{ID: program.SubsidiaryID})
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: fileAdmissionResponse{
			Students:   students,
			Subsidiary: subsidiary,
			Program:    program,
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

type exportAdmissionRequest struct {
	From  uint `json:"from"`
	To    uint `json:"to"`
	State bool `json:"state"`
}
type exportModels struct {
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

func ExportAdmission(c echo.Context) error {
	// Get data request
	request := exportAdmissionRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	var total uint
	exportModelss := make([]exportModels, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.student_id, admissions.program_id, admissions.state, students.dni , students.full_name, users.id as user_id, users.email, users.avatar").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN users on students.user_id = users.id").
		Where("admissions.year >= ? AND admissions.year <= ? AND admissions.state = ?", request.From, request.To, request.State).
		Order("admissions.id asc").Scan(&exportModelss).
		Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// CREATE EXCEL FILE
	excel := excelize.NewFile()

	// Create new sheet
	sheet1 := excel.NewSheet("Sheet1")

	// Set header values
	excel.SetCellValue("Sheet1", "A1", "ID")
	excel.SetCellValue("Sheet1", "B1", "Apellidos y nombre")
	excel.SetCellValue("Sheet1", "C1", "Fecha de nacimiento")
	excel.SetCellValue("Sheet1", "D1", "Sexo")
	excel.SetCellValue("Sheet1", "E1", "Año")

	//  Set values in excel file
	for i := 0; i < len(exportModelss); i++ {
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), exportModelss[i].ID)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), exportModelss[i].FullName)
	}

	// Default active sheet
	excel.SetActiveSheet(sheet1)

	// save file
	err := excel.SaveAs("temp/admission.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	// Return object
	return c.File("temp/admission.xlsx")
}

func ExportAdmissionByIds(c echo.Context) error {
	// Get data request
	admissions := make([]models.Admission, 0)
	if err := c.Bind(&admissions); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// CREATE EXCEL FILE
	excel := excelize.NewFile()

	// Create new sheet
	sheet1 := excel.NewSheet("Sheet1")

	// Set header values
	excel.SetCellValue("Sheet1", "A1", "ID")
	excel.SetCellValue("Sheet1", "B1", "Apellidos y nombre")
	excel.SetCellValue("Sheet1", "C1", "Fecha de nacimiento")
	excel.SetCellValue("Sheet1", "D1", "Sexo")
	excel.SetCellValue("Sheet1", "E1", "Año")

	// Query all students
	for key, admission := range admissions {
		// Query get admission all data --- current admission
		DB.First(&admission, models.Admission{ID: admission.ID})

		// Query get student all data
		student := models.Student{}
		DB.First(&student, models.Student{ID: admission.StudentID})

		// Set values in excel file
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", key+2), admission.ID)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", key+2), student.FullName)
	}

	// Default active sheet
	excel.SetActiveSheet(sheet1)

	// save file
	err := excel.SaveAs("temp/admission.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	// Return object
	return c.File("temp/admission.xlsx")
}

