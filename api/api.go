package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/controller"
	"github.com/paulantezana/review/controller/admissioncontroller"
	"github.com/paulantezana/review/controller/coursescontroller"
	"github.com/paulantezana/review/controller/institutecontroller"
	"github.com/paulantezana/review/controller/librarycontroller"
	"github.com/paulantezana/review/controller/monitoringcontroller"
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

	// Check Login
	ar.POST("/login/check", controller.LoginCheck)

	// Global settings
	ar.POST("/setting/global", controller.GetGlobalSettings)
	ar.PUT("/setting/update", controller.UpdateSetting)
	ar.POST("/setting/upload/logo", controller.UploadLogoSetting)
	ar.GET("/setting/download/logo", controller.DownloadLogoSetting)
	ar.POST("/setting/upload/ministry", controller.UploadMinistrySetting)
	ar.GET("/setting/download/ministry", controller.DownloadMinistrySetting)
	ar.GET("/setting/download/national/emblem", controller.DownloadNationalEmblemSetting)

	// ============================================================================
	//   Institutional controller api
	// subsidiary
	ar.POST("/subsidiary/all", institutecontroller.GetSubsidiaries)
	ar.POST("/subsidiary/all/tree", institutecontroller.GetSubsidiariesTree)
	ar.POST("/subsidiary/by/id", institutecontroller.GetSubsidiaryByID)
	ar.POST("/subsidiary/create", institutecontroller.CreateSubsidiary)
	ar.PUT("/subsidiary/update", institutecontroller.UpdateSubsidiary)
	ar.DELETE("/subsidiary/delete", institutecontroller.DeleteSubsidiary)

	// subsidiary - user
	ar.POST("/subsidiary/user/all/by/user", institutecontroller.GetSubsidiariesUserByUserID)
	ar.POST("/subsidiary/user/all/by/user/license", institutecontroller.GetSubsidiariesUserByUserIDLicense)
	ar.PUT("/subsidiary/user/update", institutecontroller.UpdateSubsidiariesUserByUserID)

	// Program
	ar.POST("/program/all", institutecontroller.GetPrograms)
	ar.POST("/program/create", institutecontroller.CreateProgram)
	ar.PUT("/program/update", institutecontroller.UpdateProgram)

	// Program - user
	ar.POST("/program/user/all/by/user", institutecontroller.GetProgramsUserByUserID)
	ar.POST("/program/user/all/by/user/license", institutecontroller.GetProgramsUserByUserIDLicense)
	ar.PUT("/program/user/update", institutecontroller.UpdateProgramsUserByUserID)

	// Program
	ar.POST("/semester/all", institutecontroller.GetSemesters)
	ar.POST("/semester/create", institutecontroller.CreateSemester)
	ar.PUT("/semester/update", institutecontroller.UpdateSemester)
	ar.DELETE("/semester/delete", institutecontroller.DeleteSemester)

	// Student
	ar.POST("/student/all", institutecontroller.GetStudents)
	ar.POST("/student/create", institutecontroller.CreateStudent)
	ar.PUT("/student/update", institutecontroller.UpdateStudent)
	ar.DELETE("/student/delete", institutecontroller.DeleteStudent)
	ar.POST("/student/search", institutecontroller.GetStudentSearch)
	ar.GET("/student/download/template", institutecontroller.GetTempUploadStudent)
	ar.POST("/student/upload/template", institutecontroller.SetTempUploadStudent)
	ar.POST("/student/detail/by/id", institutecontroller.GetStudentDetailByID)
	ar.POST("/student/detail/by/dni", institutecontroller.GetStudentDetailByDNI)

	// Teacher
	ar.POST("/teacher/all", institutecontroller.GetTeachers)
	ar.POST("/teacher/create", institutecontroller.CreateTeacher)
	ar.PUT("/teacher/update", institutecontroller.UpdateTeacher)
	ar.DELETE("/teacher/delete", institutecontroller.DeleteTeacher)
	ar.POST("/teacher/search", institutecontroller.GetTeacherSearch)
	ar.GET("/teacher/download/template", institutecontroller.GetTempUploadTeacher)
	ar.POST("/teacher/upload/template", institutecontroller.SetTempUploadTeacher)
	ar.GET("/teacher/export/all", institutecontroller.ExportAllTeachers)

	// Teacher Program
	ar.POST("/teacher/program/all", institutecontroller.GetTeacherProgramsByProgram)
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

	// ---------------------------------------------------------------------------
	//      Monitoring routes ----------------------------------------------------

	// poll
	ar.POST("/monitoring/poll/all", monitoringcontroller.GetPollsPaginate)
	ar.POST("/monitoring/poll/by/id", monitoringcontroller.GetPollByID)
	ar.POST("/monitoring/poll/create", monitoringcontroller.CreatePoll)
	ar.PUT("/monitoring/poll/update", monitoringcontroller.UpdatePoll)
	ar.DELETE("/monitoring/poll/delete", monitoringcontroller.DeletePoll)

	// Question
	ar.POST("/monitoring/question/all", monitoringcontroller.GetQuestions)
	ar.POST("/monitoring/question/create", monitoringcontroller.CreateQuestions)
	ar.PUT("/monitoring/question/update", monitoringcontroller.UpdateQuestion)
	ar.DELETE("/monitoring/question/delete", monitoringcontroller.DeleteQuestion)

	// Type questions
	ar.POST("/monitoring/type/question/all", monitoringcontroller.GetTypeQuestions)
	ar.DELETE("/monitoring/multiple/question/delete", monitoringcontroller.DeleteMultipleQuestion)

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
	ar.POST("/library/book/create", librarycontroller.CreateBook)
	ar.PUT("/library/book/update", librarycontroller.UpdateBook)
	ar.DELETE("/library/book/delete", librarycontroller.DeleteBook)
	ar.POST("/library/book/by/id", librarycontroller.GetBookByID)
	ar.POST("/library/book/upload/avatar", librarycontroller.UploadAvatarBook)
	ar.POST("/library/book/upload/pdf", librarycontroller.UploadPdfBook)

	// ---------------------------------------------------------------------------
	//      Admission routes -----------------------------------------------------
	// Admission
	ar.POST("/admission/admission/paginate", admissioncontroller.GetAdmissionsPaginate)
	ar.POST("/admission/admission/paginate/exam", admissioncontroller.GetAdmissionsPaginateExam)
	ar.POST("/admission/admission/create", admissioncontroller.CreateAdmission)
	ar.POST("/admission/admission/cancel", admissioncontroller.CancelAdmission)
	ar.PUT("/admission/admission/update", admissioncontroller.UpdateAdmission)

	// ---------------------------------------------------------------------------
	//      External api -----------------------------------------------------
	// RENIEC
	ar.POST("/external/reniec", controller.Reniec)

	// Web Site services
	ar.POST("/web/setting/by/student", controller.GetStudentSettings)
}
