package main

import (
	"github.com/paulantezana/review/endpoint"
	"github.com/paulantezana/review/migration"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/provider"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Migration database
	migration.Migrate()

	// Configuration cor
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT, echo.OPTIONS},
	}))

	// Assets
	static := e.Group("/static")
	static.Static("", "static")

	// Root router success
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Stating Server OK")
	})

	// Sting API services
	endpoint.PublicApi(e)
	endpoint.ProtectedApi(e)
	endpoint.PublicWs(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = provider.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}
