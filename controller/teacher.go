package controller

import (
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
    "io"
    "net/http"
    "os"
    "strings"
)

func GetTeachers(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // Get connection
    db := config.GetConnection()
    defer db.Close()

    // Pagination calculate
    if request.CurrentPage == 0 {
        request.CurrentPage = 1
    }
    offset := request.Limit*request.CurrentPage - request.Limit

    // Execute instructions
    var total uint
    teachers := make([]models.Teacher, 0)

    // Query in database
    if err := db.Where("lower(first_name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%",currentUser.ProgramID).
        Or("lower(last_name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", currentUser.ProgramID).
        Or("dni LIKE ? AND program_id = ?", "%"+request.Search+"%", currentUser.ProgramID).
        Order("id asc").
        Offset(offset).Limit(request.Limit).Find(&teachers).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
            return err
    }

    // Type response
    // 0 = all data
    // 1 = minimal data
    if request.Type == 1 {
        customTeacher := make([]models.Teacher, 0)
        for _, teacher := range teachers {
            customTeacher = append(customTeacher, models.Teacher{
                ID:       teacher.ID,
                FirstName: teacher.FirstName,
            })
        }
        return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
            Success:     true,
            Data:        customTeacher,
            Total:       total,
            CurrentPage: request.CurrentPage,
        })
    }
    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        teachers,
        Total:       total,
        CurrentPage: request.CurrentPage,
    })
}


func GetTeacherSearch(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // Get connection
    db := config.GetConnection()
    defer db.Close()

    // Execute instructions
    teachers := make([]models.Teacher, 0)
    if err := db.Where("lower(full_name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", currentUser.ProgramID).
        Or("dni LIKE ? AND program_id = ?", "%"+request.Search+"%", currentUser.ProgramID).
        Limit(10).Find(&teachers).Error; err != nil {
        return err
    }

    customTeachers := make([]models.Teacher, 0)
    for _, teacher := range teachers {
        customTeachers = append(customTeachers, models.Teacher{
            ID:       teacher.ID,
            FirstName: teacher.FirstName,
            DNI:      teacher.DNI,
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    customTeachers,
    })
}

func CreateTeacher(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    teacher := models.Teacher{}
    if err := c.Bind(&teacher); err != nil {
        return err
    }
    teacher.ProgramID = currentUser.ProgramID

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert teachers in database
    if err := db.Create(&teacher).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    teacher.ID,
        Message: fmt.Sprintf("El estudiante %s se registro correctamente", teacher.FirstName),
    })
}

func UpdateTeacher(c echo.Context) error {
    // Get data request
    teacher := models.Teacher{}
    if err := c.Bind(&teacher); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Update teacher in database
    rows := db.Model(&teacher).Update(teacher).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", teacher.ID),
        })
    }

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

    // Validation teacher exist
    if db.First(&teacher).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se encontró el registro con id %d", teacher.ID),
        })
    }

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
    return c.File("templates/uploadTeacherTemplate.xlsx")
}

func SetTempUploadTeacher(c echo.Context) error {
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
    teachers := make([]models.Teacher, 0)
    ignoreCols := 1

    // Get all the rows in the Sheet1.
    rows := xlsx.GetRows("Sheet1")
    for k, row := range rows {
        if k >= ignoreCols {
            teachers = append(teachers, models.Teacher{
                DNI:      strings.TrimSpace(row[0]),
                FirstName: strings.TrimSpace(row[1]),
                Phone:    strings.TrimSpace(row[3]),
                //State:    true,
            })
        }
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert teachers in database
    tr := db.Begin()
    for _, teacher := range teachers {
        if err := tr.Create(&teacher).Error; err != nil {
            tr.Rollback()
            return c.JSON(http.StatusOK, utilities.Response{
                Success: false,
                Message: fmt.Sprintf("Ocurrió un error al insertar el alumno %s con "+
                    "DNI: %s es posible que este alumno ya este en la base de datos o los datos son incorrectos, "+
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

