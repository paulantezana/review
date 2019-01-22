package librarycontroller

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func CreateVote(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	vote := models.Vote{}
	if err := c.Bind(&vote); err != nil {
		return err
	}
	vote.UserID = currentUser.ID

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Validate
	currentVote := models.Vote{
		UserID:    currentUser.ID,
		CommentID: vote.CommentID,
	}
	DB.Where(&currentVote).First(&currentVote)

	// If not exist
	if currentVote.ID == 0 {
		// Insert books in database
		if err := DB.Create(&vote).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Update votes in comments
		if err := updateCommentVotes(vote.CommentID, vote.Value, false); err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		c.JSON(http.StatusCreated, utilities.Response{
			Success: true,
			Message: fmt.Sprintf("Voto registrado"),
		})
		// If exist and update vote false to true OR true to false
	} else if currentVote.Value != vote.Value {
		currentVote.Value = vote.Value

		// Insert books in database
		if err := DB.Save(&vote).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Update votes in comments
		if err := updateCommentVotes(vote.CommentID, vote.Value, true); err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		c.JSON(http.StatusOK, utilities.Response{
			Success: true,
			Message: fmt.Sprintf("Voto actualizado"),
		})
	}

	// Return response error
	return c.JSON(http.StatusCreated, utilities.Response{
		Message: fmt.Sprintf("Este voto ya está registrado"),
	})
}

// update votes count in table comments
func updateCommentVotes(commentID uint, vote bool, isUpdate bool) (err error) {
	comment := models.Comment{}

	DB := config.GetConnection()
	defer DB.Close()

	rows := DB.First(&comment, commentID).RowsAffected

	if rows > 0 {
		if vote {
			comment.Votes++
			if isUpdate {
				comment.Votes++
			}
		} else {
			comment.Votes--
			if isUpdate {
				comment.Votes--
			}
		}
		DB.Save(&comment)
	} else {
		err = errors.New(fmt.Sprintf("No se encontró un registro de comentario para asignarle el voto"))
	}
	return
}
