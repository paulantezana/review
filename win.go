package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/review/endpoint"
	"github.com/paulantezana/review/migration"
	"github.com/paulantezana/review/provider"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	// Hosts
	hosts := map[string]*Host{}

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = provider.GetConfig().Server.Port
	}

	//----------------------------------------------------------------------------------------
	// Home
	//----------------------------------------------------------------------------------------
	home := echo.New()
	home.Use(middleware.Logger())
	home.Use(middleware.Recover())

	hosts["localhost:"+port] = &Host{home}

	home.Static("/", "client/home/public")
	//home.File("/","client/home/index.html")

	//----------------------------------------------------------------------------------------
	// API
	//----------------------------------------------------------------------------------------
	api := echo.New()
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	hosts["api.localhost:"+port] = &Host{api}

	// Migration database
	migration.Migrate()

	// Configuration cor
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT, echo.OPTIONS},
	}))

	// Assets
	static := api.Group("/static")
	static.Static("", "static")

	// Root router success
	api.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Stating Server OK")
	})

	// Sting API services
	endpoint.PublicApi(api)
	endpoint.ProtectedApi(api)
    endpoint.CoreApi(api)
	endpoint.PublicWs(api)

	//----------------------------------------------------------------------------------------
	// Administration
	//----------------------------------------------------------------------------------------

	admin := echo.New()
	admin.Use(middleware.Logger())
	admin.Use(middleware.Recover())

	hosts["administracion.localhost:"+port] = &Host{admin}

	admin.Static("/", "client/admin")

	//----------------------------------------------------------------------------------------
	// Student
	//----------------------------------------------------------------------------------------

	student := echo.New()
	student.Use(middleware.Logger())
	student.Use(middleware.Recover())

	hosts["alumno.localhost:"+port] = &Host{student}

	student.Static("/", "client/student")

	//----------------------------------------------------------------------------------------
	// Teacher
	//----------------------------------------------------------------------------------------

	teacher := echo.New()
	teacher.Use(middleware.Logger())
	teacher.Use(middleware.Recover())

	hosts["profesor.localhost:"+port] = &Host{teacher}

	teacher.Static("/", "client/student")

	// ----------------------------------------------------------------------------------------
	// Server Start
	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}
		return
	})

	// Open Bowser
	var err error
	urlOpen := fmt.Sprintf("http://localhost:%s", provider.GetConfig().Server.Port)
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", urlOpen).Run()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", urlOpen).Run()
	case "darwin":
		err = exec.Command("open", urlOpen).Run()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		println(err)
	}

	// Start Server
	e.Logger.Fatal(e.Start(":" + port))
}
