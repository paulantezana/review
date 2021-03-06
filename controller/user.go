package controller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/paulantezana/review/models"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
)

type loginProgramLicense struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type loginDataResponse struct {
	User     interface{}   `json:"user"`
	Token    interface{}   `json:"token"`
	Licenses []licenseUser `json:"licenses"`
}

// Login login app
func Login(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
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

	// Exception users student and invited
	if !(user.RoleID >= 1 && user.RoleID <= 4) {
		return c.NoContent(http.StatusForbidden)
	}

	// Prepare response data
	user.Password = ""
	user.Key = ""

	// Insert new Session
	session := models.Session{
		UserID:       user.ID,
		LastActivity: time.Now(),
	}
	if err := DB.Create(&session).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		},
	})
}

type loginStudent struct {
	Programs []loginProgramLicense `json:"programs"`
	User     models.User           `json:"user"`
	Token    interface{}           `json:"token"`
}

// Login by student
func LoginStudent(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
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
	student := models.Student{}
	if err := DB.First(&student, models.Student{UserID: user.ID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query programs
	loginProgramLicenses := make([]loginProgramLicense, 0)
	if err := DB.Debug().Raw("SELECT id, name FROM programs WHERE id "+
		"IN (SELECT program_id FROM student_programs WHERE student_id = ?)", student.ID).
		Scan(&loginProgramLicenses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Insert new Session
	session := models.Session{
		UserID:       user.ID,
		LastActivity: time.Now(),
	}
	if err := DB.Create(&session).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// get token key
	token := utilities.GenerateJWT(user)

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", user.UserName),
		Data: loginStudent{
			Token:    token,
			Programs: loginProgramLicenses,
			User:     user,
		},
	})
}

// Login login check
func LoginUserCheck(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
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
		Message: fmt.Sprintf("Verificación exitosa %s", user.UserName),
		Data:    user,
	})
}

// Login password login check
func LoginPasswordCheck(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)

	// Hash password
	if DB.Where("id = ? and password = ?", user.ID, pwd).First(&user).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Contraseña incorrecta")})
	}

	// Check state user
	if !user.State {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("No tiene permisos para realizar ningún tipo de acción.")})
	}

	// Prepare response data
	user.Password = ""
	user.Key = ""

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("La verificación de su contraseña fue exitosamente. usuario: %s", user.UserName),
		Data:    user.ID,
	})
}

