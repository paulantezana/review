package migration

import (
	"crypto/sha256"
	"fmt"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/models/admissionmodel"
	"github.com/paulantezana/review/models/coursemodel"
	"github.com/paulantezana/review/models/institutemodel"
	"github.com/paulantezana/review/models/librarymodel"
	"github.com/paulantezana/review/models/messengermodel"
	"github.com/paulantezana/review/models/monitoringmodel"
	"github.com/paulantezana/review/models/reviewmodel"
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
		&institutemodel.Subsidiary{},
		&institutemodel.SubsidiaryUser{},
		&institutemodel.Program{},
		&institutemodel.ProgramUser{},
		&institutemodel.Semester{},
		&institutemodel.Module{},
		&institutemodel.ModuleSemester{},
		&institutemodel.Unity{},

		&institutemodel.StudentStatus{},
		&institutemodel.Student{},
		&institutemodel.StudentHistory{},
		&institutemodel.StudentProgram{},

		&institutemodel.Teacher{},
		&institutemodel.TeacherAction{},
		&institutemodel.TeacherProgram{},

		// Admission
		&admissionmodel.Admission{},
		&admissionmodel.AdmissionPayment{},
		&admissionmodel.Payment{},

		// Review
		&reviewmodel.Review{},
		&reviewmodel.ReviewDetail{},
		&reviewmodel.Company{},

		// Migration certification
		&coursemodel.Course{},
		&coursemodel.CourseStudent{},
		&coursemodel.CourseExam{},

		// Migration monitoring
		&monitoringmodel.Answer{},
		&monitoringmodel.AnswerDetail{},
		&monitoringmodel.MultipleQuestion{},
		&monitoringmodel.Poll{},
		&monitoringmodel.Question{},
		&monitoringmodel.TypeQuestion{},

		// Libraries
		&librarymodel.Category{},
		&librarymodel.Book{},
		&librarymodel.Reading{},
		&librarymodel.Comment{},
		&librarymodel.Like{},

		// Messenger model
		&messengermodel.Group{},
		&messengermodel.Message{},
		&messengermodel.MessageRecipient{},
		&messengermodel.ReminderFrequency{},
		&messengermodel.Session{},
		&messengermodel.UserGroup{},
	)
	// General =================================================================
	db.Model(&models.User{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")

	// Institutional ===========================================================
	db.Model(&institutemodel.SubsidiaryUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.SubsidiaryUser{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.Program{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.ProgramUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.ProgramUser{}).AddForeignKey("program_id", "programs(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.Semester{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.ModuleSemester{}).AddForeignKey("semester_id", "semesters(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.ModuleSemester{}).AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.Unity{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.Unity{}).AddForeignKey("semester_id", "semesters(id)", "RESTRICT", "RESTRICT")

	db.Model(&institutemodel.Student{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.Student{}).AddForeignKey("student_status_id", "student_status(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.StudentHistory{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.StudentHistory{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.StudentProgram{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.StudentProgram{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")

	db.Model(&institutemodel.Teacher{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.TeacherAction{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")
	db.Model(&institutemodel.TeacherProgram{}).AddForeignKey("teacher_id", "teachers(id)", "CASCADE", "CASCADE")
	db.Model(&institutemodel.TeacherProgram{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")

	// Admission ===============================================================
	db.Model(&admissionmodel.Admission{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&admissionmodel.Admission{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&admissionmodel.Admission{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&admissionmodel.AdmissionPayment{}).AddForeignKey("admission_id", "admissions(id)", "RESTRICT", "RESTRICT")
	//db.Model(&admissionmodel.Payment{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "RESTRICT", "RESTRICT")

	// Reviews =================================================================
	db.Model(&reviewmodel.Review{}).AddForeignKey("creator_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&reviewmodel.Review{}).AddForeignKey("student_program_id", "student_programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&reviewmodel.Review{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")
	db.Model(&reviewmodel.Review{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")

	db.Model(&reviewmodel.ReviewDetail{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")
	db.Model(&reviewmodel.ReviewDetail{}).AddForeignKey("review_id", "reviews(id)", "RESTRICT", "RESTRICT")

	// Certification ===========================================================
	db.Model(&coursemodel.CourseStudent{}).AddForeignKey("course_id", "courses(id)", "RESTRICT", "RESTRICT")
	db.Model(&coursemodel.CourseStudent{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&coursemodel.CourseExam{}).AddForeignKey("course_student_id", "course_students(id)", "RESTRICT", "RESTRICT")

	// Monitoring ==============================================================
	db.Model(&monitoringmodel.Poll{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoringmodel.Answer{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoringmodel.Answer{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoringmodel.AnswerDetail{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")
	db.Model(&monitoringmodel.AnswerDetail{}).AddForeignKey("answer_id", "answers(id)", "RESTRICT", "RESTRICT")

	db.Model(&monitoringmodel.Question{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&monitoringmodel.Question{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")

	db.Model(&monitoringmodel.MultipleQuestion{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")

	// Libraries ===========================================================
	db.Model(&librarymodel.Book{}).AddForeignKey("category_id", "categories(id)", "RESTRICT", "RESTRICT")
	db.Model(&librarymodel.Reading{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&librarymodel.Reading{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")
	db.Model(&librarymodel.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&librarymodel.Comment{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")
	db.Model(&librarymodel.Like{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&librarymodel.Like{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")

	// Messenger ===========================================================
	db.Model(&messengermodel.MessageRecipient{}).AddForeignKey("recipient_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&messengermodel.MessageRecipient{}).AddForeignKey("recipient_group_id", "user_groups(id)", "RESTRICT", "RESTRICT")
	db.Model(&messengermodel.MessageRecipient{}).AddForeignKey("message_id", "messages(id)", "RESTRICT", "RESTRICT")
	db.Model(&messengermodel.UserGroup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&messengermodel.UserGroup{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	db.Model(&messengermodel.Message{}).AddForeignKey("creator_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&messengermodel.Message{}).AddForeignKey("reminder_frequency_id", "reminder_frequencies(id)", "RESTRICT", "RESTRICT")

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
	status := institutemodel.StudentStatus{}
	db.First(&status)
	if status.ID == 0 {
		status1 := institutemodel.StudentStatus{Name: "No asignado"}
		status2 := institutemodel.StudentStatus{Name: "Postulante"}
		status3 := institutemodel.StudentStatus{Name: "Exonerado"}
		status4 := institutemodel.StudentStatus{Name: "Trasladado"}
		status5 := institutemodel.StudentStatus{Name: "Rechazado"}
		status6 := institutemodel.StudentStatus{Name: "Aprobado"}
		status7 := institutemodel.StudentStatus{Name: "Prematriculado"}
		status8 := institutemodel.StudentStatus{Name: "Matriculado"}
		status9 := institutemodel.StudentStatus{Name: "Expulsado"}
		status10 := institutemodel.StudentStatus{Name: "Egresado"}
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
	tpq := monitoringmodel.TypeQuestion{}
	db.First(&tpq)

	if tpq.ID == 0 {
		// Create Models
		tq1 := monitoringmodel.TypeQuestion{Name: "Respuesta breve"}          // 1 = Simple input
		tq2 := monitoringmodel.TypeQuestion{Name: "Párrafo"}                  // 2 = TextArea input
		tq3 := monitoringmodel.TypeQuestion{Name: "Opción múltiple"}          // 3 = Radio input
		tq4 := monitoringmodel.TypeQuestion{Name: "Casillas de verificación"} // 4 = Checkbox input

		// Insert in Database
		db.Create(&tq1)
		db.Create(&tq2)
		db.Create(&tq3)
		db.Create(&tq4)
	}
}
