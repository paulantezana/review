package controller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/utilities"
    "io/ioutil"
    "net/http"
    "strings"
)

type reniecRequest struct {
    DNI string `json:"dni"`
}

type reniecResponse struct {
    DNI string `json:"dni"`
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    FullName string `json:"full_name"`
} 

func Reniec(c echo.Context) error {
    request := reniecRequest{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    url := "http://aplicaciones007.jne.gob.pe/srop_publico/Consulta/Afiliado/GetNombresCiudadano?DNI=" + request.DNI

    req, _ := http.NewRequest("GET", url, nil)

    //req.Header.Add("cache-control", "no-cache")
    //req.Header.Add("postman-token", "077270b3-72f2-29c0-2875-741a4d6cabd3")

    res, _ := http.DefaultClient.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    data := strings.Split(strings.ToLower(string(body)),"|")
    lastName := strings.ToUpper(fmt.Sprintf("%s %s",data[0],data[1]))
    reniecResponse := reniecResponse{
        DNI: request.DNI,
        FirstName: strings.Title(data[2]),
        LastName: lastName,
        FullName: fmt.Sprintf("%s, %s",lastName,strings.Title(data[2])),
    }

    return c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Data: reniecResponse,
    })
}
