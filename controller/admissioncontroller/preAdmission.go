package admissioncontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

func GetPreAdmission(c echo.Context) error {
	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	admissionSettings := make([]models.AdmissionSetting, 0)
	DB.Where("pre_end_date >= ? AND pre_enabled = true", time.Now()).Find(&admissionSettings)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admissionSettings,
	})
}

func GetPreAdmissionById(c echo.Context) error {
    admissionSetting := models.AdmissionSetting{}
    if err := c.Bind(&admissionSetting); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query
    DB.Where("pre_end_date >= ? AND pre_start_date <= ? AND id = ? AND pre_enabled = true", time.Now(), time.Now(), admissionSetting.ID).First(&admissionSetting)
    if admissionSetting.ID == 0 {
        // Return response
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("Cerrado"),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
       Success: true,
       Data:    admissionSetting,
    })
}

type savePreAdmissionRequest struct {
	Student models.Student `json:"student"`
	User    models.User    `json:"user"`
    AdmissionSettingID uint `json:"admission_setting_id"`
}

func SavePreAdmission(c echo.Context) error {
	// Get data request
	request := savePreAdmissionRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

    // Validate DNI
    if !utilities.ValidateDni(request.Student.DNI) {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Número de dni no valido")})
    }

	// Validate required parameters
    if request.AdmissionSettingID == 0 {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("No se especifico el proceso de admisión este campo es requerido.")})
    }

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// start transaction
	TX := DB.Begin()

	// Query student
    DB.First(&request.Student, models.Student{DNI: request.Student.DNI})

	// Validate if exist student
	if request.Student.ID == 0 {
		// has password new user account
		cc := sha256.Sum256([]byte(request.Student.DNI + "ST"))
		pwd := fmt.Sprintf("%x", cc)

		// Insert user in database
		request.User.UserName = request.Student.DNI + "ST"
		request.User.Password = pwd
		request.User.RoleID = 5

		if err := TX.Create(&request.User).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Insert student in database
		request.Student.UserID = request.User.ID
		request.Student.StudentStatusID = 1
		if err := TX.Create(&request.Student).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Insert student history
		studentHistory := models.StudentHistory{
			StudentID:   request.Student.ID,
			UserID:      request.User.ID,
			Description: "Sus datos fueron creados desde el proceso de pre inscripcion de admision",
			Date:        time.Now(),
			Type:        1,
		}
		if err := TX.Create(&studentHistory).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		// Update data
		TX.Model(&request.Student).Update(request.Student)
		TX.Model(&request.User).Update(request.User)

		// Query data user
		DB.First(&request.User, models.User{ID: request.User.ID})

		// Insert student history
		studentHistory := models.StudentHistory{
			StudentID:   request.Student.ID,
			UserID:      request.User.ID,
			Description: "Sus datos fueron modificados desde el proceso de pre inscripcion de admision",
			Date:        time.Now(),
			Type:        1,
		}
		if err := TX.Create(&studentHistory).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Register pre admission
    preAdmission := models.PreAdmission{
        StudentID: request.Student.ID,
        AdmissionSettingID: request.AdmissionSettingID,
    }

    // Validate
    DB.First(&preAdmission,preAdmission)
    if preAdmission.ID >= 1 {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("El estudiante %s ya esta registrado en el proceso de preadmisión.", request.Student.FullName),
        })
    }

    // Create pre admission
    if err := TX.Create(&preAdmission).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    request,
		Message: fmt.Sprintf("El estudiante %s se registro correctamente", request.Student.FullName),
	})
}
