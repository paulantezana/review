package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
)

type userTop struct {
	ID       uint   `json:"-"`
	UserName string `json:"user_name"`
	Top      uint   `json:"top"`
}

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
		return err
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