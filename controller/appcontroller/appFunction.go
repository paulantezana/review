package appcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetAppModuleFunctions(c echo.Context) error {
    // Get data request
    appModuleFunction := models.AppModuleFunction{}
    if err := c.Bind(&appModuleFunction); err != nil {
        return err
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Execute instructions
    appModuleFunctions := make([]models.AppModuleFunction, 0)
    if err := DB.Find(&appModuleFunctions).Order("id desc").
        Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    appModuleFunctions,
    })
}

func GetAppModuleFunctionByID(c echo.Context) error {
    // Get data request
    appModuleFunction := models.AppModuleFunction{}
    if err := c.Bind(&appModuleFunction); err != nil {
        return err
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Execute instructions
    if err := DB.First(&appModuleFunction, appModuleFunction.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    appModuleFunction,
    })
}

func CreateAppModuleFunction(c echo.Context) error {
    // Get data request
    appModuleFunction := models.AppModuleFunction{}
    if err := c.Bind(&appModuleFunction); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()


    // Create new appModuleFunction
    if err := DB.Create(&appModuleFunction).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    appModuleFunction.ID,
        Message: fmt.Sprintf("El appModuleFunctiona de estudios %s se registro exitosamente", appModuleFunction.Name),
    })
}

func UpdateAppModuleFunction(c echo.Context) error {
    // Get data request
    appModuleFunction := models.AppModuleFunction{}
    if err := c.Bind(&appModuleFunction); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Update appModuleFunction in database
    rows := DB.Model(&appModuleFunction).Update(appModuleFunction).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", appModuleFunction.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    appModuleFunction.ID,
        Message: fmt.Sprintf("Los datos del appModuleFunctiona de estudios %s se actualizaron correctamente", appModuleFunction.Name),
    })
}
