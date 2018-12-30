package institutecontroller

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models/institutemodel"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

type programUserResponse struct {
    ID           uint `json:"id"`
    UserID       uint `json:"user_id"`
    ProgramID uint `json:"program_id"`
    License      bool `json:"license"`
    Name         string `json:"name"`
}

type programUserRequest struct {
    UserID uint `json:"user_id"`
    SubsidiaryID uint `json:"subsidiary_id"`
}

// Get all programs licenses by user
func GetProgramsUserByUserID(c echo.Context) error {
    // Get data request
    request := programUserRequest{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query Programs
    programs := make([]institutemodel.Subsidiary, 0)
    if err := DB.Raw("SELECT * FROM programs WHERE id NOT IN (SELECT program_id  FROM program_users WHERE user_id = ?) AND subsidiary_id = ?", request.UserID, request.SubsidiaryID).
        Scan(&programs).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Start Transaction
    TR := DB.Begin()

    // Insert SubsidiaryUsers
    for _, program := range programs {
        programUser := institutemodel.ProgramUser{
            UserID:       request.UserID,
            ProgramID: program.ID,
        }
        if err := TR.Create(&programUser).Error; err != nil {
            TR.Rollback()
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
    }

    // End Transaction
    TR.Commit()

    // Query SubsidiaryUsers
    programUsers := make([]programUserResponse, 0)
    if err := DB.Table("program_users").
        Select("program_users.id, program_users.user_id, program_users.program_id, program_users.license, programs.name").
        Joins("INNER JOIN programs ON programs.id = program_users.program_id").
        Where("program_users.user_id = ? AND programs.subsidiary_id = ?", request.UserID,request.SubsidiaryID).
        Scan(&programUsers).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Response data
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    programUsers,
    })
}

func GetProgramsUserByUserIDLicense(c echo.Context) error {
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query
    programUsers := make([]programUserResponse, 0)
    if err := DB.Table("program_users").
        Select("program_users.id, program_users.user_id, program_users.program_id, program_users.license, programs.name").
        Joins("INNER JOIN programs ON programs.id = program_users.program_id").
        Where("program_users.user_id = ? AND license = TRUE", currentUser.ID).
        Scan(&programUsers).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Response data
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    programUsers,
    })
}

func UpdateProgramsUserByUserID(c echo.Context) error {
    // Get data request
    programUsers := make([]institutemodel.ProgramUser, 0)
    if err := c.Bind(&programUsers); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Update in Database
    for _, programUser := range programUsers {
        if err := DB.Model(programUser).UpdateColumn("license", programUser.License).Error; err != nil {
            return err
        }
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: "OK",
    })
}

