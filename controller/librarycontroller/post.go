package librarycontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetPostsPaginate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	posts := make([]models.Post, 0)

	// Query in database
	if len(request.IDs) == 0 {
		if err := DB.Table("posts").Select("posts.* ").
			Joins("INNER JOIN categories ON posts.category_id = categories.id").
			Where("lower(posts.name) PostLike lower(?) AND categories.program_id = ?", "%"+request.Search+"%", request.ProgramID).
			Order("id desc").
			Offset(offset).Limit(request.Limit).Scan(&posts).
			Offset(-1).Limit(-1).Count(&total).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Where("lower(name) PostLike lower(?) AND category_id in (?)", "%"+request.Search+"%", request.IDs).
			Order("id desc").
			Offset(offset).Limit(request.Limit).Find(&posts).
			Offset(-1).Limit(-1).Count(&total).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Query post postComments count
	for i := range posts {
		DB.Model(&models.PostComment{}).
			Where("post_id = ?", posts[i].ID).
			Count(&posts[i].Detail.Comments)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        posts,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetPostsPaginateByPostReading(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	posts := make([]models.Post, 0)

	// Query in database
	if len(request.IDs) == 0 {
		if err := DB.Where("lower(name) PostLike lower(?)", "%"+request.Search+"%").
			Order("id desc").
			Offset(offset).Limit(request.Limit).Find(&posts).
			Offset(-1).Limit(-1).Count(&total).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Where("lower(name) PostLike lower(?) AND category_id in (?)", "%"+request.Search+"%", request.IDs).
			Order("id desc").
			Offset(offset).Limit(request.Limit).Find(&posts).
			Offset(-1).Limit(-1).Count(&total).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// validates
	for i := range posts {
		// Query post postComments count
		DB.Model(&models.PostComment{}).
			Where("post_id = ?", posts[i].ID).
			Count(&posts[i].Detail.Comments)

			// Average start
		bStarts := make([]models.BStarts, 0)
		if err := DB.Raw("SELECT users.user_name, postLikes.stars FROM postLikes "+
			"INNER JOIN users ON postLikes.user_id = users.id "+
			"WHERE postLikes.post_id = ? LIMIT 15", posts[i].ID).
			Scan(&bStarts).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		posts[i].Detail.Starts = bStarts

		// has postLike
		postLike := models.PostLike{
			UserID: currentUser.ID,
			PostID: posts[i].ID,
		}
		DB.Where(&postLike).First(&postLike)
		posts[i].Detail.StartValue = postLike.Stars
		if postLike.ID >= 1 {
			posts[i].Detail.HasStart = 1
		} else {
			posts[i].Detail.HasStart = 0
		}
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        posts,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetPostByID(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	post := models.Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	if err := DB.First(&post, post.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// PostComments count
	DB.Model(&models.PostComment{}).
		Where("post_id = ?", post.ID).
		Count(&post.Detail.Comments)

	// Query start
	bStarts := make([]models.BStarts, 0)
	if err := DB.Raw("SELECT users.user_name, postLikes.stars FROM postLikes "+
		"INNER JOIN users ON postLikes.user_id = users.id "+
		"WHERE postLikes.post_id = ? LIMIT 15", post.ID).
		Scan(&bStarts).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	post.Detail.Starts = bStarts

	// has postLike
	postLike := models.PostLike{
		UserID: currentUser.ID,
		PostID: post.ID,
	}
	DB.Where(&postLike).First(&postLike)
	post.Detail.StartValue = postLike.Stars
	if postLike.ID >= 1 {
		post.Detail.HasStart = 1
	} else {
		post.Detail.HasStart = 0
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    post,
	})
}

func GetPostByIDPostReading(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	post := models.Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	if err := DB.First(&post, post.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create postReadings
	postReading := models.PostReading{
		UserID: currentUser.ID,
		PostID: post.ID,
		Date:   time.Now(),
	}
	if err := DB.Create(&postReading).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Update table post
	post.Views++
	if err := DB.Save(&post).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// PostComments count
	DB.Model(&models.PostComment{}).
		Where("post_id = ?", post.ID).
		Count(&post.Detail.Comments)

	// Query start
	bStarts := make([]models.BStarts, 0)
	if err := DB.Raw("SELECT users.user_name, postLikes.stars FROM postLikes "+
		"INNER JOIN users ON postLikes.user_id = users.id "+
		"WHERE postLikes.post_id = ? LIMIT 15", post.ID).
		Scan(&bStarts).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	post.Detail.Starts = bStarts

	// has postLike
	postLike := models.PostLike{
		UserID: currentUser.ID,
		PostID: post.ID,
	}
	DB.Where(&postLike).First(&postLike)
	post.Detail.StartValue = postLike.Stars
	if postLike.ID >= 1 {
		post.Detail.HasStart = 1
	} else {
		post.Detail.HasStart = 0
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    post,
	})
}

func CreatePost(c echo.Context) error {
	// Get data request
	post := models.Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Insert posts in database
	if err := db.Create(&post).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    post.ID,
		Message: fmt.Sprintf("El libro %s se registro correctamente", post.Title),
	})
}

func UpdatePost(c echo.Context) error {
	// Get data request
	post := models.Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update post in database
	rows := db.Model(&post).Update(post).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", post.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    post.ID,
		Message: fmt.Sprintf("Los datos del libro %s se actualizaron correctamente", post.Title),
	})
}

func DeletePost(c echo.Context) error {
	// Get data request
	post := models.Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Delete post in database
	if err := db.Delete(&post).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    post.ID,
		Message: fmt.Sprintf("El libro %s se elimino correctamente", post.Title),
	})
}

// UploadAvatarUser function upload avatar user
func UploadAvatarPost(c echo.Context) error {
	// Read form fields
	idPost := c.FormValue("id")
	post := models.Post{}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&post, "id = ?", idPost).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", idPost),
		})
	}

	// Source
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	ccc := sha256.Sum256([]byte(string(post.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/posts/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	post.Avatar = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&post).Update(post).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    post.ID,
		Message: fmt.Sprintf("El avatar del libro %s, se subió correctamente", post.Title),
	})
}

// UploadAvatarUser function upload avatar user
func UploadPdfPost(c echo.Context) error {
	// Read form fields
	idPost := c.FormValue("id")
	post := models.Post{}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&post, "id = ?", idPost).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", idPost),
		})
	}

	// Source
	file, err := c.FormFile("pdf")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	ccc := sha256.Sum256([]byte(string(post.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/posts/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&post).Update(post).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    post.ID,
		Message: fmt.Sprintf("El avatar del libro %s, se subió correctamente", post.Title),
	})
}
