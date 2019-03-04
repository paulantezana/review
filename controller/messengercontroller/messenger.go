package messengercontroller

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// Use in user chat
type userShort struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"` //
	Avatar      string `json:"avatar"`
}

// Use in list chat messages scroll reverse
type chatMessage struct {
	ID          uint        `json:"id"`
	Body        string      `json:"body"`
	BodyType    uint8       `json:"body_type"` // 0 = plain string || 1 == file
	FilePath    string      `json:"file_path"`
	IsRead      bool        `json:"is_read"`
	CreatedAt   time.Time   `json:"created_at"`
	CreatorID   uint        `json:"-"`
	RecipientID uint        `json:"-"`
	Mode        string      `json:"mode"`
	ReID        uint        `json:"-"`
	Recipient   userShort   `json:"recipient, omitempty"`
	Creator     userShort   `json:"creator, omitempty"`
	Reads       []userShort `json:"reads, omitempty"`
}

type lastMessage struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Mode      string    `json:"mode"` // user || group
	IsRead    bool      `json:"is_read"`
	Contact   userShort `json:"contact"`
	CreatorID uint      `json:"creator_id"`
}
type timeSlice []lastMessage

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	return p[i].CreatedAt.After(p[j].CreatedAt)
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

var Melody *melody.Melody

func init() {
	Melody = melody.New()
	Melody.Config.MaxMessageSize = 1024 * 1024 * 1024
}

