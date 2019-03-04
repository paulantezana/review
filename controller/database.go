package controller

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/review/config"
	"github.com/paulantezana/review/utilities"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

func BackupDB(c echo.Context) error {
	con := config.GetConfig()

	// Hash current time
	cc := sha256.Sum256([]byte(time.Now().String()))
	cd := fmt.Sprintf("%x", cc)

	// Execute command
	out, err := exec.Command("pg_dump", "--host", con.Database.Server, "--port", con.Database.Port, "--username", con.Database.User, "--file", fmt.Sprintf("static/backup/database/%s.sql", cd), con.Database.Database).Output()
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response json
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se genero un archivo de copia de seguridad de la base de datos actual. %s", out),
	})
}

func BackupDBList(c echo.Context) error {
	// Query folder database
	fls, err := ioutil.ReadDir("static/backup/database/")
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Prepare struct
	files := make([]utilities.File, 0)
	for _, f := range fls {
		files = append(files, utilities.File{
			Name: f.Name(),
			Size: uint(f.Size()),
			Date: f.ModTime(),
		})
	}

	// Response json
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    files,
	})
}
