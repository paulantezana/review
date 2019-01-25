package institutecontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
	"net/http"
)

func GetSubsidiaries(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	subsidiaries := make([]models.Subsidiary, 0)
	if err := db.Order("id desc").Find(&subsidiaries).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    subsidiaries,
	})
}

type SubsidiariesTree struct {
	ID       uint                     `json:"id" gorm:"primary_key"`
	Name     string                   `json:"name"`
	Programs []models.Program `json:"programs"`
}

func GetSubsidiariesTree(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute Query
	subsidiariesTree := make([]SubsidiariesTree, 0)
	if err := db.Table("subsidiaries").Select("id, name").
		Scan(&subsidiariesTree).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query programs
	for k, subsidiary := range subsidiariesTree {
		programs := make([]models.Program, 0)
		if err := db.Find(&programs, models.Program{SubsidiaryID: subsidiary.ID}).Order("id desc").Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		subsidiariesTree[k].Programs = programs
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    subsidiariesTree,
	})
}

func GetSubsidiaryByID(c echo.Context) error {
	// Get data request
	subsidiary := models.Subsidiary{}
	if err := c.Bind(&subsidiary); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&subsidiary, subsidiary.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    subsidiary,
	})
}

func CreateSubsidiary(c echo.Context) error {
	// Get data request
	subsidiary := models.Subsidiary{}
	if err := c.Bind(&subsidiary); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Create new subsidiary
	if err := DB.Create(&subsidiary).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiary.ID,
		Message: fmt.Sprintf("the subsidiary %s successfully registered", subsidiary.Name),
	})
}

func UpdateSubsidiary(c echo.Context) error {
	// Get data request
	subsidiary := models.Subsidiary{}
	if err := c.Bind(&subsidiary); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update subsidiary in database
	rows := db.Model(&subsidiary).Update(subsidiary).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("the subsidiary %s could not update", subsidiary.Name),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiary.ID,
		Message: fmt.Sprintf("The data of the subsidiary %s was updated correctly", subsidiary.Name),
	})
}

func DeleteSubsidiary(c echo.Context) error {
	// Get data request
	subsidiary := models.Subsidiary{}
	if err := c.Bind(&subsidiary); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete teacher in database
	if err := db.Delete(&subsidiary).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    subsidiary.ID,
		Message: fmt.Sprintf("The subsidiary %s was successfully deleted", subsidiary.Name),
	})
}

func UpdateMainSubsidiary(c echo.Context) error {
	// Get data request
	subsidiary := models.Subsidiary{}
	if err := c.Bind(&subsidiary); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	//  all subsidiaries main = false
	if err := DB.Exec("UPDATE subsidiaries SET main = false").Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	//  current subsidiary main = true
	if err := DB.Model(subsidiary).UpdateColumn("main", true).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiary.ID,
		Message: fmt.Sprintf("The data of the subsidiary %s was updated correctly", subsidiary.Name),
	})
}
