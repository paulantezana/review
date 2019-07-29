package monitoringcontroller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
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
	DB := provider.GetConnection()
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
	DB := provider.GetConnection()
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
	case "subsidiary":
		DB.Raw("SELECT id, name FROM subsidiaries "+
			"WHERE lower(name) LIKE lower(?) LIMIT 5", "%"+searchRequest.Search+"%").
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
	DB := provider.GetConnection()
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

func validateRestrictions(tableName string, user models.User) ([]uint, error) {
	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Query all filters by quiz_diplomats
	monitoringFilters := make([]models.MonitoringFilter, 0)
	if err := DB.Find(&monitoringFilters, models.MonitoringFilter{Table: tableName}).Error; err != nil {
		return nil, err
	}
	for k := range monitoringFilters {
		if err := DB.Model(&monitoringFilters[k]).Related(&monitoringFilters[k].Details).Error; err != nil {
			return nil, err
		}
	}

	// Query current student
	currentStudent := models.Student{}
	if err := DB.First(&currentStudent, models.Student{UserID: user.ID}).Error; err != nil {
		return nil, err
	}

	// loop query
	tableIDS := make([]uint, 0)
	for _, monitoringFilter := range monitoringFilters {
		switch monitoringFilter.Type {
		case "student":
			for _, detail := range monitoringFilter.Details {
				// validation
				if detail.ReferenceID == currentStudent.ID {
					tableIDS = append(tableIDS, monitoringFilter.TableID)
				}
			}
		case "program":
			// create new slice
			referIDS := make([]uint, 0)
			for _, detail := range monitoringFilter.Details {
				referIDS = append(referIDS, detail.ReferenceID)
			}
			// Query database programs
			studentPrograms := make([]models.StudentProgram, 0)
			DB.Where("id in (?)", referIDS).Find(&studentPrograms)

			// Loop validations
			for _, pro := range studentPrograms {
				// validation
				if pro.StudentID == currentStudent.ID {
					tableIDS = append(tableIDS, monitoringFilter.TableID)
				}
			}
		case "subsidiary":
			// Create new slices
			referIDS := make([]uint, 0)
			for _, detail := range monitoringFilter.Details {
				referIDS = append(referIDS, detail.ReferenceID)
			}

			// Query database subsidiaries
			sStudentIDS := make([]utilities.Counter, 0)
			DB.Raw("SELECT student_programs.student_id as id FROM student_programs "+
				"INNER JOIN programs ON student_programs.program_id = programs.id "+
				"INNER JOIN subsidiaries ON programs.subsidiary_id = subsidiaries.id "+
				"WHERE subsidiaries.id IN (?)", referIDS).Scan(&sStudentIDS)

			// Loop validations
			for _, sub := range sStudentIDS {
				// validation
				if sub.ID == currentStudent.ID {
					tableIDS = append(tableIDS, monitoringFilter.TableID)
				}
			}
			//default:
			//    innerIDS = append(innerIDS, monitoringFilter.TableID)
		}
	}

	// Remove duplicates
	tableIDS = utilities.RemoveDuplicates(tableIDS)

	// return ids
	return tableIDS, nil
}
