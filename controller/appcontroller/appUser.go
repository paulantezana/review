package appcontroller

import (
    "crypto/sha256"
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

type loginCoreResponse struct {
    AppUser     interface{} `json:"app_user"`
    Token    interface{}   `json:"token"`
}

// Login login app
func Login(c echo.Context) error {
    // Get data request
    appUser := models.AppUser{}
    if err := c.Bind(&appUser); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Hash password
    cc := sha256.Sum256([]byte(appUser.Password))
    pwd := fmt.Sprintf("%x", cc)

    // login with profile
    if DB.Where("user_name = ? and password = ?", appUser.UserName, pwd).First(&appUser).RecordNotFound() {
        if DB.Where("email = ? and password = ?", appUser.UserName, pwd).First(&appUser).RecordNotFound() {
            return c.JSON(http.StatusOK, utilities.Response{
                Message: "El nombre de usuario o contraseña es incorecta app",
            })
        }
    }

    // Check state appUser
    if !appUser.State {
        return c.NoContent(http.StatusForbidden)
    }

    // Prepare response data
    appUser.Password = ""
    appUser.TempKey = ""


    // get token key
    token := utilities.GenerateCoreJWT(appUser)

    // Login success
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: fmt.Sprintf("Bienvenido al sistema %s", appUser.UserName),
        Data: loginCoreResponse{
            AppUser:  appUser,
            Token: token,
        },
    })
}