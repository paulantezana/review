---
title: "Modelos del sistema de admisión"
date: "2019-27-02"
---


## AdmissionSetting
En esta tabla se almacena la informacion de cada proceso de admision ya sea de cada año o de una temporada
los datos almacenados en esta tabla luego son usado spor la tabla `Admission` para realizar los calculos pertinentes
```go
type AdmissionSetting struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	VacantByProgram uint   `json:"vacant_by_program"`
	Year            uint   `json:"year"`
	Seats           uint   `json:"seats"`
	Description     string `json:"description"`
	ShowInWeb       bool   `json:"show_in_web"`

	SubsidiaryID uint        `json:"subsidiary_id"`
	Admissions   []Admission `json:"admissions, omitempty"`
}
```
### Campos
- **VacantByProgram(requerido)** Campo para ingresar el numero maximo de vacantes por porgrama de estudios
- **Year(requerido)** Campo para ingresar el año actual
- **Seats(requerido)** Campo para ingresar el numero de asientos por clase
- **Description(requerido)** Alguna descripcion extra sobre este porceso de admision
- **ShowInWeb(opcional)** Cuando ya se tenga los resultados de las exámenes que se rindió se puede habilitar  este campo a `true` para mostrar los resultados de manera publica como en un sitio web 
    - `true` : Los resultados de la evaluación estarán de manera pública
    - `false` : Los resultado seran privados
    - Valor por defecto `true`
- **SubsidiaryID(requerido)** ID de la filial en el que se aperturo el proceso de admision
- **Admissions(opcional)** Referencia a la tabla `Admission`

## Admission
Se usa para almacenar los datos del proceso de admission de todo los años es importante
que este tabla este relaciondo con la tabla `AdmissionSetting` ya que hace una copia de 
algunos campos ademas le sirve para hacer calculos como la clase y los asientos

```go
type Admission struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	ExamNote      float32   `json:"exam_note"`
	ExamDate      time.Time `json:"exam_date"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`
	Classroom     uint      `json:"classroom"`
	Seat          uint      `json:"seat"`
	State         bool      `json:"state" gorm:"default:'true'"`

	StudentID          uint `json:"student_id"`
	ProgramID          uint `json:"program_id"`
	UserID             uint `json:"user_id"`
	AdmissionSettingID uint `json:"admission_setting_id"`
}
```

### Campos
- **Observation(opcional)**  Campo para agregar alguna observación de la admision
- **Exonerated(opcional)**  Para especificar si el alumno es exonerado o no
    - false : alumno regular
    - true  : Indica que el alumno es **Exonerado**
    - valor por defecto false
- **ExamNote(opcional)**  Resultadi de la evaluacion del examen que rindio el alumno
- **ExamDate(opcional)**  Fecha en que se rindio el examen
- **AdmissionDate(requerido)** Fecha en que el alumno se inscrivio al proceso de admision
- **Year(requerido)** Campo para ingresar el año actual del porceso de admision. Este campo se debe hacer una copia de la tabla padre AdmissionSetting su campo Year 
- **Classroom(requerido)** Campo para ingresar la clase en la cual redira el examen
- **Seat(requerido)** Campo para ingresar el lugar exacto donde se sentara el alumno para redir su examen
- **State(opcional)** Estado de la admision - Utilize este campo para anular la admision
    - false : admision anulada
    - true  : alumno en proceso de admision
    - Valor por defecto true
- **StudentID(requerido)** ID del alumno a la que se esta registrando en el proceso de admision
- **ProgramID(requerido)** ID del programa de estudios a la cual esta postulando
- **UserID(requerido)** ID del usuario que esta registrando este registro
- **AdmissionSettingID(requerido)** ID de la configuracion de la admision