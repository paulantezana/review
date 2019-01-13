package institutecontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/institutemodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetPrograms(c echo.Context) error {
	// Get data request
	program := institutemodel.Program{}
	if err := c.Bind(&program); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	programs := make([]institutemodel.Program, 0)
	if err := DB.Where("subsidiary_id = ?", program.SubsidiaryID).Find(&programs).Order("id desc").
		Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    programs,
	})
}

func GetProgramByID(c echo.Context) error {
	// Get data request
	program := institutemodel.Program{}
	if err := c.Bind(&program); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	if err := DB.First(&program, program.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    program,
	})
}

type createProgramRequest struct {
	Name         string `json:"name"`
	Level        string `json:"level"`
	SubsidiaryID uint   `json:"subsidiary_id"`

	DNI       string `json:"dni"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

func CreateProgram(c echo.Context) error {
	// Get data request
	request := createProgramRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// ------------------------------------
	// Starting transaction
	// ------------------------------------
	TR := DB.Begin()
	program := institutemodel.Program{
		Name:         request.Name,
		Level:        request.Level,
		SubsidiaryID: request.SubsidiaryID,
	}
	// Create new program
	if err := TR.Create(&program).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create user
	cc := sha256.Sum256([]byte(request.Password))
	pwd := fmt.Sprintf("%x", cc)

	user := models.User{
		UserName: request.UserName,
		Password: pwd,
		Email:    request.Email,
		RoleID:   3,
		Freeze:   true,
	}
	if err := TR.Create(&user).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create program user - relation
	programUser := institutemodel.ProgramUser{
		UserID:    user.ID,
		ProgramID: program.ID,
		License:   true,
	}
	if err := TR.Create(&programUser).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create teacher
	teacher := institutemodel.Teacher{
		DNI:       request.DNI,
		FirstName: request.FirstName,
		ProgramID: program.ID,
		UserID:    user.ID,
	}
	if err := TR.Create(&teacher).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create Relation
	teacherProgram := institutemodel.TeacherProgram{
		ProgramID: program.ID,
		TeacherID: teacher.ID,
		Type:      "career",
		ByDefault: true,
	}
	if err := TR.Create(&teacherProgram).Error; err != nil {
		TR.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	TR.Commit()
	// ------------------------------------
	// End Transaction
	// ------------------------------------

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    program.ID,
		Message: fmt.Sprintf("El programa de estudios %s se registro exitosamente", program.Name),
	})
}

func UpdateProgram(c echo.Context) error {
	// Get data request
	program := institutemodel.Program{}
	if err := c.Bind(&program); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Update program in database
	rows := DB.Model(&program).Update(program).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", program.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    program.ID,
		Message: fmt.Sprintf("Los datos del programa de estudios %s se actualizaron correctamente", program.Name),
	})
}
