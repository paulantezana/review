package messengercontroller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

func GetGroupsScroll(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	counter := utilities.Counter{}

	// Query groups
	groups := make([]models.Group, 0)
	if err := DB.Debug().Raw("SELECT * FROM groups WHERE id IN "+
		"( SELECT group_id FROM user_groups WHERE user_id = ? ) "+
		"LIMIT ? OFFSET ?", currentUser.ID, request.Limit, offset).
		Scan(&groups).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	if err := DB.Raw("SELECT count(*) FROM groups WHERE id IN "+
		"( SELECT group_id FROM user_groups WHERE user_id = ? )", currentUser.ID).
		Scan(&counter).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Validate scroll
	var hasMore = false
	if request.CurrentPage < 10 {
		if request.Limit*request.CurrentPage < counter.Count {
			hasMore = true
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.ResponseScroll{
		Success:     true,
		Data:        groups,
		HasMore:     hasMore,
		CurrentPage: request.CurrentPage,
	})
}

func GetGroupByID(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Execute instructions
	if err := DB.First(&group, group.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// related gro users
	DB.Model(&group).Related(&group.UserGroups)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group,
	})
}

func CreateGroup(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}
	fmt.Println("PASO")

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// add data
	group.Date = time.Now()
	group.UserGroups = append(group.UserGroups, models.UserGroup{UserID: currentUser.ID})

	// Insert courses in database
	if err := DB.Create(&group).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", group.Name),
	})
}

func UpdateGroup(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	//
	group.Date = time.Now()

	// Insert courses in database
	rows := DB.Model(&group).Update(&group).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %s", group.Name),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", group.Name),
	})
}

func DisabelGroup(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	//
	group.Date = time.Now()

	// Disable
	if err := DB.Model(&group).UpdateColumn("is_active", false).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", group.Name),
	})
}
