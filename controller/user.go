package controller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/messengermodel"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

type loginProgramLicense struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type loginSubsidiaryLicense struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type loginLicenses struct {
	Programs     []loginProgramLicense    `json:"programs"`
	Subsidiaries []loginSubsidiaryLicense `json:"subsidiaries"`
}

type loginDataResponse struct {
	User     interface{}   `json:"user"`
	Token    interface{}   `json:"token"`
	Licenses loginLicenses `json:"licenses"`
}

// Login login app
func Login(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)

	// Validate user and email
	if user.RoleID == 0 {
		// login without using the profile
		if DB.Where("user_name = ? and password = ?", user.UserName, pwd).First(&user).RecordNotFound() {
			if DB.Where("email = ? and password = ?", user.UserName, pwd).First(&user).RecordNotFound() {
				return c.JSON(http.StatusOK, utilities.Response{
					Message: "El nombre de usuario o contraseña es incorecta",
				})
			}
		}
	} else {
		// login with profile
		if DB.Where("user_name = ? and password = ? and role_id = ?", user.UserName, pwd, user.RoleID).First(&user).RecordNotFound() {
			if DB.Where("email = ? and password = ? and role_id = ?", user.UserName, pwd, user.RoleID).First(&user).RecordNotFound() {
				return c.JSON(http.StatusOK, utilities.Response{
					Message: "El nombre de usuario o contraseña es incorecta",
				})
			}
		}
	}

	// Check state user
	if !user.State {
		return c.NoContent(http.StatusForbidden)
	}

	// Prepare response data
	user.Password = ""
	user.Key = ""

	// Query licenses
	loginProgramLicenses := make([]loginProgramLicense, 0)
	if err := DB.Table("program_users").
		Select("programs.id, programs.name").
		Joins("INNER JOIN programs on program_users.program_id = programs.id").
		Where("program_users.user_id = ? AND program_users.license = true", user.ID).
		Scan(&loginProgramLicenses).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}
	loginSubsidiaryLicenses := make([]loginSubsidiaryLicense, 0)
	if err := DB.Table("subsidiary_users").
		Select("subsidiaries.id, subsidiaries.name").
		Joins("INNER JOIN subsidiaries on subsidiary_users.subsidiary_id = subsidiaries.id").
		Where("subsidiary_users.user_id = ? AND subsidiary_users.license = true", user.ID).
		Scan(&loginSubsidiaryLicenses).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Insert new Session
	session := messengermodel.Session{
		UserName:     user.UserName,
		LastActivity: time.Now(),
	}
	if err := DB.Create(&session).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// get token key
	token := utilities.GenerateJWT(user)

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", user.UserName),
		Data: loginDataResponse{
			User:  user,
			Token: token,
			Licenses: loginLicenses{
				Programs:     loginProgramLicenses,
				Subsidiaries: loginSubsidiaryLicenses,
			},
		},
	})
}

type loginStudent struct {
	User  interface{} `json:"user"`
	Token interface{} `json:"token"`
}

// Login by student
func LoginStudent(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)

	// Validate user and email
	if DB.Where("user_name = ? and password = ? and role_id = ?", user.UserName, pwd, 5).First(&user).RecordNotFound() {
		if DB.Where("email = ? and password = ? and role_id = ?", user.UserName, pwd, 5).First(&user).RecordNotFound() {
			return c.JSON(http.StatusOK, utilities.Response{
				Message: fmt.Sprintf("No lo sé rik parece falso"),
			})
		}
	}

	// Check state user
	if !user.State {
		return c.NoContent(http.StatusForbidden)
	}

	// Prepare response data
	user.Password = ""
	user.Key = ""

	// Query student
	//student := institutemodel.Student{}
	//DB.First(&student,institutemodel.Student{UserID: user.ID})

	// Query program by student
	//programs := make([]institutemodel.Program,0)
	//if err := DB.Table("student_programs").
	//    Select("programs.id, programs.name").
	//    Joins("INNER JOIN programs on student_programs.program_id = programs.id").
	//    Where("student_programs.student_id = ?", student.ID).
	//    Scan(&programs).Error; err != nil {
	//        return c.NoContent(http.StatusInternalServerError)
	//}

	// Insert new Session
	session := messengermodel.Session{
		UserName:     user.UserName,
		LastActivity: time.Now(),
	}
	if err := DB.Create(&session).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// get token key
	token := utilities.GenerateJWT(user)

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", user.UserName),
		Data: loginStudent{
			Token: token,
			User:  user,
		},
	})
}

// Login login check
func LoginCheck(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Hash password
	if DB.Where("id = ?", user.ID).First(&user).RecordNotFound() {
		return c.NoContent(http.StatusForbidden)
	}

	// Check state user
	if !user.State {
		return c.NoContent(http.StatusForbidden)
	}

	// Prepare response data
	user.Password = ""
	user.Key = ""

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", user.UserName),
		Data:    user,
	})
}

