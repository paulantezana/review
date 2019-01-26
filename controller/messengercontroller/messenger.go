package messengercontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

type chatMessage struct {
	Body        string    `json:"body"`
	IsRead      bool      `json:"is_read"`
	CreatorID   uint      `json:"creator_id"`
	Date        time.Time `json:"date"`
	RecipientId uint      `json:"-"`
	ReID        uint      `json:"-"`
}

type userMessage struct {
	ID           uint          `json:"id"`
	UserName     string        `json:"user_name"` //
	Avatar       string        `json:"avatar"`
	LastActivity time.Time     `json:"last_activity"`
	LastMessages []chatMessage `json:"last_messages"`
}

func GetUsersMessageScroll(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query users
	users := make([]userMessage, 0)
	if err := DB.Table("users").Select("id, user_name, avatar").
		Scan(&users).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	for i := range users {
		// Find las activity
		session := models.Session{}
		DB.First(&session, models.Session{UserID: users[i].ID})
		users[i].LastActivity = session.LastActivity

		// Query semesters
		chatMessage := make([]chatMessage, 0)
		if err := DB.Table("messages").
			Select("messages.body, message_recipients.is_read, messages.creator_id, messages.date").
			Joins("INNER JOIN message_recipients ON messages.id = message_recipients.message_id").
			Where("messages.creator_id = ? AND message_recipients.recipient_id = ?", users[i].ID, currentUser.ID).
			Or("messages.creator_id = ? AND message_recipients.recipient_id = ?", currentUser.ID, users[i].ID).
			Limit(1).
			Order("messages.id desc").
			Scan(&chatMessage).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		users[i].LastMessages = chatMessage
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    users,
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
	requestUser := models.User{}
	if err := c.Bind(&requestUser); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query semesters
	chatMessages := make([]chatMessage, 0)
	if err := DB.Table("messages").
		Select("messages.body, message_recipients.is_read, messages.creator_id, messages.date, message_recipients.id as re_id,  message_recipients.recipient_id  ").
		Joins("INNER JOIN message_recipients ON messages.id = message_recipients.message_id").
		Where("messages.creator_id = ? AND message_recipients.recipient_id = ?", requestUser.ID, currentUser.ID).
		Or("messages.creator_id = ? AND message_recipients.recipient_id = ?", currentUser.ID, requestUser.ID).
		Scan(&chatMessages).Error; err != nil {
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

	// Read message
	DB.Model(models.MessageRecipient{}).Where("id in (?)", rIds).Update(models.MessageRecipient{IsRead: true})

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    chatMessages,
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
