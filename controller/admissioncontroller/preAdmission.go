package admissioncontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

func GetPreAdmissionsPaginate(c echo.Context) error {
	// Get data request
	request := admissionsPaginateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Pagination calculate
	offset := request.validate()

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	var total uint
	admissionsPaginateResponses := make([]admissionsPaginateResponse, 0)
	if err := DB.Debug().Table("pre_admissions").
		Select("pre_admissions.*, students.dni , students.full_name, users.id as user_id, users.email, users.avatar").
		Joins("INNER JOIN students ON pre_admissions.student_id = students.id").
		Joins("INNER JOIN users on students.user_id = users.id").
		Where("students.dni LIKE ? AND pre_admissions.admission_setting_id = ? AND pre_admissions.program_id = ?", "%"+request.Search+"%", request.AdmissionSettingID, request.ProgramID).
		Or("lower(students.full_name) LIKE lower(?) AND pre_admissions.admission_setting_id = ? AND pre_admissions.program_id = ?", "%"+request.Search+"%", request.AdmissionSettingID, request.ProgramID).
		Order("pre_admissions.id desc").
		Offset(offset).Limit(request.Limit).Scan(&admissionsPaginateResponses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.ResponsePaginate{
		Success:     true,
		Data:        admissionsPaginateResponses,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetPreAdmission(c echo.Context) error {
	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	admissionSettings := make([]models.AdmissionSetting, 0)
	DB.Where("pre_end_date >= ? AND pre_enabled = true", time.Now()).Find(&admissionSettings)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    admissionSettings,
	})
}

func GetPreAdmissionPrograms(c echo.Context) error {
	// Get data request
	admissionSetting := models.AdmissionSetting{}
	if err := c.Bind(&admissionSetting); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query settings
	DB.First(&admissionSetting, models.AdmissionModality{ID: admissionSetting.ID})

	// Query programs
	programs := make([]models.Program, 0)
	if err := DB.Where("subsidiary_id = ?", admissionSetting.SubsidiaryID).Find(&programs).Order("id desc").
		Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    programs,
	})
}

func GetPreAdmissionModalities(c echo.Context) error {
    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Query programs
    admissionModalities := make([]models.AdmissionModality, 0)
    if err := DB.Select("id, name").Find(&admissionModalities).Order("id desc").
        Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    admissionModalities,
    })
}

func GetPreAdmissionById(c echo.Context) error {
	// Get data request
	admissionSetting := models.AdmissionSetting{}
	if err := c.Bind(&admissionSetting); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
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
	Student      models.Student      `json:"student"`
	User         models.User         `json:"user"`
	PreAdmission models.PreAdmission `json:"pre_admission"`
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
	if request.PreAdmission.AdmissionSettingID == 0 {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("No se especifico el proceso de admisión este campo es requerido.")})
	}

	// get connection
	DB := provider.GetConnection()
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
		StudentID:          request.Student.ID,
		AdmissionSettingID: request.PreAdmission.AdmissionSettingID,
		ProgramID:          request.PreAdmission.ProgramID,
	}

	// Validate
	DB.First(&preAdmission, preAdmission)
	if preAdmission.ID >= 1 {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Success: true,
			Data:    request,
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
