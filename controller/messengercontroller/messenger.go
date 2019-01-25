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

func GetUsersMessage(c echo.Context) error {
    // Get user token authenticate
    //user := c.Get("user").(*jwt.Token)
    //claims := user.Claims.(*utilities.Claim)
    //currentUser := claims.User

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query users
    users := make([]models.User,0)
    if err := DB.Find(&users).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }



    // Return response
    return c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Data: users,
    })
}

type createMessageRequest struct {
    RecipientID      uint `json:"recipient_id"`
    Body string `json:"body"`
    ParentID uint `json:"parent_id"`
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
        Body:request.Body,
        Date: time.Now(),
        CreatorID: currentUser.ID,
    }
    if err := TX.Create(&message).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // create message recipient
    recipient := models.MessageRecipient{
        RecipientID: request.RecipientID,
        MessageID: message.ID,
    }

    if err := TX.Create(&recipient).Error; err != nil {
        TX.Rollback()
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Commit transaction
    TX.Commit()

    // Return response
    return c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Message: "OK",
    })
}
