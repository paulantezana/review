package librarycontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/olahol/melody"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"golang.org/x/net/websocket"
	"net/http"
)

var Melody *melody.Melody

func init() {
	Melody = melody.New()
	Melody.Config.MaxMessageSize = 1024 * 1024 * 1024
}

func GetCommentsAll(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	comment := models.Comment{}
	if err := c.Bind(&comment); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	comments := make([]models.Comment, 0)

	// Query in database
	if err := DB.Where("parent_id = 0 AND book_id = ?", comment.BookID).Order("id asc").Find(&comments).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// find users by comment
	vote := models.Vote{}
	for i := range comments {
		DB.Model(&comments[i]).Related(&comments[i].User)
		comments[i].User[0].Password = ""
		comments[i].User[0].Key = ""
		comments[i].Children = commentGetChildren(comments[i].ID)

		// Find votes if Has Vote
		vote.CommentID = comments[i].ID
		vote.UserID = currentUser.ID
		if count := DB.Where(&vote).Find(&vote).RowsAffected; count > 0 {
			if vote.Value {
				comments[i].HasVote = 1
			} else {
				comments[i].HasVote = -1
			}
		}
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success: true,
		Data:    comments,
	})
}

func CreateComment(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	comment := models.Comment{}
	if err := c.Bind(&comment); err != nil {
		return err
	}
	comment.UserID = currentUser.ID

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Insert books in database
	if err := DB.Create(&comment).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find current user
	if err := DB.Model(&comment).Related(&comment.User).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	comment.User[0].Password = ""
	comment.User[0].Key = ""

	// Serialize struct to json
	json, err := json.Marshal(&utilities.SocketResponse{
		Type: "create",
		Data: comment,
	})

	// websocket
	origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/api/v1/ws/comment", config.GetConfig().Server.Port)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write(json); err != nil {
		log.Fatal(err)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    comment.ID,
		Message: fmt.Sprintf("El comentario %d se registro correctamente", comment.ID),
	})
}

func UpdateComment(c echo.Context) error {
	// Get data request
	comment := models.Comment{}
	if err := c.Bind(&comment); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update category in database
	rows := db.Model(&comment).Update(&comment).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", comment.ID),
		})
	}

	// Find data
	if err := db.First(&comment, comment.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Serialize struct to json
	json, err := json.Marshal(&utilities.SocketResponse{
		Type: "update",
		Data: comment,
	})

	// websocket
	origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/api/v1/ws/comment", config.GetConfig().Server.Port)

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
		Data:    comment.ID,
		Message: fmt.Sprintf("Los datos del curso %s se actualizaron correctamente", comment.ID),
	})
}

func DeleteComment(c echo.Context) error {
	// Get data request
	comment := models.Comment{}
	if err := c.Bind(&comment); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Find data
	if err := db.First(&comment, comment.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Delete book in database
	if err := db.Delete(&comment).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Delete comments children
	if comment.ID >= 1 {
		if err := db.Delete(models.Comment{}, "parent_id = ?", comment.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Empty data
	comment.Body = ""

	// Serialize struct to json
	json, err := json.Marshal(&utilities.SocketResponse{
		Type: "delete",
		Data: comment,
	})

	// websocket
	origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/api/v1/ws/comment", config.GetConfig().Server.Port)

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
		Data:    comment.ID,
		Message: fmt.Sprintf("El curso %s se elimino correctamente", comment.ID),
	})
}

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
		if err := updateCommentVotes(vote.CommentID, vote.Value); err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Find data
		comment := models.Comment{}
		if err := DB.First(&comment, vote.CommentID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		comment.HasVote = 1

		// Serialize struct to json
		json, err := json.Marshal(&utilities.SocketResponse{
			Type: "update",
			Data: comment,
		})

		// websocket
		origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
		url := fmt.Sprintf("ws://localhost:%s/api/v1/ws/comment", config.GetConfig().Server.Port)

		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := ws.Write(json); err != nil {
			log.Fatal(err)
		}

		// Return response
		return c.JSON(http.StatusCreated, utilities.Response{
			Success: true,
			Message: fmt.Sprintf("Voto registrado"),
		})
		// If exist and update vote false to true OR true to false
	} else if currentVote.Value != vote.Value {
		currentVote.Value = vote.Value

		// Insert books in database
		if err := DB.Save(&currentVote).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Update votes in comments
		if err := updateCommentVotes(vote.CommentID, vote.Value); err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Find data
		comment := models.Comment{}
		if err := DB.First(&comment, vote.CommentID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		if vote.Value {
			comment.HasVote = 1
		} else {
			comment.HasVote = -1
		}

		// Serialize struct to json
		json, err := json.Marshal(&utilities.SocketResponse{
			Type: "update",
			Data: comment,
		})

		// websocket
		origin := fmt.Sprintf("http://localhost:%s/", config.GetConfig().Server.Port)
		url := fmt.Sprintf("ws://localhost:%s/api/v1/ws/comment", config.GetConfig().Server.Port)

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
			Message: fmt.Sprintf("Voto actualizado"),
		})
	}

	// Return response error
	return c.JSON(http.StatusCreated, utilities.Response{
		Message: fmt.Sprintf("Este voto ya está registrado"),
	})
}

// update votes count in table comments
func updateCommentVotes(commentID uint, vote bool) (err error) {
	comment := models.Comment{}

	DB := config.GetConnection()
	defer DB.Close()

	rows := DB.First(&comment, commentID).RowsAffected

	if rows > 0 {
		if vote {
			comment.Votes++
		} else {
			comment.Votes--
		}
		DB.Save(&comment)
	} else {
		err = errors.New(fmt.Sprintf("No se encontró un registro de comentario para asignarle el voto"))
	}
	return
}

func commentGetChildren(id uint) (children []models.Comment) {
	DB := config.GetConnection()
	defer DB.Close()

	DB.Where("parent_id = ?", id).Find(&children)
	for i := range children {
		DB.Model(&children[i]).Related(&children[i].User)
		children[i].User[0].Password = ""
		children[i].User[0].Key = ""
		if children[i].ParentID >= 1 {
			children[i].Children = commentGetChildren(children[i].ID)
		}
	}
	return
}
