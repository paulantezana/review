package provider

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Essalud struct {
	DNI             string
	ApellidoPaterno string
	Nombres         string
	FechaNacimiento string
	Sexo            string
	ApellidoMaterno string
}

type EssaludContainer struct {
	DatosPerson []Essalud
}

func EssaludQuery(dni string) (EssaludContainer, error) {
	container := EssaludContainer{}

	url := "https://ww1.essalud.gob.pe/sisep/postulante/postulante/postulante_obtenerDatosPostulante.htm"
	payload := strings.NewReader("strDni=" + dni)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
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
	if err != nil {
		return container, err
	}

	// string to struct
	err = json.Unmarshal([]byte(string(body)), &container)
	if err != nil {
		return container, err
	}

	// return response success
	return container, nil
}

type JnePerson struct {
	Name      string
	FirstName string
	LastName  string
}

func JneQuery(dni string) (JnePerson, error) {
	jnePerson := JnePerson{}

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
