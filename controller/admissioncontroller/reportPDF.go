package admissioncontroller

import (
    "crypto/sha256"
    "fmt"
    "github.com/jung-kurt/gofpdf"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/utilities"
    "net/http"
    "strings"
    "time"
)

func GetPDFAdmissionStudentLicense(c echo.Context) error {
    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Settings
    pageMargin := 19.0

    // Create PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(pageMargin,pageMargin,pageMargin)
    pdf.AddPage()

    // Set file name
    cc := sha256.Sum256([]byte(time.Now().String()))
    pwd := fmt.Sprintf("%x", cc)
    fileName :=  fmt.Sprintf("static/rpe/%s.pdf",pwd)

    // Save file
    err := pdf.OutputFileAndClose(fileName)
    if err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data: fileName,
    })
}

func GetPDFAdmissionStudentList(c echo.Context) error {
    // Get data request
    admission := models.Admission{}
    if err := c.Bind(&admission); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Settings
    pageMargin := 19.0

    // Create PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(pageMargin,pageMargin,pageMargin)
    pdf.AddPage()

    // Set file name
    cc := sha256.Sum256([]byte(time.Now().String()))
    pwd := fmt.Sprintf("%x", cc)
    fileName :=  fmt.Sprintf("static/rpe/%s.pdf",pwd)

    // Save file
    err := pdf.OutputFileAndClose(fileName)
    if err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data: fileName,
    })
}

func GetPDFAdmissionStudentFile(c echo.Context) error {
    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Find settings
    setting := models.Setting{}
    if err := db.First(&setting).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    pageMargin := 19.0

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(pageMargin,pageMargin,pageMargin)
    pdf.AddUTF8Font("Calibri", "", "static/font/Calibri_Regular.ttf")
    pdf.AddUTF8Font("Calibri", "B", "static/font/Calibri_Bold.ttf")
    pdf.AddUTF8Font("Calibri", "I", "static/font/Calibri_Italic.ttf")
    pdf.AddUTF8Font("Calibri", "BI", "static/font/Calibri_Bold_Italic.ttf")
    pdf.AddUTF8Font("Calibri", "L", "static/font/Calibri_Light.ttf")
    pdf.AddUTF8Font("Calibri", "LI", "static/font/Calibri_Light_Italic.ttf")

    leftMargin, _, rightMargin, _ := pdf.GetMargins()
    pageWidth, _ := pdf.GetPageSize()
    pageWidth -= leftMargin + rightMargin

    // Header
    pdf.SetHeaderFunc(func() {
        //pdf.SetY(19)
        pdf.Image("static/ministrySmall.jpg", pageMargin, pageMargin - 8, 70, 0, false, "", 0, "")
        pdf.Image(setting.Logo, (pageWidth + leftMargin) - 12, pageMargin - 8, 12, 0, false, "", 0, "")

        pdf.SetFont("Calibri", "B", 13)
        pdf.WriteAligned(pageWidth,13, strings.ToUpper(setting.Prefix),"C")
        pdf.Ln(5)
        pdf.SetFont("Calibri", "B", 16)
        pdf.WriteAligned(pageWidth,16,fmt.Sprintf("%s",strings.ToUpper(setting.Institute)),"C")
        pdf.Ln(5)

        pdf.SetFont("Calibri", "", 7)
        pdf.WriteAligned(pageWidth / 2,16,fmt.Sprintf("Resolución de Creación %s",setting.ResolutionAuthorization),"C")
        pdf.WriteAligned(pageWidth / 2,16,fmt.Sprintf("Resolución de Revalidación %s",setting.ResolutionRenovation),"C")
        pdf.Ln(5)

        pdf.SetLineWidth(0.3)
        pdf.Line(leftMargin,pdf.GetY() + 5,pageWidth + leftMargin,pdf.GetY() + 5)
    })

    pdf.SetFooterFunc(func() {
        // Position at 1.5 cm from bottom
        pdf.SetY(-15)
        // Arial italic 8
        pdf.SetFont("Calibri", "I", 8)
        // Text color in gray
        pdf.SetTextColor(128, 128, 128)
        // Page number
        pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
            "", 0, "C", false, 0, "")
    })

    pdf.AddPage()

    pdf.SetY(pdf.GetY() + 10)
    pdf.SetFont("Calibri", "B", 13)
    pdf.WriteAligned(pageWidth, 13, "FICHA DE INSCRIPCIÓN DEL POSTULANTE","C")
    pdf.Ln(15)

    pdf.SetFont("Calibri", "", 10)
    pdf.CellFormat(15, 6, "AÑO", "1", 0, "C", false, 0, "")
    pdf.CellFormat(30, 6, "2019", "1", 0, "C", false, 0, "")


    // Set file name
    cc := sha256.Sum256([]byte(time.Now().String()))
    pwd := fmt.Sprintf("%x", cc)
    fileName :=  fmt.Sprintf("static/rpe/%s.pdf",pwd)

    // Save file
    err := pdf.OutputFileAndClose(fileName)
    if err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data: fileName,
    })
}