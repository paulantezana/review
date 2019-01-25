package institutecontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
	"net/http"
)

type teacherProgramResponse struct {
	ID        uint   `json:"id"`
	TeacherID uint   `json:"teacher_id"`
	ProgramID uint   `json:"program_id"`
	ByDefault bool   `json:"by_default"`
	Type      string `json:"type"` // cross // career

	DNI       string `json:"dni"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

func GetTeacherProgramByProgram(c echo.Context) error {
	// Get data request
	teacherProgram := models.TeacherProgram{}
	if err := c.Bind(&teacherProgram); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	teachers := make([]teacherProgramResponse, 0)
	if err := DB.Table("teachers").
		Select("teacher_programs.id, teacher_programs.teacher_id, teacher_programs.program_id, teacher_programs.by_default, teacher_programs.type, teachers.dni, teachers.first_name, teachers.last_name").
		Joins("INNER JOIN teacher_programs ON teacher_programs.teacher_id = teachers.id").
		Order("teacher_programs.id desc").
		Where("teacher_programs.program_id = ?", teacherProgram.ProgramID).
		Scan(&teachers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    teachers,
	})
}

func CreateTeacherProgram(c echo.Context) error {
	// Get data request
	teacherProgram := models.TeacherProgram{}
	if err := c.Bind(&teacherProgram); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Validation
	tValidate := make([]models.TeacherProgram, 0)
	if err := DB.Find(&tValidate, models.TeacherProgram{
		TeacherID: teacherProgram.TeacherID,
		ProgramID: teacherProgram.ProgramID,
	}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	if len(tValidate) > 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El profesor con el id = %d ya esta asignado en este programa de estudios", teacherProgram.TeacherID),
		})
	}

	// Create new teacherProgram
	if err := DB.Create(&teacherProgram).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    teacherProgram.ID,
		Message: fmt.Sprintf("El teacherPrograma de estudios %d se registro exitosamente", teacherProgram.ID),
	})
}

func UpdateTeacherProgram(c echo.Context) error {
	// Get data request
	teacherProgram := models.TeacherProgram{}
	if err := c.Bind(&teacherProgram); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update teacherProgram in database
	rows := db.Model(&teacherProgram).Update(teacherProgram).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", teacherProgram.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    teacherProgram.ID,
		Message: fmt.Sprintf("Los datos del teacherPrograma de estudios %s se actualizaron correctamente", teacherProgram.ID),
	})
}

func DeleteTeacherProgram(c echo.Context) error {
	// Get data request
	teacherProgram := models.TeacherProgram{}
	if err := c.Bind(&teacherProgram); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&teacherProgram).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    teacherProgram.ID,
		Message: fmt.Sprintf("The teacherProgram %s was successfully deleted", teacherProgram.ID),
	})
}
