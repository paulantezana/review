package migration

import (
	"crypto/sha256"
	"fmt"
    "github.com/paulantezana/review/models/coursemodel"

    "github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/monitoring"
)

// migration function
func Migrate() {
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

		// Migration certification
		&coursemodel.Course{},
		&coursemodel.CourseStudent{},
		&coursemodel.CourseExam{},

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

	// Certification
    db.Model(&coursemodel.CourseStudent{}).AddForeignKey("course_id", "courses(id)", "RESTRICT", "RESTRICT")
    db.Model(&coursemodel.CourseStudent{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
    db.Model(&coursemodel.CourseExam{}).AddForeignKey("course_student_id", "course_students(id)", "RESTRICT", "RESTRICT")

	// Monitoring
	db.Model(&monitoring.Poll{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoring.AnswerDetail{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")
	db.Model(&monitoring.AnswerDetail{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoring.AnswerDetail{}).AddForeignKey("answer_id", "answers(id)", "RESTRICT", "RESTRICT")

	db.Model(&monitoring.Question{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoring.Question{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")

	db.Model(&monitoring.MultipleQuestion{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")

	// -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	usr := models.User{}
	db.First(&usr)

	// Validate
	if usr.ID == 0 {
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
		db.Create(&user)
	}

	// First Setting
	prm := models.Setting{}
	db.First(&prm)

	// Validate
	if prm.ID == 0 {
		co := models.Setting{
			Prefix:          "INSTITUTO DE EDUCACIÓN SUPERIOR TECNOLÓGICO PÚBLICO",
			PrefixShortName: "I.E.S.T.P.",
			Institute:       "SEDNA",
			NationalEmblem:  "static/nationalEmblem.jpg",
			Logo:            "static/logo.png",
			Ministry:        "static/ministry.jpg",
		}
		// Insert in database
		db.Create(&co)
	}

	// ====================================================
	// -- Insert Type Quiestions
	tpq := monitoring.TypeQuestion{}
	db.First(&tpq)

	if tpq.ID == 0 {
		// Create Models
		tq1 := monitoring.TypeQuestion{Name: "Respuesta breve"}          // 1 = Simple input
		tq2 := monitoring.TypeQuestion{Name: "Párrafo"}                  // 2 = TextArea input
		tq3 := monitoring.TypeQuestion{Name: "Opción múltiple"}          // 3 = Radio input
		tq4 := monitoring.TypeQuestion{Name: "Casillas de verificación"} // 4 = Checkbox input

		// Insert in Database
		db.Create(&tq1)
		db.Create(&tq2)
		db.Create(&tq3)
		db.Create(&tq4)
	}
}
