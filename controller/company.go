package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetCompanies(c echo.Context) error {
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
	companies := make([]models.Company, 0)

	// Query in database
	if err := db.Where("nombre_o_razon_social LIKE ?", "%"+request.Search+"%").
		Order("id asc").
		Offset(offset).Limit(request.Limit).Find(&companies).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        companies,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetCompanySearch(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	companies := make([]models.Company, 0)
	if err := db.Where("lower(nombre_o_razon_social) LIKE lower(?)", "%"+request.Search+"%").
		Limit(10).Find(&companies).Error; err != nil {
		return err
	}

	customCompanies := make([]models.Company, 0)
	for _, student := range companies {
		customCompanies = append(customCompanies, models.Company{
			ID:                 student.ID,
			NombreORazonSocial: student.NombreORazonSocial,
			RUC:                student.RUC,
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customCompanies,
	})
}

func CreateCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert companies in database
	if err := db.Create(&company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("La empresa %s se registro correctamente", company.NombreORazonSocial),
	})
}

func UpdateCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update company in database
	rows := db.Model(&company).Update(company).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", company.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("Los datos del la empresa %s se actualizaron correctamente", company.NombreORazonSocial),
	})
}

func DeleteCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation company exist
	if db.First(&company).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", company.ID),
		})
	}

	// Delete company in database
	if err := db.Delete(&company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("La empresa %s se elimino correctamente", company.NombreORazonSocial),
	})
}
