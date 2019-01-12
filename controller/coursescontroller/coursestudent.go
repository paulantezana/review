package coursescontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models/coursemodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetCourseStudentsPaginate(c echo.Context) error {
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
	courseStudents := make([]coursemodel.CourseStudent, 0)

	// Query in database
	if err := db.Where("course_id = ?", request.ID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&courseStudents).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return err
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
	courseStudent := coursemodel.CourseStudent{}
	if err := c.Bind(&courseStudent); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

    // Validation
    VCStudent := coursemodel.CourseStudent{}
    DB.First(&VCStudent,coursemodel.CourseStudent{DNI:courseStudent.DNI,CourseID:courseStudent.CourseID})
    if VCStudent.ID != 0 {
        return c.JSON(http.StatusOK,utilities.Response{
            Message: fmt.Sprintf("El estudiante %s ya está matriculado en este curso",courseStudent.FullName),
        })
    }

    // Insert courseStudents in database
	if err := DB.Create(&courseStudent).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
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
	courseStudent := coursemodel.CourseStudent{}
	if err := c.Bind(&courseStudent); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update course in database
	rows := db.Model(&courseStudent).Update(courseStudent).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
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
	courseStudent := coursemodel.CourseStudent{}
	if err := c.Bind(&courseStudent); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete course in database
	if err := db.Delete(&courseStudent).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    courseStudent.ID,
		Message: fmt.Sprintf("La participante %s se elimino correctamente", courseStudent.FullName),
	})
}
