package admissioncontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/admissionmodel"
	"github.com/paulantezana/review/models/institutemodel"
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
	Year        uint   `json:"year"`
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
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

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
		Where("students.dni LIKE ? AND admissions.year = ? AND admissions.program_id = ?", "%"+request.Search+"%", request.Year, request.ProgramID).
		Or("lower(students.full_name) LIKE lower(?) AND admissions.year = ? AND admissions.program_id = ?", "%"+request.Search+"%", request.Year, request.ProgramID).
		Order("admissions.id desc").
		Offset(offset).Limit(request.Limit).Scan(&admissionsPaginateResponses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
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

func GetAdmissionsByID(c echo.Context) error {
	// Get data request
	admission := admissionmodel.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query admission
	DB.First(&admission)

	// Query Student
	student := institutemodel.Student{}
	DB.First(&student, institutemodel.Student{ID: admission.StudentID})

	// Query user
	user := models.User{}
	DB.First(&user, models.User{ID: student.UserID})

	// Reset response
	user.Password = ""
	user.Key = ""

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data: createAdmissionRequest{
			Student:   student,
			Admission: admission,
			User:      user,
		},
	})
}

func GetAdmissionsPaginateExam(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

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
		Where("students.dni LIKE ? AND admissions.year = ? AND admissions.program_id = ? AND admissions.state = true", "%"+request.Search+"%", request.Year, request.ProgramID).
		Or("lower(students.full_name) LIKE lower(?) AND admissions.year = ? AND admissions.program_id = ? AND admissions.state = true", "%"+request.Search+"%", request.Year, request.ProgramID).
		Order("admissions.id desc").
		Offset(offset).Limit(request.Limit).Scan(&admissionsPaginateExamResponses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
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

type createAdmissionRequest struct {
	Student   institutemodel.Student   `json:"student"`
	Admission admissionmodel.Admission `json:"admission"`
	User      models.User              `json:"user"`
}

type countValidate struct {
	Count uint
}

func CreateAdmission(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := createAdmissionRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Init vars
	currentYear := uint(time.Now().Year())

	// Validation

	countV := countValidate{}
	if err := DB.Raw("SELECT count(*) as count FROM admissions WHERE student_id IN (SELECT id FROM students WHERE dni = ?) AND year = ? AND state = true", request.Student.DNI, currentYear).
		Scan(&countV).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if countV.Count >= 1 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El estudiante %s ya esta registrado en el proceso de admision del año %d", request.Student.FullName, currentYear),
		})
	}

	// start transaction
	TX := DB.Begin()

	// find if exist student
	st := institutemodel.Student{}
	DB.First(&st, institutemodel.Student{DNI: request.Student.DNI})

	// Validate if exist student
	if st.ID == 0 {
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
		// Set new student ID
		request.Admission.StudentID = request.Student.ID
	} else {
		// Set current student ID
		request.Admission.StudentID = st.ID
		request.Student.ID = st.ID
		request.Student.UserID = st.UserID
		request.User.ID = st.UserID

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

	// Insert admission
	request.Admission.AdmissionDate = time.Now()
	request.Admission.Year = currentYear
	request.Admission.UserID = currentUser.ID
	if err := TX.Create(&request.Admission).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	//  Update relations   StudentProgram by default
	TX.Exec("UPDATE student_programs SET by_default = false WHERE student_id = ?", request.Student.ID)

	// Insert new Relation program and student
	studentProgram := institutemodel.StudentProgram{
		StudentID: request.Student.ID,
		ProgramID: request.Admission.ProgramID,
		ByDefault: true,
	}
	if err := TX.Create(&studentProgram).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert student history
	studentHistory := institutemodel.StudentHistory{
		StudentID:   request.Student.ID,
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

	// Reset Keys and fields
	request.User.Password = ""
	request.User.Key = ""

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    request,
		Message: fmt.Sprintf("El estudiante %s se registro correctamente", request.Student.FullName),
	})
}

func UpdateAdmission(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	request := createAdmissionRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update student
	rows := db.Model(&request.Student).Update(request.Student).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", request.Student.ID),
		})
	}

	// Update student
	rows = db.Model(&request.User).Update(request.User).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", request.User.ID),
		})
	}

	// Update admission
	rows = db.Model(&request.Admission).Update(request.Admission).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", request.Admission.ID),
		})
	}

	// Query student
	db.First(&request.Student, institutemodel.Student{ID: request.Student.ID})
	db.First(&request.Admission, admissionmodel.Admission{ID: request.Admission.ID})
	db.First(&request.User, models.User{ID: request.User.ID})

	// Reset Keys and fields
	request.User.Password = ""
	request.User.Key = ""

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    request,
		Message: fmt.Sprintf("Los datos del admision %d se actualizaron correctamente", request.Admission.ID),
	})
}

func CancelAdmission(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	admission := admissionmodel.Admission{}
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
	studentHistory := institutemodel.StudentHistory{
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

func UpdateExamAdmission(c echo.Context) error {
	// Get data request
	admission := admissionmodel.Admission{}
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
	Students   []institutemodel.Student  `json:"students"`
	Subsidiary institutemodel.Subsidiary `json:"subsidiary"`
	Program    institutemodel.Program    `json:"program"`
}

func FileAdmission(c echo.Context) error {
	// Get data request
	admissions := make([]admissionmodel.Admission, 0)
	if err := c.Bind(&admissions); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	students := make([]institutemodel.Student, 0)

	// Query all students
	for _, admission := range admissions {
		// Query get admission all data --- current admission
		DB.First(&admission, admissionmodel.Admission{ID: admission.ID})

		// Query get student all data
		student := institutemodel.Student{}
		DB.First(&student, institutemodel.Student{ID: admission.StudentID})

		// Append array
		students = append(students, student)
	}

	// Query program
	program := institutemodel.Program{}
	subsidiary := institutemodel.Subsidiary{}
	if len(admissions) >= 1 {
		DB.First(&program, institutemodel.Program{ID: admissions[0].ProgramID})
		DB.First(&subsidiary, institutemodel.Subsidiary{ID: program.SubsidiaryID})
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

type exportAdmissionRequest struct {
	From  uint `json:"from"`
	To    uint `json:"to"`
	State bool `json:"state"`
}
type exportAdmissionModel struct {
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
	exportAdmissionModels := make([]exportAdmissionModel, 0)
	if err := DB.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.student_id, admissions.program_id, admissions.state, students.dni , students.full_name, users.id as user_id, users.email, users.avatar").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN users on students.user_id = users.id").
		Where("admissions.year >= ? AND admissions.year <= ? AND admissions.state = ?", request.From, request.To, request.State).
		Order("admissions.id asc").Scan(&exportAdmissionModels).
		Count(&total).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
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
	for i := 0; i < len(exportAdmissionModels); i++ {
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), exportAdmissionModels[i].ID)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), exportAdmissionModels[i].FullName)
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
	admissions := make([]admissionmodel.Admission, 0)
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
		DB.First(&admission, admissionmodel.Admission{ID: admission.ID})

		// Query get student all data
		student := institutemodel.Student{}
		DB.First(&student, institutemodel.Student{ID: admission.StudentID})

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
