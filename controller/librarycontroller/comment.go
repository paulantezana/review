package librarycontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/olahol/melody"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"golang.org/x/net/websocket"
	"net/http"
)

var Melody *melody.Melody

func init() {
	Melody = melody.New()
	Melody.Config.MaxMessageSize = 1024 * 1024 * 1024
}

func GetPostCommentsAll(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	postComment := models.PostComment{}
	if err := c.Bind(&postComment); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	postComments := make([]models.PostComment, 0)

	// Query in database
	if err := DB.Where("parent_id = 0 AND book_id = ?", postComment.PostID).Order("id asc").Find(&postComments).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// find users by postComment
	postVote := models.PostVote{}
	for i := range postComments {
		DB.Model(&postComments[i]).Related(&postComments[i].User)
		postComments[i].User[0].Password = ""
		postComments[i].User[0].TempKey = ""
		postComments[i].Children = postCommentGetChildren(postComments[i].ID)

		// Find postVotes if Has PostVote
		postVote.PostCommentID = postComments[i].ID
		postVote.UserID = currentUser.ID
		if count := DB.Where(&postVote).Find(&postVote).RowsAffected; count > 0 {
			if postVote.Value {
				postComments[i].HasPostVote = 1
			} else {
				postComments[i].HasPostVote = -1
			}
		}
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success: true,
		Data:    postComments,
	})
}

func CreatePostComment(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	postComment := models.PostComment{}
	if err := c.Bind(&postComment); err != nil {
		return err
	}
	postComment.UserID = currentUser.ID

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Insert books in database
	if err := DB.Create(&postComment).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Find current user
	if err := DB.Model(&postComment).Related(&postComment.User).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	postComment.User[0].Password = ""
	postComment.User[0].TempKey = ""

	// Serialize struct to json
	json, err := json.Marshal(&utilities.SocketResponse{
		Type: "create",
		Data: postComment,
	})

	// websocket
	origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/postComment", provider.GetConfig().Server.Port)

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
		Data:    postComment.ID,
		Message: fmt.Sprintf("El comentario %d se registro correctamente", postComment.ID),
	})
}

func UpdatePostComment(c echo.Context) error {
	// Get data request
	postComment := models.PostComment{}
	if err := c.Bind(&postComment); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update category in database
	rows := db.Model(&postComment).Update(&postComment).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", postComment.ID),
		})
	}

	// Find data
	if err := db.First(&postComment, postComment.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Serialize struct to json
	json, err := json.Marshal(&utilities.SocketResponse{
		Type: "update",
		Data: postComment,
	})

	// websocket
	origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/postComment", provider.GetConfig().Server.Port)

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
		Data:    postComment.ID,
		Message: fmt.Sprintf("Los datos del curso %s se actualizaron correctamente", postComment.ID),
	})
}

func DeletePostComment(c echo.Context) error {
	// Get data request
	postComment := models.PostComment{}
	if err := c.Bind(&postComment); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Find data
	if err := db.First(&postComment, postComment.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Delete book in database
	if err := db.Delete(&postComment).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Delete postComments children
	if postComment.ID >= 1 {
		if err := db.Delete(models.PostComment{}, "parent_id = ?", postComment.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Empty data
	postComment.Body = ""

	// Serialize struct to json
	json, err := json.Marshal(&utilities.SocketResponse{
		Type: "delete",
		Data: postComment,
	})

	// websocket
	origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
	url := fmt.Sprintf("ws://localhost:%s/ws/postComment", provider.GetConfig().Server.Port)

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
		Data:    postComment.ID,
		Message: fmt.Sprintf("El curso %s se elimino correctamente", postComment.ID),
	})
}

func CreatePostVote(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	postVote := models.PostVote{}
	if err := c.Bind(&postVote); err != nil {
		return err
	}
	postVote.UserID = currentUser.ID

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate
	currentPostVote := models.PostVote{
		UserID:    currentUser.ID,
		PostCommentID: postVote.PostCommentID,
	}
	DB.Where(&currentPostVote).First(&currentPostVote)

	// If not exist
	if currentPostVote.ID == 0 {
		// Insert books in database
		if err := DB.Create(&postVote).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Update postVotes in postComments
		if err := updatePostCommentPostVotes(postVote.PostCommentID, postVote.Value); err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Find data
		postComment := models.PostComment{}
		if err := DB.First(&postComment, postVote.PostCommentID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		postComment.HasPostVote = 1

		// Serialize struct to json
		json, err := json.Marshal(&utilities.SocketResponse{
			Type: "update",
			Data: postComment,
		})

		// websocket
		origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
		url := fmt.Sprintf("ws://localhost:%s/ws/postComment", provider.GetConfig().Server.Port)

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
		// If exist and update postVote false to true OR true to false
	} else if currentPostVote.Value != postVote.Value {
		currentPostVote.Value = postVote.Value

		// Insert books in database
		if err := DB.Save(&currentPostVote).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Update postVotes in postComments
		if err := updatePostCommentPostVotes(postVote.PostCommentID, postVote.Value); err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Find data
		postComment := models.PostComment{}
		if err := DB.First(&postComment, postVote.PostCommentID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		if postVote.Value {
			postComment.HasPostVote = 1
		} else {
			postComment.HasPostVote = -1
		}

		// Serialize struct to json
		json, err := json.Marshal(&utilities.SocketResponse{
			Type: "update",
			Data: postComment,
		})

		// websocket
		origin := fmt.Sprintf("http://localhost:%s/", provider.GetConfig().Server.Port)
		url := fmt.Sprintf("ws://localhost:%s/ws/postComment", provider.GetConfig().Server.Port)

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

// update postVotes count in table postComments
func updatePostCommentPostVotes(postCommentID uint, postVote bool) (err error) {
	postComment := models.PostComment{}

	DB := provider.GetConnection()
	defer DB.Close()

	rows := DB.First(&postComment, postCommentID).RowsAffected

	if rows > 0 {
		if postVote {
			postComment.PostVotes++
		} else {
			postComment.PostVotes--
		}
		DB.Save(&postComment)
	} else {
		err = errors.New(fmt.Sprintf("No se encontró un registro de comentario para asignarle el voto"))
	}
	return
}

func postCommentGetChildren(id uint) (children []models.PostComment) {
	DB := provider.GetConnection()
	defer DB.Close()

	DB.Where("parent_id = ?", id).Find(&children)
	for i := range children {
		DB.Model(&children[i]).Related(&children[i].User)
		children[i].User[0].Password = ""
		children[i].User[0].TempKey = ""
		if children[i].ParentID >= 1 {
			children[i].Children = postCommentGetChildren(children[i].ID)
		}
	}
	return
}
