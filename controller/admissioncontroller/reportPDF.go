package admissioncontroller

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/utilities"
	"gopkg.in/ahmetb/go-linq.v3"
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
	DB := provider.GetConnection()
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

	// Query Subsidiary
	subsidiary := models.Subsidiary{}
	if err := DB.First(&subsidiary, models.Subsidiary{ID: request.SubsidiaryID}).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Settings
	pageMargin := 9.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddUTF8Font("Calibri", "", "static/font/Calibri_Regular.ttf")
	pdf.AddUTF8Font("Calibri", "B", "static/font/Calibri_Bold.ttf")
	pdf.AddUTF8Font("Calibri", "I", "static/font/Calibri_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "BI", "static/font/Calibri_Bold_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "L", "static/font/Calibri_Light.ttf")
	pdf.AddUTF8Font("Calibri", "LI", "static/font/Calibri_Light_Italic.ttf")

	// Settings
	leftMargin, topMargin, rightMargin, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	pageWidth -= leftMargin + rightMargin
	fontFamilyName := "Calibri"
	gutter := 2.0

	// Init
	pdf.AddPage()
	pdf.SetFont(fontFamilyName, "B", 10)

	gCols := 2.0
	cCol := 0.0
	cRow := 0.0
	padding := 4.0

	for _, license := range admissionLicences {
		w := (pageWidth - (gutter * (gCols - 1))) / gCols
		h := 56.0

		x := leftMargin + (w * cCol) + (gutter * cCol)
		y := topMargin + (h * cRow) + (gutter * cRow)

		limitRow := pageHeight > (y + h + topMargin)
		if !limitRow {
			pdf.AddPage()

			// Reset row col
			cRow = 0.0
			cCol = 0.0

			// Redefine X Y
			x = leftMargin + (w * cCol) + (gutter * cCol)
			y = topMargin + (h * cRow) + (gutter * cRow)
		}

		// Background
		if utilities.FileExist("static/backgroundMachupicchu.jpg") {
			pdf.Image("static/backgroundMachupicchu.jpg", x, y, w, h, false, "", 0, "")
		}

		// Header Images
		pdf.Image(setting.NationalEmblem, x+padding, y+padding, 9, 0, false, "", 0, "")
		pdf.Image(setting.Logo, (x+w)-(9+padding), y+4, 9, 0, false, "", 0, "")

		// Header content
		pdf.SetFont(fontFamilyName, "B", 10)
		pdf.SetXY(x+9+padding, y+padding)
		headerString := fmt.Sprintf("%s \n %s - %s CARNET DEL POSTULANTE", setting.Prefix, setting.Institute, subsidiary.District)
		headerString = strings.ToUpper(headerString)
		pdf.MultiCell(w-(18+(padding*2)), 3.5, headerString, "", "C", false)

		// Profile
		if utilities.FileExist(license.Avatar) {
			pdf.Image(license.Avatar, x+padding, y+15+padding, 25, 0, false, "", 0, "")
		}

		// Content
		pdf.SetXY(x+padding+27, y+padding+15)
		pdf.Cell(30, 3.5, "DNI:")
		pdf.Cell(30, 3.5, "ADMISION:")

		pdf.SetXY(x+padding+27, y+padding+18.5)
		pdf.SetFont(fontFamilyName, "", 10)
		pdf.Cell(30, 3.5, license.DNI)
		pdf.Cell(30, 3.5, fmt.Sprintf("%d", license.Year))

		// Name
		pdf.SetXY(x+padding+27, y+padding+25)
		pdf.SetFontStyle("B")
		pdf.Cell(30, 3.5, "APELLIDOS Y NOMBRES:")

		pdf.SetXY(x+padding+27, y+padding+28.5)
		pdf.SetFontStyle("")
		pdf.Cell(30, 3.5, license.FullName)

		// Program
		pdf.SetXY(x+padding+27, y+padding+35)
		pdf.SetFontStyle("B")
		pdf.Cell(30, 3.5, "PROGRAMA DE ESTUDIOS:")

		pdf.SetXY(x+padding+27, y+padding+38.5)
		pdf.SetFontStyle("")
		pdf.Cell(30, 3.5, license.Program)

		// Position
		pdf.SetXY(x+padding+27, y+padding+45)
		pdf.SetFontStyle("B")
		pdf.Cell(10, 3.5, "AULA: ")
		pdf.SetFontStyle("")
		pdf.Cell(15, 3.5, fmt.Sprintf("%d", license.Classroom))
		pdf.SetFontStyle("B")
		pdf.Cell(17, 3.5, "NUMERO: ")
		pdf.SetFontStyle("")
		pdf.Cell(15, 3.5, fmt.Sprintf("%d", license.Seat))

		// Rect
		pdf.SetFillColor(200, 200, 200)
		pdf.SetDrawColor(200, 200, 200)
		pdf.Rect(x, y, w, h, "")

		// Set new params
		if cCol < (gCols - 1) {
			cCol += 1
		} else {
			cCol = 0.0
			cRow += 1
		}
	}

	// Set file name
	cc := sha256.Sum256([]byte(fmt.Sprintf("%d-%d", request.IDs, currentUser.ID)))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/admission/%s.pdf", pwd)

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

