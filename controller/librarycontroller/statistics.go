package librarycontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

type libraryCount struct {
	Books      uint `json:"books"`
	Students   uint `json:"students"`
	Categories uint `json:"categories"`
	Readings   uint `json:"readings"`
}

func LibraryCounts(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Create struct
	counts := libraryCount{}

	// Count all books
	DB.Table("books").Count(&counts.Books)

	// Count all students
	DB.Table("students").Count(&counts.Students)

	// Count all categories
	DB.Table("categories").Count(&counts.Categories)

	// Count all readings
	DB.Table("readings").Count(&counts.Readings)

	// return data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    counts,
	})
}

type top10Reading struct {
	ID       uint   `json:"id"`
	DNI      string `json:"dni"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Count    uint   `json:"count"`
}

func Top10ReadingByStudent(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query count reading
	counter := make([]utilities.Counter, 0)
	if err := DB.Raw("SELECT user_id as id, count(book_id) as count FROM readings " +
		"GROUP BY user_id ORDER BY count DESC LIMIT 10").Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// other
	top10Readings := make([]top10Reading, 0)
	for _, cou := range counter {
		// Query student
		st := models.Student{}
		DB.First(&st, models.Student{UserID: cou.ID})

		// Query user
		us := models.User{}
		DB.First(&us, models.User{ID: cou.ID})

		// struct
		top := top10Reading{
			ID:       st.ID,
			DNI:      st.DNI,
			FullName: st.FullName,
			Avatar:   us.Avatar,
			Count:    cou.Count,
		}
		top10Readings = append(top10Readings, top)
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    top10Readings,
		Message: fmt.Sprintf("Voto actualizado"),
	})
}

func Top10ReadingByProgram(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	//
	counter := make([]utilities.Counter, 0)
	DB.Raw("SELECT user_id as id, count(book_id) FROM readings " +
		"GROUP BY user_id LIMIT 10").Scan(&counter)

	// other
	top10Readings := make([]top10Reading, 0)
	for _, cou := range counter {
		// Query student
		st := models.Student{}
		DB.First(&st, models.Student{ID: cou.ID})

		// Query user
		us := models.User{}
		DB.First(&us, models.User{ID: st.UserID})

		// struct
		top := top10Reading{
			ID:       st.ID,
			DNI:      st.DNI,
			FullName: st.FullName,
			Avatar:   us.Avatar,
			Count:    cou.Count,
		}
		top10Readings = append(top10Readings, top)
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    top10Readings,
		Message: fmt.Sprintf("Voto actualizado"),
	})
}

type topReadingByBookResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Count    uint   `json:"count"`
}

func TopReadingByBook(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Readings
	readings := make([]topReadingByBookResponse, 0)
	DB.Raw("SELECT books.id, books.name, count(readings.user_id) FROM readings " +
		"INNER JOIN books ON readings.book_id = books.id " +
		"GROUP BY books.id " +
		"LIMIT 20").Scan(&readings)

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    readings,
		Message: fmt.Sprintf("Voto actualizado"),
	})
}

type lastCommentsResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Votes     uint32    `json:"votes"`
	HasVote   int8      `json:"has_vote"` // if current user has vote
	Body      string    `json:"body"`
	User      struct {
		ID       uint   `json:"id"`
		UserName string `json:"user_name"`
		Avatar   string `json:"avatar"`
	} `json:"user"`
}

func LastComments(c echo.Context) error {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query last comments
	comments := make([]models.PostComment, 0)
	if err := DB.Limit(10).Order("created_at DESC").Find(&comments).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// for comments
	lastCommentsResponses := make([]lastCommentsResponse, 0)
	for _, comment := range comments {
		// Query user
		user := models.User{}
		DB.First(&user, models.User{ID: comment.UserID})

		lastC := lastCommentsResponse{
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Votes:     comment.PostVotes,
			HasVote:   comment.HasPostVote,
		}
		lastC.User.ID = user.ID
		lastC.User.UserName = user.UserName
		lastC.User.Avatar = user.Avatar
		lastCommentsResponses = append(lastCommentsResponses, lastC)
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    lastCommentsResponses,
	})
}
