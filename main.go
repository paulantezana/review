package main

import (
	"crypto/sha256"
	"fmt"
    "github.com/paulantezana/review/models/monitoring"
    "os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/api"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Migration database
	migration()

	// Configuration cor
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// Assets
	static := e.Group("/static")
	static.Static("", "static")

	// Sting API services
	api.PublicApi(e)
	api.ProtectedApi(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = config.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}

// migration function
func migration() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.User{},
		&models.Module{},
		&models.Review{},
		&models.ReviewDetail{},
		&models.Program{},
		&models.Student{},
		&models.Teacher{},
		&models.Company{},
		&models.Setting{},

		// Migration monitoring
        &monitoring.Answer{},
        &monitoring.AnswerDetail{},
        &monitoring.MultipleQuestion{},
        &monitoring.Poll{},
        &monitoring.Question{},
        &monitoring.TypeQuestion{},
	)
	db.Model(&models.Review{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.ReviewDetail{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ReviewDetail{}).AddForeignKey("review_id", "reviews(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Student{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
    db.Model(&models.Student{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Teacher{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
    db.Model(&models.Teacher{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Module{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")


	// Monitoring
    db.Model(&monitoring.Poll{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
    db.Model(&monitoring.AnswerDetail{}).AddForeignKey("question_id", "questions(id)", "RESTRICT", "RESTRICT")
    db.Model(&monitoring.AnswerDetail{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")
    db.Model(&monitoring.AnswerDetail{}).AddForeignKey("answer_id", "answers(id)", "RESTRICT", "RESTRICT")

    db.Model(&monitoring.Question{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
    db.Model(&monitoring.Question{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")

    db.Model(&monitoring.MultipleQuestion{}).AddForeignKey("question_id", "questions(id)", "RESTRICT", "RESTRICT")

    // -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	usr := models.User{}
	db.First(&usr)

	// hash password
	cc := sha256.Sum256([]byte("FyYkbJo2W1T1"))
	pwd := fmt.Sprintf("%x", cc)

	// create model
	user := models.User{
		UserName: "sa",
		Password: pwd,
		Email:    "yoel.antezana@gmail.com",
		Profile:  "sa",
	}

	// insert database
	if usr.ID == 0 {
		db.Create(&user)
	}

	// First Setting
	prm := models.Setting{}
	db.First(&prm)
	co := models.Setting{
		Prefix:          "INSTITUTO DE EDUCACIÓN SUPERIOR TECNOLÓGICO PÚBLICO",
		PrefixShortName: "I.E.S.T.P.",
		Institute:       "SEDNA",
		Logo:            "static/logo.png",
		Ministry:        "static/ministry.jpg",
	}

	// Insert database
	if prm.ID == 0 {
		db.Create(&co)
	}
}
