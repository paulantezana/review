package coursescontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models/coursemodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetCoursesPaginate(c echo.Context) error {
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
	courses := make([]coursemodel.Course, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&courses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        courses,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetCourseByID(c echo.Context) error {
	// Get data request
	course := coursemodel.Course{}
	if err := c.Bind(&course); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&course, course.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    course,
	})
}

func CreateCourse(c echo.Context) error {
	// Get data request
	course := coursemodel.Course{}
	if err := c.Bind(&course); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert courses in database
	if err := db.Create(&course).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    course.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", course.Name),
	})
}

func UpdateCourse(c echo.Context) error {
	// Get data request
	course := coursemodel.Course{}
	if err := c.Bind(&course); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update course in database
	rows := db.Model(&course).Update(course).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", course.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    course.ID,
		Message: fmt.Sprintf("Los datos del curso %s se actualizaron correctamente", course.Name),
	})
}

func DeleteCourse(c echo.Context) error {
	// Get data request
	course := coursemodel.Course{}
	if err := c.Bind(&course); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete course in database
	if err := db.Delete(&course).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    course.ID,
		Message: fmt.Sprintf("El curso %s se elimino correctamente", course.Name),
	})
}
