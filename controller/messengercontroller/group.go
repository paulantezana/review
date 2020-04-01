package messengercontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func GetMssGroupsScroll(c echo.Context) error {
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	counter := utilities.Counter{}

	// Query mssGroups
	mssGroups := make([]models.MssGroup, 0)
	if err := DB.Raw("SELECT * FROM mssGroups WHERE id IN "+
		"( SELECT mssGroup_id FROM user_mssGroups WHERE user_id = ? AND is_active = true)   "+
		"ORDER BY id asc LIMIT ? OFFSET ? ", currentUser.ID, request.Limit, offset).
		Scan(&mssGroups).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	if err := DB.Raw("SELECT count(*) FROM mssGroups WHERE id IN "+
		"( SELECT mssGroup_id FROM user_mssGroups WHERE user_id = ? )", currentUser.ID).
		Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// query las messages
	lastMssMessages := make([]lastMssMessage, 0)
	for _, mssGroup := range mssGroups {
		// Find last message
		lastMssMessageByMssGroup := make([]lastMssMessage, 0)
		if err := DB.Debug().Table("mssGroup_messages").
			Select("mssGroup_messages.id, mssGroup_messages.body, mssGroup_messages.created_at, mssGroup_messages.creator_id").
			Joins("INNER JOIN mssGroup_message_recipients ON mssGroup_messages.id = mssGroup_message_recipients.message_id").
			Where("mssGroup_message_recipients.recipient_mssGroup_id = ?", mssGroup.ID).
			Limit(1).
			Order("mssGroup_messages.id DESC").
			Scan(&lastMssMessageByMssGroup).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// struct response message
		lastMssMessage := lastMssMessage{
			Mode: "user",
			Contact: userShort{
				ID:     mssGroup.ID,
				Name:   mssGroup.Name,
				Avatar: mssGroup.Avatar,
			},
		}

		// Set messages
		if len(lastMssMessageByMssGroup) >= 1 {
			lastMssMessage.Body = lastMssMessageByMssGroup[0].Body
			lastMssMessage.CreatedAt = lastMssMessageByMssGroup[0].CreatedAt
			lastMssMessage.IsRead = lastMssMessageByMssGroup[0].IsRead
			lastMssMessage.CreatorID = lastMssMessageByMssGroup[0].CreatorID
		}

		// Add last message
		lastMssMessages = append(lastMssMessages, lastMssMessage)
	}

	// Order By date
	lastMssMessagesSorted := make(timeSlice, 0, len(lastMssMessages))
	for _, lasM := range lastMssMessages {
		lastMssMessagesSorted = append(lastMssMessagesSorted, lasM)
	}
	sort.Sort(lastMssMessagesSorted)

	// Validate scroll
	var hasMore = false
	if request.CurrentPage < 10 {
		if request.Limit*request.CurrentPage < counter.Count {
			hasMore = true
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.ResponseScroll{
		Success:     true,
		Data:        lastMssMessagesSorted,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

type userMssGroupResponse struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active" gorm:"default:'true'"`
	IsAdmin   bool      `json:"is_admin"`

	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type mssGroupResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`

	Users []userMssGroupResponse `json:"users"`
}

func GetMssGroupByID(c echo.Context) error {
	// Get data request
	mssGroup := models.MssGroup{}
	if err := c.Bind(&mssGroup); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	mssGroupResponse := mssGroupResponse{}
	if err := DB.Raw("SELECT * FROM mssGroups WHERE id = ?", mssGroup.ID).Scan(&mssGroupResponse).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	userMssGroupResponses := make([]userMssGroupResponse, 0)
	if err := DB.Table("user_mssGroups").
		Select("user_mssGroups.id, user_mssGroups.created_at, user_mssGroups.is_active, user_mssGroups.is_admin, user_mssGroups.user_id, users.user_name as name, users.avatar").
		Joins("INNER JOIN users ON user_mssGroups.user_id = users.id").
		Where("user_mssGroups.mssGroup_id = ?", mssGroup.ID).
		Order("user_mssGroups.id asc").
		Scan(&userMssGroupResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	mssGroupResponse.Users = userMssGroupResponses

	// Query current student Name
	for i := range userMssGroupResponses {
		student := models.Student{}
		DB.First(&student, models.Student{UserID: userMssGroupResponses[i].UserID})
		if student.ID >= 1 {
			userMssGroupResponses[i].Name = student.FullName
		} else {
			teacher := models.Teacher{}
			DB.First(&teacher, models.Teacher{UserID: userMssGroupResponses[i].UserID})
			if teacher.ID >= 1 {
				userMssGroupResponses[i].Name = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
			}
		}

	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    mssGroupResponse,
	})
}

func CreateMssGroup(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	mssGroup := models.MssGroup{}
	if err := c.Bind(&mssGroup); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// add data
	mssGroup.Date = time.Now()
	mssGroup.MssUserGroups = append(mssGroup.MssUserGroups, models.MssUserGroup{
		UserID:  currentUser.ID,
		IsAdmin: true,
	})

	// Insert courses in database
	if err := DB.Create(&mssGroup).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    mssGroup.ID,
		Message: fmt.Sprintf("El canal %s se registro correctamente", mssGroup.Name),
	})
}

func AddUsers(c echo.Context) error {
	// Get data request
	mssUserGroups := make([]models.MssUserGroup, 0)
	if err := c.Bind(&mssUserGroups); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Start transaction
	TX := DB.Begin()

	// Add new users
	count := 0
	for _, uMssGroup := range mssUserGroups {
		if err := TX.Create(&uMssGroup).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		count++
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se añadieron %d usuarios.", count),
	})
}

func UpdateMssGroup(c echo.Context) error {
	// Get data request
	mssGroup := models.MssGroup{}
	if err := c.Bind(&mssGroup); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	//
	mssGroup.Date = time.Now()

	// Insert courses in database
	rows := DB.Model(&mssGroup).Update(&mssGroup).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %s", mssGroup.Name),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    mssGroup.ID,
		Message: fmt.Sprintf("Los datos del canal %s se modificarón correctamente.", mssGroup.Name),
	})
}

// Enable or disable mssGroup
func IsActiveMssGroup(c echo.Context) error {
	// Get data request
	mssGroup := models.MssGroup{}
	if err := c.Bind(&mssGroup); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Disable
	if err := DB.Model(&mssGroup).UpdateColumn("is_active", mssGroup.IsActive).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    mssGroup.ID,
		Message: fmt.Sprintf("Los datos del canal se modificaron correctamente."),
	})
}

// Enable or disable user in mssGroup
func IsActiveUserMssGroup(c echo.Context) error {
	// Get data request
	userMssGroup := models.MssUserGroup{}
	if err := c.Bind(&userMssGroup); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Disable
	if err := DB.Model(&userMssGroup).Where("user_id = ? AND mssGroup_id = ?", userMssGroup.UserID, userMssGroup.MssGroupID).
		UpdateColumn("is_active", userMssGroup.IsActive).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    userMssGroup.ID,
		Message: fmt.Sprintf("Los datos del canal se modificaron correctamente."),
	})
}

// UploadAvatarUser function upload avatar user
func UploadAvatarMssGroup(c echo.Context) error {
	// Read form fields
	idMssGroup := c.FormValue("id")

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	mssGroup := models.MssGroup{}
	if db.First(&mssGroup, "id = ?", idMssGroup).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", mssGroup.ID),
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
	ccc := sha256.Sum256([]byte(string(mssGroup.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/profiles/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	mssGroup.Avatar = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&mssGroup).Update(mssGroup).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    mssGroup,
		Message: fmt.Sprintf("El avatar del canal %s, se subió correctamente", mssGroup.Name),
	})
}
