package endpoint

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/controller"
    "github.com/paulantezana/review/controller/Institutioncontroller"
    "github.com/paulantezana/review/controller/admissioncontroller"
    "github.com/paulantezana/review/controller/appcontroller"
    "github.com/paulantezana/review/controller/institutecontroller"
	"github.com/paulantezana/review/controller/librarycontroller"
	"github.com/paulantezana/review/controller/messengercontroller"
	"github.com/paulantezana/review/controller/monitoringcontroller"
	"github.com/paulantezana/review/controller/reviewcontroller"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/utilities"
)

// CoreApi functions
func CoreApi(e *echo.Echo) {
    cor := e.Group("/api/v1/core")

    // Configure middleware with the custom claims type
    con := middleware.JWTConfig{
        Claims:     &utilities.CoreClaim{},
        SigningKey: []byte(provider.GetConfig().Server.Key),
    }
    cor.Use(middleware.JWTWithConfig(con))

    cor.POST("/institution/paginate", institutecontroller.GetInstitutions)
    cor.POST("/institution/by/id", institutecontroller.GetInstitutionByID)
    cor.POST("/institution/create", institutecontroller.CreateInstitution)
    cor.PUT("/institution/update", institutecontroller.UpdateInstitution)

    cor.POST("/app/by/id", appcontroller.GetAppByID)
    cor.PUT("/app/update", appcontroller.UpdateApp)

    cor.POST("/app/module/paginate", appcontroller.GetAppModules)
    cor.POST("/app/module/by/id", appcontroller.GetAppModuleByID)
    cor.POST("/app/module/create", appcontroller.CreateAppModule)
    cor.PUT("/app/module/update", appcontroller.UpdateAppModule)

    cor.POST("/app/function/paginate", appcontroller.GetAppModuleFunctions)
    cor.POST("/app/function/by/id", appcontroller.GetAppModuleFunctionByID)
    cor.POST("/app/function/create", appcontroller.CreateAppModuleFunction)
    cor.PUT("/app/function/update", appcontroller.UpdateAppModuleFunction)
}

// PublicApi function public urls
func PublicApi(e *echo.Echo) {
	pb := e.Group("/api/v1/public")

    pb.POST("/core/user/login", appcontroller.Login)

	pb.POST("/user/login", controller.Login)
	pb.POST("/user/login/student", controller.LoginStudent)
	pb.POST("/user/forgot/search", controller.ForgotSearch)
	pb.POST("/user/forgot/validate", controller.ForgotValidate)
	pb.POST("/user/forgot/change", controller.ForgotChange)
	pb.POST("/library/paginate", controller.ForgotChange)
	pb.POST("/library/by/id", controller.ForgotChange)

	// Global
	pb.GET("/setting", controller.GetSetting)

	pb.GET("/subsidiaries", institutecontroller.GetSubsidiaries)
	pb.GET("/subsidiaries/detail", admissioncontroller.GetSubsidiariesDetail)

	// Admission
	pb.POST("/admission/results", admissioncontroller.GetAdmissionExamAllResults)
	pb.POST("/admission/results/by/id", admissioncontroller.GetAdmissionExamResultsById)
	pb.POST("/admission/results/by/program/id", admissioncontroller.GetAdmissionExamResultsByProgramId)
	pb.POST("/admission/brochure", admissioncontroller.GetAdmissionExamAllResults) // Prospecto
	pb.POST("/admission/pre/admission/save", admissioncontroller.SavePreAdmission)
	pb.POST("/admission/pre/admission/get", admissioncontroller.GetPreAdmission)
	pb.POST("/admission/pre/admission/by/id", admissioncontroller.GetPreAdmissionById)
	pb.POST("/admission/pre/admission/programs", admissioncontroller.GetPreAdmissionPrograms)
	pb.POST("/admission/pre/admission/modalities", admissioncontroller.GetPreAdmissionModalities)

	// modalities
	pb.GET("/admission/modalities", admissioncontroller.GetModalities)
	pb.POST("/admission/modalities/by/id", admissioncontroller.GetModalityById)

	pb.POST("/external/dni", controller.GetStudentByDni)
	pb.POST("/external/ruc", controller.GetStudentByDni)


	// APP
    pb.POST("/app", appcontroller.GetAppByID)
    pb.POST("/app/module", appcontroller.GetAppModules)
}

