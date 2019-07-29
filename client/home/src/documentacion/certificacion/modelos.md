---
title: "Modelos del sistema de certificación"
date: "2019-27-02"
---

## Course
Se usa para almacenar los CURSOS
```go
type Course struct {
	ID                      uint      `json:"id" gorm:"primary_key"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Price                   float32   `json:"price"`
	ResolutionAuthorization string    `json:"resolution_authorization"`
	DateExam                time.Time `json:"date_exam"`

	CourseStudents []CourseStudent `json:"course_students"`
}
```
### Campos
- **Name(requerido)** Campo para registrar el nombre del curso
- **Description(opcional)** Campo para registrar la descripción corta del curso que se llevara
- **StartDate(requerido)** Campo para registrar la fecha de inicio del curso
- **EndDate(requerido)** Campo para registrar fecha de finalización del curso
- **Price(opcional)** Campo para registrar el precio que se establece al curso
- **ResolutionAuthorization(opcional)** Campo para registrar la resolución de autorización del curso
- **DateExam(opcional)** Fecha estimado para el examen final del curso


## CourseStudent
Se usa para almacenar los alumnos que se registraron en algún curso en especifico

```go
type CourseStudent struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	DNI      string  `json:"dni" gorm:" type:varchar(15)"`
	FullName string  `json:"full_name" gorm:"type:varchar(250)"`
	Phone    string  `json:"phone" gorm:"type:varchar(32)"`
	State    bool    `json:"state" gorm:"default:'true'"`
	Gender   string  `json:"gender"`
	Year     uint    `json:"year"`
	Payment  float32 `json:"payment"`
	Note     float32 `json:"note"`

	StudentID uint `json:"student_id"`

	CourseID  uint `json:"course_id"`
	ProgramID uint `json:"program_id"`

	CourseExams []CourseExam `json:"course_exams"`
}
```

### Campos
- **DNI(requerido)** Campo para ingresar el DNI del alumno
- **FullName(requerido)** Campo para ingresar el nombre completo del alumno se recomienda usar el siguiente formato APELLIDO APELLIDO, Nombre
- **Phone(requerido)** Teléfono o número celular del alumno
- **State(requerido)** Estado del alumno
- **Gender(requerido)** Sexo del alumno F = femenino, M = masculino
- **Year(opcional)** *Aun sin uso*
- **Payment(requerido)** Monto que está pagando al inscribirse al curso este dato se debe hacer una copia de la tabla `Course`
- **Note(opcional)** Resultado del examen que rindió el alumno en el curso
- **StudentID(opcional)** ID del estudiante *aun esta sin uso*
- **CourseID(requerido)** ID del curso a la que se está inscribiendo el alumno
- **ProgramID(requerido)** ID programa de estudios a la que pertenece el alumno
- **CourseExams(opcional)** *aun sin uso*

## CourseExam
*aun sin uso*

```go
type CourseExam struct {
	ID   uint      `json:"id" gorm:"primary_key"`
	Note float32   `json:"note"`
	Date time.Time `json:"date"`

	CourseStudentID uint `json:"course_student_id"`
}
```