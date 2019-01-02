package admissioncontroller

import (
    "crypto/sha256"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/models/admissionmodel"
    "github.com/paulantezana/review/models/institutemodel"
    "github.com/paulantezana/review/utilities"
    "net/http"
    "time"
)

func GetAdmissionsPaginate(c echo.Context) error {
    // Get user token authenticate
    //user := c.Get("user").(*jwt.Token)
    //claims := user.Claims.(*utilities.Claim)
    //currentUser := claims.User

    // Get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Execute instructions
    admissions := make([]admissionmodel.Admission, 0)
    DB.Find(&admissions)

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    admissions,
    })
}

type createAdmissionRequest struct {
    Student institutemodel.Student `json:"student"`
    Admission admissionmodel.Admission `json:"admission"`
}
func CreateAdmission(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    request := createAdmissionRequest{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // start transaction
    TX := DB.Begin()

    // has password new user account
    cc := sha256.Sum256([]byte(request.Student.DNI + "ST"))
    pwd := fmt.Sprintf("%x", cc)

    // Insert user in database
    userAccount := models.User{
        UserName: request.Student.DNI + "ST",
        Password: pwd,
        RoleID:   5,
    }
    if err := TX.Create(&userAccount).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Insert student in database
    request.Student.UserID = userAccount.ID
    request.Student.StudentStatusID = 2
    if err := TX.Create(&request.Student).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Insert admission
    request.Admission.StudentID = request.Student.ID
    request.Admission.AdmissionDate = time.Now()
    request.Admission.Year = uint(time.Now().Year())
    request.Admission.UserID = currentUser.ID
    if err := TX.Create(&request.Admission).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Insert student history
    studentHistory := institutemodel.StudentHistory{
        StudentID: request.Student.ID,
        UserID: currentUser.ID,
        Description: fmt.Sprintf("Admision"),
        Date: time.Now(),
    }
    if err := TX.Create(&studentHistory).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Commit transaction
    TX.Commit()

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    request.Student.ID,
        Message: fmt.Sprintf("El estudiante %s se registro correctamente", request.Student.FullName),
    })
}
