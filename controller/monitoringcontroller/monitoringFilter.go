package monitoringcontroller

import (
    "fmt"
    "github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"net/http"
)

func GetMonitoringFilterQuery(c echo.Context) error {
	// Get data request
	monitoringFilter := models.MonitoringFilter{}
	if err := c.Bind(&monitoringFilter); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Query
	DB.FirstOrCreate(&monitoringFilter, models.MonitoringFilter{
		Table:   monitoringFilter.Table,
		TableID: monitoringFilter.TableID,
	})

	// Query details
    if err := DB.Model(&monitoringFilter).Related(&monitoringFilter.Details).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    monitoringFilter,
	})
}

type searchRequest struct {
	Type   string `json:"type"`
	Search string `json:"search"`
}
type searchResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetMonitoringFilterSearch(c echo.Context) error {
	// Get data request
	searchRequest := searchRequest{}
	if err := c.Bind(&searchRequest); err != nil {
		return err
	}

	// Get connection
	DB := config.GetConnection()
	defer DB.Close()

	// Search
	searchResponse := make([]searchResponse, 0)
	switch searchRequest.Type {
	case "program":
		DB.Raw("SELECT id, name FROM programs "+
			"WHERE lower(name) LIKE lower(?) LIMIT 5", "%"+searchRequest.Search+"%").
			Scan(&searchResponse)
	case "student":
		DB.Raw("SELECT id, full_name as name FROM students "+
			"WHERE lower(full_name) LIKE lower(?) OR dni LIKE ? LIMIT 5", "%"+searchRequest.Search+"%", "%"+searchRequest.Search+"%").
			Scan(&searchResponse)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    searchResponse,
	})
}

// SaveMonitoringFilter -
func SaveMonitoringFilter(c echo.Context) error {
    // Get data request
    monitoringFilter := models.MonitoringFilter{}
    if err := c.Bind(&monitoringFilter); err != nil {
        return err
    }

    // Get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Reset all
    if err := DB.Where("monitoring_filter_id = ?", monitoringFilter.ID).Delete(&models.MonitoringFilterDetail{}).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Save
    if err := DB.Save(&monitoringFilter).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    monitoringFilter.ID,
        Message: fmt.Sprintf("El filtro %d se registro correctamente", monitoringFilter.ID),
    })
}