// ForgotSearch function forgot user search
func ForgotSearch(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validations
	if err := DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("Tu búsqueda no arrojó ningún resultado. Vuelve a intentarlo con otros datos."),
		})
	}

	// Generate key validation
	key := (int)(rand.Float32() * 10000000)
	user.Key = fmt.Sprint(key)

	// Update database
	if err := DB.Model(&user).Update(user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query Database Get Settings
	con := models.Setting{}
	DB.First(&con)

	// SEND EMAIL get html template
	t, err := template.ParseFiles("templates/email.html")
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// SEND EMAIL new buffer
	buf := new(bytes.Buffer)
	err = t.Execute(buf, user)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// SEND EMAIL
	err = provider.SendEmail(
		con.PrefixShortName+" "+con.Institute,
		user.Email,
		fmt.Sprint(key)+" es el código de recuperación de tu cuenta",
		buf.String(),
	)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Update key
	if err := db.Model(&user).UpdateColumn("key", "").Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	users := make([]models.User, 0)

	// Find users
	if err := DB.Where("user_name LIKE ? AND role_id >= ?", "%"+request.Search+"%", currentUser.RoleID).
		Order("id desc").Offset(offset).Limit(request.Limit).Find(&users).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// find
	for i := range users {
		student := models.Student{}
		DB.First(&student, models.Student{UserID: users[i].ID})
		if student.ID >= 1 {
			users[i].UserName = student.FullName
		} else {
			teacher := models.Teacher{}
			DB.First(&teacher, models.Teacher{UserID: users[i].ID})
			if teacher.ID >= 1 {
				users[i].UserName = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
			}
		}
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

// GetUsers function get all users
type searchUsersResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}

func SearchUsers(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Find users
	users := make([]searchUsersResponse, 0)
	if err := DB.Raw("SELECT id, user_name, avatar FROM users "+
		"WHERE lower(user_name) LIKE lower(?) "+
		"OR id IN (SELECT user_id FROM teachers WHERE lower(first_name) LIKE lower(?) LIMIT 20) "+
		"OR id IN (SELECT user_id FROM students WHERE lower(full_name) LIKE lower(?) LIMIT 20) "+
		"LIMIT 30", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%").Scan(&users).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Queries
	newUsers := make([]searchUsersResponse, 0)
	for i := range users {
		// ignore current user
		if users[i].ID == currentUser.ID {
			continue
		}

		// search user
		nUser := searchUsersResponse{
			ID:       users[i].ID,
			UserName: users[i].UserName,
			Avatar:   users[i].Avatar,
		}

		// Query current student Name
		student := models.Student{}
		DB.First(&student, models.Student{UserID: nUser.ID})
		if student.ID >= 1 {
			nUser.UserName = student.FullName
		} else {
			teacher := models.Teacher{}
			DB.First(&teacher, models.Teacher{UserID: nUser.ID})
			if teacher.ID >= 1 {
				nUser.UserName = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
			}
		}

		// append child
		newUsers = append(newUsers, nUser)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    newUsers,
	})
}

// GetUserByID function get user by id
func GetUserByID(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Default empty values
	if user.RoleID == 0 {
		user.RoleID = 6
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Insert user in database
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !user.State {
		if err := db.Model(user).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
	db := provider.GetConnection()
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
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
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
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
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
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

type licenseUser struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	UserID       uint   `json:"user_id"`
	SubsidiaryID uint   `json:"subsidiary_id"`
	Name         string `json:"name"`
	License      bool   `json:"license"`
	Programs     []struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		ProgramID uint   `json:"program_id"`
		License   bool   `json:"license"`
		Name      string `json:"name"`
	} `json:"programs"`
}

type licenseUserResponse struct {
	Licenses []licenseUser `json:"licenses"`
	User     models.User   `json:"user"`
}

func GetLicensePostLogin(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query licenses
	licenseUsers := make([]licenseUser, 0)
	licenseUsersFilter := make([]licenseUser, 0)
	if err := DB.Table("subsidiary_users").
		Select("subsidiary_users.id, subsidiary_users.user_id, subsidiary_users.subsidiary_id, subsidiary_users.license, subsidiaries.name").
		Joins("INNER JOIN subsidiaries ON subsidiaries.id = subsidiary_users.subsidiary_id").
		Where("subsidiary_users.user_id = ?", currentUser.ID).
		Scan(&licenseUsers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	for i, subsidiaryR := range licenseUsers {
		DB.Table("program_users").
			Select("program_users.id, program_users.user_id, program_users.program_id, program_users.license, programs.name").
			Joins("INNER JOIN programs ON programs.id = program_users.program_id").
			Where("program_users.user_id = ? AND programs.subsidiary_id = ? AND license = true", currentUser.ID, subsidiaryR.SubsidiaryID).
			Scan(&licenseUsers[i].Programs)

		if licenseUsers[i].License == true || len(licenseUsers[i].Programs) >= 1 {
			licenseUsersFilter = append(licenseUsersFilter, licenseUsers[i])
		}
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data: licenseUserResponse{
			Licenses: licenseUsersFilter,
			User:     currentUser,
		},
	})
}

func GetLicenseUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query Subsidiaries
	subsidiaries := make([]models.Subsidiary, 0)
	if err := DB.Raw("SELECT * FROM subsidiaries WHERE id NOT IN (SELECT subsidiary_id FROM subsidiary_users WHERE user_id = ?)", user.ID).
		Scan(&subsidiaries).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Start Transaction
	TR := DB.Begin()

	// Insert SubsidiaryUsers
	for _, subsidiary := range subsidiaries {
		subsidiaryUser := models.SubsidiaryUser{
			UserID:       user.ID,
			SubsidiaryID: subsidiary.ID,
		}
		if err := TR.Create(&subsidiaryUser).Error; err != nil {
			TR.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Get all subsidiary users
	subsidiaryUsers := make([]models.SubsidiaryUser, 0)
	if err := DB.Where("user_id = ?", user.ID).
		Find(&subsidiaryUsers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// loop subsidiaryUsers
	for _, subsidiaryUser := range subsidiaryUsers {
		// Query Programs
		programs := make([]models.Program, 0)
		if err := DB.Raw("SELECT * FROM programs WHERE id NOT IN (SELECT program_id  FROM program_users WHERE user_id = ? AND subsidiary_user_id = ?) AND subsidiary_id = ?", user.ID, subsidiaryUser.ID, subsidiaryUser.SubsidiaryID).
			Scan(&programs).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Insert SubsidiaryUsers
		for _, program := range programs {
			programUser := models.ProgramUser{
				UserID:           user.ID,
				ProgramID:        program.ID,
				SubsidiaryUserID: subsidiaryUser.ID,
			}
			if err := TR.Create(&programUser).Error; err != nil {
				TR.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
		}
	}

	// End Transaction
	TR.Commit()

	// Query SubsidiaryUsers
	licenseUsers := make([]licenseUser, 0)
	if err := DB.Table("subsidiary_users").
		Select("subsidiary_users.id, subsidiary_users.user_id, subsidiary_users.subsidiary_id, subsidiary_users.license, subsidiaries.name").
		Joins("INNER JOIN subsidiaries ON subsidiaries.id = subsidiary_users.subsidiary_id").
		Where("subsidiary_users.user_id = ?", user.ID).
		Scan(&licenseUsers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	for i, subsidiaryR := range licenseUsers {
		DB.Table("program_users").
			Select("program_users.id, program_users.user_id, program_users.program_id, program_users.license, programs.name").
			Joins("INNER JOIN programs ON programs.id = program_users.program_id").
			Where("program_users.user_id = ? AND programs.subsidiary_id = ?", user.ID, subsidiaryR.SubsidiaryID).
			Scan(&licenseUsers[i].Programs)
	}

	// Response data
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    licenseUsers,
	})
}
