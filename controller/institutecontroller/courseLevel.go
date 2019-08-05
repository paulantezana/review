package institutecontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetCourseLevelPaginate(c echo.Context) error {
    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // Get connection
    db := provider.GetConnection()
    defer db.Close()

    // Pagination calculate
    offset := request.Validate()

    // Execute instructions
    var total uint
    courseLevels := make([]models.CourseLevel, 0)

    // Query in database
    if err := db.Where("lower(name) LIKE lower(?) AND course_id = ?", "%"+request.Search+"%", request.CourseID).
        Order("id desc").
        Offset(offset).Limit(request.Limit).Find(&courseLevels).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        courseLevels,
        Total:       total,
        CurrentPage: request.CurrentPage,
        Limit:       request.Limit,
    })
}


func GetAllCourseLevel(c echo.Context) error {
    // Get data request
    courseLevel := models.CourseLevel{}
    if err := c.Bind(&courseLevel); err != nil {
        return err
    }

    // Get connection
    db := provider.GetConnection()
    defer db.Close()

    // Execute instructions
    courseLevels := make([]models.CourseLevel, 0)

    // Query in database
    if err := db.Find(&courseLevels, &courseLevel).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    courseLevels,
    })
}

func CreateCourseLevel(c echo.Context) error {
    // Get data request
    courseLevel := models.CourseLevel{}
    if err := c.Bind(&courseLevel); err != nil {
        return err
    }
    //courseLevel.ProgramID = currentUser.ProgramID

    // get connection
    db := provider.GetConnection()
    defer db.Close()

    // Insert courseLevels in database
    if err := db.Create(&courseLevel).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    courseLevel.ID,
        Message: fmt.Sprintf("El modulo %s se registro correctamente", courseLevel.Name),
    })
}

func UpdateCourseLevel(c echo.Context) error {
    // Get data request
    courseLevel := models.CourseLevel{}
    if err := c.Bind(&courseLevel); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Update courseLevel in database
    rows := DB.Model(&courseLevel).Update(courseLevel).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", courseLevel.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    courseLevel.ID,
        Message: fmt.Sprintf("Los datos del modulo %s se actualizaron correctamente", courseLevel.Name),
    })
}

func DeleteCourseLevel(c echo.Context) error {
    // Get data request
    courseLevel := models.CourseLevel{}
    if err := c.Bind(&courseLevel); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Delete courseLevel in database
    if err := DB.Delete(&courseLevel).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    courseLevel.ID,
        Message: fmt.Sprintf("El modulo %s se elimino correctamente", courseLevel.Name),
    })
}

