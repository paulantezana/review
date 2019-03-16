package endpoint

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/controller"
	"github.com/paulantezana/review/controller/admissioncontroller"
	"github.com/paulantezana/review/controller/coursescontroller"
	"github.com/paulantezana/review/controller/institutecontroller"
	"github.com/paulantezana/review/controller/librarycontroller"
	"github.com/paulantezana/review/controller/messengercontroller"
	"github.com/paulantezana/review/controller/monitoringcontroller"
	"github.com/paulantezana/review/controller/publiccontroller"
	"github.com/paulantezana/review/controller/reviewcontroller"
	"github.com/paulantezana/review/utilities"
)

// PublicApi function public urls
func PublicApi(e *echo.Echo) {
	pb := e.Group("/api/v1/public")
	pb.POST("/user/login", controller.Login)
	pb.POST("/user/login/student", controller.LoginStudent)
	pb.POST("/user/forgot/search", controller.ForgotSearch)
	pb.POST("/user/forgot/validate", controller.ForgotValidate)
	pb.POST("/user/forgot/change", controller.ForgotChange)
	pb.POST("/library/paginate", controller.ForgotChange)
	pb.POST("/library/by/id", controller.ForgotChange)

	pb.GET("/subsidiaries", institutecontroller.GetSubsidiaries)
	pb.POST("/programs", institutecontroller.GetPrograms)
	pb.POST("/admission/exam/results", publiccontroller.GetAdmissionExamResults)
}

