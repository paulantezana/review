package institutecontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetInstitutions(c echo.Context) error {
    // Get data request
    institution := models.Institution{}
    if err := c.Bind(&institution); err != nil {
        return err
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Execute instructions
    institutions := make([]models.Institution, 0)
    if err := DB.Find(&institutions).Order("id desc").
        Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    institutions,
    })
}

func GetInstitutionByID(c echo.Context) error {
    // Get data request
    institution := models.Institution{}
    if err := c.Bind(&institution); err != nil {
        return err
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Execute instructions
    if err := DB.First(&institution, institution.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    institution,
    })
}

func CreateInstitution(c echo.Context) error {
    // Get data request
    institution := models.Institution{}
    if err := c.Bind(&institution); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // ------------------------------------
    // Starting transaction
    // ------------------------------------
    TR := DB.Begin()

    // Create new institution
    if err := TR.Create(&institution).Error; err != nil {
        TR.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    TR.Commit()
    // ------------------------------------
    // End Transaction
    // ------------------------------------

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    institution.ID,
        Message: fmt.Sprintf("El institutiona de estudios %s se registro exitosamente", institution.Institute),
    })
}

func UpdateInstitution(c echo.Context) error {
    // Get data request
    institution := models.Institution{}
    if err := c.Bind(&institution); err != nil {
        return err
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Update institution in database
    rows := DB.Model(&institution).Update(institution).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", institution.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    institution.ID,
        Message: fmt.Sprintf("Los datos del institutiona de estudios %s se actualizaron correctamente", institution.Institute),
    })
}
