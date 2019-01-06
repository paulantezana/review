package controller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/review/config"
    "github.com/paulantezana/review/models"
    "github.com/paulantezana/review/models/institutemodel"
    "github.com/paulantezana/review/utilities"
    "io/ioutil"
    "net/http"
    "strings"
)

type reniecRequest struct {
    DNI string `json:"dni"`
}

type reniecResponse struct {
    Student institutemodel.Student `json:"student"`
    User models.User `json:"user"`
    Exist bool `json:"exist"`
} 

func Reniec(c echo.Context) error {
    request := reniecRequest{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Search student
    student := institutemodel.Student{}
    DB.First(&student,institutemodel.Student{DNI: request.DNI})
    reniecResponse := reniecResponse{
        Student: student,
    }

    // Validation
    if student.ID == 0 {
        url := "http://aplicaciones007.jne.gob.pe/srop_publico/Consulta/Afiliado/GetNombresCiudadano?DNI=" + request.DNI
        req, _ := http.NewRequest("GET", url, nil)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("Error en la consulta a la Reniec")} )
        }
        defer res.Body.Close()

        body, _ := ioutil.ReadAll(res.Body)

        // Split string
        data := strings.Split(strings.ToLower(string(body)),"|")
        lastName := strings.ToUpper(fmt.Sprintf("%s %s",data[0],data[1]))
        firstName := strings.Title(data[2])

        // fill data
        reniecResponse.Student.DNI = request.DNI
        reniecResponse.Student.FullName = fmt.Sprintf("%s, %s",lastName,firstName)
    }else {
        // Find User
        user:= models.User{}
        if err := DB.First(&user,student.UserID).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{ Message: fmt.Sprintf("%s", err)} )
        }
        user.Password = ""
        user.Key = ""

        reniecResponse.Exist = true
        reniecResponse.User = user
    }

    return c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Data: reniecResponse,
    })
}
