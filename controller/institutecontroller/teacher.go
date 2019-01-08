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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dgrijalva/jwt-go"
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
	teachers := make([]institutemodel.Teacher, 0)

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
		teacherProgram := institutemodel.TeacherProgram{}
		db.First(&teacherProgram, institutemodel.TeacherProgram{
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
    teachers := make([]institutemodel.Teacher, 0)

    // Query in database
    DB.Raw("SELECT * FROM teachers " +
        "WHERE id IN (SELECT teacher_id FROM teacher_programs where program_id = ?) " +
        "AND (lower(first_name) LIKE lower(?) OR lower(last_name) LIKE lower(?) OR dni LIKE ?) ORDER BY id desc " +
        "OFFSET ? LIMIT ?",
        request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%",offset,request.Limit).Scan(&teachers)

    // Query students count total
    DB.Raw("SELECT * FROM teachers " +
        "WHERE id IN (SELECT teacher_id FROM teacher_programs where program_id = ?) " +
        "AND (lower(first_name) LIKE lower(?) OR lower(last_name) LIKE lower(?) OR dni LIKE ?) ",
        request.ProgramID, "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%").Scan(&total)

    // Get type teacher
    for k, teacher := range teachers {
        teacherProgram := institutemodel.TeacherProgram{}
        DB.First(&teacherProgram, institutemodel.TeacherProgram{
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
        DB.Raw("SELECT * FROM teachers " +
            "WHERE id IN (SELECT teacher_id FROM teacher_programs where program_id = ?) " +
            "AND (lower(first_name) LIKE lower(?) OR lower(last_name) LIKE lower(?) OR dni LIKE ?) ORDER BY id desc " +
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
	teacher := institutemodel.Teacher{}
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
	teacherPrograms := make([]institutemodel.TeacherProgram, 0)
	teacherPrograms = append(teacherPrograms, institutemodel.TeacherProgram{
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
	teacher := institutemodel.Teacher{}
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
	teacherProgram := institutemodel.TeacherProgram{
		ProgramID: teacher.ProgramID,
		Type:      teacher.Type,
	}
	if err := TR.Debug().Model(&institutemodel.TeacherProgram{}).Where("teacher_id = ? AND by_default = true", teacher.ID).
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
	teacher := institutemodel.Teacher{}
	if err := c.Bind(&teacher); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&teacher).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    teacher.ID,
		Message: fmt.Sprintf("El estudiante %s se elimino correctamente", teacher.FirstName),
	})
}

func GetTempUploadTeacher(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Return file sa
	if currentUser.RoleID == 2 || currentUser.RoleID == 1 {
		// Get data request
		subsidiary := institutemodel.Subsidiary{}
		if err := c.Bind(&subsidiary); err != nil {
			return err
		}

		// Get template excel
		fileDir := "templates/templateTeacherSA.xlsx"
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
		if err := db.Find(&programs, institutemodel.Program{SubsidiaryID: subsidiary.ID}).
			Order("id desc").Error; err != nil {
			return err
		}

		xlsx.SetCellValue("ProgramIDS", "A1", "ID")
		xlsx.SetCellValue("ProgramIDS", "B1", "Programa De Estudios")

		// Clear cells
		for i := 0; i < 100; i++ {
			xlsx.SetCellValue("ProgramIDS", fmt.Sprintf("A%d", i+2), "")
			xlsx.SetCellValue("ProgramIDS", fmt.Sprintf("B%d", i+2), "")
		}

		// Fill data in cells
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
	return c.File("templates/templateTeacher.xlsx")
}

func SetTempUploadTeacher(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

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
	teachers := make([]institutemodel.Teacher, 0)
	ignoreCols := 1

	// Get all the rows in the Sheet1.
	rows := xlsx.GetRows("teacher")
	for k, row := range rows {
		if k >= ignoreCols {

			// Validate required fields
			if row[0] == "" {
				break
			}

			// program id
			var currentProgram uint
			//currentProgram = currentUser.DefaultProgramID

			if currentProgram == 0 {
				u, _ := strconv.ParseUint(strings.TrimSpace(row[12]), 0, 32)
				currentProgram = uint(u)
			}

			// Create model teacherPrograms
			teacherPrograms := make([]institutemodel.TeacherProgram, 0)
			teacherPrograms = append(teacherPrograms, institutemodel.TeacherProgram{
				ProgramID: currentProgram,
				Type:      "career",
			})

			// Create AND Append model Teacher
			teachers = append(teachers, institutemodel.Teacher{
				DNI:             strings.TrimSpace(row[0]),
				LastName:        strings.TrimSpace(row[1]),
				FirstName:       strings.TrimSpace(row[2]),
				Gender:          strings.TrimSpace(row[4]),
				Address:         strings.TrimSpace(row[5]),
				Phone:           strings.TrimSpace(row[6]),
				WorkConditions:  strings.TrimSpace(row[7]),
				EducationLevel:  strings.TrimSpace(row[8]),
				Specialty:       strings.TrimSpace(row[11]),
				TeacherPrograms: teacherPrograms,
			})
		}
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert teachers in database
	tr := db.Begin()
	for _, teacher := range teachers {
		// has password new user account
		cc := sha256.Sum256([]byte(teacher.DNI + "TA"))
		pwd := fmt.Sprintf("%x", cc)

		// New Account
		userAccount := models.User{
			UserName: teacher.DNI + "TA",
			Password: pwd,
			RoleID:   4,
		}

		// Insert user in database
		if err := tr.Create(&userAccount).Error; err != nil {
			tr.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: fmt.Sprintf("%s", err),
			})
		}

		// Insert teacher in database
		teacher.UserID = userAccount.ID
		if err := tr.Create(&teacher).Error; err != nil {
			tr.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: fmt.Sprintf("Ocurrió un error al insertar el profesor %s con "+
					"DNI: %s es posible que este profesor ya este en la base de datos o los datos son incorrectos, "+
					"Error: %s, no se realizo ninguna cambio en la base de datos", teacher.FirstName, teacher.DNI, err),
			})
		}
	}
	tr.Commit()

	// Response success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se guardo %d registros den la base de datos", len(teachers)),
	})
}

func ExportAllTeachers(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query in database
	teachers := make([]institutemodel.Teacher, 0)
	if err := db.Order("id asc").Find(&teachers).Error; err != nil {
		return err
	}

	// Create excel file
	xlsx := excelize.NewFile()

	// Create a new sheet.
	index := xlsx.NewSheet("Sheet1")

	// Set value of a cell.
	xlsx.SetCellValue("Sheet1", "A1", "DNI")
	xlsx.SetCellValue("Sheet1", "B1", "Apellidos")
	xlsx.SetCellValue("Sheet1", "C1", "Nombres")
	xlsx.SetCellValue("Sheet1", "D1", "Fecha Nacimiento")
	xlsx.SetCellValue("Sheet1", "E1", "Genero")
	xlsx.SetCellValue("Sheet1", "F1", "Direccion")
	xlsx.SetCellValue("Sheet1", "G1", "Telefono")
	xlsx.SetCellValue("Sheet1", "H1", "Condicion Laboral")
	xlsx.SetCellValue("Sheet1", "I1", "Nivel de educacion")
	xlsx.SetCellValue("Sheet1", "J1", "Fecha ingreso")
	xlsx.SetCellValue("Sheet1", "K1", "Fecha retiro")
	xlsx.SetCellValue("Sheet1", "L1", "Especialidad")

	currentRow := 2
	for k, teacher := range teachers {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow+k), teacher.DNI)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow+k), teacher.LastName)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow+k), teacher.FirstName)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", currentRow+k), teacher.BirthDate.Format("01/02/2006"))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", currentRow+k), teacher.Gender)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", currentRow+k), teacher.Address)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", currentRow+k), teacher.Phone)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", currentRow+k), teacher.WorkConditions)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("I%d", currentRow+k), teacher.EducationLevel)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("J%d", currentRow+k), teacher.AdmissionDate.Format("01/02/2006"))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("K%d", currentRow+k), teacher.RetirementDate.Format("01/02/2006"))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("L%d", currentRow+k), teacher.Specialty)
	}

	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)

	// Save xlsx file by the given path.
	err := xlsx.SaveAs("temp/allTeachers.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	// Response file excel
	return c.File("temp/allTeachers.xlsx")
}
