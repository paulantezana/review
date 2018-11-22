package monitoringcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models/monitoring"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetPollsPaginate(c echo.Context) error {
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
    companies := make([]monitoring.Poll, 0)

    // Query in database
    if err := db.Where("lower(name_social_reason) LIKE lower(?)", "%"+request.Search+"%").
        Order("id asc").
        Offset(offset).Limit(request.Limit).Find(&companies).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        companies,
        Total:       total,
        CurrentPage: request.CurrentPage,
    })
}

func CreatePoll(c echo.Context) error {
    // Get data request
    poll := monitoring.Poll{}
    if err := c.Bind(&poll); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert companies in database
    if err := db.Create(&poll).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    poll.ID,
        Message: fmt.Sprintf("La empresa %s se registro correctamente", poll.Name),
    })
}

func UpdatePoll(c echo.Context) error {
    // Get data request
    poll := monitoring.Poll{}
    if err := c.Bind(&poll); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Update poll in database
    rows := db.Model(&poll).Update(poll).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", poll.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    poll.ID,
        Message: fmt.Sprintf("Los datos del la empresa %s se actualizaron correctamente", poll.Name),
    })
}

func DeletePoll(c echo.Context) error {
    // Get data request
    poll := monitoring.Poll{}
    if err := c.Bind(&poll); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Validation poll exist
    if db.First(&poll).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se encontró el registro con id %d", poll.ID),
        })
    }

    // Delete poll in database
    if err := db.Delete(&poll).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    poll.ID,
        Message: fmt.Sprintf("La empresa %s se elimino correctamente", poll.Name),
    })
}


