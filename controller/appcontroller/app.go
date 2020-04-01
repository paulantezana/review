package appcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetAppByID(c echo.Context) error {
    DB := provider.GetConnection()
    defer DB.Close()

    app := models.App{}
    if err := DB.First(&app).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    app,
    })
}

func UpdateApp(c echo.Context) error {
    // Get data request
    app := models.App{}
    if err := c.Bind(&app); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Update app in database
    rows := DB.Model(&app).Update(app).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", app.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    app.ID,
        Message: fmt.Sprintf("Los datos del app %s se actualizaron correctamente", app.Name),
    })
}

