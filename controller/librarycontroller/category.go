package librarycontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models/librarymodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetCategoriesPaginate(c echo.Context) error {
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
	categories := make([]librarymodel.Category, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&categories).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        categories,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetCategoriesAll(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	categories := make([]librarymodel.Category, 0)

	// Query in database
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    categories,
	})
}

func GetCategoryByID(c echo.Context) error {
	// Get data request
	category := librarymodel.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&category, category.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    category,
	})
}

func CreateCategory(c echo.Context) error {
	// Get data request
	category := librarymodel.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert categories in database
	if err := db.Create(&category).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    category.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", category.Name),
	})
}

func UpdateCategory(c echo.Context) error {
	// Get data request
	category := librarymodel.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update category in database
	rows := db.Model(&category).Update(category).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", category.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    category.ID,
		Message: fmt.Sprintf("Los datos del curso %s se actualizaron correctamente", category.Name),
	})
}

func DeleteCategory(c echo.Context) error {
	// Get data request
	category := librarymodel.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete category in database
	if err := db.Delete(&category).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    category.ID,
		Message: fmt.Sprintf("El curso %s se elimino correctamente", category.Name),
	})
}
