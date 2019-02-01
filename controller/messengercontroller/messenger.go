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
    "strconv"
    "time"
)

type chatMessageShort struct {
	Body        string
	IsRead      bool
	CreatorID   uint
	Date        time.Time 
	RecipientId uint
	ReID        uint
}

type userShort struct {
	ID          uint   `json:"id"`
	Name    string `json:"name"` //
	Avatar      string `json:"avatar"`
	UserGroupID uint   `json:"-"`
}

type chatMessage struct {
	ID          uint       `json:"id"`
	Body        string     `json:"body"`
	BodyType    uint8      `json:"body_type"` // 0 = plain string || 1 == file
	FilePath    string     `json:"file_path"`
	IsRead      bool       `json:"is_read"`
	Date        time.Time  `json:"date"`
	CreatorID   uint       `json:"-"`
	RecipientID uint       `json:"-"`
	ReID        uint       `json:"-"`
	Creator     userShort   `json:"creator, omitempty"`
	Reads       []userShort `json:"reads, omitempty"`
}

type lastMessage struct {
	Body    string    `json:"body"`
	Date    time.Time `json:"date"`
	Mode    string    `json:"mode"` // user || group
	IsRead  bool      `json:"is_read"`
	Contact userShort  `json:"contact"`
    CreatorID uint `json:"creator_id"`
}
type timeSlice []lastMessage

func (p timeSlice) Len() int {
    return len(p)
}

func (p timeSlice) Less(i, j int) bool {
    return p[i].Date.After(p[j].Date)
}

func (p timeSlice) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
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

	// Users
	lastMessages := make([]lastMessage,0)
	for i := range users {
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

		// struct
        lastMessage := lastMessage{
            Body: chatMessageShort[0].Body,
            Date: chatMessageShort[0].Date,
            IsRead: chatMessageShort[0].IsRead,
            Mode: "user",
            CreatorID: chatMessageShort[0].CreatorID,
            Contact: userShort{
                ID: users[i].ID,
                Name: users[i].UserName,
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

func GetMessagesGroup(c echo.Context) error {
	//// Get user token authenticate
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*utilities.Claim)
	//currentUser := claims.User

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
	if err := DB.Raw("SELECT id, body, body_type, file_path, creator_id, date FROM messages WHERE id "+
		" IN ( "+
		"   SELECT message_id FROM message_recipients WHERE recipient_group_id "+
		" IN (SELECT id FROM user_groups WHERE group_id = ?) "+
		" ) ORDER BY messages.id desc "+
		" OFFSET ? LIMIT ?", request.GroupID, offset, request.Limit).Scan(&chatMessages).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if err := DB.Raw("SELECT count(*) FROM messages WHERE id "+
		" IN ( "+
		"   SELECT message_id FROM message_recipients WHERE recipient_group_id "+
		" IN (SELECT id FROM user_groups WHERE group_id = ?) "+
		" ) ", request.GroupID).Scan(&counter).Error; err != nil {
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

		//// Find reads this message
		//reads := make([]userShot, 0)
		//DB.Raw("SELECT users.id, users.user_name, users.avatar, user_groups.id as user_group_id FROM message_recipients " +
		//    "INNER JOIN user_groups ON user_groups.id = message_recipients.recipient_group_id " +
		//    "INNER JOIN users on user_groups.user_id = users.id " +
		//    "WHERE message_recipients.message_id = ?", chatMessages[i].ID).
		//        Scan(&reads)
		//if len(userShots) == 0 {
		//    return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usuario no encontrado")})
		//}
		//chatMessages[i].Reads = reads
	}

	// Validate scroll
	var hasMore = false
	if request.CurrentPage < 10 {
		if request.Limit*request.CurrentPage < counter.Count {
			hasMore = true
		}
	}

	// Read message
	//DB.Model(models.MessageRecipient{}).Where("id in (?)", rIds).Update(models.MessageRecipient{IsRead: true})

	// Return response data scroll reverse
	return c.JSON(http.StatusOK, utilities.ResponseScroll{
		Success:     true,
		Data:        chatMessages,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
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
	mode := c.FormValue("mode")

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Convert string to int
    rID, err := strconv.ParseUint(recipientID, 0, 32)
    if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Valida if user is active
    if mode == "group" {
        userGroup := models.UserGroup{}
        if err := DB.First(&userGroup,models.UserGroup{ UserID: currentUser.ID, GroupID: uint(rID) }) .Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
        if !userGroup.IsActive {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usted está bloqueado en este grupo")})
        }
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
		Date:      time.Now(),
		CreatorID: currentUser.ID,
	}
	if err := TX.Create(&message).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// if is user
	if mode == "user" {
		recipient := models.MessageRecipient{
			RecipientID: uint(rID),
			MessageID:   message.ID,
		}
		if err := TX.Create(&recipient).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// if is group
	if mode == "group" {
		userGroup := make([]models.UserGroup, 0)
		DB.Find(&userGroup, models.UserGroup{GroupID: uint(rID)})
		for _, uGroup := range userGroup {
			recipient := models.MessageRecipient{
				RecipientGroupID: uGroup.ID,
				MessageID:        message.ID,
			}
			if err := TX.Create(&recipient).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
		}
	}

	// Commit transaction
	TX.Commit()

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

	// Valida if user is active
    if request.Mode == "group" {
        userGroup := models.UserGroup{}
        if err := DB.First(&userGroup,models.UserGroup{ UserID: currentUser.ID, GroupID: request.RecipientID }) .Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
        if !userGroup.IsActive {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Usted está bloqueado en este grupo")})
        }
    }

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

	// Create message recipient if USER
	if request.Mode == "user" {
		recipient := models.MessageRecipient{
			RecipientID: request.RecipientID,
			MessageID:   message.ID,
		}

		if err := TX.Create(&recipient).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Create Recipient IF GROUP
	if request.Mode == "group" {
		userGroup := make([]models.UserGroup, 0)
		DB.Find(&userGroup, models.UserGroup{GroupID: request.RecipientID})
		for _, uGroup := range userGroup {
			recipient := models.MessageRecipient{
				RecipientGroupID: uGroup.ID,
				MessageID:        message.ID,
			}
			if err := TX.Create(&recipient).Error; err != nil {
				TX.Rollback()
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
		}
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
	})
}
