package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/api"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/models"
	"os"
    "crypto/sha256"
    "fmt"
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

func migration() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.User{},
		&models.Module{},
		&models.Representative{},
		&models.Review{},
		&models.ReviewDetail{},
		&models.Setting{},
		&models.Student{},
		&models.Company{},
	)
	db.Model(&models.Review{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Review{}).AddForeignKey("module_id", "modules(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.ReviewDetail{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.ReviewDetail{}).AddForeignKey("review_id", "reviews(id)", "RESTRICT", "RESTRICT")

    // -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	usr := models.User{}
	db.First(&usr)
	// hash password
    cc := sha256.Sum256([]byte("admin"))
    pwd := fmt.Sprintf("%x", cc)
    // create model
	user := models.User{
	    UserName: "admin",
	    Password: pwd,
	    Email: "yoel.antezana@gmail.com",
    }
    // insert database
    if usr.ID == 0 {
        db.Create(&user)
    }

    // First Setting
    cg := models.Setting{}
    db.First(&cg)
    co := models.Setting{ ItemTable: 10 }
    // Insert database
    if cg.ID == 0 {
        db.Create(&co)
    }
}
