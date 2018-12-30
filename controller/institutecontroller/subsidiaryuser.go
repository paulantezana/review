package institutecontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/institutemodel"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

type subsidiaryUserResponse struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	UserID       uint   `json:"user_id"`
	SubsidiaryID uint   `json:"subsidiary_id"`
	Name         string `json:"name"`
	License      bool   `json:"license"`
}

func GetSubsidiariesUserByUserID(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query Subsidiaries
	subsidiaries := make([]institutemodel.Subsidiary, 0)
	if err := DB.Raw("SELECT * FROM subsidiaries WHERE id NOT IN (SELECT subsidiary_id FROM subsidiary_users WHERE user_id = ?)", user.ID).
		Scan(&subsidiaries).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Start Transaction
	TR := DB.Begin()

	// Insert SubsidiaryUsers
	for _, subsidiary := range subsidiaries {
		subsidiaryUser := institutemodel.SubsidiaryUser{
			UserID:       user.ID,
			SubsidiaryID: subsidiary.ID,
		}
		if err := TR.Create(&subsidiaryUser).Error; err != nil {
			TR.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// End Transaction
	TR.Commit()

	// Query SubsidiaryUsers
	subsidiaryUsers := make([]subsidiaryUserResponse, 0)
	if err := DB.Table("subsidiary_users").
		Select("subsidiary_users.id, subsidiary_users.user_id, subsidiary_users.subsidiary_id, subsidiary_users.license, subsidiaries.name").
		Joins("INNER JOIN subsidiaries ON subsidiaries.id = subsidiary_users.subsidiary_id").
		Where("subsidiary_users.user_id = ?", user.ID).
		Scan(&subsidiaryUsers).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiaryUsers,
	})
}

func GetSubsidiariesUserByUserIDLicense(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// ss
	subsidiaryUsers := make([]subsidiaryUserResponse, 0)
	if err := DB.Table("subsidiary_users").
		Select("subsidiary_users.id, subsidiary_users.user_id, subsidiary_users.subsidiary_id, subsidiary_users.license, subsidiaries.name").
		Joins("INNER JOIN subsidiaries ON subsidiaries.id = subsidiary_users.subsidiary_id").
		Where("subsidiary_users.user_id = ? AND license = TRUE", user.ID).
		Scan(&subsidiaryUsers).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiaryUsers,
	})
}

func UpdateSubsidiariesUserByUserID(c echo.Context) error {
	// Get data request
	subsidiaryUsers := make([]institutemodel.SubsidiaryUser, 0)
	if err := c.Bind(&subsidiaryUsers); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Update in Database
	for _, subsidiaryUser := range subsidiaryUsers {
		if err := DB.Model(subsidiaryUser).UpdateColumn("license", subsidiaryUser.License).Error; err != nil {
			return err
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}