// ProtectedApi function protected urls
func ProtectedApi(e *echo.Echo) {
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(config.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// System
	ar.POST("/system/database/backup", controller.BackupDB)
	ar.POST("/system/database/backup/list", controller.BackupDBList)

	// Check Loginreport
	ar.POST("/login/user/check", controller.LoginUserCheck)
	ar.POST("/login/password/check", controller.LoginPasswordCheck)

	// Global settings
	ar.POST("/setting/global", controller.GetGlobalSettings)
	ar.PUT("/setting/update", controller.UpdateSetting)
	ar.POST("/setting/upload/logo", controller.UploadLogoSetting)
	ar.GET("/setting/download/logo", controller.DownloadLogoSetting)
	ar.POST("/setting/upload/ministry", controller.UploadMinistrySetting)
	ar.POST("/setting/download/file", controller.DownloadFile)
	ar.GET("/setting/download/ministry", controller.DownloadMinistrySetting)
	ar.GET("/setting/download/ministry/small", controller.DownloadMinistrySmallSetting)
	ar.GET("/setting/download/national/emblem", controller.DownloadNationalEmblemSetting)

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
	ar.POST("/subsidiary/user/all/by/user/license", institutecontroller.GetSubsidiariesUserByUserIDLicense)
	ar.PUT("/subsidiary/user/update", institutecontroller.UpdateSubsidiariesUserByUserID)

	// Program
	ar.POST("/program/all", institutecontroller.GetPrograms)
	ar.POST("/program/by/id", institutecontroller.GetProgramByID)
	ar.POST("/program/create", institutecontroller.CreateProgram)
	ar.PUT("/program/update", institutecontroller.UpdateProgram)

	// Program - user
	ar.POST("/program/user/all/by/user", institutecontroller.GetProgramsUserByUserID)
	ar.POST("/program/user/all/by/user/license", institutecontroller.GetProgramsUserByUserIDLicense)
	ar.POST("/program/user/all/by/student/license", institutecontroller.GetProgramsUserByStudentIDLicense)
	ar.PUT("/program/user/update", institutecontroller.UpdateProgramsUserByUserID)

	// Program
	ar.POST("/semester/all", institutecontroller.GetSemesters)
	ar.POST("/semester/create", institutecontroller.CreateSemester)
	ar.PUT("/semester/update", institutecontroller.UpdateSemester)
	ar.DELETE("/semester/delete", institutecontroller.DeleteSemester)

	// Student
	ar.POST("/student/paginate", institutecontroller.GetStudentsPaginate)
	ar.POST("/student/paginate/by/subsidiary", institutecontroller.GetStudentsPaginateBySubsidiary)
	ar.POST("/student/paginate/program", institutecontroller.GetStudentsPaginateByProgram)
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
	ar.POST("/review/acta", reviewcontroller.GetActaReview)
	ar.POST("/review/cons", reviewcontroller.GetConstReview)
	ar.POST("/review/consolidate", reviewcontroller.GetConsolidateReview)

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

	// Statistic
	ar.POST("/statistic/top/users", controller.TopUsers)
	ar.POST("/statistic/top/student/whit/reviews", controller.TopStudentsWithReview)

	// Review Detail
	ar.POST("/reviewDetail/by/review", reviewcontroller.GetReviewsDetailByReview)
	ar.DELETE("/reviewDetail/delete", reviewcontroller.DeleteReviewDetail)

	// ---------------------------------------------------------------------------
	//      Certificate routes ----------------------------------------------------
	ar.POST("/course/all", coursescontroller.GetCoursesPaginate)
	ar.POST("/course/create", coursescontroller.CreateCourse)
	ar.PUT("/course/update", coursescontroller.UpdateCourse)
	ar.DELETE("/course/delete", coursescontroller.DeleteCourse)
	ar.POST("/course/by/id", coursescontroller.GetCourseByID)

	ar.POST("/course/student/all", coursescontroller.GetCourseStudentsPaginate)
	ar.POST("/course/student/create", coursescontroller.CreateCourseStudent)
	ar.PUT("/course/student/update", coursescontroller.UpdateCourseStudent)
	ar.DELETE("/course/student/delete", coursescontroller.DeleteCourseStudent)
	ar.POST("/course/student/act", coursescontroller.ActCourseStudent)
	ar.POST("/course/student/download/template/by/subsidiary", coursescontroller.GetTempUploadCourseStudentBySubsidiary)
	ar.POST("/course/student/upload/template/by/subsidiary", coursescontroller.SetTempUploadStudentBySubsidiary)

	ar.POST("/review/promotion/const", reviewcontroller.GetConstGraduated)
	ar.POST("/review/promotion/certificate", reviewcontroller.GetCertGraduated)
	ar.POST("/review/promotion/certificate/module", reviewcontroller.GetCertModule)

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

	ar.POST("/monitoring/quiz/answer/last", monitoringcontroller.GetLastQuizAnswer)
	ar.POST("/monitoring/quiz/answer/create", monitoringcontroller.CreateQuizAnswer)
	ar.POST("/monitoring/quiz/answer/create/detail", monitoringcontroller.CreateQuizAnswerDetail)

	// ---------------------------------------------------------------------------
	//      Book routes ----------------------------------------------------
	// category
	ar.POST("/library/category/paginate", librarycontroller.GetCategoriesPaginate)
	ar.POST("/library/category/all", librarycontroller.GetCategoriesAll)
	ar.POST("/library/category/create", librarycontroller.CreateCategory)
	ar.PUT("/library/category/update", librarycontroller.UpdateCategory)
	ar.DELETE("/library/category/delete", librarycontroller.DeleteCategory)
	ar.POST("/library/category/by/id", librarycontroller.GetCategoryByID)

	// book
	ar.POST("/library/book/paginate", librarycontroller.GetBooksPaginate)
	ar.POST("/library/book/paginate/reading", librarycontroller.GetBooksPaginateByReading)
	ar.POST("/library/book/like", librarycontroller.CreateLike)
	ar.POST("/library/book/create", librarycontroller.CreateBook)
	ar.PUT("/library/book/update", librarycontroller.UpdateBook)
	ar.DELETE("/library/book/delete", librarycontroller.DeleteBook)
	ar.POST("/library/book/by/id", librarycontroller.GetBookByID)
	ar.POST("/library/book/by/id/reading", librarycontroller.GetBookByIDReading)
	ar.POST("/library/book/upload/avatar", librarycontroller.UploadAvatarBook)
	ar.POST("/library/book/upload/pdf", librarycontroller.UploadPdfBook)

	// Comments
	ar.POST("/library/comment/all", librarycontroller.GetCommentsAll)
	ar.POST("/library/comment/create", librarycontroller.CreateComment)
	ar.POST("/library/comment/vote", librarycontroller.CreateVote)
	ar.PUT("/library/comment/update", librarycontroller.UpdateComment)
	ar.DELETE("/library/comment/delete", librarycontroller.DeleteComment)

	// ---------------------------------------------------------------------------
	//      Admission routes -----------------------------------------------------
	// Admission setting
	ar.POST("/admission/setting/all", admissioncontroller.GetAdmissionSettings)
	ar.POST("/admission/setting/by/id", admissioncontroller.GetAdmissionSettingByID)
	ar.PUT("/admission/setting/default/in/web", admissioncontroller.ShowInWebAdmissionSetting)
	ar.POST("/admission/setting/create", admissioncontroller.CreateAdmissionSetting)
	ar.PUT("/admission/setting/update", admissioncontroller.UpdateAdmissionSetting)
	ar.DELETE("/admission/setting/delete", admissioncontroller.DeleteAdmissionSetting)

	// Admission
	ar.POST("/admission/admission/paginate", admissioncontroller.GetAdmissionsPaginate)
	ar.POST("/admission/admission/by/id", admissioncontroller.GetAdmissionsByID)
	ar.POST("/admission/admission/paginate/exam", admissioncontroller.GetAdmissionsPaginateExam)
	ar.POST("/admission/admission/create", admissioncontroller.CreateAdmission)
	ar.POST("/admission/admission/update/student", admissioncontroller.UpdateStudentAdmission)
	ar.POST("/admission/admission/cancel", admissioncontroller.CancelAdmission)
	ar.PUT("/admission/admission/update", admissioncontroller.UpdateAdmission)
	ar.PUT("/admission/admission/update/exam", admissioncontroller.UpdateExamAdmission)
	ar.POST("/admission/admission/df/file", admissioncontroller.FileAdmissionDF)
	ar.POST("/admission/admission/df/license", admissioncontroller.LicenseAdmissionDF)
	ar.POST("/admission/admission/df/list", admissioncontroller.ListAdmissionDF)
	ar.POST("/admission/admission/next/classroom", admissioncontroller.GetNextClassroomAdmission)
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
	ar.POST("/messenger/message/user/scroll", messengercontroller.GetUsersMessageScroll)
	ar.POST("/messenger/message/create", messengercontroller.CreateMessage)
	ar.POST("/messenger/message/create/by/group", messengercontroller.CreateGroupMessage)
	ar.POST("/messenger/message/create/upload/file", messengercontroller.CreateMessageFileUpload)
	ar.POST("/messenger/message/create/upload/file/by/group", messengercontroller.CreateMessageFileUploadByGroup)
	ar.POST("/messenger/message/by/user", messengercontroller.GetMessages)
	ar.POST("/messenger/message/by/group", messengercontroller.GetMessagesByGroup)
	ar.POST("/messenger/message/unread", messengercontroller.UnreadMessages)
	ar.POST("/messenger/group/scroll", messengercontroller.GetGroupsScroll)
	ar.POST("/messenger/group/by/id", messengercontroller.GetGroupByID)
	ar.POST("/messenger/group/create", messengercontroller.CreateGroup)
	ar.POST("/messenger/group/upload/avatar", messengercontroller.UploadAvatarGroup)
	ar.POST("/messenger/group/add/users", messengercontroller.AddUsers)
	ar.PUT("/messenger/group/update", messengercontroller.UpdateGroup)
	ar.POST("/messenger/group/is/active", messengercontroller.IsActiveGroup)
	ar.POST("/messenger/group/user/is/active", messengercontroller.IsActiveUserGroup)

	// ---------------------------------------------------------------------------
	//      External api -----------------------------------------------------
	// RENIEC
	ar.POST("/external/reniec", controller.Reniec)

	// ---------------------------------------------------------------------------
	//      Student routes -----------------------------------------------------
	ar.POST("/setting/global/student", controller.GetStudentSettings)
}
