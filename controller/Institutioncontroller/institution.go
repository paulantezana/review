package Institutioncontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"io"
	"net/http"
	"os"
	"path"
)

// GlobalInstitutions struct
type gInstitutionsResponse struct {
	UserRoles      []models.UserRole     `json:"userRoles"`
	Institution    models.Institution    `json:"institution"`
	User       models.User       `json:"user"`
	Program    models.Program    `json:"program"`
	Subsidiary models.Subsidiary `json:"subsidiary"`
}

// GetGlobalInstitutions function
func GetInstitutionSetting(c echo.Context) error {
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

	// Find institutions
	con := models.Institution{}
	db.First(&con,models.Institution{ID: currentUser.InstitutionID})

	// Find institutions
	userRoles := make([]models.UserRole, 0)
	if err := db.Where("id >= ?", currentUser.UserRoleID).Find(&userRoles).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Fond subsidiary
	subsidiary := models.Subsidiary{}
	db.Where("id = ?", request.SubsidiaryID).First(&subsidiary)

	// Find program
	program := models.Program{}
	if currentUser.UserRoleID == 3 {
		db.Where("id = ?", request.ProgramID).First(&program)
	}

	// Set object response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
		Data: gInstitutionsResponse{
			User:       currentUser,
			Institution:    con,
			UserRoles:      userRoles,
			Program:    program,
			Subsidiary: subsidiary,
		},
	})
}

type studentInstitutionsResponse struct {
	Institution models.Institution `json:"institution"`
	User    models.User    `json:"user"`
	Student models.Student `json:"student"`
	Message string         `json:"message"`
	Success bool           `json:"success"`
}

// GetGlobalInstitutions function
func GetStudentInstitutions(c echo.Context) error {
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
	user.TempKey = ""

	// Find institutions
	institution := models.Institution{}
	if err := db.First(&institution).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// find student
	student := models.Student{}
	if err := db.First(&student, models.Student{UserID: user.ID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Set object response
	return c.JSON(http.StatusOK, studentInstitutionsResponse{
		User:    user,
		Institution: institution,
		Student: student,
		Success: true,
		Message: "OK",
	})
}

// UpdateInstitution function update institutions
func UpdateInstitution(c echo.Context) error {
	// Get data request
	con := models.Institution{}
	if err := c.Bind(&con); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation first data
	var exist uint
	db.Model(&models.Institution{}).Count(&exist)

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

// UploadLogoInstitution function upload logo institutions
func UploadLogoInstitution(c echo.Context) error {
	// Read form fields
	idInstitution := c.FormValue("id")
	institution := models.Institution{}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&institution, "id = ?", idInstitution).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", institution.ID),
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
	institution.Logo = logoSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&institution).Update(institution).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    institution.ID,
		Message: fmt.Sprintf("El logo se guardo satisafactoriamente"),
	})
}

// DownloadLogoInstitution function dowloand logo institutions
func DownloadLogoInstitution(c echo.Context) error {
	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	institution := models.Institution{}
	if db.First(&institution).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", institution.ID),
		})
	}
	return c.File(institution.Logo)
}

// UploadMinistryInstitution function upload logo institutions
func UploadMinistryInstitution(c echo.Context) error {
	// Read form fields
	idInstitution := c.FormValue("id")
	institution := models.Institution{}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&institution, "id = ?", idInstitution).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", institution.ID),
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
	institution.Ministry = ministrySRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&institution).Update(institution).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    institution.ID,
		Message: "OK",
	})
}

// DownloadMinistryInstitution function dowloand logo institutions
func DownloadMinistryInstitution(c echo.Context) error {
	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	institution := models.Institution{}
	if db.First(&institution).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", institution.ID),
		})
	}
	return c.File(institution.Ministry)
}

// DownloadMinistrySmallInstitution function dowloand logo institutions
func DownloadMinistrySmallInstitution(c echo.Context) error {
	return c.File("static/ministrySmall.jpg")
}

// DownloadNationalEmblemInstitution function download logo institutions
func DownloadNationalEmblemInstitution(c echo.Context) error {
	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	institution := models.Institution{}
	if db.First(&institution).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", institution.ID),
		})
	}
	return c.File(institution.NationalEmblem)
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
