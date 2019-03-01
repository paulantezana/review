package controller

import (
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
    cd := time.Now().Format("01-02-2006")
	out, err := exec.Command("pg_dump","--host", con.Database.Server, "--port",con.Database.Port,"--username",con.Database.User,"--file",fmt.Sprintf("static/backup/database/%s_%s.sql",con.Database.Database,cd),con.Database.Database).Output()
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se genero un archivo de copia de seguridad de la base de datos actual. %s", out),
	})
}

func BackupDBList(c echo.Context) error {
    fls, err := ioutil.ReadDir("static/backup/database/")
    if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }
    files := make([]utilities.File,0)
    for _, f := range fls {
        files = append(files, utilities.File{
            Name: f.Name(),
            Size: uint(f.Size()),
            Date: f.ModTime(),
        })
    }
    return  c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Data: files,
    })
}
