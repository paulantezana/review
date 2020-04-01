package reviewcontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
	"net/http"
	"time"
)

func GetPDFReviewStudentActa(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	// Set file name
	cc := sha256.Sum256([]byte(time.Now().String()))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/%s.pdf", pwd)

	// Save file
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    fileName,
	})
}

func GetPDFReviewStudentConst(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	// Set file name
	cc := sha256.Sum256([]byte(time.Now().String()))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/%s.pdf", pwd)

	// Save file
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    fileName,
	})
}

func GetPDFReviewStudentConsolidate(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	// Set file name
	cc := sha256.Sum256([]byte(time.Now().String()))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/%s.pdf", pwd)

	// Save file
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    fileName,
	})
}

// Certificate
func GetPDFReviewStudentConstGraduated(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	// Set file name
	cc := sha256.Sum256([]byte(time.Now().String()))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/%s.pdf", pwd)

	// Save file
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    fileName,
	})
}

func GetPDFReviewStudentCertGraduated(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	// Set file name
	cc := sha256.Sum256([]byte(time.Now().String()))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/%s.pdf", pwd)

	// Save file
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    fileName,
	})
}

func GetPDFReviewStudentCertModule(c echo.Context) error {
	// Get data request
	review := models.Review{}
	if err := c.Bind(&review); err != nil {
		return err
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	// Set file name
	cc := sha256.Sum256([]byte(time.Now().String()))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/%s.pdf", pwd)

	// Save file
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    fileName,
	})
}
