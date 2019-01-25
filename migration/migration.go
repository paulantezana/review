package migration

import (
	"crypto/sha256"
	"fmt"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
)

// migration function
func Migrate() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.Role{},
		&models.User{},

		&models.Setting{},

		// Institute
		&models.Subsidiary{},
		&models.SubsidiaryUser{},
		&models.Program{},
		&models.ProgramUser{},
		&models.Semester{},
		&models.Module{},
		&models.ModuleSemester{},
		&models.Unity{},

		&models.StudentStatus{},
		&models.Student{},
		&models.StudentHistory{},
		&models.StudentProgram{},

		&models.Teacher{},
		&models.TeacherAction{},
		&models.TeacherProgram{},

		// Admission
		&models.Admission{},
		&models.AdmissionPayment{},
		&models.Payment{},

		// Review
		&models.Review{},
		&models.ReviewDetail{},
		&models.Company{},

		// Migration certification
		&models.Course{},
		&models.CourseStudent{},
		&models.CourseExam{},

		// Migration monitoring
		&models.Answer{},
		&models.AnswerDetail{},
		&models.MultipleQuestion{},
		&models.Poll{},
		&models.Question{},
		&models.TypeQuestion{},

		// Libraries
		&models.Category{},
		&models.Book{},
		&models.Reading{},
		&models.Comment{},
		&models.Like{},
		&models.Vote{},

		// Messenger model
		&models.Group{},
		&models.Message{},
		&models.MessageRecipient{},
		&models.ReminderFrequency{},
		&models.Session{},
		&models.UserGroup{},
	)
	// General =================================================================
	db.Model(&models.User{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")

	// Institutional ===========================================================
	db.Model(&models.SubsidiaryUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.SubsidiaryUser{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "CASCADE", "CASCADE")
	db.Model(&models.Program{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ProgramUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.ProgramUser{}).AddForeignKey("program_id", "programs(id)", "CASCADE", "CASCADE")
	db.Model(&models.Semester{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ModuleSemester{}).AddForeignKey("semester_id", "semesters(id)", "CASCADE", "CASCADE")
	db.Model(&models.ModuleSemester{}).AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE")
	db.Model(&models.Unity{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Unity{}).AddForeignKey("semester_id", "semesters(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Student{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Student{}).AddForeignKey("student_status_id", "student_status(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.StudentHistory{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.StudentHistory{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.StudentProgram{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.StudentProgram{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Teacher{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.TeacherAction{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.TeacherProgram{}).AddForeignKey("teacher_id", "teachers(id)", "CASCADE", "CASCADE")
	db.Model(&models.TeacherProgram{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")

	// Admission ===============================================================
	db.Model(&models.Admission{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admission{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admission{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.AdmissionPayment{}).AddForeignKey("admission_id", "admissions(id)", "RESTRICT", "RESTRICT")
	//db.Model(&models.Payment{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "RESTRICT", "RESTRICT")

	// Reviews =================================================================
	db.Model(&models.Review{}).AddForeignKey("creator_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("student_program_id", "student_programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.ReviewDetail{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ReviewDetail{}).AddForeignKey("review_id", "reviews(id)", "RESTRICT", "RESTRICT")

	// Certification ===========================================================
	db.Model(&models.CourseStudent{}).AddForeignKey("course_id", "courses(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CourseStudent{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CourseExam{}).AddForeignKey("course_student_id", "course_students(id)", "RESTRICT", "RESTRICT")

	// Monitoring ==============================================================
	db.Model(&models.Poll{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Answer{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Answer{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.AnswerDetail{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")
	db.Model(&models.AnswerDetail{}).AddForeignKey("answer_id", "answers(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Question{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Question{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.MultipleQuestion{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")

	// Libraries ===========================================================
	db.Model(&models.Book{}).AddForeignKey("category_id", "categories(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Reading{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.Reading{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")
	db.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.Comment{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")
	db.Model(&models.Like{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.Like{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")

	// Messenger ===========================================================
	//db.Model(&models.MessageRecipient{}).AddForeignKey("recipient_id", "users(id)", "RESTRICT", "RESTRICT")
	//db.Model(&models.MessageRecipient{}).AddForeignKey("recipient_group_id", "user_groups(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MessageRecipient{}).AddForeignKey("message_id", "messages(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.UserGroup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.UserGroup{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Message{}).AddForeignKey("creator_id", "users(id)", "RESTRICT", "RESTRICT")
	//db.Model(&models.Message{}).AddForeignKey("reminder_frequency_id", "reminder_frequencies(id)", "RESTRICT", "RESTRICT")

	// -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	role := models.Role{}
	db.First(&role)

	// Validate
	if role.ID == 0 {
		role1 := models.Role{Name: "Director@"}     // Global Level
		role2 := models.Role{Name: "Administrador"} // Filial Level
		role3 := models.Role{Name: "Coordinador"}   // Program level
		role4 := models.Role{Name: "Profesor"}      // Teacher
		role5 := models.Role{Name: "Estudiante"}    // Student
		role6 := models.Role{Name: "Invitado"}      // Invited level
		db.Create(&role1).Create(&role2).Create(&role3)
		db.Create(&role4).Create(&role5).Create(&role6)
	}

	// -------------------------------------------------------------
	// Insert State students ---------------------------------------
	status := models.StudentStatus{}
	db.First(&status)
	if status.ID == 0 {
		status1 := models.StudentStatus{Name: "No asignado"}
		status2 := models.StudentStatus{Name: "Postulante"}
		status3 := models.StudentStatus{Name: "Exonerado"}
		status4 := models.StudentStatus{Name: "Trasladado"}
		status5 := models.StudentStatus{Name: "Rechazado"}
		status6 := models.StudentStatus{Name: "Aprobado"}
		status7 := models.StudentStatus{Name: "Prematriculado"}
		status8 := models.StudentStatus{Name: "Matriculado"}
		status9 := models.StudentStatus{Name: "Expulsado"}
		status10 := models.StudentStatus{Name: "Egresado"}
		db.Create(&status1).Create(&status2).Create(&status3).Create(&status4).Create(&status5)
		db.Create(&status6).Create(&status7).Create(&status8).Create(&status9).Create(&status10)
	}

	// -------------------------------------------------------------
	// Insert user -------------------------------------------------
	usr := models.User{}
	db.First(&usr)

	// Validate
	if usr.ID == 0 {
		// hash password
		cc := sha256.Sum256([]byte("sa"))
		pwd := fmt.Sprintf("%x", cc)

		// create model
		user := models.User{
			UserName: "sa",
			Password: pwd,
			Email:    "yoel.antezana@gmail.com",
			RoleID:   1,
			Freeze:   true,
		}
		db.Create(&user)
	}

	// =====================================================
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
			Logo:            "static/logo.jpg",
			Ministry:        "static/ministry.jpg",
		}
		// Insert in database
		db.Create(&co)
	}

	// ====================================================
	// -- Insert Type Quiestions
	tpq := models.TypeQuestion{}
	db.First(&tpq)

	if tpq.ID == 0 {
		// Create Models
		tq1 := models.TypeQuestion{Name: "Respuesta breve"}          // 1 = Simple input
		tq2 := models.TypeQuestion{Name: "Párrafo"}                  // 2 = TextArea input
		tq3 := models.TypeQuestion{Name: "Opción múltiple"}          // 3 = Radio input
		tq4 := models.TypeQuestion{Name: "Casillas de verificación"} // 4 = Checkbox input

		// Insert in Database
		db.Create(&tq1)
		db.Create(&tq2)
		db.Create(&tq3)
		db.Create(&tq4)
	}
}
