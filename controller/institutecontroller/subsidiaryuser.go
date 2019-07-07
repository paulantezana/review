package institutecontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
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

func SubsidiariesUserUpdate(c echo.Context) error {
    // Get data request
    subsidiaryUser := models.SubsidiaryUser{}
    if err := c.Bind(&subsidiaryUser); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Update module in database
    DB.Model(&subsidiaryUser).Where("id = ?", subsidiaryUser.ID).UpdateColumn("license", subsidiaryUser.License)

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    subsidiaryUser.ID,
        Message: fmt.Sprintf("Los datos del se actualizaron correctamente"),
    })
}

// get all subsidiaries by user id
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
	subsidiaries := make([]models.Subsidiary, 0)
	if err := DB.Raw("SELECT * FROM subsidiaries WHERE id NOT IN (SELECT subsidiary_id FROM subsidiary_users WHERE user_id = ?)", user.ID).
		Scan(&subsidiaries).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Start Transaction
	TR := DB.Begin()

	// Insert SubsidiaryUsers
	for _, subsidiary := range subsidiaries {
		subsidiaryUser := models.SubsidiaryUser{
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiaryUsers,
	})
}

// Get subsidiaries license by user id
func GetSubsidiariesUserByUserIDLicense(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// ss
	subsidiaryUsers := make([]subsidiaryUserResponse, 0)
	if err := DB.Table("subsidiary_users").
		Select("subsidiary_users.id, subsidiary_users.user_id, subsidiary_users.subsidiary_id, subsidiary_users.license, subsidiaries.name").
		Joins("INNER JOIN subsidiaries ON subsidiaries.id = subsidiary_users.subsidiary_id").
		Where("subsidiary_users.user_id = ? AND license = TRUE", currentUser.ID).
		Scan(&subsidiaryUsers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    subsidiaryUsers,
	})
}