// ProtectedApi function protected urls
func ProtectedApi(e *echo.Echo) {
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(provider.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// System
	ar.POST("/system/database/backup", controller.BackupDB)
	ar.POST("/system/database/backup/list", controller.BackupDBList)

	// Check Loginreport
	ar.POST("/login/user/check", controller.LoginUserCheck)
	ar.POST("/login/password/check", controller.LoginPasswordCheck)

	// Global settings
	ar.POST("/setting/global", Institutioncontroller.GetInstitutionSetting)
	ar.PUT("/setting/update", Institutioncontroller.UpdateInstitution)
	ar.POST("/setting/upload/logo", Institutioncontroller.UploadLogoInstitution)
	ar.GET("/setting/download/logo", Institutioncontroller.DownloadLogoInstitution)
	ar.POST("/setting/upload/ministry", Institutioncontroller.UploadMinistryInstitution)
	ar.POST("/setting/download/file", Institutioncontroller.DownloadFile)
	ar.GET("/setting/download/ministry", Institutioncontroller.DownloadMinistryInstitution)
	ar.GET("/setting/download/ministry/small", Institutioncontroller.DownloadMinistrySmallInstitution)
	ar.GET("/setting/download/national/emblem", Institutioncontroller.DownloadNationalEmblemInstitution)

	// ============================================================================
	//   Institutional controller api
	// subsidiary
	ar.POST("/subsidiary/all", institutecontroller.GetSubsidiaries)
	ar.POST("/subsidiary/all/tree", institutecontroller.GetSubsidiariesTree)
	ar.POST("/subsidiary/by/id", institutecontroller.GetSubsidiaryByID)
	ar.POST("/subsidiary/create", institutecontroller.CreateSubsidiary)
	ar.PUT("/subsidiary/update", institutecontroller.UpdateSubsidiary)
	ar.PUT("/subsidiary/update/main", institutecontroller.UpdateMainSubsidiary)
	ar.DELETE("/subsidiary/delete", institutecontroller.DeleteSubsidiary)

	// subsidiary - user
	ar.POST("/subsidiary/user/all/by/user", institutecontroller.GetSubsidiariesUserByUserID)
	ar.PUT("/subsidiary/user/update", institutecontroller.SubsidiariesUserUpdate)
	ar.POST("/subsidiary/user/all/by/user/license", institutecontroller.GetSubsidiariesUserByUserIDLicense)

	// Program
	ar.POST("/program/all", institutecontroller.GetPrograms)
	ar.POST("/program/all/by/license", institutecontroller.GetProgramsByLicense)
	ar.POST("/program/by/id", institutecontroller.GetProgramByID)
	ar.POST("/program/create", institutecontroller.CreateProgram)
	ar.PUT("/program/update", institutecontroller.UpdateProgram)

	// Program - user
	ar.POST("/program/user/all/by/user", institutecontroller.GetProgramsUserByUserID)
	ar.PUT("/program/user/update", institutecontroller.ProgramsUserUpdate)
	ar.POST("/program/user/all/by/student/license", institutecontroller.GetProgramsUserByStudentIDLicense)

	// Program
	ar.POST("/semester/all", institutecontroller.GetSemesters)
	ar.POST("/semester/by/id", institutecontroller.GetSemesterByID)
	ar.POST("/semester/create", institutecontroller.CreateSemester)
	ar.PUT("/semester/update", institutecontroller.UpdateSemester)
	ar.DELETE("/semester/delete", institutecontroller.DeleteSemester)

	// Student
	ar.POST("/student/paginate", institutecontroller.GetStudentsPaginate)
	ar.POST("/student/paginate/by/subsidiary", institutecontroller.GetStudentsPaginateBySubsidiary)
	ar.POST("/student/paginate/program", institutecontroller.GetStudentsPaginateByProgram)
	ar.POST("/student/paginate/by/license", institutecontroller.GetStudentsPaginateByLicense)
	ar.POST("/student/history", institutecontroller.GetStudentHistory)
	ar.POST("/student/programs", institutecontroller.GetStudentPrograms)
	ar.POST("/student/create", institutecontroller.CreateStudent)
	ar.PUT("/student/update", institutecontroller.UpdateStudent)
	ar.DELETE("/student/delete", institutecontroller.DeleteStudent)
	ar.POST("/student/search", institutecontroller.GetStudentSearch)
	ar.POST("/student/download/template/by/subsidiary", institutecontroller.GetTempUploadStudentBySubsidiary)
	ar.POST("/student/upload/template/by/subsidiary", institutecontroller.SetTempUploadStudentBySubsidiary)
	ar.POST("/student/download/template/by/program", institutecontroller.GetTempUploadStudentByProgram)
	ar.POST("/student/upload/template/by/program", institutecontroller.SetTempUploadStudentByProgram)
	ar.POST("/student/by/id", institutecontroller.GetStudentByID)
	ar.POST("/student/by/dni", institutecontroller.GetStudentByDNI)

	// Teacher
	ar.POST("/teacher/all", institutecontroller.GetTeachers)
	ar.POST("/teacher/paginate/program", institutecontroller.GetTeachersPaginateByProgram)
	ar.POST("/teacher/create", institutecontroller.CreateTeacher)
	ar.PUT("/teacher/update", institutecontroller.UpdateTeacher)
	ar.DELETE("/teacher/delete", institutecontroller.DeleteTeacher)
	ar.POST("/teacher/search", institutecontroller.GetTeacherSearch)
	ar.POST("/teacher/search/program", institutecontroller.GetTeacherSearchProgram)
	ar.POST("/teacher/download/template/by/subsidiary", institutecontroller.GetTempUploadTeacherBySubsidiary)
	ar.POST("/teacher/upload/template/by/subsidiary", institutecontroller.SetTempUploadTeacherBySubsidiary)
	ar.POST("/teacher/download/template/by/program", institutecontroller.GetTempUploadTeacherByProgram)
	ar.POST("/teacher/upload/template/by/program", institutecontroller.SetTempUploadTeacherByProgram)
	ar.GET("/teacher/export/all", institutecontroller.ExportAllTeachers)

	// Teacher Program
	ar.POST("/teacher/program/all", institutecontroller.GetTeacherProgramByProgram)
	ar.POST("/teacher/program/create", institutecontroller.CreateTeacherProgram)
	ar.PUT("/teacher/program/update", institutecontroller.UpdateTeacherProgram)
	ar.DELETE("/teacher/program/delete", institutecontroller.DeleteTeacherProgram)

	// Module
	ar.POST("/module/all", institutecontroller.GetModules)
	ar.POST("/module/create", institutecontroller.CreateModule)
	ar.PUT("/module/update", institutecontroller.UpdateModule)
	ar.DELETE("/module/delete", institutecontroller.DeleteModule)
	ar.POST("/module/search", institutecontroller.GetModuleSearch)

	// Course
	ar.POST("/course/paginate", institutecontroller.GetCoursesPaginate)
	ar.POST("/course/by/id", institutecontroller.GetCourseByID)
	ar.POST("/course/create", institutecontroller.CreateCourse)
	ar.PUT("/course/update", institutecontroller.UpdateCourse)
	ar.DELETE("/course/delete", institutecontroller.DeleteCourse)

	// Course level
	ar.POST("/course/level/paginate", institutecontroller.GetCourseLevelPaginate)
	ar.POST("/course/level/all", institutecontroller.GetAllCourseLevel)
	ar.POST("/course/level/create", institutecontroller.CreateCourseLevel)
	ar.PUT("/course/level/update", institutecontroller.UpdateCourseLevel)
	ar.DELETE("/course/level/delete", institutecontroller.DeleteCourseLevel)

	// Course level
	ar.POST("/course/node/paginate", institutecontroller.GetCourseNodePaginate)
	ar.POST("/course/node/create", institutecontroller.CreateCourseNode)
	ar.PUT("/course/node/update", institutecontroller.UpdateCourseNode)
	ar.DELETE("/course/node/delete", institutecontroller.DeleteCourseNode)

	// Company
	ar.POST("/company/all", reviewcontroller.GetCompanies)
	ar.POST("/company/create", reviewcontroller.CreateCompany)
	ar.PUT("/company/update", reviewcontroller.UpdateCompany)
	ar.DELETE("/company/delete", reviewcontroller.DeleteCompany)
	ar.DELETE("/company/delete/multiple", reviewcontroller.MultipleDeleteCompany)
	ar.POST("/company/search", reviewcontroller.GetCompanySearch)
	ar.GET("/company/download/template", reviewcontroller.GetTempUploadCompany)
	ar.POST("/company/upload/template", reviewcontroller.SetTempUploadCompany)
	ar.GET("/company/export/all", reviewcontroller.ExportAllCompanies)

	// Review
	ar.POST("/review/all", reviewcontroller.GetReviews)
	ar.POST("/review/create", reviewcontroller.CreateReview)
	ar.PUT("/review/update", reviewcontroller.UpdateReview)
	ar.DELETE("/review/delete", reviewcontroller.DeleteReview)

	ar.POST("/review/report/pdf/acta", reviewcontroller.GetPDFReviewStudentActa)
	ar.POST("/review/report/pdf/cons", reviewcontroller.GetPDFReviewStudentConst)
	ar.POST("/review/report/pdf/consolidate", reviewcontroller.GetPDFReviewStudentConsolidate)

	// User
	ar.POST("/user/all", controller.GetUsers)
	ar.POST("/user/search", controller.SearchUsers)
	ar.POST("/user/create", controller.CreateUser)
	ar.PUT("/user/update", controller.UpdateUser)
	ar.DELETE("/user/delete", controller.DeleteUser)
	ar.POST("/user/by/id", controller.GetUserByID)
	ar.POST("/user/upload/avatar", controller.UploadAvatarUser)
	ar.POST("/user/reset/password", controller.ResetPasswordUser)
	ar.POST("/user/change/password", controller.ChangePasswordUser)
	ar.POST("/user/licenses", controller.GetLicenseUser)
	ar.POST("/user/login/licenses", controller.GetLicensePostLogin)

	// Statistic
	ar.POST("/statistic/top/users", controller.TopUsers)
	ar.POST("/statistic/top/student/whit/reviews", controller.TopStudentsWithReview)

	// Review Detail
	ar.POST("/reviewDetail/by/review", reviewcontroller.GetReviewsDetailByReview)
	ar.DELETE("/reviewDetail/delete", reviewcontroller.DeleteReviewDetail)

	// ---------------------------------------------------------------------------
	//      Certificate routes ----------------------------------------------------

	//ar.POST("/course/all", coursescontroller.GetCoursesPaginate)
	//ar.POST("/course/create", coursescontroller.CreateCourse)
	//ar.PUT("/course/update", coursescontroller.UpdateCourse)
	//ar.DELETE("/course/delete", coursescontroller.DeleteCourse)
	//ar.POST("/course/by/id", coursescontroller.GetCourseByID)
	//
	//ar.POST("/course/student/all", coursescontroller.GetLanguageCourseStudentsPaginate)
	//ar.POST("/course/student/create", coursescontroller.CreateLanguageCourseStudent)
	//ar.PUT("/course/student/update", coursescontroller.UpdateLanguageCourseStudent)
	//ar.DELETE("/course/student/delete", coursescontroller.DeleteLanguageCourseStudent)
	//ar.POST("/course/student/act", coursescontroller.ActLanguageCourseStudent)
	//ar.POST("/course/student/download/template/by/subsidiary", coursescontroller.GetTempUploadLanguageCourseStudentBySubsidiary)
	//ar.POST("/course/student/upload/template/by/subsidiary", coursescontroller.SetTempUploadStudentBySubsidiary)

	ar.POST("/review/report/pdf/promotion/const", reviewcontroller.GetPDFReviewStudentConstGraduated)
	ar.POST("/review/report/pdf/promotion/certificate", reviewcontroller.GetPDFReviewStudentCertGraduated)
	ar.POST("/review/report/pdf/promotion/certificate/module", reviewcontroller.GetPDFReviewStudentCertModule)

	// ---------------------------------------------------------------------------
	//      Monitoring routes ----------------------------------------------------

	// poll
	ar.POST("/monitoring/poll/paginate", monitoringcontroller.GetPollsPaginate)
	ar.POST("/monitoring/poll/paginate/student", monitoringcontroller.GetPollsPaginateStudent)
	ar.POST("/monitoring/poll/by/id", monitoringcontroller.GetPollByID)
	ar.POST("/monitoring/poll/create", monitoringcontroller.CreatePoll)
	ar.PUT("/monitoring/poll/update", monitoringcontroller.UpdatePoll)
	ar.PUT("/monitoring/poll/update/state", monitoringcontroller.UpdateStatePoll)
	ar.DELETE("/monitoring/poll/delete", monitoringcontroller.DeletePoll)

	// quiz
	ar.POST("/monitoring/quiz/paginate", monitoringcontroller.GetQuizzesPaginate)
	ar.POST("/monitoring/quiz/all/by/diplomat", monitoringcontroller.GetQuizzesAllByDiplomat)
	ar.POST("/monitoring/quiz/all/by/diplomat/student", monitoringcontroller.GetQuizzesAllByDiplomatStudent)
	ar.POST("/monitoring/quiz/paginate/student", monitoringcontroller.GetQuizzesPaginateStudent)
	ar.POST("/monitoring/quiz/by/id", monitoringcontroller.GetQuizByID)
	ar.POST("/monitoring/quiz/create", monitoringcontroller.CreateQuiz)
	ar.PUT("/monitoring/quiz/update", monitoringcontroller.UpdateQuiz)
	ar.PUT("/monitoring/quiz/update/state", monitoringcontroller.UpdateStateQuiz)
	ar.DELETE("/monitoring/quiz/delete", monitoringcontroller.DeleteQuiz)

	// Question
	ar.POST("/monitoring/question/all", monitoringcontroller.GetQuestions)
	ar.POST("/monitoring/question/save", monitoringcontroller.SaveQuestions)
	ar.PUT("/monitoring/question/update", monitoringcontroller.UpdateQuestion)
	ar.DELETE("/monitoring/question/delete", monitoringcontroller.DeleteQuestion)

	// Quiz Diplomat
	ar.POST("/monitoring/quiz/diplomat/paginate", monitoringcontroller.GetQuizDiplomatPaginate)
	ar.POST("/monitoring/quiz/diplomat/paginate/student", monitoringcontroller.GetQuizDiplomatPaginateStudent)
	ar.POST("/monitoring/quiz/diplomat/by/id", monitoringcontroller.GetQuizDiplomatByID)
	ar.POST("/monitoring/quiz/diplomat/create", monitoringcontroller.CreateQuizDiplomat)
	ar.PUT("/monitoring/quiz/diplomat/update", monitoringcontroller.UpdateQuizDiplomat)
	//ar.PUT("/monitoring/quiz/diplomat/update/state", monitoringcontroller.UpdateQuizDiplomat)
	ar.DELETE("/monitoring/quiz/diplomat/delete", monitoringcontroller.DeleteQuizDiplomat)

	// Quiz Question
	ar.POST("/monitoring/quiz/question/all", monitoringcontroller.GetQuizQuestions)
	ar.POST("/monitoring/quiz/question/navigate", monitoringcontroller.GetQuizQuestionsNavigate)
	ar.POST("/monitoring/quiz/question/save", monitoringcontroller.SaveQuizQuestions)
	ar.PUT("/monitoring/quiz/question/update", monitoringcontroller.UpdateQuizQuestion)
	ar.DELETE("/monitoring/quiz/question/delete", monitoringcontroller.DeleteQuizQuestion)

	// Type questions
	ar.POST("/monitoring/type/question/all", monitoringcontroller.GetTypeQuestions)

	// Multiple question
	ar.DELETE("/monitoring/multiple/question/delete", monitoringcontroller.DeleteMultipleQuestion)

	// Answer
	ar.POST("/monitoring/answer/create", monitoringcontroller.CreateAnswer)
	ar.POST("/monitoring/answer/summary", monitoringcontroller.GetAnswerSummary)
	ar.POST("/monitoring/answer/navigate", monitoringcontroller.GetAnswerNavigate)
	ar.POST("/monitoring/answer/export/excel", monitoringcontroller.ExportExcelAnswers)

	// Quiz Answer
	ar.POST("/monitoring/quiz/answer/last", monitoringcontroller.GetLastQuizAnswer)
	ar.POST("/monitoring/quiz/answer/create", monitoringcontroller.CreateQuizAnswer)
	ar.POST("/monitoring/quiz/answer/time/finish", monitoringcontroller.TimeFinishQuizAnswer)
	ar.POST("/monitoring/quiz/answer/create/detail", monitoringcontroller.CreateQuizAnswerDetail)
	ar.POST("/monitoring/quiz/answer/analyze/by/student", monitoringcontroller.GetAnalyzeQuizAnswerByStudent)

	// Filters restrictions
	ar.POST("/monitoring/filter/query", monitoringcontroller.GetMonitoringFilterQuery)
	ar.POST("/monitoring/filter/search", monitoringcontroller.GetMonitoringFilterSearch)
	ar.POST("/monitoring/filter/save", monitoringcontroller.SaveMonitoringFilter)

	// ---------------------------------------------------------------------------
	//      Book routes ----------------------------------------------------
	// category
	ar.POST("/library/category/paginate", librarycontroller.GetCategoriesPaginate)
	ar.POST("/library/category/all", librarycontroller.GetCategoriesAll)
	ar.POST("/library/category/create", librarycontroller.CreatePostCategory)
	ar.PUT("/library/category/update", librarycontroller.UpdatePostCategory)
	ar.DELETE("/library/category/delete", librarycontroller.DeletePostCategory)
	ar.POST("/library/category/by/id", librarycontroller.GetPostCategoryByID)

	// book
	ar.POST("/library/book/paginate", librarycontroller.GetPostsPaginate)
	ar.POST("/library/book/paginate/reading", librarycontroller.GetPostsPaginateByPostReading)
	ar.POST("/library/book/like", librarycontroller.CreatePostLike)
	ar.POST("/library/book/create", librarycontroller.CreatePost)
	ar.PUT("/library/book/update", librarycontroller.UpdatePost)
	ar.DELETE("/library/book/delete", librarycontroller.DeletePost)
	ar.POST("/library/book/by/id", librarycontroller.GetPostByID)
	ar.POST("/library/book/by/id/reading", librarycontroller.GetPostByIDPostReading)
	ar.POST("/library/book/upload/avatar", librarycontroller.UploadAvatarPost)
	ar.POST("/library/book/upload/pdf", librarycontroller.UploadPdfPost)

	// Comments
	ar.POST("/library/comment/all", librarycontroller.GetPostCommentsAll)
	ar.POST("/library/comment/create", librarycontroller.CreatePostComment)
	ar.POST("/library/comment/vote", librarycontroller.CreatePostVote)
	ar.PUT("/library/comment/update", librarycontroller.UpdatePostComment)
	ar.DELETE("/library/comment/delete", librarycontroller.DeletePostComment)

	// Statics
	ar.POST("/library/statics/counts", librarycontroller.LibraryCounts)
	ar.POST("/library/statics/top/reading/by/student", librarycontroller.Top10ReadingByStudent)
	ar.POST("/library/statics/top/reading/by/program", librarycontroller.Top10ReadingByProgram)
	ar.POST("/library/statics/top/reading/by/book", librarycontroller.TopReadingByBook)
	ar.POST("/library/statics/last/comments", librarycontroller.LastComments)

	// ---------------------------------------------------------------------------
	//      Admission routes -----------------------------------------------------
	// Admission setting
	ar.POST("/admission/setting/all", admissioncontroller.GetAdmissionSettings)
	ar.POST("/admission/setting/by/id", admissioncontroller.GetAdmissionSettingByID)
	ar.POST("/admission/setting/create", admissioncontroller.CreateAdmissionSetting)
	ar.PUT("/admission/setting/update", admissioncontroller.UpdateAdmissionSetting)
	ar.DELETE("/admission/setting/delete", admissioncontroller.DeleteAdmissionSetting)

	// Admission setting
	ar.POST("/admission/modality/all", admissioncontroller.GetModalities)
	//ar.POST("/admission/modality/by/id", admissioncontroller.GetAdmissionSettingByID)
	ar.POST("/admission/modality/create", admissioncontroller.CreateModality)
	ar.PUT("/admission/modality/update", admissioncontroller.UpdateModality)
	ar.DELETE("/admission/modality/delete", admissioncontroller.DeleteModality)

	// Admission
	ar.POST("/admission/admission/paginate", admissioncontroller.GetAdmissionsPaginate)
	ar.POST("/admission/admission/by/id", admissioncontroller.GetAdmissionsByID)
	ar.POST("/admission/admission/paginate/exam", admissioncontroller.GetAdmissionsPaginateExam)
	ar.POST("/admission/admission/create", admissioncontroller.CreateAdmission)
	ar.POST("/admission/admission/update/student", admissioncontroller.UpdateStudentAdmission)
	ar.POST("/admission/admission/cancel", admissioncontroller.CancelAdmission)
	ar.PUT("/admission/admission/update", admissioncontroller.UpdateAdmission)
	ar.PUT("/admission/admission/update/exam", admissioncontroller.UpdateExamAdmission)
	ar.POST("/admission/admission/next/classroom", admissioncontroller.GetNextClassroomAdmission)

	ar.POST("/admission/pre/admission/paginate", admissioncontroller.GetPreAdmissionsPaginate)

	// Admission report PDF
	ar.POST("/admission/admission/report/pdf/file", admissioncontroller.GetPDFAdmissionStudentFile)
	ar.POST("/admission/admission/report/pdf/license", admissioncontroller.GetPDFAdmissionStudentLicense)
	ar.POST("/admission/admission/report/pdf/list", admissioncontroller.GetPDFAdmissionStudentList)

	// Admission export excel
	ar.POST("/admission/admission/export", admissioncontroller.ExportAdmission)
	ar.POST("/admission/admission/export/by/ids", admissioncontroller.ExportAdmissionByIds)
	ar.POST("/admission/admission/export/exam/results", admissioncontroller.ExportAdmissionExamResults)
	// Admission reports
	ar.POST("/admission/admission/report/general/by/settings", admissioncontroller.ReportAdmissionGeneral)

	// Payment
	ar.POST("/admission/payment/all", admissioncontroller.GetPayments)
	ar.POST("/admission/payment/create", admissioncontroller.CreatePayment)
	ar.PUT("/admission/payment/update", admissioncontroller.UpdatePayment)
	ar.DELETE("/admission/payment/delete", admissioncontroller.DeletePayment)

	// ---------------------------------------------------------------------------
	//      Messenger api -----------------------------------------------------
	ar.POST("/messenger/message/user/scroll", messengercontroller.GetUsersMssMessageScroll)
	ar.POST("/messenger/message/create", messengercontroller.CreateMssMessage)
	ar.POST("/messenger/message/create/by/group", messengercontroller.CreateMssGroupMessage)
	ar.POST("/messenger/message/create/upload/file", messengercontroller.CreateMssMessageFileUpload)
	ar.POST("/messenger/message/create/upload/file/by/group", messengercontroller.CreateMssMessageFileUploadByGroup)
	ar.POST("/messenger/message/by/user", messengercontroller.GetMssMessages)
	ar.POST("/messenger/message/by/group", messengercontroller.GetMssMessagesByGroup)
	ar.POST("/messenger/message/unread", messengercontroller.UnreadMssMessages)
	ar.POST("/messenger/group/scroll", messengercontroller.GetMssGroupsScroll)
	ar.POST("/messenger/group/by/id", messengercontroller.GetMssGroupByID)
	ar.POST("/messenger/group/create", messengercontroller.CreateMssGroup)
	ar.POST("/messenger/group/upload/avatar", messengercontroller.UploadAvatarMssGroup)
	ar.POST("/messenger/group/add/users", messengercontroller.AddUsers)
	ar.PUT("/messenger/group/update", messengercontroller.UpdateMssGroup)
	ar.POST("/messenger/group/is/active", messengercontroller.IsActiveMssGroup)
	ar.POST("/messenger/group/user/is/active", messengercontroller.IsActiveUserMssGroup)

	// ---------------------------------------------------------------------------
	//      External api -----------------------------------------------------
	ar.POST("/external/dni", controller.GetStudentByDni)
	ar.POST("/external/ruc", controller.GetStudentByDni)

	// ---------------------------------------------------------------------------
	//      Student routes -----------------------------------------------------
	//ar.POST("/setting/global/student", controller.GetStudentSettings)
}
