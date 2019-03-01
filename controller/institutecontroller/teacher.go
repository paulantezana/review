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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

func GetTeachers(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

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
	teachers := make([]models.Teacher, 0)

	if err := db.Where("lower(first_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("lower(last_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("dni LIKE ?", "%"+request.Search+"%").
		Order("id asc").
		Offset(offset).Limit(request.Limit).Find(&teachers).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get type teacher
	for k, teacher := range teachers {
		teacherProgram := models.TeacherProgram{}
		db.First(&teacherProgram, models.TeacherProgram{
			TeacherID: teacher.ID,
		})
		teachers[k].Type = teacherProgram.Type
		teachers[k].ProgramID = teacherProgram.ProgramID
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        teachers,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetTeachersPaginateByProgram(c echo.Context) error {
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
	teachers := make([]models.Teacher, 0)

	// Query in database
	DB.Raw("SELECT * FROM teachers "+
		"WHERE id IN (SELECT teacher_id FROM teacher_programs where program_id = ?) "+
		"AND (lower(first_name) LIKE lower(?) OR lower(last_name) LIKE lower(?) OR dni LIKE ?) ORDER BY id desc "+
		"OFFSET ? LIMIT ?",
		request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%", offset, request.Limit).Scan(&teachers)

	// Query students count total
	DB.Raw("SELECT * FROM teachers "+
		"WHERE id IN (SELECT teacher_id FROM teacher_programs where program_id = ?) "+
		"AND (lower(first_name) LIKE lower(?) OR lower(last_name) LIKE lower(?) OR dni LIKE ?) ",
		request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%").Scan(&total)

	// Get type teacher
	for k, teacher := range teachers {
		teacherProgram := models.TeacherProgram{}
		DB.First(&teacherProgram, models.TeacherProgram{
			TeacherID: teacher.ID,
		})
		teachers[k].Type = teacherProgram.Type
		teachers[k].ProgramID = teacherProgram.ProgramID
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        teachers,
		Total:       total.Count,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

type teacherSearchResponse struct {
	ID        uint   `json:"id"`
	DNI       string `json:"dni"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

func GetTeacherSearch(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	teachers := make([]teacherSearchResponse, 0)

	if request.Search != "" {
		if err := DB.Table("teachers").Select("id, dni, first_name, last_name").Where("lower(last_name) LIKE lower(?)", "%"+request.Search+"%").
			Or("lower(first_name) LIKE lower(?)", "%"+request.Search+"%").
			Or("lower(dni) LIKE lower(?)", "%"+request.Search+"%").
			Limit(5).Find(&teachers).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    teachers,
	})
}

func GetTeacherSearchProgram(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	teachers := make([]teacherSearchResponse, 0)

	// Search teachers
	if request.Search != "" {
		DB.Raw("SELECT * FROM teachers "+
			"WHERE id IN (SELECT teacher_id FROM teacher_programs where program_id = ?) "+
			"AND (lower(first_name) LIKE lower(?) OR lower(last_name) LIKE lower(?) OR dni LIKE ?) ORDER BY id desc "+
			"LIMIT 5",
			request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%").Scan(&teachers)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    teachers,
	})
}

func CreateTeacher(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

	// Get data request
	teacher := models.Teacher{}
	if err := c.Bind(&teacher); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	TR := DB.Begin()

	// has password new user account
	cc := sha256.Sum256([]byte(teacher.DNI + "TA"))
	pwd := fmt.Sprintf("%x", cc)

	// Insert user in database
	userAccount := models.User{
		UserName: teacher.DNI + "TA",
		Password: pwd,
		RoleID:   4,
	}
	if err := TR.Create(&userAccount).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert teachers in database
	teacherPrograms := make([]models.TeacherProgram, 0)
	teacherPrograms = append(teacherPrograms, models.TeacherProgram{
		ProgramID: teacher.ProgramID,
		Type:      teacher.Type,
		ByDefault: true,
	})

	teacher.UserID = userAccount.ID
	teacher.TeacherPrograms = teacherPrograms
	if err := TR.Create(&teacher).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TR.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    teacher.ID,
		Message: fmt.Sprintf("El profesor %s se registro correctamente", teacher.FirstName),
	})
}

func UpdateTeacher(c echo.Context) error {
	// Get data request
	teacher := models.Teacher{}
	if err := c.Bind(&teacher); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	TR := DB.Begin()

	// Update teacher in database
	if err := TR.Model(&teacher).Update(teacher).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", teacher.ID),
		})
	}

	// Update teacher program
	teacherProgram := models.TeacherProgram{
		ProgramID: teacher.ProgramID,
		Type:      teacher.Type,
	}
	if err := TR.Debug().Model(&models.TeacherProgram{}).Where("teacher_id = ? AND by_default = true", teacher.ID).
		Update(teacherProgram).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", teacher.ID),
		})
	}

	// Commit transaction
	TR.Commit()

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    teacher.ID,
		Message: fmt.Sprintf("Los datos del estudiante %s se actualizaron correctamente", teacher.FirstName),
	})
}

func DeleteTeacher(c echo.Context) error {
	// Get data request
	teacher := models.Teacher{}
	if err := c.Bind(&teacher); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&teacher).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    teacher.ID,
		Message: fmt.Sprintf("El estudiante %s se elimino correctamente", teacher.FirstName),
	})
}

func GetTempUploadTeacherBySubsidiary(c echo.Context) error {
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
	fileDir := "templates/templateTeacherSubsidiary.xlsx"
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

func GetTempUploadTeacherByProgram(c echo.Context) error {
	// Return file excel
	return c.File("templates/templateTeacher.xlsx")
}

// Set SetTempUploadTeacherBySubsidiary upload teacher
func SetTempUploadTeacherBySubsidiary(c echo.Context) error {
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

	// Get all the rows in the Sheet1.
	rows := excel.GetRows("Teacher")
	for k, row := range rows {

		if k >= ignoreCols {

			// Validate required fields
			if row[0] == "" || row[1] == "" {
				break
			}

			// program id
			u, _ := strconv.ParseUint(strings.TrimSpace(row[0]), 0, 32)
			currentProgram := uint(u)

			// Create model teacherPrograms
			teacherPrograms := make([]models.TeacherProgram, 0)
			teacherPrograms = append(teacherPrograms, models.TeacherProgram{
				ProgramID: currentProgram,
				Type:      "career",
			})

			// Create model teacher
			teacher := models.Teacher{
				DNI:       strings.TrimSpace(row[1]),
				LastName:  strings.TrimSpace(row[2]),
				FirstName: strings.TrimSpace(row[3]),
				Phone:     strings.TrimSpace(row[4]),
				Gender:    strings.TrimSpace(row[6]),

				Address:        strings.TrimSpace(row[8]),
				WorkConditions: strings.TrimSpace(row[9]),
				EducationLevel: strings.TrimSpace(row[10]),
				Specialty:      strings.TrimSpace(row[13]),
			}

			// has password new user account
			cc := sha256.Sum256([]byte(teacher.DNI + "TA"))
			pwd := fmt.Sprintf("%x", cc)

			// New Account
			userAccount := models.User{
				UserName: teacher.DNI + "TA",
				Email:    strings.TrimSpace(row[4]),
				Password: pwd,
				RoleID:   5,
			}

			// Insert user in database
			if err := TX.Create(&userAccount).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
			teacher.UserID = userAccount.ID // Set new user id

			// Create teacher
			if err := TX.Create(&teacher).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{
					Success: false,
					Message: fmt.Sprintf("Ocurrió un error al insertar el profesor %s con "+
						"DNI: %s es posible que este profesor ya este en la base de datos o los datos son incorrectos, "+
						"Error: %s, no se realizo ninguna cambio en la base de datos", teacher.FirstName, teacher.DNI, err),
				})
			}

			// Relation student
			teacherProgram := models.TeacherProgram{
				ProgramID: currentProgram,
				TeacherID: teacher.ID,
			}
			if err := TX.Create(&teacherProgram).Error; err != nil {
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
		Message: fmt.Sprintf("Se guardo %d registros den la base de datos", counter),
	})
}

func SetTempUploadTeacherByProgram(c echo.Context) error {
	// Response success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		//Message: fmt.Sprintf("Se guardo %d registros den la base de datos",'s'),
	})
}

func ExportAllTeachers(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query in database
	teachers := make([]models.Teacher, 0)
	if err := db.Order("id asc").Find(&teachers).Error; err != nil {
		return err
	}

	// Create excel file
	excel := excelize.NewFile()

	// Create a new sheet.
	index := excel.NewSheet("Sheet1")

	// Set value of a cell.
	excel.SetCellValue("Sheet1", "A1", "DNI")
	excel.SetCellValue("Sheet1", "B1", "Apellidos")
	excel.SetCellValue("Sheet1", "C1", "Nombres")
	excel.SetCellValue("Sheet1", "D1", "Fecha Nacimiento")
	excel.SetCellValue("Sheet1", "E1", "Genero")
	excel.SetCellValue("Sheet1", "F1", "Direccion")
	excel.SetCellValue("Sheet1", "G1", "Telefono")
	excel.SetCellValue("Sheet1", "H1", "Condicion Laboral")
	excel.SetCellValue("Sheet1", "I1", "Nivel de educacion")
	excel.SetCellValue("Sheet1", "J1", "Fecha ingreso")
	excel.SetCellValue("Sheet1", "K1", "Fecha retiro")
	excel.SetCellValue("Sheet1", "L1", "Especialidad")

	currentRow := 2
	for k, teacher := range teachers {
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow+k), teacher.DNI)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow+k), teacher.LastName)
		excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow+k), teacher.FirstName)
		excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", currentRow+k), teacher.BirthDate.Format("01/02/2006"))
		excel.SetCellValue("Sheet1", fmt.Sprintf("E%d", currentRow+k), teacher.Gender)
		excel.SetCellValue("Sheet1", fmt.Sprintf("F%d", currentRow+k), teacher.Address)
		excel.SetCellValue("Sheet1", fmt.Sprintf("G%d", currentRow+k), teacher.Phone)
		excel.SetCellValue("Sheet1", fmt.Sprintf("H%d", currentRow+k), teacher.WorkConditions)
		excel.SetCellValue("Sheet1", fmt.Sprintf("I%d", currentRow+k), teacher.EducationLevel)
		excel.SetCellValue("Sheet1", fmt.Sprintf("J%d", currentRow+k), teacher.AdmissionDate.Format("01/02/2006"))
		excel.SetCellValue("Sheet1", fmt.Sprintf("K%d", currentRow+k), teacher.RetirementDate.Format("01/02/2006"))
		excel.SetCellValue("Sheet1", fmt.Sprintf("L%d", currentRow+k), teacher.Specialty)
	}

	// Set active sheet of the workbook.
	excel.SetActiveSheet(index)

	// Save excel file by the given path.
	err := excel.SaveAs("temp/allTeachers.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	// Response file excel
	return c.File("temp/allTeachers.xlsx")
}
