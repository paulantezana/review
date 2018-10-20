package models

// Company struct
type Company struct {
	ID                     uint   `json:"id" gorm:"primary_key"`
	RUC                    string `json:"ruc"  gorm:"type:varchar(11); unique; not null"`
	NombreORazonSocial     string `json:"nombre_o_razon_social"`
	EstadoDelContribuyente string `json:"estado_del_contribuyente"`
	CondicionDeDomicilio   string `json:"condicion_de_domicilio"`
	Ubigeo                 string `json:"ubigeo"`
	TipoDeVia              string `json:"tipo_de_via"`
	NombreDeVia            string `json:"nombre_de_via"`
	CodigoDeZona           string `json:"codigo_de_zona"`
	TipoDeZona             string `json:"tipo_de_zona"`
	Numero                 string `json:"numero"`
	Interior               string `json:"interior"`
	Lote                   string `json:"lote"`
	Dpto                   string `json:"dpto"`
	Manzana                string `json:"manzana"`
	Kilometro              string `json:"kilometro"`
	Departamento           string `json:"departamento"`
	Provincia              string `json:"provincia"`
	Distrito               string `json:"distrito"`
	Direccion              string `json:"direccion"`
	DireccionCompleta      string `json:"direccion_completa"`
	UltimaActualizacion    string `json:"ultima_actualizacion"`

	ReviewDetails []ReviewDetail `json:"review_details"`
}

func (Company) TableName() string {
	return "companies"
}
