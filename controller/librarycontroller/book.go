package librarycontroller

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
    "time"
)

func GetBooksPaginate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	books := make([]models.Book, 0)

	// Query in database
	if err := DB.Where("lower(name) LIKE lower(?) AND category_id in (?)", "%"+request.Search+"%",request.IDs).
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&books).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

    // Query book comments count
    for i := range books {
        DB.Model(&models.Comment{}).
            Where("book_id = ?", books[i].ID).
            Count(&books[i].CommentCount)
    }

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        books,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetBookByID(c echo.Context) error {
	// Get data request
	book := models.Book{}
	if err := c.Bind(&book); err != nil {
		return err
	}

	// Get connection
    DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	if err := DB.First(&book, book.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

    // Comments count
    DB.Model(&models.Comment{}).
        Where("book_id = ?", book.ID).
        Count(&book.CommentCount)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    book,
	})
}

func GetBookByIDReading(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    book := models.Book{}
    if err := c.Bind(&book); err != nil {
        return err
    }

    // Get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Execute instructions
    if err := DB.First(&book, book.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Comments count
    DB.Model(&models.Comment{}).
        Where("book_id = ?", book.ID).
        Count(&book.CommentCount)

    // Create readings
    reading := models.Reading{
        UserID: currentUser.ID,
        BookID: book.ID,
        Date: time.Now(),
    }
    if err := DB.Create(&reading).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Update table book
    book.Views++
    if err := DB.Save(&book).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    book,
    })
}

func CreateBook(c echo.Context) error {
	// Get data request
	book := models.Book{}
	if err := c.Bind(&book); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert books in database
	if err := db.Create(&book).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    book.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", book.Name),
	})
}

func UpdateBook(c echo.Context) error {
	// Get data request
	book := models.Book{}
	if err := c.Bind(&book); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update book in database
	rows := db.Model(&book).Update(book).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", book.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    book.ID,
		Message: fmt.Sprintf("Los datos del curso %s se actualizaron correctamente", book.Name),
	})
}

func DeleteBook(c echo.Context) error {
	// Get data request
	book := models.Book{}
	if err := c.Bind(&book); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete book in database
	if err := db.Delete(&book).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    book.ID,
		Message: fmt.Sprintf("El curso %s se elimino correctamente", book.Name),
	})
}

// UploadAvatarUser function upload avatar user
func UploadAvatarBook(c echo.Context) error {
	// Read form fields
	idBook := c.FormValue("id")
	book := models.Book{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&book, "id = ?", idBook).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", idBook),
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
	ccc := sha256.Sum256([]byte(string(book.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/books/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	book.Avatar = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&book).Update(book).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    book.ID,
		Message: fmt.Sprintf("El avatar del libro %s, se subió correctamente", book.Name),
	})
}

// UploadAvatarUser function upload avatar user
func UploadPdfBook(c echo.Context) error {
	// Read form fields
	idBook := c.FormValue("id")
	book := models.Book{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&book, "id = ?", idBook).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", idBook),
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
	ccc := sha256.Sum256([]byte(string(book.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/books/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	book.Pdf = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&book).Update(book).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    book.ID,
		Message: fmt.Sprintf("El avatar del libro %s, se subió correctamente", book.Name),
	})
}
