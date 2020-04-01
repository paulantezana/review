package messengercontroller

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
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
	ID     uint   `json:"id"`
	Name   string `json:"name"` //
	Avatar string `json:"avatar"`
}

// Use in list chat mssMessages scroll reverse
type chatMssMessage struct {
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

type lastMssMessage struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Mode      string    `json:"mode"` // user || group
	IsRead    bool      `json:"is_read"`
	Contact   userShort `json:"contact"`
	CreatorID uint      `json:"creator_id"`
}
type timeSlice []lastMssMessage

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

// Get all users width mssMessages
func GetUsersMssMessageScroll(c echo.Context) error {
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

	// Query users
	users := make([]models.User, 0)
	if err := DB.Raw("SELECT * FROM users "+
		"WHERE  id IN ( SELECT creator_id FROM mssMessages "+
		"INNER JOIN mssMessage_recipients ON mssMessages.id = mssMessage_recipients.mssMessage_id "+
		"WHERE mssMessage_recipients.recipient_id = ? "+
		") OR id IN ( SELECT recipient_id FROM mssMessage_recipients "+
		"INNER JOIN mssMessages ON mssMessage_recipients.mssMessage_id = mssMessages.id "+
		"WHERE creator_id = ?) "+
		"OFFSET ? LIMIT ?", currentUser.ID, currentUser.ID, offset, request.Limit).Scan(&users).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if err := DB.Raw("SELECT count(*) FROM users "+
		"WHERE  id IN ( SELECT creator_id FROM mssMessages "+
		"INNER JOIN mssMessage_recipients ON mssMessages.id = mssMessage_recipients.mssMessage_id "+
		"WHERE mssMessage_recipients.recipient_id = ? "+
		") OR id IN ( SELECT recipient_id FROM mssMessage_recipients "+
		"INNER JOIN mssMessages ON mssMessage_recipients.mssMessage_id = mssMessages.id "+
		"WHERE creator_id = ?)", currentUser.ID, currentUser.ID).Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query last mssMessages
	lastMssMessages := make([]lastMssMessage, 0)
	for i := range users {
		// Query mssMessages
		lastMssMessageByUser := make([]lastMssMessage, 0)
		if err := DB.Table("mssMessages").
			Select("mssMessages.body, mssMessage_recipients.is_read, mssMessages.creator_id, mssMessages.created_at").
			Joins("INNER JOIN mssMessage_recipients ON mssMessages.id = mssMessage_recipients.mssMessage_id").
			Where("mssMessages.creator_id = ? AND mssMessage_recipients.recipient_id = ?", users[i].ID, currentUser.ID).
			Or("mssMessages.creator_id = ? AND mssMessage_recipients.recipient_id = ?", currentUser.ID, users[i].ID).
			Limit(1).
			Order("mssMessages.id desc").
			Scan(&lastMssMessageByUser).Error; err != nil {
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

		// struct response mssMessage
		lastMssMessage := lastMssMessage{
			Body:      lastMssMessageByUser[0].Body,
			CreatedAt: lastMssMessageByUser[0].CreatedAt,
			IsRead:    lastMssMessageByUser[0].IsRead,
			Mode:      "user",
			CreatorID: lastMssMessageByUser[0].CreatorID,
			Contact: userShort{
				ID:     users[i].ID,
				Name:   users[i].UserName,
				Avatar: users[i].Avatar,
			},
		}
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

// Get mssMessages by group
func GetMssMessagesByGroup(c echo.Context) error {
	// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

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

	// Query chatMssMessage scroll
	chatMssMessages := make([]chatMssMessage, 0)
	if err := DB.Raw("SELECT group_mssMessages.id, group_mssMessages.body, group_mssMessages.body_type, group_mssMessages.file_path, group_mssMessages.created_at, group_mssMessages.creator_id  FROM group_mssMessages "+
		"INNER JOIN group_mssMessage_recipients ON group_mssMessages.id = group_mssMessage_recipients.mssMessage_id "+
		"WHERE group_mssMessage_recipients.recipient_group_id = ? "+
		"GROUP BY group_mssMessages.id, group_mssMessages.body, group_mssMessages.body_type, group_mssMessages.created_at, group_mssMessages.created_at "+
		"ORDER BY group_mssMessages.created_at DESC "+
		" OFFSET ? LIMIT ?", request.MssGroupID, offset, request.Limit).Scan(&chatMssMessages).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	DB.Raw("SELECT count(*) FROM group_mssMessages "+
		"INNER JOIN group_mssMessage_recipients ON group_mssMessages.id = group_mssMessage_recipients.mssMessage_id "+
		"WHERE group_mssMessage_recipients.recipient_group_id = ? "+
		"GROUP BY group_mssMessages.id", request.MssGroupID).Scan(&counter)

	// find user creator info
	for i := range chatMssMessages {
		userShots := make([]userShort, 0)
		DB.Raw("SELECT * FROM users WHERE id = ?", chatMssMessages[i].CreatorID).Scan(&userShots)
		if len(userShots) == 0 {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usuario no encontrado")})
		}
		chatMssMessages[i].Creator = userShots[0]
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
		Data:        chatMssMessages,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

// Get mssMessages by user
func GetMssMessages(c echo.Context) error {
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

	// Query semesters
	var total uint
	chatMssMessages := make([]chatMssMessage, 0)
	if err := DB.Table("mssMessages").
		Select("mssMessages.body, mssMessages.body_type, mssMessages.file_path, mssMessage_recipients.is_read, mssMessages.creator_id, mssMessages.created_at, "+
			"mssMessage_recipients.id as re_id,  mssMessage_recipients.recipient_id  ").
		Joins("INNER JOIN mssMessage_recipients ON mssMessages.id = mssMessage_recipients.mssMessage_id").
		Where("mssMessages.creator_id = ? AND mssMessage_recipients.recipient_id = ?", request.UserID, currentUser.ID).
		Or("mssMessages.creator_id = ? AND mssMessage_recipients.recipient_id = ?", currentUser.ID, request.UserID).
		Order("mssMessages.id desc").Limit(request.Limit).Offset(offset).
		Scan(&chatMssMessages).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Get ids and read true
	var rIds = make([]uint, 0)
	for i, m := range chatMssMessages {
		if chatMssMessages[i].RecipientID == currentUser.ID {
			rIds = append(rIds, m.ReID)
			chatMssMessages[i].IsRead = true
		}

		// Find data creator
		userShots := make([]userShort, 0)
		DB.Raw("SELECT * FROM users WHERE id = ?", chatMssMessages[i].CreatorID).Scan(&userShots)
		if len(userShots) == 0 {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usuario no encontrado")})
		}
		chatMssMessages[i].Creator = userShots[0]
	}

	// Validate scroll
	var hasMore = false
	if request.CurrentPage < 10 {
		if request.Limit*request.CurrentPage < total {
			hasMore = true
		}
	}

	// Read mssMessage
	DB.Model(models.MssMessageRecipient{}).Where("id in (?)", rIds).Update(models.MssMessageRecipient{IsRead: true})

	// Return response data scroll reverse
	return c.JSON(http.StatusOK, utilities.ResponseScroll{
		Success:     true,
		Data:        chatMssMessages,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

func CreateMssMessageFileUpload(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	recipientID := c.FormValue("recipient_id")

	// get connection
	DB := provider.GetConnection()
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

	// create struct mssMessage
	mssMessage := models.MssMessage{
		Body:      file.Filename,
		BodyType:  1,
		FilePath:  fileSRC,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&mssMessage).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// if is user
	recipient := models.MssMessageRecipient{
		RecipientID: uint(rID),
		MssMessageID:   mssMessage.ID,
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
	chatMssMessage := chatMssMessage{
		ID:        mssMessage.ID,
		Body:      mssMessage.Body,
		BodyType:  mssMessage.BodyType,
		FilePath:  mssMessage.FilePath,
		CreatedAt: mssMessage.CreatedAt,
		Recipient: userShort{
			ID:     userRecipient.ID,
			Name:   userRecipient.Name,
			Avatar: userRecipient.Avatar,
		},
		Creator: userShort{
			ID:     currentUser.ID,
			Name:   currentUser.UserName,
			Avatar: currentUser.Avatar,
		},
	}

	json, err := json.Marshal(&utilities.SocketResponse{
		Type:   "chat",
		Action: "create",
		Data:   chatMssMessage,
	})

	// Socket
	origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/chat", provider.GetConfig().Server.Port)

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

func CreateMssMessageFileUploadByGroup(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	recipientID := c.FormValue("recipient_id")

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Convert string to int
	rID, err := strconv.ParseUint(recipientID, 0, 32)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Valida if user is active
	mssUserGroup := models.MssUserGroup{}
	if err := DB.First(&mssUserGroup, models.MssUserGroup{UserID: currentUser.ID, MssGroupID: uint(rID)}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !mssUserGroup.IsActive {
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

	// create struct mssMessage
	mssGroupMessage := models.MssGroupMessage{
		Body:      file.Filename,
		BodyType:  1,
		FilePath:  fileSRC,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&mssGroupMessage).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// create recipient mssMessage
	mssUserGroups := make([]models.MssUserGroup, 0)
	DB.Find(&mssUserGroups, models.MssUserGroup{MssGroupID: uint(rID)})
	for _, uGroup := range mssUserGroups {
		recipient := models.MssGroupMessageRecipient{
			RecipientGroupID: uint(rID),
			RecipientID:      uGroup.UserID,
			MssGroupMessageID:        mssGroupMessage.ID,
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
	chatMssMessage := chatMssMessage{
		ID:        mssGroupMessage.ID,
		Body:      mssGroupMessage.Body,
		BodyType:  mssGroupMessage.BodyType,
		FilePath:  mssGroupMessage.FilePath,
		CreatedAt: mssGroupMessage.CreatedAt,
		Mode:      "group",
		Recipient: userShort{
			ID:     groupRecipient.ID,
			Name:   groupRecipient.Name,
			Avatar: groupRecipient.Avatar,
		},
		Creator: userShort{
			ID:     currentUser.ID,
			Name:   currentUser.UserName,
			Avatar: currentUser.Avatar,
		},
	}

	json, err := json.Marshal(&utilities.SocketResponse{
		Type:   "chat",
		Action: "create",
		Data:   chatMssMessage,
	})

	// Socket
	origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/chat", provider.GetConfig().Server.Port)

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

type createMssMessageRequest struct {
	RecipientID uint   `json:"recipient_id"`
	Body        string `json:"body"`
	Mode        string `json:"mode"` // user || group
}

// Create mssMessage by user
func CreateMssMessage(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := createMssMessageRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Start transaction
	TX := DB.Begin()

	// create struct mssMessage
	mssMessage := models.MssMessage{
		Body:      request.Body,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&mssMessage).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create mssMessage recipient if USER
	recipient := models.MssMessageRecipient{
		RecipientID: request.RecipientID,
		MssMessageID:   mssMessage.ID,
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
	chatMssMessage := chatMssMessage{}
	if request.Mode == "user" {
		chatMssMessage.ID = mssMessage.ID
		chatMssMessage.Body = mssMessage.Body
		chatMssMessage.BodyType = mssMessage.BodyType
		chatMssMessage.FilePath = mssMessage.FilePath
		chatMssMessage.CreatedAt = mssMessage.CreatedAt
		chatMssMessage.Mode = request.Mode
		chatMssMessage.Recipient = userShort{
			ID:     userRecipient.ID,
			Name:   userRecipient.Name,
			Avatar: userRecipient.Avatar,
		}
		chatMssMessage.Creator = userShort{
			ID:     currentUser.ID,
			Name:   currentUser.UserName,
			Avatar: currentUser.Avatar,
		}
	}

	json, err := json.Marshal(&utilities.SocketResponse{
		Type:   "chat",
		Action: "create",
		Data:   chatMssMessage,
	})

	// Socket
	origin := fmt.Sprintf("%s:%s", provider.GetConfig().Server.Host, provider.GetConfig().Server.Port)
	url := fmt.Sprintf("%s:%s/ws/chat", provider.GetConfig().Server.Socket, provider.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 1 == %s", err)})
	}
	if _, err := ws.Write(json); err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 2 == %s", err)})
	}

	// Send chat Notices
	// Websocket las notices
	getUnreadMssMessages(models.User{ID: request.RecipientID}, true)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
        Message: "OK",
	})
}

// Create mssMessage by group
func CreateMssGroupMessage(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := createMssMessageRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Valida if user is active
	mssUserGroup := models.MssUserGroup{}
	if err := DB.First(&mssUserGroup, models.MssUserGroup{UserID: currentUser.ID, MssGroupID: request.RecipientID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !mssUserGroup.IsActive {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usted está bloqueado en este grupo")})
	}

	// Start transaction
	TX := DB.Begin()

	// create struct mssGroupMessage
	mssGroupMessage := models.MssGroupMessage{
		Body:      request.Body,
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&mssGroupMessage).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create Recipient mssMessage by group
	mssUserGroups := make([]models.MssUserGroup, 0)
	DB.Find(&mssUserGroups, models.MssUserGroup{MssGroupID: request.RecipientID})
	for _, uGroup := range mssUserGroups {
		recipient := models.MssGroupMessageRecipient{
			RecipientGroupID: request.RecipientID,
			RecipientID:      uGroup.UserID,
			MssGroupMessageID:        mssGroupMessage.ID,
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
	chatMssMessage := chatMssMessage{
		ID:        mssGroupMessage.ID,
		Body:      mssGroupMessage.Body,
		BodyType:  mssGroupMessage.BodyType,
		FilePath:  mssGroupMessage.FilePath,
		CreatedAt: mssGroupMessage.CreatedAt,
		Mode:      request.Mode,
		Recipient: userShort{
			ID:     groupRecipient.ID,
			Name:   groupRecipient.Name,
			Avatar: groupRecipient.Avatar,
		},
		Creator: userShort{
			ID:     currentUser.ID,
			Name:   currentUser.UserName,
			Avatar: currentUser.Avatar,
		},
	}

	json, err := json.Marshal(&utilities.SocketResponse{
		Type:   "chat",
		Action: "create",
		Data:   chatMssMessage,
	})

	// Socket
	origin := fmt.Sprintf("%s:%s/", provider.GetConfig().Server.Host, provider.GetConfig().Server.Port)
	url := fmt.Sprintf("%s:%s/ws/chat", provider.GetConfig().Server.Socket, provider.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 1 == %s", err)})
	}
	if _, err := ws.Write(json); err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("ERROR 2 == %s", err)})
	}

	// Send chat Notices
	// Websocket las notices
	getUnreadMssMessages(models.User{ID: request.RecipientID}, true)

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
        Message: "OK",
	})
}

func UnreadMssMessages(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get  unread last mssMessages
	notices := getUnreadMssMessages(models.User{ID: currentUser.ID}, false)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
        Message: "OK",
		Data:    notices,
	})
}

func getUnreadMssMessages(u models.User, socket bool) []utilities.Notice {
	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// query
	lastMssMessages := make([]lastMssMessage, 0)
	if err := DB.Table("mssMessages").
		Select("mssMessages.body, mssMessage_recipients.is_read, mssMessages.creator_id, mssMessages.created_at").
		Joins("INNER JOIN mssMessage_recipients ON mssMessages.id = mssMessage_recipients.mssMessage_id").
		Where("mssMessage_recipients.recipient_id = ? AND mssMessage_recipients.is_read = false", u.ID).
		Limit(1).
		Order("mssMessages.id desc").
		Scan(&lastMssMessages).Error; err != nil {
		log.Fatal(err)
	}

	notices := make([]utilities.Notice, 0)
	for i := range lastMssMessages {
		user := models.User{}
		if err := DB.First(&user, models.User{ID: lastMssMessages[i].CreatorID}).Error; err != nil {
			log.Fatal(err)
		}

		notice := utilities.Notice{
			ID:          lastMssMessages[i].CreatorID,
			Title:       user.UserName,
			Avatar:      user.Avatar,
			Description: lastMssMessages[i].Body,
			Date:        lastMssMessages[i].CreatedAt,
			RecipientID: u.ID,
			Type:        "mssMessage",
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
		origin := fmt.Sprintf("%s:%s/", provider.GetConfig().Server.Host, provider.GetConfig().Server.Port)
		url := fmt.Sprintf("%s:%s/ws/chat", provider.GetConfig().Server.Socket, provider.GetConfig().Server.Port)

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
