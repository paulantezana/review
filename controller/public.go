package controller

import (
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

// GetGlobalSettings function
func GetSetting(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Find settings
	con := models.Setting{}
	db.First(&con)

	// Set object response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "OK",
		Data:    con,
	})
}
