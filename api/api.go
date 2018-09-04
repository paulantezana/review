package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/controller"
	"github.com/paulantezana/review/utilities"
)

func PublicApi(e *echo.Echo) {
	pb := e.Group("/api/v1/public")
	pb.POST("/user/login", controller.Login)
	pb.POST("/user/forgot/search", controller.ForgotSearch)
	pb.POST("/user/forgot/validate", controller.ForgotValidate)
	pb.POST("/user/forgot/change", controller.ForgotChange)
}

func ProtectedApi(e *echo.Echo) {
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(config.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// Global settings
	ar.POST("/setting/global", controller.GetGlobalSettings)
	ar.PUT("/setting/update", controller.UpdateSetting)
	ar.POST("/setting/upload/logo", controller.UploadLogoSetting)
	ar.GET("/setting/download/logo", controller.DownloadLogoSetting)

    // Program
    ar.POST("/program/all", controller.GetPrograms)
    ar.POST("/program/create", controller.CreateProgram)
    ar.PUT("/program/update", controller.UpdateProgram)

	// Student
	ar.POST("/student/all", controller.GetStudents)
	ar.POST("/student/create", controller.CreateStudent)
	ar.PUT("/student/update", controller.UpdateStudent)
	ar.DELETE("/student/delete", controller.DeleteStudent)
	ar.POST("/student/search", controller.GetStudentSearch)
	ar.GET("/student/download/template", controller.GetTempUploadStudent)
	ar.POST("/student/upload/template", controller.SetTempUploadStudent)
	ar.POST("/student/detail/by/id", controller.GetStudentDetailByID)

    // Student
    ar.POST("/teacher/all", controller.GetTeachers)
    ar.POST("/teacher/create", controller.CreateTeacher)
    ar.PUT("/teacher/update", controller.UpdateTeacher)
    ar.DELETE("/teacher/delete", controller.DeleteTeacher)
    ar.POST("/teacher/search", controller.GetTeacherSearch)
    ar.GET("/teacher/download/template", controller.GetTempUploadTeacher)
    ar.POST("/teacher/upload/template", controller.SetTempUploadTeacher)

	// Module
	ar.POST("/module/all", controller.GetModules)
	ar.POST("/module/create", controller.CreateModule)
	ar.PUT("/module/update", controller.UpdateModule)
	ar.DELETE("/module/delete", controller.DeleteModule)
	ar.POST("/module/search", controller.GetModuleSearch)

	// Company
	ar.POST("/company/all", controller.GetCompanies)
	ar.POST("/company/create", controller.CreateCompany)
	ar.PUT("/company/update", controller.UpdateCompany)
	ar.DELETE("/company/delete", controller.DeleteCompany)
	ar.POST("/company/search", controller.GetCompanySearch)

	// Review
	ar.POST("/review/all", controller.GetReviews)

	ar.POST("/review/create", controller.CreateReview)
	ar.PUT("/review/update", controller.UpdateReview)
	ar.DELETE("/review/delete", controller.DeleteReview)

	// User
	ar.POST("/user/all", controller.GetUsers)
	ar.POST("/user/create", controller.CreateUser)
	ar.PUT("/user/update", controller.UpdateUser)
	ar.DELETE("/user/delete", controller.DeleteUser)
	ar.POST("/user/by/id", controller.GetUserByID)
	ar.POST("/user/upload/avatar", controller.UploadAvatarUser)
	ar.POST("/user/reset/password", controller.ResetPasswordUser)
	ar.POST("/user/change/password", controller.ChangePasswordUser)

	// Review Detail
	ar.POST("/reviewdetail/all", controller.GetReviewsDetail)
}
