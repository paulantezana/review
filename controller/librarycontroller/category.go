package librarycontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
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
	db := provider.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	categories := make([]models.PostCategory, 0)

	// Query in database
	if err := db.Where("lower(name) LIKE lower(?) AND program_id = ?", "%"+request.Search+"%", request.ProgramID).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&categories).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
	request := models.PostCategory{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	categories := make([]models.PostCategory, 0)

	// Query in database
	if err := db.Where("program_id = ?", request.ProgramId).Order("id desc").Find(&categories).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    categories,
	})
}

func GetPostCategoryByID(c echo.Context) error {
	// Get data request
	postCategory := models.PostCategory{}
	if err := c.Bind(&postCategory); err != nil {
		return err
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&postCategory, postCategory.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    postCategory,
	})
}

func CreatePostCategory(c echo.Context) error {
	// Get data request
	postCategory := models.PostCategory{}
	if err := c.Bind(&postCategory); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Insert categories in database
	if err := db.Create(&postCategory).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    postCategory.ID,
		Message: fmt.Sprintf("La categoria %s se registro correctamente", postCategory.Name),
	})
}

func UpdatePostCategory(c echo.Context) error {
	// Get data request
	postCategory := models.PostCategory{}
	if err := c.Bind(&postCategory); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update postCategory in database
	rows := db.Model(&postCategory).Update(postCategory).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", postCategory.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    postCategory.ID,
		Message: fmt.Sprintf("Los datos de la categoria %s se actualizaron correctamente", postCategory.Name),
	})
}

func DeletePostCategory(c echo.Context) error {
	// Get data request
	postCategory := models.PostCategory{}
	if err := c.Bind(&postCategory); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Delete postCategory in database
	if err := db.Delete(&postCategory).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    postCategory.ID,
		Message: fmt.Sprintf("La categoria con el id %d se elimino correctamente", postCategory.ID),
	})
}
