package migration

import (
    "crypto/sha256"
    "fmt"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/provider"
    "time"
)

// migration function
func Migrate() {
	db := provider.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.App{},
		&models.AppModule{},
		&models.AppModuleFunction{},
		&models.AppUser{},

		// Institutions
		&models.Institution{},
		&models.InstitutionModule{},

		// Authorization
		&models.User{},
		&models.UserSubsidiary{},
		&models.UserProgram{},
		&models.UserRole{},
		&models.UserRoleModule{},
		&models.UserRoleFunction{},
		&models.UserSession{},
		&models.UserAuthorizationType{},

		// Institute
		&models.Subsidiary{},
		&models.Program{},
		&models.Semester{},
		&models.Module{},
		&models.ModuleSemester{},
		&models.Course{},
		&models.CourseLevel{},
		&models.CourseNode{},

		&models.StudentStatus{},
		&models.Student{},
		&models.StudentHistory{},
		&models.StudentProgram{},

		&models.Teacher{},
		&models.TeacherAction{},
		&models.TeacherProgram{},

		// Admission
		&models.AdmissionSetting{},
		&models.Admission{},
		&models.PreAdmission{},
		&models.AdmissionModality{},
		&models.AdmissionPayment{},
		&models.Payment{},

		// Review
		&models.Review{},
		&models.ReviewDetail{},
		&models.Company{},

		// Migration certification
		//&models.Course{},
		//&models.CourseStudent{},
		//&models.CourseExam{},

		// Migration monitoring
		&models.TypeQuestion{},
		&models.Poll{},
		&models.Question{},
		&models.MultipleQuestion{},
		&models.Answer{},
		&models.AnswerDetail{},
		&models.QuizDiplomat{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.MultipleQuizQuestion{},
		&models.QuizAnswer{},
		&models.QuizAnswerDetail{},
		&models.MonitoringFilter{},
		&models.MonitoringFilterDetail{},

		// Libraries
		&models.Post{},
		&models.PostType{},
		&models.PostFile{},
		&models.PostCategory{},
		&models.PostReading{},
		&models.PostComment{},
		&models.PostLike{},
		&models.PostVote{},

		// Messenger model
		&models.MssGroup{},
		&models.MssMessage{},
		&models.MssMessageRecipient{},
		&models.MssReminderFrequency{},
		&models.MssUserGroup{},
		&models.MssGroupMessage{},
		&models.MssGroupMessageRecipient{},
	)
	// Institutional

	// General =================================================================
	db.Model(&models.User{}).AddForeignKey("user_role_id", "user_roles(id)", "RESTRICT", "RESTRICT")

	// Authorization ===========================================================
	db.Model(&models.UserSubsidiary{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.UserSubsidiary{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "CASCADE", "CASCADE")
	db.Model(&models.UserProgram{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.UserProgram{}).AddForeignKey("program_id", "programs(id)", "CASCADE", "CASCADE")
	db.Model(&models.UserProgram{}).AddForeignKey("user_subsidiary_id", "user_subsidiaries(id)", "CASCADE", "CASCADE")

	// Institutional ===========================================================
	db.Model(&models.Program{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Semester{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ModuleSemester{}).AddForeignKey("semester_id", "semesters(id)", "CASCADE", "CASCADE")
	db.Model(&models.ModuleSemester{}).AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE")
	db.Model(&models.Course{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Course{}).AddForeignKey("semester_id", "semesters(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CourseLevel{}).AddForeignKey("course_id", "courses(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CourseNode{}).AddForeignKey("course_level_id", "course_levels(id)", "RESTRICT", "RESTRICT")

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
	db.Model(&models.AdmissionSetting{}).AddForeignKey("subsidiary_id", "subsidiaries(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admission{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admission{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admission{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admission{}).AddForeignKey("admission_setting_id", "admission_settings(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PreAdmission{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PreAdmission{}).AddForeignKey("admission_setting_id", "admission_settings(id)", "RESTRICT", "RESTRICT")
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
	//db.Model(&models.CourseStudent{}).AddForeignKey("course_id", "courses(id)", "RESTRICT", "RESTRICT")
	//db.Model(&models.CourseStudent{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	//db.Model(&models.CourseExam{}).AddForeignKey("course_student_id", "course_students(id)", "RESTRICT", "RESTRICT")

	// Monitoring ==============================================================
	db.Model(&models.Poll{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Question{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Question{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MultipleQuestion{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")
	db.Model(&models.Answer{}).AddForeignKey("poll_id", "polls(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Answer{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.AnswerDetail{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "RESTRICT")
	db.Model(&models.AnswerDetail{}).AddForeignKey("answer_id", "answers(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MonitoringFilterDetail{}).AddForeignKey("monitoring_filter_id", "monitoring_filters(id)", "CASCADE", "CASCADE")

	// Quiz diplomat
	db.Model(&models.QuizDiplomat{}).AddForeignKey("program_id", "programs(id)", "RESTRICT", "RESTRICT")

	// Quiz
	db.Model(&models.QuizQuestion{}).AddForeignKey("quiz_id", "quizzes(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.QuizQuestion{}).AddForeignKey("type_question_id", "type_questions(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MultipleQuizQuestion{}).AddForeignKey("quiz_question_id", "multiple_quiz_questions(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.QuizAnswer{}).AddForeignKey("quiz_id", "quizzes(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.QuizAnswer{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.QuizAnswerDetail{}).AddForeignKey("quiz_question_id", "quiz_questions(id)", "CASCADE", "RESTRICT")
	db.Model(&models.QuizAnswerDetail{}).AddForeignKey("quiz_answer_id", "quiz_answers(id)", "RESTRICT", "RESTRICT")

	// Libraries ===========================================================
	db.Model(&models.Post{}).AddForeignKey("post_category_id", "post_categories(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Post{}).AddForeignKey("post_type_id", "post_types(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.PostReading{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.PostReading{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
	db.Model(&models.PostComment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.PostComment{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
	db.Model(&models.PostLike{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.PostLike{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")

	// Messenger =================================s==========================
	db.Model(&models.MssMessageRecipient{}).AddForeignKey("recipient_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MssMessageRecipient{}).AddForeignKey("mss_message_id", "mss_messages(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MssGroupMessageRecipient{}).AddForeignKey("mss_group_message_id", "mss_group_messages(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MssUserGroup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MssUserGroup{}).AddForeignKey("mss_group_id", "mss_groups(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MssMessage{}).AddForeignKey("creator_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.MssGroupMessage{}).AddForeignKey("creator_id", "users(id)", "RESTRICT", "RESTRICT")
	//db.Model(&models.Message{}).AddForeignKey("reminder_frequency_id", "reminder_frequencies(id)", "RESTRICT", "RESTRICT")

	// -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	role := models.UserRole{}
	db.First(&role)

	// Validate
	if role.ID == 0 {
		role1 := models.UserRole{Name: "Director@", IsMain: true, State: true}     // Global Level
		role2 := models.UserRole{Name: "Administrador", IsMain: true, State: true} // Filial Level
		role3 := models.UserRole{Name: "Coordinador", IsMain: true, State: true}   // Program level
		role4 := models.UserRole{Name: "Profesor", IsMain: true, State: true}      // Teacher
		role5 := models.UserRole{Name: "Estudiante", IsMain: true, State: true}    // Student
		role6 := models.UserRole{Name: "Invitado", IsMain: true, State: true}      // Invited level
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
    // Insert State Post Types ---------------------------------------
    postType := models.PostType{}
    db.First(&postType)
    if postType.ID == 0 {
        postType1 := models.PostType{Name: "Libro"}
        db.Create(&postType1)
    }

    // -------------------------------------------------------------
    // Insert State Post Types ---------------------------------------
    appUser := models.AppUser{}
    db.First(&appUser)
    if appUser.ID == 0 {
        cc := sha256.Sum256([]byte("admin"))
        pwd := fmt.Sprintf("%x", cc)
        user := models.AppUser{
            UserName: "admin",
            Password: pwd,
            State:true,
            Avatar: "static/apps/logo.png",
        }
        db.Create(&user)
    }

	// ====================================================
	// -- Insert Type Quiestions
	tpq := models.TypeQuestion{}
	db.First(&tpq)

	if tpq.ID == 0 {
		// Create Models
		tq1 := models.TypeQuestion{Name: "Respuesta breve"}   // 1 = Simple input
		tq2 := models.TypeQuestion{Name: "Párrafo"}           // 2 = TextArea input
		tq3 := models.TypeQuestion{Name: "Una respuesta"}     // 3 = Radio input
		tq4 := models.TypeQuestion{Name: "Varias respuestas"} // 4 = Checkbox input

		// Insert in Database
		db.Create(&tq1).Create(&tq2).Create(&tq3).Create(&tq4)
	}

    // ====================================================
    // -- Insert App
    appData := models.App{}
    db.First(&appData)

    if appData.ID == 0 {
        appData1 := models.App{
            Name: "Respuesta breve",
            Version: "0.0.1",
            LastUpdate: time.Now(),
        }
        db.Create(&appData1)
    }
}
