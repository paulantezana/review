package controller

import (
    "encoding/json"
    "fmt"
    "github.com/paulantezana/review/models"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

type essalud struct {
    DNI string
    ApellidoPaterno string
    Nombres string
    FechaNacimiento string
    Sexo string
    ApellidoMaterno string
}

type essaludContainer struct {
    DatosPerson []essalud
}

func essaludQuery(dni string) (essaludContainer, error)  {
    container := essaludContainer{}

    url := "https://ww1.essalud.gob.pe/sisep/postulante/postulante/postulante_obtenerDatosPostulante.htm"
    payload := strings.NewReader("strDni=" + dni)

    req, err := http.NewRequest("POST", url, payload)
    if err != nil  {
        return container, err
    }
    req.Header.Add("content-type", "application/x-www-form-urlencoded")

    res, err := http.DefaultClient.Do(req)
    if res != nil {
        defer res.Body.Close()
    }
    if err != nil {
        return container, err
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil  {
        return container, err
    }

    // string to struct
    err = json.Unmarshal([]byte(string(body)), &container)
    if err != nil  {
        return container, err
    }

    // return response success
    return container, nil
}

type jnePerson struct {
    Name string
    FirstName string
    LastName string
}

func jneQuery(dni string) (jnePerson, error)  {
    jnePerson := jnePerson{}

    url := "http://aplicaciones007.jne.gob.pe/srop_publico/Consulta/Afiliado/GetNombresCiudadano?DNI=" + dni
    req, _ := http.NewRequest("GET", url, nil)

    res, err := http.DefaultClient.Do(req)
    if res != nil {
        defer res.Body.Close()
    }
    if err != nil {
        return jnePerson, err
    }

    body, _ := ioutil.ReadAll(res.Body)

    // Split string
    data := strings.Split(strings.ToLower(string(body)), "|")

    jnePerson.Name = strings.Title(data[2])
    jnePerson.FirstName = strings.Title(data[0])
    jnePerson.LastName = strings.Title(data[1])

    // return response success
    return jnePerson, nil
}

func Dni(dni string) (models.Student, error) {
    student := models.Student{}
    student.DNI = dni

    // Query essalud
    container, err := essaludQuery(dni)

    // Validate
    if err == nil && len(container.DatosPerson) >= 1 {
        first := strings.ToUpper(container.DatosPerson[0].ApellidoPaterno)
        last := strings.ToUpper(container.DatosPerson[0].ApellidoMaterno)
        name := strings.Title(strings.ToLower(container.DatosPerson[0].Nombres))
        student.FullName = fmt.Sprintf("%s %s, %s",first, last, name)
        if container.DatosPerson[0].Sexo == "2" {
            student.Gender = "M"
        } else {
            student.Gender = "F"
        }
        layout := "02/01/2006"
        t, _ := time.Parse(layout, container.DatosPerson[0].FechaNacimiento)
        student.BirthDate = t
    }else {
        // Query jne
        data, _ := jneQuery(student.DNI)
        student.FullName =  fmt.Sprintf("%s %s, %s", strings.ToUpper(data.LastName), strings.ToUpper(data.FirstName), strings.Title(data.Name))
    }

    // Response data
    return student, nil
}
