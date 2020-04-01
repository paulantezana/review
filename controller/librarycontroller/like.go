package librarycontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func CreatePostLike(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	postLike := models.PostLike{}
	if err := c.Bind(&postLike); err != nil {
		return err
	}
	postLike.UserID = currentUser.ID

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate
	currentPostLike := models.PostLike{
		UserID: currentUser.ID,
		PostID: postLike.PostID,
	}
	DB.Where(&currentPostLike).First(&currentPostLike)

	// If not exist
	if currentPostLike.ID == 0 {
		// Insert books in database
		if err := DB.Create(&postLike).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		//

		// Response data
		return c.JSON(http.StatusCreated, utilities.Response{
			Success: true,
			Message: fmt.Sprintf("Voto registrado"),
		})
		// If exist and update vote false to true OR true to false
	} else if currentPostLike.Stars != postLike.Stars {
		currentPostLike.Stars = postLike.Stars

		// Insert books in database
		if err := DB.Save(&currentPostLike).Error; err != nil {
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
