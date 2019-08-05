package institutecontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetCourseNodePaginate(c echo.Context) error {
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
    courseNodes := make([]models.CourseNode, 0)

    // Query in database
    if err := db.Where("lower(name) LIKE lower(?) AND course_level_id = ?", "%"+request.Search+"%", request.CourseLevelID).
        Order("id desc").
        Offset(offset).Limit(request.Limit).Find(&courseNodes).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        courseNodes,
        Total:       total,
        CurrentPage: request.CurrentPage,
        Limit:       request.Limit,
    })
}

func CreateCourseNode(c echo.Context) error {
    // Get data request
    courseNode := models.CourseNode{}
    if err := c.Bind(&courseNode); err != nil {
        return err
    }
    //courseNode.ProgramID = currentUser.ProgramID

    // get connection
    db := provider.GetConnection()
    defer db.Close()

    // Insert courseNodes in database
    if err := db.Create(&courseNode).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    courseNode.ID,
        Message: fmt.Sprintf("El modulo %s se registro correctamente", courseNode.Name),
    })
}

func UpdateCourseNode(c echo.Context) error {
    // Get data request
    courseNode := models.CourseNode{}
    if err := c.Bind(&courseNode); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Update courseNode in database
    rows := DB.Model(&courseNode).Update(courseNode).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", courseNode.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    courseNode.ID,
        Message: fmt.Sprintf("Los datos del modulo %s se actualizaron correctamente", courseNode.Name),
    })
}

func DeleteCourseNode(c echo.Context) error {
    // Get data request
    courseNode := models.CourseNode{}
    if err := c.Bind(&courseNode); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Delete courseNode in database
    if err := DB.Delete(&courseNode).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    courseNode.ID,
        Message: fmt.Sprintf("El modulo %s se elimino correctamente", courseNode.Name),
    })
}