// Get all users width messages
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
	users := make([]models.User, 0)
	if err := DB.Raw("SELECT * FROM users "+
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

	// Query last messages
	lastMessages := make([]lastMessage, 0)
	for i := range users {
		// Query messages
		lastMessageByUser := make([]lastMessage, 0)
		if err := DB.Table("messages").
			Select("messages.body, message_recipients.is_read, messages.creator_id, messages.date").
			Joins("INNER JOIN message_recipients ON messages.id = message_recipients.message_id").
			Where("messages.creator_id = ? AND message_recipients.recipient_id = ?", users[i].ID, currentUser.ID).
			Or("messages.creator_id = ? AND message_recipients.recipient_id = ?", currentUser.ID, users[i].ID).
			Limit(1).
			Order("messages.id desc").
			Scan(&lastMessageByUser).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Query current Full Name
		student := models.Student{}
		DB.First(&student, models.Student{UserID: users[i].ID}) // In student
		if student.ID >= 1 {
			users[i].UserName = student.FullName
		} else {
			teacher := models.Teacher{}
			DB.First(&teacher, models.Teacher{UserID: users[i].ID}) // In teacher
			if teacher.ID >= 1 {
				users[i].UserName = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
			}
		}

		// struct response message
		lastMessage := lastMessage{
			Body:      lastMessageByUser[0].Body,
			CreatedAt: lastMessageByUser[0].CreatedAt,
			IsRead:    lastMessageByUser[0].IsRead,
			Mode:      "user",
			CreatorID: lastMessageByUser[0].CreatorID,
			Contact: userShort{
				ID:     users[i].ID,
				Name:   users[i].UserName,
				Avatar: users[i].Avatar,
			},
		}
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

// Get messages by group
func GetMessagesByGroup(c echo.Context) error {
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

	// Query chatMessage scroll
	chatMessages := make([]chatMessage, 0)
	if err := DB.Debug().Raw("SELECT group_messages.id, group_messages.body, group_messages.body_type, group_messages.file_path, group_messages.created_at, group_messages.creator_id  FROM group_messages " +
        "INNER JOIN group_message_recipients ON group_messages.id = group_message_recipients.message_id " +
        "WHERE group_message_recipients.recipient_group_id = ? " +
        "GROUP BY group_messages.id, group_messages.body, group_messages.body_type, group_messages.created_at, group_messages.created_at " +
        "ORDER BY group_messages.created_at DESC " +
		" OFFSET ? LIMIT ?", request.GroupID, offset, request.Limit).Scan(&chatMessages).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if err := DB.Raw("SELECT count(*) FROM group_messages "+
        "INNER JOIN group_message_recipients ON group_messages.id = group_message_recipients.message_id " +
        "WHERE group_message_recipients.recipient_group_id = ? " +
        "GROUP BY group_messages.id, group_messages.body, group_messages.body_type, group_messages.created_at, group_messages.created_at " +
		" ", request.GroupID).Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// find user creator info
	for i := range chatMessages {
		userShots := make([]userShort, 0)
		DB.Raw("SELECT * FROM users WHERE id = ?", chatMessages[i].CreatorID).Scan(&userShots)
		if len(userShots) == 0 {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usuario no encontrado")})
		}
		chatMessages[i].Creator = userShots[0]
	}

	// Validate scroll
	var hasMore = false
	if request.CurrentPage < 10 {
		if request.Limit*request.CurrentPage < counter.Count {
			hasMore = true
		}
	}

	// Return response data scroll reverse
	return c.JSON(http.StatusOK, utilities.ResponseScroll{
		Success:     true,
		Data:        chatMessages,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

// Get messages by user
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
		Select("messages.body, messages.body_type, messages.file_path, message_recipients.is_read, messages.creator_id, messages.date, "+
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
		if chatMessages[i].RecipientID == currentUser.ID {
			rIds = append(rIds, m.ReID)
			chatMessages[i].IsRead = true
		}

		// Find data creator
		userShots := make([]userShort, 0)
		DB.Raw("SELECT * FROM users WHERE id = ?", chatMessages[i].CreatorID).Scan(&userShots)
		if len(userShots) == 0 {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usuario no encontrado")})
		}
		chatMessages[i].Creator = userShots[0]
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
	recipientID := c.FormValue("recipient_id")

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Convert string to int
	rID, err := strconv.ParseUint(recipientID, 0, 32)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

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

	// Start transaction
	TX := DB.Begin()

	// create struct message
	message := models.Message{
		Body:      file.Filename,
		BodyType:  1,
		FilePath:  fileSRC,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&message).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// if is user
	recipient := models.MessageRecipient{
		RecipientID: uint(rID),
		MessageID:   message.ID,
	}
	if err := TX.Create(&recipient).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Find recipient user detail
	userRecipient := userShort{}
    DB.Raw("SELECT id, user_name as name, avatar FROM users WHERE id = ? LIMIT 1", uint(rID)).Scan(&userRecipient)

	// Socket init send data
	chatMessage := chatMessage{
		ID : message.ID,
		Body : message.Body,
		BodyType : message.BodyType,
		FilePath : message.FilePath,
		CreatedAt : message.CreatedAt,
		Recipient : userShort{
			ID:     userRecipient.ID,
			Name:   userRecipient.Name,
			Avatar: userRecipient.Avatar,
		},
		Creator : userShort{
			ID:     currentUser.ID,
			Name:   currentUser.UserName,
			Avatar: currentUser.Avatar,
		},
    }

	json, err := json.Marshal(&utilities.SocketResponse{
		Type:   "chat",
		Action: "create",
		Data:   chatMessage,
	})

	// Socket
	origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/chat", config.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write(json); err != nil {
		log.Fatal(err)
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}

func CreateMessageFileUploadByGroup(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	recipientID := c.FormValue("recipient_id")

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Convert string to int
	rID, err := strconv.ParseUint(recipientID, 0, 32)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Valida if user is active
	userGroup := models.UserGroup{}
	if err := DB.First(&userGroup, models.UserGroup{UserID: currentUser.ID, GroupID: uint(rID)}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !userGroup.IsActive {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usted está bloqueado en este grupo")})
	}

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

	// Start transaction
	TX := DB.Begin()

	// create struct message
	groupMessage := models.GroupMessage{
		Body:      file.Filename,
		BodyType:  1,
		FilePath:  fileSRC,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&groupMessage).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// create recipient message
	userGroups := make([]models.UserGroup, 0)
	DB.Find(&userGroups, models.UserGroup{GroupID: uint(rID)})
	for _, uGroup := range userGroups {
		recipient := models.GroupMessageRecipient{
			RecipientGroupID: uint(rID),
			RecipientID:      uGroup.UserID,
			MessageID:        groupMessage.ID,
		}
		if err := TX.Create(&recipient).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Commit transaction
	TX.Commit()

    // Find recipient group detail
    groupRecipient := userShort{}
    DB.Raw("SELECT * FROM groups WHERE id = ? LIMIT 1", uint(rID)).Scan(&groupRecipient)

    //Socket init send data
    chatMessage := chatMessage{
        ID : groupMessage.ID,
        Body : groupMessage.Body,
        BodyType : groupMessage.BodyType,
        FilePath : groupMessage.FilePath,
        CreatedAt : groupMessage.CreatedAt,
        Mode : "group",
        Recipient : userShort{
            ID:     groupRecipient.ID,
            Name:   groupRecipient.Name,
            Avatar: groupRecipient.Avatar,
        },
        Creator : userShort{
            ID:     currentUser.ID,
            Name:   currentUser.UserName,
            Avatar: currentUser.Avatar,
        },
    }

	json, err := json.Marshal(&utilities.SocketResponse{
	   Type:   "chat",
	   Action: "create",
	   Data:   chatMessage,
	})

	// Socket
	origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/chat", config.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
	   log.Fatal(err)
	}
	if _, err := ws.Write(json); err != nil {
	   log.Fatal(err)
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}

type createMessageRequest struct {
	RecipientID uint   `json:"recipient_id"`
	Body        string `json:"body"`
	Mode        string `json:"mode"` // user || group
}

// Create message by user
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
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&message).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create message recipient if USER
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

	// Find recipient user detail
	userRecipient := userShort{}
	if request.Mode == "user" {
		DB.Raw("SELECT id, user_name as name, avatar FROM users WHERE id = ? LIMIT 1", request.RecipientID).Scan(&userRecipient)
	}

	// Socket init send data
	chatMessage := chatMessage{}
	if request.Mode == "user" {
		chatMessage.ID = message.ID
		chatMessage.Body = message.Body
		chatMessage.BodyType = message.BodyType
		chatMessage.FilePath = message.FilePath
		chatMessage.CreatedAt = message.CreatedAt
		chatMessage.Mode = request.Mode
		chatMessage.Recipient = userShort{
			ID:     userRecipient.ID,
			Name:   userRecipient.Name,
			Avatar: userRecipient.Avatar,
		}
		chatMessage.Creator = userShort{
			ID:     currentUser.ID,
			Name:   currentUser.UserName,
			Avatar: currentUser.Avatar,
		}
	}

	json, err := json.Marshal(&utilities.SocketResponse{
		Type:   "chat",
		Action: "create",
		Data:   chatMessage,
	})

	// Socket
	origin := fmt.Sprintf("%s:%s/", config.GetConfig().Server.Host, config.GetConfig().Server.Port)
	url := fmt.Sprintf("%s:%s/ws/chat", config.GetConfig().Server.Socket, config.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 1 == %s", err)})
	}
	if _, err := ws.Write(json); err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 2 == %s", err)})
	}

	// Send chat Notices
	// Websocket las notices
	getUnreadMessages(models.User{ID: request.RecipientID}, true)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}

// Create message by group
func CreateGroupMessage(c echo.Context) error {
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

	// Valida if user is active
	userGroup := models.UserGroup{}
	if err := DB.First(&userGroup, models.UserGroup{UserID: currentUser.ID, GroupID: request.RecipientID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !userGroup.IsActive {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usted está bloqueado en este grupo")})
	}

	// Start transaction
	TX := DB.Begin()

	// create struct groupMessage
	groupMessage := models.GroupMessage{
		Body:      request.Body,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&groupMessage).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create Recipient message by group
	userGroups := make([]models.UserGroup, 0)
	DB.Find(&userGroups, models.UserGroup{GroupID: request.RecipientID})
	for _, uGroup := range userGroups {
		recipient := models.GroupMessageRecipient{
			RecipientGroupID: request.RecipientID,
			RecipientID:      uGroup.UserID,
			MessageID:        groupMessage.ID,
		}
		if err := TX.Create(&recipient).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Commit transaction
	TX.Commit()

	// Find recipient group detail
	groupRecipient := userShort{}
    DB.Raw("SELECT * FROM groups WHERE id = ? LIMIT 1", request.RecipientID).Scan(&groupRecipient)

	//Socket init send data
	chatMessage := chatMessage{
        ID : groupMessage.ID,
        Body : groupMessage.Body,
        BodyType : groupMessage.BodyType,
        FilePath : groupMessage.FilePath,
        CreatedAt : groupMessage.CreatedAt,
        Mode : request.Mode,
        Recipient : userShort{
            ID:     groupRecipient.ID,
            Name:   groupRecipient.Name,
            Avatar: groupRecipient.Avatar,
        },
        Creator : userShort{
            ID:     currentUser.ID,
            Name:   currentUser.UserName,
            Avatar: currentUser.Avatar,
        },
    }

	json, err := json.Marshal(&utilities.SocketResponse{
	   Type:   "chat",
	   Action: "create",
	   Data:   chatMessage,
	})

	// Socket
	origin := fmt.Sprintf("%s:%s/", config.GetConfig().Server.Host, config.GetConfig().Server.Port)
	url := fmt.Sprintf("%s:%s/ws/chat", config.GetConfig().Server.Socket, config.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
	   return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 1 == %s", err)})
	}
	if _, err := ws.Write(json); err != nil {
	   return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 2 == %s", err)})
	}

	// Send chat Notices
	// Websocket las notices
	getUnreadMessages(models.User{ID: request.RecipientID}, true)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}

func UnreadMessages(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get  unread last messages
	notices := getUnreadMessages(models.User{ID: currentUser.ID}, false)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
		Data:    notices,
	})
}

func getUnreadMessages(u models.User, socket bool) []utilities.Notice {
	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// query
    lastMessages := make([]lastMessage, 0)
	if err := DB.Table("messages").
		Select("messages.body, message_recipients.is_read, messages.creator_id, messages.date").
		Joins("INNER JOIN message_recipients ON messages.id = message_recipients.message_id").
		Where("message_recipients.recipient_id = ? AND message_recipients.is_read = false", u.ID).
		Limit(1).
		Order("messages.id desc").
		Scan(&lastMessages).Error; err != nil {
		log.Fatal(err)
	}

	notices := make([]utilities.Notice, 0)
	for i := range lastMessages {
		user := models.User{}
		if err := DB.First(&user, models.User{ID: lastMessages[i].CreatorID}).Error; err != nil {
			log.Fatal(err)
		}

		notice := utilities.Notice{
			ID:          lastMessages[i].CreatorID,
			Title:       user.UserName,
			Avatar:      user.Avatar,
			Description: lastMessages[i].Body,
			Date:        lastMessages[i].CreatedAt,
			RecipientID: u.ID,
			Type:        "message",
		}

		notices = append(notices, notice)
	}

	// Socket prepare data
	if socket {
		json, err := json.Marshal(&utilities.SocketResponse{
			Type:   "notice",
			Action: "info",
			Data:   notices,
		})

		// Socket
		origin := fmt.Sprintf("%s:%s/", config.GetConfig().Server.Host, config.GetConfig().Server.Port)
		url := fmt.Sprintf("%s:%s/ws/chat", config.GetConfig().Server.Socket, config.GetConfig().Server.Port)

		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := ws.Write(json); err != nil {
			log.Fatal(err)
		}
	}

	return notices
}
