package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"io"
	"net/http"
	"os"
	"path"
)

// GlobalSettings struct
type gSettingsResponse struct {
	Roles      []models.Role     `json:"roles"`
	Setting    models.Setting    `json:"setting"`
	User       models.User       `json:"user"`
	Program    models.Program    `json:"program"`
	Subsidiary models.Subsidiary `json:"subsidiary"`
}

// GetGlobalSettings function
func GetGlobalSettings(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Find settings
	con := models.Setting{}
	db.First(&con)

	// Find settings
	roles := make([]models.Role, 0)
	if err := db.Where("id >= ?", currentUser.RoleID).Find(&roles).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Fond subsidiary
	subsidiary := models.Subsidiary{}
	db.Where("id = ?", request.SubsidiaryID).First(&subsidiary)

	// Find program
	program := models.Program{}
	if currentUser.RoleID == 3 {
		db.Where("id = ?", request.ProgramID).First(&program)
	}

	// Set object response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
		Data: gSettingsResponse{
			User:       currentUser,
			Setting:    con,
			Roles:      roles,
			Program:    program,
			Subsidiary: subsidiary,
		},
	})
}

type studentSettingsResponse struct {
	Setting models.Setting `json:"setting"`
	User    models.User    `json:"user"`
	Student models.Student `json:"student"`
	Message string         `json:"message"`
	Success bool           `json:"success"`
}

// GetGlobalSettings function
func GetStudentSettings(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	user.Password = ""
	user.Key = ""

	// Find settings
	setting := models.Setting{}
	if err := db.First(&setting).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// find student
	student := models.Student{}
	if err := db.First(&student, models.Student{UserID: user.ID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Set object response
	return c.JSON(http.StatusOK, studentSettingsResponse{
		User:    user,
		Setting: setting,
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
	db := provider.GetConnection()
	defer db.Close()

	// Validation first data
	var exist uint
	db.Model(&models.Setting{}).Count(&exist)

	// Insert provider in database
	if exist == 0 {
		if err := db.Create(&con).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Update con in database
	if err := db.Model(&con).Update(con).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response provider
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    con.ID,
		Message: fmt.Sprintf("Los datos se guardarón satisafactoriamente"),
	})
}

// UploadLogoSetting function upload logo settings
func UploadLogoSetting(c echo.Context) error {
	// Read form fields
	idSetting := c.FormValue("id")
	setting := models.Setting{}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&setting, "id = ?", idSetting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
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
		Message: fmt.Sprintf("El logo se guardo satisafactoriamente"),
	})
}

// DownloadLogoSetting function dowloand logo settings
func DownloadLogoSetting(c echo.Context) error {
	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
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
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&setting, "id = ?", idSetting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}
	return c.File(setting.Ministry)
}

// DownloadMinistrySmallSetting function dowloand logo settings
func DownloadMinistrySmallSetting(c echo.Context) error {
	return c.File("static/ministrySmall.jpg")
}

// DownloadNationalEmblemSetting function download logo settings
func DownloadNationalEmblemSetting(c echo.Context) error {
	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}
	return c.File(setting.NationalEmblem)
}

func DownloadFile(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Return request
	return c.File(request.Search)
}
