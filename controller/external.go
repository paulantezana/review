package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

type reniecRequest struct {
	DNI string `json:"dni"`
}

type reniecResponse struct {
	Student models.Student `json:"student"`
	User    models.User    `json:"user"`
	Exist   bool           `json:"exist"`
}

func GetStudentByDni(c echo.Context) error {
	request := reniecRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Validate DNI
	if !utilities.ValidateDni(request.DNI) {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Número de dni no valido")})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Search student
	student := models.Student{}
	DB.First(&student, models.Student{DNI: request.DNI})
	reniecResponse := reniecResponse{
		Student: student,
	}

	// Validation
	if student.ID == 0 {
		student, err := Dni(request.DNI)
		if err == nil {
			reniecResponse.Student = student
		}
	} else {
		// Find User
		user := models.User{}
		if err := DB.First(&user, student.UserID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		user.Password = ""
		user.TempKey = ""

		reniecResponse.Exist = true
		reniecResponse.User = user
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    reniecResponse,
	})
}
