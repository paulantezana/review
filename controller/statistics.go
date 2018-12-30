package controller

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/institutemodel"
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

type userTop struct {
	ID       uint   `json:"-"`
	UserName string `json:"user_name"`
	Top      uint   `json:"top"`
}

// TopUsers top all users
func TopUsers(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query database top 5
	userTops := make([]userTop, 0)
	if err := db.Table("users").
		Select("users.id, users.user_name, count(users.id) as top").
		Group("users.id, users.user_name").
		Order("top desc").
		Limit(15).
		Scan(&userTops).Error; err != nil {
		return c.JSON(http.StatusOK,utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Total registers
	var total uint
	db.Model(models.User{}).Count(&total)

	return c.JSON(http.StatusOK, utilities.ResponsePaginate{
		Success: true,
		Data:    userTops,
		Total:   total,
	})
}

type studentsWithReviewResponse struct {
	Message                   string `json:"message"`
	Success                   bool   `json:"success"`
	Students                  uint   `json:"students"`
	Reviews                   uint   `json:"reviews"`
	PercentagePositiveReviews uint   `json:"percentage_positive_reviews"`
	PercentageNegativeReviews uint   `json:"percentage_negative_reviews"`
}

// TopStudents top all count students
func TopStudentsWithReview(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// All modules
	var countModules uint
	if err := db.Model(&institutemodel.Module{}).Where("program_id = ?", currentUser.DefaultProgramID).Count(&countModules).Error; err != nil {
        return c.JSON(http.StatusOK,utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// All students
	var countStudents uint
	if err := db.Model(&institutemodel.Student{}).Where("program_id = ?", currentUser.DefaultProgramID).Count(&countStudents).Error; err != nil {
        return c.JSON(http.StatusOK,utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// All revisions
	countAllRevisions := countModules * countStudents
	var countReviews uint
	if err := db.Table("reviews").Joins("INNER JOIN students on reviews.student_id = students.id").
		Where("students.program_id = ?", currentUser.DefaultProgramID).Count(&countReviews).Error; err != nil {
        return c.JSON(http.StatusOK,utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	percentageP := uint(0)
	percentageN := uint(0)

	if countAllRevisions > 0 {
		percentageP = uint((countReviews * 100) / countAllRevisions)
		percentageN = uint(100 - percentageP)
	}

	return c.JSON(http.StatusOK, studentsWithReviewResponse{
		Success:                   true,
		Students:                  countStudents,
		Reviews:                   countReviews,
		PercentagePositiveReviews: percentageP,
		PercentageNegativeReviews: percentageN,
	})
}
