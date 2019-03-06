package messengercontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func GetGroupsScroll(c echo.Context) error {
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
	DB := config.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	counter := utilities.Counter{}

	// Query groups
	groups := make([]models.Group, 0)
	if err := DB.Raw("SELECT * FROM groups WHERE id IN "+
		"( SELECT group_id FROM user_groups WHERE user_id = ? AND is_active = true)   "+
		"ORDER BY id asc LIMIT ? OFFSET ? ", currentUser.ID, request.Limit, offset).
		Scan(&groups).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	if err := DB.Raw("SELECT count(*) FROM groups WHERE id IN "+
		"( SELECT group_id FROM user_groups WHERE user_id = ? )", currentUser.ID).
		Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// query las messages
	lastMessages := make([]lastMessage, 0)
	for _, group := range groups {
		// Find last message
		lastMessageByGroup := make([]lastMessage, 0)
		if err := DB.Debug().Table("group_messages").
			Select("group_messages.id, group_messages.body, group_messages.created_at, group_messages.creator_id").
			Joins("INNER JOIN group_message_recipients ON group_messages.id = group_message_recipients.message_id").
			Where("group_message_recipients.recipient_group_id = ?", group.ID).
			Limit(1).
			Order("group_messages.id DESC").
			Scan(&lastMessageByGroup).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// struct response message
		lastMessage := lastMessage{
			Mode: "user",
			Contact: userShort{
				ID:     group.ID,
				Name:   group.Name,
				Avatar: group.Avatar,
			},
		}

		// Set messages
		if len(lastMessageByGroup) >= 1 {
			lastMessage.Body = lastMessageByGroup[0].Body
			lastMessage.CreatedAt = lastMessageByGroup[0].CreatedAt
			lastMessage.IsRead = lastMessageByGroup[0].IsRead
			lastMessage.CreatorID = lastMessageByGroup[0].CreatorID
		}

		// Add last message
		lastMessages = append(lastMessages, lastMessage)
	}

	// Order By date
	lastMessagesSorted := make(timeSlice, 0, len(lastMessages))
	for _, lasM := range lastMessages {
		lastMessagesSorted = append(lastMessagesSorted, lasM)
	}
	sort.Sort(lastMessagesSorted)

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
		Data:        lastMessagesSorted,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

type userGroupResponse struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active" gorm:"default:'true'"`
	IsAdmin   bool      `json:"is_admin"`

	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type groupResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`

	Users []userGroupResponse `json:"users"`
}

func GetGroupByID(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	groupResponse := groupResponse{}
	if err := DB.Raw("SELECT * FROM groups WHERE id = ?", group.ID).Scan(&groupResponse).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	userGroupResponses := make([]userGroupResponse, 0)
	if err := DB.Table("user_groups").
		Select("user_groups.id, user_groups.created_at, user_groups.is_active, user_groups.is_admin, user_groups.user_id, users.user_name as name, users.avatar").
		Joins("INNER JOIN users ON user_groups.user_id = users.id").
		Where("user_groups.group_id = ?", group.ID).
		Order("user_groups.id asc").
		Scan(&userGroupResponses).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	groupResponse.Users = userGroupResponses

	// Query current student Name
	for i := range userGroupResponses {
		student := models.Student{}
		DB.First(&student, models.Student{UserID: userGroupResponses[i].UserID})
		if student.ID >= 1 {
			userGroupResponses[i].Name = student.FullName
		} else {
			teacher := models.Teacher{}
			DB.First(&teacher, models.Teacher{UserID: userGroupResponses[i].UserID})
			if teacher.ID >= 1 {
				userGroupResponses[i].Name = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
			}
		}

	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    groupResponse,
	})
}

func CreateGroup(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// add data
	group.Date = time.Now()
	group.UserGroups = append(group.UserGroups, models.UserGroup{
		UserID:  currentUser.ID,
		IsAdmin: true,
	})

	// Insert courses in database
	if err := DB.Create(&group).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("El grupo %s se registro correctamente", group.Name),
	})
}

func AddUsers(c echo.Context) error {
	// Get data request
	userGroups := make([]models.UserGroup, 0)
	if err := c.Bind(&userGroups); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Start transaction
	TX := DB.Begin()

	// Add new users
	count := 0
	for _, uGroup := range userGroups {
		if err := TX.Create(&uGroup).Error; err != nil {
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

func UpdateGroup(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	//
	group.Date = time.Now()

	// Insert courses in database
	rows := DB.Model(&group).Update(&group).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %s", group.Name),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("Los datos del grupo %s se desactivo correctamente.", group.Name),
	})
}

// Enable or disable group
func IsActiveGroup(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Disable
	if err := DB.Model(&group).UpdateColumn("is_active", group.IsActive).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("Los datos del grupo se modificaron correctamente."),
	})
}

// Enable or disable user in group
func IsActiveUserGroup(c echo.Context) error {
	// Get data request
	userGroup := models.UserGroup{}
	if err := c.Bind(&userGroup); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Disable
	if err := DB.Model(&userGroup).Where("user_id = ? AND group_id = ?", userGroup.UserID, userGroup.GroupID).
		UpdateColumn("is_active", userGroup.IsActive).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    userGroup.ID,
		Message: fmt.Sprintf("Los datos del grupo se modificaron correctamente."),
	})
}

// UploadAvatarUser function upload avatar user
func UploadAvatarGroup(c echo.Context) error {
	// Read form fields
	idGroup := c.FormValue("id")

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	group := models.Group{}
	if db.First(&group, "id = ?", idGroup).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", group.ID),
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
	ccc := sha256.Sum256([]byte(string(group.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/profiles/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	group.Avatar = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&group).Update(group).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    group,
		Message: fmt.Sprintf("El avatar del grupo %s, se subió correctamente", group.Name),
	})
}
