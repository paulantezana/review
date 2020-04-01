package appcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetAppModules(c echo.Context) error {
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Pagination calculate
    offset := request.Validate()

    // Execute instructions
    var total uint
    appModules := make([]models.AppModule, 0)
    if err := DB.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
        Order("id asc").
        Offset(offset).Limit(request.Limit).Find(&appModules).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.ResponsePaginate{
        Success:     true,
        Data:        appModules,
        Total:       total,
        CurrentPage: request.CurrentPage,
        Limit:       request.Limit,
    })
}

func GetAppModuleByID(c echo.Context) error {
    // Get data request
    appModule := models.AppModule{}
    if err := c.Bind(&appModule); err != nil {
        return err
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Execute instructions
    if err := DB.First(&appModule, appModule.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    appModule,
    })
}

func CreateAppModule(c echo.Context) error {
    // Get data request
    appModule := models.AppModule{}
    if err := c.Bind(&appModule); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Create new appModule
    if err := DB.Create(&appModule).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    appModule.ID,
        Message: fmt.Sprintf("El appModulea de estudios %s se registro exitosamente", appModule.Name),
    })
}

func UpdateAppModule(c echo.Context) error {
    // Get data request
    appModule := models.AppModule{}
    if err := c.Bind(&appModule); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Update appModule in database
    rows := DB.Model(&appModule).Update(appModule).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", appModule.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    appModule.ID,
        Message: fmt.Sprintf("Los datos del appModulea de estudios %s se actualizaron correctamente", appModule.Name),
    })
}