type admissionList struct {
	ID            uint
	Observation   string
	Exonerated    bool
	AdmissionDate time.Time
	Year          uint
	Classroom     uint
	Seat          uint

	DNI      string
	FullName string
	Program  string
}

func GetPDFAdmissionStudentList(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	admission := models.Admission{}
	if err := c.Bind(&admission); err != nil {
		return err
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Find settings
	setting := models.Setting{}
	if err := db.First(&setting).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query all students
	admissionLists := make([]admissionList, 0)
	if err := db.Table("admissions").
		Select("admissions.id, admissions.observation, admissions.exonerated, admissions.admission_date, admissions.year, admissions.classroom, admissions.seat, "+
			"students.dni, students.full_name, programs.name as program ").
		Joins("INNER JOIN students ON admissions.student_id = students.id").
		Joins("INNER JOIN programs ON admissions.program_id = programs.id").
		Where("admissions.admission_setting_id = ? AND admissions.state = true", admission.AdmissionSettingID).
		Scan(&admissionLists).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Settings
	pageMargin := 19.0

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin+4, pageMargin)
	pdf.AddUTF8Font("Calibri", "", "static/font/Calibri_Regular.ttf")
	pdf.AddUTF8Font("Calibri", "B", "static/font/Calibri_Bold.ttf")
	pdf.AddUTF8Font("Calibri", "I", "static/font/Calibri_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "BI", "static/font/Calibri_Bold_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "L", "static/font/Calibri_Light.ttf")
	pdf.AddUTF8Font("Calibri", "LI", "static/font/Calibri_Light_Italic.ttf")

	// Settings
	leftMargin, topMargin, rightMargin, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	pageWidth -= leftMargin + rightMargin
	fontFamilyName := "Calibri"
	//gutter := 2.0

	// Header
	pdf.SetHeaderFunc(func() {
		clearTop := 12.7

		pdf.Image(setting.NationalEmblem, leftMargin, topMargin-clearTop, 12, 0, false, "", 0, "")
		pdf.Image(setting.Logo, (pageWidth+leftMargin)-12, topMargin-clearTop, 12, 0, false, "", 0, "")

		pdf.SetY(topMargin - clearTop)
		pdf.SetFont(fontFamilyName, "B", 13)
		pdf.WriteAligned(pageWidth, 5, strings.ToUpper(setting.Prefix), "C")
		pdf.Ln(5)

		pdf.SetFontSize(18)
		pdf.WriteAligned(pageWidth, 8, fmt.Sprintf("%s", strings.ToUpper(setting.Institute)), "C")
		pdf.Ln(9)

		pdf.SetFont(fontFamilyName, "", 7)
		pdf.WriteAligned(pageWidth, 2, fmt.Sprintf("Resolución de Creación %s - Resolución de Revalidación %s", setting.ResolutionAuthorization, setting.ResolutionRenovation), "C")
		pdf.Ln(2)

		pdf.SetLineWidth(0.3)
		pdf.Line(leftMargin, pdf.GetY()+2, pageWidth+leftMargin, pdf.GetY()+2)
		pdf.Ln(5)
	})

	// Add Page
	pdf.AddPage()

	// Group By Classroom
	aList := linq.From(admissionLists).GroupBy(
		func(i interface{}) interface{} {
			return i.(admissionList).Classroom
		},
		func(i interface{}) interface{} {
			return i.(admissionList)
		},
	)

	// Table
	next := aList.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		gr := item.(linq.Group)

		pdf.SetFont(fontFamilyName, "", 9)
		pdf.WriteAligned(pageWidth, 3, "AULA", "C")
		pdf.Ln(3)

		pdf.SetFont(fontFamilyName, "B", 18)
		pdf.WriteAligned(pageWidth, 8, fmt.Sprintf("%d", gr.Key), "C")
		pdf.Ln(8)

		// Table header
		pdf.SetFont(fontFamilyName, "B", 10)
		pdf.SetFillColor(230, 230, 230)
		pdf.CellFormat(10, 7, "Nº", "1", 0, "C", true, 0, "")
		pdf.CellFormat(20, 7, "DNI", "1", 0, "C", true, 0, "")
		pdf.CellFormat(80, 7, "APELLIDOS Y NOMBRES", "1", 0, "", true, 0, "")
		pdf.CellFormat(42, 7, "PROGRAMA", "1", 0, "", true, 0, "")
		pdf.CellFormat(10, 7, "AULA", "1", 0, "C", true, 0, "")
		pdf.CellFormat(10, 7, "Nº", "1", 0, "C", true, 0, "")
		pdf.Ln(-1)

		for i, value := range gr.Group {
			admission := value.(admissionList)

			// Table Body
			pdf.SetFont(fontFamilyName, "", 10)
			pdf.SetFillColor(255, 255, 255)
			pdf.CellFormat(10, 7, fmt.Sprintf("%d", i), "1", 0, "C", false, 0, "")
			pdf.CellFormat(20, 7, admission.DNI, "1", 0, "C", false, 0, "")
			pdf.CellFormat(80, 7, admission.FullName, "1", 0, "", false, 0, "")
			pdf.CellFormat(42, 7, admission.Program, "1", 0, "", false, 0, "")
			pdf.CellFormat(10, 7, fmt.Sprintf("%d", admission.Classroom), "1", 0, "C", false, 0, "")
			pdf.CellFormat(10, 7, fmt.Sprintf("%d", admission.Seat), "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		// Add New Page
		pdf.AddPage()
	}

	// Table
	nextTo := aList.Iterate()
	for item, ok := nextTo(); ok; item, ok = nextTo() {
		gr := item.(linq.Group)

		pdf.SetFont(fontFamilyName, "", 9)
		pdf.WriteAligned(pageWidth, 3, "AULA", "C")
		pdf.Ln(3)

		pdf.SetFont(fontFamilyName, "B", 18)
		pdf.WriteAligned(pageWidth, 8, fmt.Sprintf("%d", gr.Key), "C")
		pdf.Ln(8)

		gCols := 2.0
		cCol := 0.0
		cRow := 0.0
		padding := 2.5
		gutter := 2.0
		topMarginAux := topMargin + 24

		for _, value := range gr.Group {
			admission := value.(admissionList)

			w := (pageWidth - (gutter * (gCols - 1))) / gCols
			h := 30.0

			x := leftMargin + (w * cCol) + (gutter * cCol)
			y := topMarginAux + (h * cRow) + (gutter * cRow)

			limitRow := pageHeight > (y + h + topMargin)
			if !limitRow {
				pdf.AddPage()

				// Reset row col
				cRow = 0.0
				cCol = 0.0

				// Redefine X Y
				x = leftMargin + (w * cCol) + (gutter * cCol)
				y = topMarginAux + (h * cRow) + (gutter * cRow)
			}

			// Classroom and seat
			pdf.SetXY(x+padding, y+padding)
			pdf.SetFont(fontFamilyName, "B", 20)
			pdf.Cell(15, 8, fmt.Sprintf("%d - %d", admission.Classroom, admission.Seat))

			pdf.SetXY(x+padding, y+padding+8)
			pdf.SetFont(fontFamilyName, "", 10)
			pdf.Cell(15, 3.5, fmt.Sprintf("%s", admission.Program))
			pdf.SetXY(x+padding, y+padding+11.5)
			pdf.Cell(15, 3.5, fmt.Sprintf("%s", admission.DNI))

			// Line
			pdf.Line(x+padding, y+padding+17, x+(w-padding), y+padding+17)

			// Student Full Name
			pdf.SetXY(x+padding, y+padding+20)
			pdf.SetFont(fontFamilyName, "B", 11)
			pdf.Cell(15, 3.5, fmt.Sprintf("%s", admission.FullName))

			// Raw Rect
			pdf.Rect(x, y, w, h, "")

			// Set new params
			if cCol < (gCols - 1) {
				cCol += 1
			} else {
				cCol = 0.0
				cRow += 1
			}
		}
		// Add New Page
		pdf.AddPage()
	}

	// Set file name
	cc := sha256.Sum256([]byte(fmt.Sprintf("%d-%d", admission.ID, currentUser.ID)))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/admission/%s.pdf", pwd)

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

func GetPDFAdmissionStudentFile(c echo.Context) error {
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
	db := provider.GetConnection()
	defer db.Close()

	// Find settings
	setting := models.Setting{}
	if err := db.First(&setting).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	pageMargin := 19.0

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddUTF8Font("Calibri", "", "static/font/Calibri_Regular.ttf")
	pdf.AddUTF8Font("Calibri", "B", "static/font/Calibri_Bold.ttf")
	pdf.AddUTF8Font("Calibri", "I", "static/font/Calibri_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "BI", "static/font/Calibri_Bold_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "L", "static/font/Calibri_Light.ttf")
	pdf.AddUTF8Font("Calibri", "LI", "static/font/Calibri_Light_Italic.ttf")

	// Settings
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	pageWidth -= leftMargin + rightMargin

	// Header
	pdf.SetHeaderFunc(func() {
		//pdf.SetY(19)
		pdf.Image("static/ministrySmall.jpg", pageMargin, pageMargin-8, 70, 0, false, "", 0, "")
		pdf.Image(setting.Logo, (pageWidth+leftMargin)-12, pageMargin-8, 12, 0, false, "", 0, "")

		pdf.SetFont("Calibri", "B", 13)
		pdf.WriteAligned(pageWidth, 13, strings.ToUpper(setting.Prefix), "C")
		pdf.Ln(5)
		pdf.SetFont("Calibri", "B", 16)
		pdf.WriteAligned(pageWidth, 16, fmt.Sprintf("%s", strings.ToUpper(setting.Institute)), "C")
		pdf.Ln(5)

		pdf.SetFont("Calibri", "", 7)
		pdf.WriteAligned(pageWidth/2, 16, fmt.Sprintf("Resolución de Creación %s", setting.ResolutionAuthorization), "C")
		pdf.WriteAligned(pageWidth/2, 16, fmt.Sprintf("Resolución de Revalidación %s", setting.ResolutionRenovation), "C")
		pdf.Ln(5)

		pdf.SetLineWidth(0.3)
		pdf.Line(leftMargin, pdf.GetY()+5, pageWidth+leftMargin, pdf.GetY()+5)
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
	pdf.WriteAligned(pageWidth, 13, "FICHA DE INSCRIPCIÓN DEL POSTULANTE", "C")
	pdf.Ln(15)

	pdf.SetFont("Calibri", "", 10)
	pdf.CellFormat(15, 6, "AÑO", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 6, "2019", "1", 0, "C", false, 0, "")

	// Set file name
	cc := sha256.Sum256([]byte(fmt.Sprintf("%d-%d", request.IDs, currentUser.ID)))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/rpe/admission/%s.pdf", pwd)

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
