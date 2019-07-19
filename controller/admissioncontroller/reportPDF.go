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

type licenceADF struct {
    ID            uint      `json:"id"`
    Observation   string    `json:"observation"`
    Exonerated    bool      `json:"exonerated"`
    AdmissionDate time.Time `json:"admission_date"`
    Year          uint      `json:"year"`
    Classroom     uint      `json:"classroom"`
    Seat          uint      `json:"seat"`

    DNI      string `json:"dni"`
    FullName string `json:"full_name"`
    Avatar   string `json:"avatar"`
    Program  string `json:"program"`
}

func GetPDFAdmissionStudentLicense(c echo.Context) error {
    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Query all students
    admissionLicences := make([]licenceADF, 0)
    if err := DB.Table("admissions").
        Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.classroom, admissions.seat, "+
            "students.dni, students.full_name, programs.name as program, "+
            "users.avatar").
        Joins("INNER JOIN students ON admissions.student_id = students.id").
        Joins("INNER JOIN users ON students.user_id = users.id").
        Joins("INNER JOIN programs ON admissions.program_id = programs.id").
        Where("admissions.id IN (?)", request.IDs).
        Scan(&admissionLicences).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Find settings
    setting := models.Setting{}
    if err := DB.First(&setting).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Settings
    pageMargin := 12.7

    // Create PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetMargins(pageMargin,pageMargin,pageMargin)
    pdf.AddUTF8Font("Calibri", "", "static/font/Calibri_Regular.ttf")
    pdf.AddUTF8Font("Calibri", "B", "static/font/Calibri_Bold.ttf")
    pdf.AddUTF8Font("Calibri", "I", "static/font/Calibri_Italic.ttf")
    pdf.AddUTF8Font("Calibri", "BI", "static/font/Calibri_Bold_Italic.ttf")
    pdf.AddUTF8Font("Calibri", "L", "static/font/Calibri_Light.ttf")
    pdf.AddUTF8Font("Calibri", "LI", "static/font/Calibri_Light_Italic.ttf")

    // Settings
    leftMargin, topMargin, rightMargin, _ := pdf.GetMargins()
    pageWidth, _ := pdf.GetPageSize()
    pageWidth -= leftMargin + rightMargin
    fontFamilyName := "Calibri"
    gutter := 2.0

    // Init
    pdf.AddPage()
    pdf.SetFont(fontFamilyName, "B", 10)

    // Background
    pdf.Image("static/backgroundPattern1.jpg", 0, 0, 110, 110, false, "", 0, "")
    //pdf.Image("static/backgroundPattern1.jpg", 110, 0, 110, 110, false, "", 0, "")
    //pdf.Image("static/backgroundPattern1.jpg", 0, 110, 110, 110, false, "", 0, "")
    //pdf.Image("static/backgroundPattern1.jpg", 110, 110, 110, 110, false, "", 0, "")

    gCols := 3.0
    cCol := 0.0
    cRow := 0.0
    for _, license := range admissionLicences {
        headerString := fmt.Sprintf("%s \n %s - Sicuani",setting.Prefix, setting.Institute )
        headerString = strings.ToUpper(headerString)
        pdf.Image(setting.NationalEmblem, leftMargin, topMargin, 9, 0, false, "", 0, "")
        pdf.SetX(leftMargin + 9 + gutter)
        pdf.MultiCell(70,3.5,headerString,"","C",false)
        pdf.Image(setting.Logo, leftMargin + 79 + (gutter * 2), topMargin, 9, 0, false, "", 0, "")

        // Profile
        if utilities.FileExist(license.Avatar) {
            pdf.Image(license.Avatar, leftMargin, topMargin + 15, 25, 0, false, "", 0, "")
        }



        rW := (pageWidth / gCols) - (gutter * (gCols - 1))
        rH := 50.0

        sX :=  leftMargin
        if cCol > 0.0 {
            sX = (rW * cCol) + leftMargin + (gutter * cCol)
        }

        sY := topMargin
        if cRow == 0.0 {
            sY = topMargin
        }else {
            sY = (rH * cRow) + topMargin + gutter
        }


        pdf.Rect(sX,sY,rW,rH,"")

        // Set new params
        if cCol < (gCols - 1) {
            cCol += 1
        }else {
            cCol = 0.0
            cRow += 1
        }
    }

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