// ForgotSearch function forgot user search
func ForgotSearch(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Validations
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("Tu búsqueda no arrojó ningún resultado. Vuelve a intentarlo con otros datos."),
		})
	}

	// Generate key validation
	key := (int)(rand.Float32() * 10000000)
	user.Key = fmt.Sprint(key)

	// Update database
	if err := db.Model(&user).Update(user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// SEND EMAIL get html template
	t, err := template.ParseFiles("utilities/email.html")
	if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// SEND EMAIL new buffer
	buf := new(bytes.Buffer)
	err = t.Execute(buf, user)
	if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// SEND EMAIL
	err = config.SendEmail(user.Email, fmt.Sprint(key)+" es el código de recuperación de tu cuenta en REVIEW", buf.String())
	if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Response success api service
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

// ForgotValidate function forgot user validate
func ForgotValidate(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validations
	if err := db.Where("id = ? AND key = ?", user.ID, user.Key).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El número %s que ingresaste no coincide con tu código de seguridad. Vuelve a intentarlo", user.Key),
		})
	}

	// Response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

// ForgotChange function forgot password change
func ForgotChange(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validate
	currentUser := models.User{}
	if err := db.Where("id = ?", user.ID).First(&currentUser).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontro ningun registro con el id %d", user.ID),
		})
	}

	// Encrypted old password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Update
	if err := db.Model(&user).Update(user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Update key
	if err := db.Model(&user).UpdateColumn("key", "").Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
		Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", currentUser.UserName),
	})
}

// GetUsers function get all users
func GetUsers(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	users := make([]models.User, 0)

	// Find users
	if err := db.Where("user_name LIKE ? AND role_id >= ?", "%"+request.Search+"%", currentUser.RoleID).
		Order("id desc").Offset(offset).Limit(request.Limit).Find(&users).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        users,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

// GetUserByID function get user by id
func GetUserByID(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user,
	})
}

// CreateUser function create new user
func CreateUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// Default empty values
	if user.RoleID == 0 {
		user.RoleID = 6
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Insert user in database
	if err := db.Create(&user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user.ID,
		Message: fmt.Sprintf("El usuario %s se registro exitosamente", user.UserName),
	})
}

// UpdateUser function update current user
func UpdateUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	aux := models.User{ID: user.ID}
	if db.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Update user in database
	if err := db.Model(&user).Update(user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}
	if !user.State {
		if err := db.Model(user).UpdateColumn("state", false).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

// DeleteUser function delete user by id
func DeleteUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validate
	db.First(&user, user.ID)
	if user.Freeze {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El usuario %s está protegido por el sistema y no se permite eliminar", user.UserName),
		})
	}

	// Delete user in database
	if err := db.Delete(&user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", user.UserName),
	})
}

// UploadAvatarUser function upload avatar user
func UploadAvatarUser(c echo.Context) error {
	// Read form fields
	idUser := c.FormValue("id")
	user := models.User{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&user, "id = ?", idUser).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Source
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	ccc := sha256.Sum256([]byte(string(user.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/profiles/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	user.Avatar = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&user).Update(user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user,
		Message: fmt.Sprintf("El avatar del usuario %s, se subió correctamente", user.UserName),
	})
}

// ResetPasswordUser function reset password
func ResetPasswordUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&user, "id = ?", user.ID).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Set new password
	cc := sha256.Sum256([]byte(fmt.Sprintf("%d%s", user.ID, user.UserName)))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Update user in database
	if err := db.Model(&user).Update(user).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err) })
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("La contraseña del usuario se cambio exitosamente. ahora su numevacontraseña es %d%s", user.ID, user.UserName),
	})
}

// ChangePasswordUser function change password user
func ChangePasswordUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	aux := models.User{ID: user.ID}
	if db.First(&aux, "id = ?", aux.ID).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", aux.ID),
		})
	}

	// Change password
	if len(user.Password) > 0 {
		// Validate empty length old password
		if len(user.OldPassword) == 0 {
			return c.JSON(http.StatusOK, utilities.Response{
				Message: "Ingrese la contraseña antigua",
			})
		}

		// Hash old password
		ccc := sha256.Sum256([]byte(user.OldPassword))
		old := fmt.Sprintf("%x", ccc)

		// validate old password
		if db.Where("password = ?", old).First(&aux).RecordNotFound() {
			return c.JSON(http.StatusOK, utilities.Response{
				Message: "La contraseña antigua es incorrecta",
			})
		}

		// Set and hash new password
		cc := sha256.Sum256([]byte(user.Password))
		pwd := fmt.Sprintf("%x", cc)
		user.Password = pwd
	}

	// Update user in database
	if err := db.Model(&user).Update(user).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", aux.UserName),
	})
}
