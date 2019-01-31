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
		"ORDER BY id asc LIMIT ? OFFSET ? ", currentUser.ID, request.Limit, offset).
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

type userGroupResponse struct {
    ID       uint      `json:"id" gorm:"primary_key"`
    Date     time.Time `json:"date"`
    IsActive bool      `json:"is_active" gorm:"default:'true'"`
    IsAdmin bool `json:"is_admin"`

    UserID  uint `json:"user_id"`
    Name string `json:"name"`
    Avatar string `json:"avatar"`
}

type groupResponse struct {
    ID       uint      `json:"id"`
    Name     string    `json:"name"`
    Avatar   string    `json:"avatar"`
    Date     time.Time `json:"date"`
    IsActive bool      `json:"is_active"`

    Users []userGroupResponse `json:"users"`
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
    groupResponse := groupResponse{}
	if err := DB.Raw("SELECT * FROM groups WHERE id = ?", group.ID).Scan(&groupResponse).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

    userGroupResponses := make([]userGroupResponse,0)
    if err := DB.Table("user_groups").
        Select("user_groups.id, user_groups.date, user_groups.is_active, user_groups.is_admin, user_groups.user_id, users.user_name as name, users.avatar").
        Joins("INNER JOIN users ON user_groups.user_id = users.id").
        Where("user_groups.group_id = ?", group.ID).
        Order("user_groups.id asc").
        Scan(&userGroupResponses).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }
    groupResponse.Users = userGroupResponses

    // Query current student Name
    for i := range userGroupResponses {
        student := models.Student{}
        DB.First(&student, models.Student{UserID: userGroupResponses[i].UserID})
        if student.ID >= 1 {
            userGroupResponses[i].Name = student.FullName
        } else {
            teacher := models.Teacher{}
            DB.First(&teacher, models.Teacher{UserID: userGroupResponses[i].UserID})
            if teacher.ID >= 1 {
                userGroupResponses[i].Name = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
            }
        }

    }

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    groupResponse,
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
	group.UserGroups = append(group.UserGroups, models.UserGroup{
	    UserID: currentUser.ID,
	    IsAdmin: true,
	})

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

func IsActiveGroup(c echo.Context) error {
	// Get data request
	group := models.Group{}
	if err := c.Bind(&group); err != nil {
		return err
	}

	// get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Disable
	if err := DB.Model(&group).UpdateColumn("is_active", group.IsActive).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    group.ID,
		Message: fmt.Sprintf("El curso %s se registro correctamente", group.Name),
	})
}

func IsActiveUserGroup(c echo.Context) error {
    // Get data request
    userGroup := models.UserGroup{}
    if err := c.Bind(&userGroup); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Disable
    if err := DB.Debug().Model(&userGroup).Where("user_id = ? AND group_id = ?",userGroup.UserID,userGroup.GroupID).
        UpdateColumn("is_active", userGroup.IsActive).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    userGroup.ID,
        Message: fmt.Sprintf("El cursos se registro correctamente"),
    })
}
