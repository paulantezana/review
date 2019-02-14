package reviewcontroller

import (
	"fmt"
	"github.com/paulantezana/review/models"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
)

func GetCompanies(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	offset := request.Validate()

	// Execute instructions
	var total uint
	companies := make([]models.Company, 0)

	// Query in database
	if err := db.Where("lower(name_social_reason) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&companies).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        companies,
		Total:       total,
		CurrentPage: request.CurrentPage,
		Limit:       request.Limit,
	})
}

func GetCompanySearch(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	companies := make([]models.Company, 0)
	if err := db.Where("lower(name_social_reason) LIKE lower(?)", "%"+request.Search+"%").
		Limit(10).Find(&companies).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	customCompanies := make([]models.Company, 0)
	for _, student := range companies {
		customCompanies = append(customCompanies, models.Company{
			ID:               student.ID,
			NameSocialReason: student.NameSocialReason,
			RUC:              student.RUC,
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customCompanies,
	})
}

func CreateCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert companies in database
	if err := db.Create(&company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("La empresa %s se registro correctamente", company.NameSocialReason),
	})
}

func UpdateCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update company in database
	rows := db.Model(&company).Update(company).RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", company.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("Los datos del la empresa %s se actualizaron correctamente", company.NameSocialReason),
	})
}

func DeleteCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Delete company in database
	if err := db.Delete(&company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("La empresa %s se elimino correctamente", company.NameSocialReason),
	})
}

func MultipleDeleteCompany(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	tx := db.Begin()
	for _, value := range request.IDs {
		company := models.Company{
			ID: value,
		}

		// Delete company in database
		if err := tx.Delete(&company).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Sel eliminaron %d registros", len(request.IDs)),
	})
}

// GetTempUploadStudent download template
func GetTempUploadCompany(c echo.Context) error {
	// Return file admin
	return c.File("templates/templateCompany.xlsx")
}

// SetTempUploadStudent set upload student
func SetTempUploadCompany(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	auxDir := "temp/" + file.Filename
	dst, err := os.Create(auxDir)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// ---------------------
	// Read File whit Excel
	// ---------------------
	excel, err := excelize.OpenFile(auxDir)
	if err != nil {
		return err
	}

	// Prepare
	companies := make([]models.Company, 0)
	ignoreCols := 5

	// Get all the rows in the student.
	rows := excel.GetRows("empresa")
	for k, row := range rows {
		if k >= ignoreCols {
			companies = append(companies, models.Company{
				RUC:              strings.TrimSpace(row[0]),
				NameSocialReason: strings.TrimSpace(row[1]),
				Address:          strings.TrimSpace(row[2]),
				Manager:          strings.TrimSpace(row[3]),
				Phone:            strings.TrimSpace(row[4]),
				CompanyType:      strings.TrimSpace(row[5]),
			})
		}
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert students in database
	tr := db.Begin()
	for _, company := range companies {
		if err := tr.Create(&company).Error; err != nil {
			tr.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{
				Message: fmt.Sprintf("Ocurrió un error al insertar al empresa %s con "+
					"RUC: %s es posible que esta empresa ya este en la base de datos o los datos son incorrectos, "+
					"Error: %s, no se realizo ninguna cambio en la base de datos", company.NameSocialReason, company.RUC, err),
			})
		}
	}
	tr.Commit()

	// Response success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se guardo %d registros den la base de datos", len(companies)),
	})
}

func ExportAllCompanies(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	companies := make([]models.Company, 0)

	// Query in database
	if err := db.Order("id asc").Find(&companies).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Create excel file
	excel := excelize.NewFile()

	// Create a new sheet.
	index := excel.NewSheet("Sheet1")

	// Set value of a cell.
	excel.SetCellValue("Sheet1", "A1", "RUC")
	excel.SetCellValue("Sheet1", "B1", "Nombre o razón social")
	excel.SetCellValue("Sheet1", "C1", "Dirección")
	excel.SetCellValue("Sheet1", "D1", "Gerente")
	excel.SetCellValue("Sheet1", "E1", "Telefono")
	excel.SetCellValue("Sheet1", "F1", "Tipo")

	currentRow := 2
	for k, company := range companies {
		excel.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow+k), company.RUC)
		excel.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow+k), company.NameSocialReason)
		excel.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow+k), company.Address)
		excel.SetCellValue("Sheet1", fmt.Sprintf("D%d", currentRow+k), company.Manager)
		excel.SetCellValue("Sheet1", fmt.Sprintf("E%d", currentRow+k), company.Phone)
		excel.SetCellValue("Sheet1", fmt.Sprintf("F%d", currentRow+k), company.CompanyType)
	}

	// Set active sheet of the workbook.
	excel.SetActiveSheet(index)

	// Save excel file by the given path.
	err := excel.SaveAs("temp/allCompanies.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	// Response file excel
	return c.File("temp/allCompanies.xlsx")
}
