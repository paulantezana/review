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
	"strconv"
	"time"
)

type chatMessageShort struct {
	Body        string    `json:"body"`
	IsRead      bool      `json:"is_read"`
	CreatorID   uint      `json:"creator_id"`
	Date        time.Time `json:"date"`
	RecipientId uint      `json:"-"`
	ReID        uint      `json:"-"`
}

type chatMessage struct {
	Body        string    `json:"body"`
	BodyType    uint8     `json:"body_type"` // 0 = plain string || 1 == file
	FilePath    string    `json:"file_path"`
	IsRead      bool      `json:"is_read"`
	CreatorID   uint      `json:"creator_id"`
	Date        time.Time `json:"date"`
	RecipientId uint      `json:"-"`
	ReID        uint      `json:"-"`
}

type userMessage struct {
	ID           uint               `json:"id"`
	UserName     string             `json:"user_name"` //
	Avatar       string             `json:"avatar"`
	LastActivity time.Time          `json:"last_activity"`
	LastMessages []chatMessageShort `json:"last_messages"`
}

func GetUsersMessageScroll(c echo.Context) error {
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

	// Query users
	users := make([]userMessage, 0)
	if err := DB.Raw("SELECT id, user_name, avatar FROM users "+
		"WHERE  id IN ( SELECT creator_id FROM messages "+
		"INNER JOIN message_recipients ON messages.id = message_recipients.message_id "+
		"WHERE message_recipients.recipient_id = ? "+
		") OR id IN ( SELECT recipient_id FROM message_recipients "+
		"INNER JOIN messages ON message_recipients.message_id = messages.id "+
		"WHERE creator_id = ?) "+
		"OFFSET ? LIMIT ?", currentUser.ID, currentUser.ID, offset, request.Limit).Scan(&users).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if err := DB.Raw("SELECT count(*) FROM users "+
		"WHERE  id IN ( SELECT creator_id FROM messages "+
		"INNER JOIN message_recipients ON messages.id = message_recipients.message_id "+
		"WHERE message_recipients.recipient_id = ? "+
		") OR id IN ( SELECT recipient_id FROM message_recipients "+
		"INNER JOIN messages ON message_recipients.message_id = messages.id "+
		"WHERE creator_id = ?)", currentUser.ID, currentUser.ID).Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Users
	for i := range users {
		// Find las activity
		session := models.Session{}
		DB.First(&session, models.Session{UserID: users[i].ID})
		users[i].LastActivity = session.LastActivity

		// Query semesters
		chatMessageShort := make([]chatMessageShort, 0)
		if err := DB.Table("messages").
			Select("messages.body, message_recipients.is_read, messages.creator_id, messages.date").
			Joins("INNER JOIN message_recipients ON messages.id = message_recipients.message_id").
			Where("messages.creator_id = ? AND message_recipients.recipient_id = ?", users[i].ID, currentUser.ID).
			Or("messages.creator_id = ? AND message_recipients.recipient_id = ?", currentUser.ID, users[i].ID).
			Limit(1).
			Order("messages.id desc").
			Scan(&chatMessageShort).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		users[i].LastMessages = chatMessageShort

		// Query current student Name
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
		Data:        users,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

type createMessageRequest struct {
	RecipientID uint   `json:"recipient_id"`
	Body        string `json:"body"`
	ParentID    uint   `json:"parent_id"`
}

// Get messages
func GetMessages(c echo.Context) error {
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

	// Query semesters
	var total uint
	chatMessages := make([]chatMessage, 0)
	if err := DB.Table("messages").
		Select("messages.body, messages.body_type, messages.file_path, message_recipients.is_read, messages.creator_id, messages.date, " +
		    "message_recipients.id as re_id,  message_recipients.recipient_id  ").
		Joins("INNER JOIN message_recipients ON messages.id = message_recipients.message_id").
		Where("messages.creator_id = ? AND message_recipients.recipient_id = ?", request.UserID, currentUser.ID).
		Or("messages.creator_id = ? AND message_recipients.recipient_id = ?", currentUser.ID, request.UserID).
		Order("messages.id desc").Limit(request.Limit).Offset(offset).
		Scan(&chatMessages).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get ids and read true
	var rIds = make([]uint, 0)
	for i, m := range chatMessages {
		if chatMessages[i].RecipientId == currentUser.ID {
			rIds = append(rIds, m.ReID)
			chatMessages[i].IsRead = true
		}
	}

	// Validate scroll
	var hasMore = false
	if request.CurrentPage < 10 {
		if request.Limit*request.CurrentPage < total {
			hasMore = true
		}
	}

	// Read message
	DB.Model(models.MessageRecipient{}).Where("id in (?)", rIds).Update(models.MessageRecipient{IsRead: true})

	// Return response data scroll reverse
	return c.JSON(http.StatusOK, utilities.ResponseScroll{
		Success:     true,
		Data:        chatMessages,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

func CreateMessageFileUpload(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	recipientIDs := c.FormValue("recipient_id")

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	ccc := sha256.Sum256([]byte(time.Now().String() + string(currentUser.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	fileSRC := "static/chat/" + name
	dst, err := os.Create(fileSRC)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Start transaction
	TX := DB.Begin()

	// create struct message
	message := models.Message{
		Body:      file.Filename,
		BodyType:  1,
		FilePath:  fileSRC,
		Date:      time.Now(),
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&message).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// create message recipient
	rIDs, err := strconv.ParseUint(recipientIDs, 0, 32)
	if err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	recipient := models.MessageRecipient{
		RecipientID: uint(rIDs),
		MessageID:   message.ID,
	}
	if err := TX.Create(&recipient).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}

func CreateMessage(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := createMessageRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Start transaction
	TX := DB.Begin()

	// create struct message
	message := models.Message{
		Body:      request.Body,
		Date:      time.Now(),
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&message).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// create message recipient
	recipient := models.MessageRecipient{
		RecipientID: request.RecipientID,
		MessageID:   message.ID,
	}

	if err := TX.Create(&recipient).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}
