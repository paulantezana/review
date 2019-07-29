---
title: "Modelos del sistema de revisión"
date: "2019-27-02"
---

## Company
Sirve para almacenar todas las empresas posibles en las que se realizan las practicas en situaciones reales de trabajo.

```go
type Company struct {
	ID               uint           `json:"id" gorm:"primary_key"`
	RUC              string         `json:"ruc"  gorm:"type:varchar(11); unique; not null"`
	NameSocialReason string         `json:"name_social_reason"`
	Address          string         `json:"address"`
	Manager          string         `json:"manager"`
	Phone            string         `json:"phone"`
	CompanyType      string         `json:"company_type"` // 1 = public || 2 = private
}
```

### Campos
- **RUC(requerido)** : Campo para registar el RUC de la empresa
- **NameSocialReason(requerido)** : Campo para registrar la razón social o nombre de la empresa
- **Address(opcional)** : Campo para registrar la dirección exacta de donde este ubicado geográficamente.
- **Manager(opcional)** : Campo para registrar el nombre del gerente o representante de la empresa.
- **Phone(opcional)** : Campo para registrar el número telefónico o cualquier otro numero de contacto
- **CompanyType(requerido)** : Campo para registrar el tipo de empresa si una publica o privada
	- `private`	: Empresa privada 
	- `public`	: Empresa publica

## Review
Se usa para almacenar todas las revisiones de toda las carreras y filiales
```go
type Review struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	ApprobationDate time.Time `json:"approbation_date"`

	ModuleId         uint `json:"module_id"`
	StudentProgramID uint `json:"student_program_id"`
	CreatorID        uint `json:"creator_id"`
	TeacherID        uint `json:"teacher_id"`

	ReviewDetails []ReviewDetail `json:"review_details"`
}
```

### Campos
- **ApprobationDate(opcional)** : Campo para registrar la fecha en la que se aprobó la revisión por su asesor.
- **ModuleId(requerido)** : ID del modulo a la se se le está haciendo la revisión correspondiente
- **StudentProgramID(opcional)** : ID del estudiante que se esta haciendo la revisión de las practicas que este este relacionado a la carrera con corresponde el `ModuleID`
- **CreatorID(opcional)** : ID de usuario quien esta registrando la nueva revisión en la base de datos.
- **TeacherID(opcional)** : ID de asesor que se le asigno al estudiante para la revisión de sus practicas modulares en situaciones reales de trabajo.
- **ReviewDetails(requerido)** : Referencia al modelo ` ReviewDetail` que contiene los desatolles de la revisión se creó una tabla de detalles en el caso de que en algún programa de estudios un modulo puede tener muchas practicas en situaciones reales de trabajo.

## ReviewDetail
Se usa para almacenar los detalles de una revisión de un modulo
```go
type ReviewDetail struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Hours     uint      `json:"hours"`
	Note      uint      `json:"note"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	ReviewID  uint `json:"review_id"`
	CompanyID uint `json:"company_id"`
}
```
- **Hours(requerido)** : Campo para registrar el número de horas que realizo sus practicas modulares
- **Note(requerido)** : Campo para registrar la nota que obtuvo durante sus prácticas modulares.
- **StartDate(requerido)** : Campo para registrar la fecha en la que inicio sus prácticas modulares
- **EndDate(requerido)** : Campo para registrar la fecha en la que finalizo sus prácticas modulares
- **ReviewID(requerido)** : ID de la revisión a la que pertenece este registro o detalle
- **CompanyID(requerido)** : ID de la empresa en la que realizo sus prácticas modulares.
