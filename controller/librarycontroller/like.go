package librarycontroller

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func CreateLike(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    like := models.Like{}
    if err := c.Bind(&like); err != nil {
        return err
    }
    like.UserID = currentUser.ID

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Validate
    currentVote := models.Like{
        UserID:    currentUser.ID,
        BookID:like.BookID,
    }
    DB.Where(&currentVote).First(&currentVote)

    // If not exist
    if currentVote.ID == 0 {
        // Insert books in database
        if err := DB.Create(&like).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }

        // Response data
        return c.JSON(http.StatusCreated, utilities.Response{
            Success: true,
            Message: fmt.Sprintf("Voto registrado"),
        })
        // If exist and update vote false to true OR true to false
    } else if currentVote.Stars != like.Stars {
        currentVote.Stars = like.Stars

        // Insert books in database
        if err := DB.Save(&currentVote).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }

        // Response data
        return c.JSON(http.StatusOK, utilities.Response{
            Success: true,
            Message: fmt.Sprintf("Voto actualizado"),
        })
    }

    // Return response error
    return c.JSON(http.StatusCreated, utilities.Response{
        Message: fmt.Sprintf("Este voto ya está registrado"),
    })
}
