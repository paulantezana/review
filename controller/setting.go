package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
)

// GlobalSettings struct
type GlobalSettings struct {
	Message string         `json:"message"`
	Success bool           `json:"success"`
	Setting models.Setting `json:"setting"`
	User    models.User    `json:"user"`
	Program models.Program `json:"program"`
}

// GetGlobalSettings function
func GetGlobalSettings(c echo.Context) error {
	// Get data request
	con := models.Setting{}
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	user.Password = ""
	user.Key = ""

	// Find settings
	db.First(&con)

	// Find program
	var program models.Program
	if user.ProgramID > 0 {
		if err := db.First(&program, user.ProgramID).Error; err != nil {
			return err
		}
	}

	// Set object response
	return c.JSON(http.StatusOK, GlobalSettings{
		User:    user,
		Setting: con,
		Program: program,
		Success: true,
		Message: "OK",
	})
}

type studentSettingsResponse struct {
	Setting models.Setting `json:"setting"`
	Program models.Program `json:"program"`
	User    models.User    `json:"user"`
	Student models.Student `json:"student"`
	Message string         `json:"message"`
	Success bool           `json:"success"`
}

// GetGlobalSettings function
func GetStudentSettings(c echo.Context) error {
	// Get data request
	setting := models.Setting{}
	program := models.Program{}
	student := models.Student{}

	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// find user
	if err := db.First(&user, user.ID).Error; err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	user.Password = ""
	user.Key = ""

	// find program
	if err := db.First(&program, user.ProgramID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("%s", err),
        })
	}

	// find student
	if err := db.Where("user_id = ?",user.ID).First(&student).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("%s", err),
        })
	}

	// Find settings
	db.First(&setting)

	// Set object response
	return c.JSON(http.StatusOK, studentSettingsResponse{
		User:    user,
		Setting: setting,
		Program: program,
		Student: student,
		Success: true,
		Message: "OK",
	})
}

// UpdateSetting function update settings
func UpdateSetting(c echo.Context) error {
	// Get data request
	con := models.Setting{}
	if err := c.Bind(&con); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation first data
	var exist uint
	db.Model(&models.Setting{}).Count(&exist)

	// Insert config in database
	if exist == 0 {
		if err := db.Create(&con).Error; err != nil {
			return err
		}
	}

	// Update con in database
	if err := db.Model(&con).Update(con).Error; err != nil {
		return err
	}

	// Response config
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    con.ID,
	})
}

// UploadLogoSetting function upload logo settings
func UploadLogoSetting(c echo.Context) error {
	// Read form fields
	idSetting := c.FormValue("id")
	setting := models.Setting{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&setting, "id = ?", idSetting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}

	// Source
	file, err := c.FormFile("logo")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	logoSRC := fmt.Sprintf("static/logo%s", path.Ext(file.Filename))
	dst, err := os.Create(logoSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	setting.Logo = logoSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&setting).Update(setting).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    setting.ID,
		Message: "OK",
	})
}

// DownloadLogoSetting function dowloand logo settings
func DownloadLogoSetting(c echo.Context) error {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}
	return c.File(setting.Logo)
}

// UploadMinistrySetting function upload logo settings
func UploadMinistrySetting(c echo.Context) error {
	// Read form fields
	idSetting := c.FormValue("id")
	setting := models.Setting{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&setting, "id = ?", idSetting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}

	// Source
	file, err := c.FormFile("ministry")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	ministrySRC := fmt.Sprintf("static/ministry%s", path.Ext(file.Filename))
	dst, err := os.Create(ministrySRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	setting.Ministry = ministrySRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&setting).Update(setting).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    setting.ID,
		Message: "OK",
	})
}

// DownloadMinistrySetting function dowloand logo settings
func DownloadMinistrySetting(c echo.Context) error {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}
	return c.File(setting.Ministry)
}

// DownloadNationalEmblemSetting function download logo settings
func DownloadNationalEmblemSetting(c echo.Context) error {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}
	return c.File(setting.NationalEmblem)
}
