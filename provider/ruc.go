package provider

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
)

type Cloud struct {
    Ruc string `json:"ruc"`
    RazonSocial string `json:"razon_social"`
    Ciiu string `json:"ciiu"`
    FechaActividad string `json:"fecha_actividad"`
    ContribuyenteCondicion string `json:"contribuyente_condicion"`
    ContribuyenteTipo string `json:"contribuyente_tipo"`
    ContribuyenteEstado string `json:"contribuyente_estado"`
    NombreComercial string `json:"nombre_comercial"`
    FechaInscripcion string `json:"fecha_inscripcion"`
    DomicilioFiscal string `json:"domicilio_fiscal"`
    SistemaEmision string `json:"sistema_emision"`
    SistemaContabilidad string `json:"sistema_contabilidad"`
    ActividadExterior string `json:"actividad_exterior"`
    EmisionElectronica string `json:"emision_electronica"`
    FechaInscripcion_ple string `json:"fecha_inscripcion_ple"`
    Oficio string `json:"oficio"`
    FechaBaja string `json:"fecha_baja"`
    RepresentanteLegal string `json:"representante_legal"`
    Empleados string `json:"empleados"`
    Locale string `json:"locale"`
}

func SunatCloud(dni string) (Cloud, error) {
    cloud := Cloud{}

    // Prepare
    url := "https://api.sunat.cloud/ruc/" + dni
    req, _ := http.NewRequest("GET", url, nil)

    // Send Query
    res, err := http.DefaultClient.Do(req)
    if res != nil {
        defer res.Body.Close()
    }
    if err != nil {
        return cloud, err
    }

    // Read
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return cloud, err
    }

    // string to struct
    err = json.Unmarshal([]byte(string(body)), &cloud)
    if err != nil {
        return cloud, err
    }

    // return response success
    return cloud, nil
}