---
title: "Modelos de la base de datos"
date: "2019-27-02"
---

## Subsidiary:
Se usa para almacenar las filiales
```go
type Subsidiary struct {
	ID                  uint   `json:"id" gorm:"primary_key"`
	Name                string `json:"name"`
	Country             string `json:"country"`
	Department          string `json:"department"`
	Province            string `json:"province"`
	District            string `json:"district"`
	TownCenter          string `json:"town_center"`
	Address             string `json:"address"`
	Main                bool   `json:"main"`
	Phone               string `json:"phone"`
}
```

### Campos
- **Name(requerido)** Campo para almacenar el nombre de la filial
- **Country(requerido)** Campo para almacenar el país donde se encuentra la filial
- **Department(requerido)** Campo para almacenar el departamento donde se encuentra la filial
- **Province(opcional)** Campo para almacenar la provincia donde se encuentra la filial
- **District(opcional)** Campo para almacenar el distrito donde se encuentra la filial
- **TownCenter(opcional)** Campo para almacenar el centro poblado donde se encuentra la filial
- **Address(requerido)** Campo para almacenar la dirección exacta donde se encuentra la filial
- **Main(opcional)** Campo para indicar cual es la filial principal en todo el registro solo debe existir uno
    - `false`   : Es filial
    - `true`    : Central
    Valor por defecto `false`
- **Phone(opcional)** Campo para almacenar el teléfono de la filial

## SubsidiaryUser
Se usa para asignarle los permisos a un usuario sobre una filial.
se tiene que crear de forma automática las referencias te todo los usuarios con todo los filiales que existan
con el valor por defecto en el campo `License` en false

```go
type SubsidiaryUser struct {
	ID           uint `json:"id" gorm:"primary_key"`
	UserID       uint `json:"user_id"`
	SubsidiaryID uint `json:"subsidiary_id"`
	License      bool `json:"license"`
}
```

### Campos
- **UserID(requerido)** ID del usuario a la que se esta asignando el permiso
- **SubsidiaryID(requerido)** ID de la filial al que el usuario podrá acceder
- **License(requerido)** Licencia si el usuario tiene permiso o no
    - `false`   : Indica que el usuario no podrá acceder a la filial
    - `true`    : el usuario tiene permiso para acceder a la filial
    - Valor por defecto `false`




## Program
Se usa para almacenar los programas de estudios de cada filial
```go
type Program struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" type:varchar(255); unique; not null"`
	Level string `json:"level"`

	SubsidiaryID uint `json:"subsidiary_id"`
}
```
### Campos
- **Name(requerido)** Campo para almacenar el nombre del programa de estudios
- **Level(requerido)** Campo para almacenar el nivel académico del programa de estudios
- **SubsidiaryID(requerido)** ID de la filial a la que pertenece el programa de estudios



## ProgramUser
Se usa para asignarle los permisos a un usuaro sobre un programa de estudios.
se tiene que crear de forma automática las referencias te todo los usuarios con todo los programas de estudios que existan
con el valor por defecto en el campo `License` en false

```go
type ProgramUser struct {
	ID        uint `json:"id" gorm:"primary_key"`
	UserID    uint `json:"user_id"`
	ProgramID uint `json:"program_id"`
	License   bool `json:"license"`
}
```

### Campos
- **UserID(requerido)** ID del usuario a la que se esta asignando el permiso
- **ProgramID(requerido)** ID del programa de estudios al que el usuario podrá acceder
- **License(requerido)** Licencia si el usuario tiene permiso o no
    - `false`   : Indica que el usuario no podrá acceder al programa de estudios
    - `true`    : el usuario tiene permiso para acceder al programa de estudios
    - Valor por defecto `false`



## Semester
Se usa para almacenar los semestres de cada programa de estudios
```go
type Semester struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"type:varchar(128); not null"`
	Sequence uint   `json:"sequence"`
	Period   string `json:"period"`
	Year     uint   `json:"year" gorm:"not null"`

	ProgramID uint `json:"program_id"`
}
```

### Campos
- **Name(requerido)** Campo para almacenar el nombre del semestre
- **Sequence(requerido)** Campo para almacenar la secuencia numérica del semestre
- **Period(opcional)** Campo para almacenar el periodo academico
- **Year(opcional)** Campo para almacenar el año del semestre
- **ProgramID(requerido)** ID del programa de estudios a la que pertenece este semestre


## Module
Se usa para alamcenar los modulos de un programa de estudios
```go
type Module struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Sequence    uint   `json:"sequence"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      uint   `json:"points"`
	Hours       uint   `json:"hours"`

	ProgramID uint `json:"program_id"`

	Semesters []ModuleSemester `json:"semesters"`
}
```

### Campos
- **Sequence(requerido)** Campo para alamcenar la secuencia numerica del modulo
- **Name(requerido)** Campo para alamcenar el nombre del modulo
- **Description(opcional)** Campo para alamcenar la descripcion completa del modulo
- **Points(requerido)** Campo para alamcenar los creditos que tiene el modulo
- **Hours(requerido)** Campo para alamcenar las horas que tiene este modulo
- **ProgramID(requerido)** ID del programa de estudios a la que pertenece el modulo
- **Semesters(opcional)** Referencia a los semestres que esta asignado el modulo


## ModuleSemester
Se usa en la relacionde muchos a muchos de las tablas `Module` y `Semester`
```go
type ModuleSemester struct {
	SemesterID uint `json:"semester_id" gorm:"primary_key"`
	ModuleID   uint `json:"module_id" gorm:"primary_key"`
}
```

### Campos
- **SemesterID(requerido)** ID del semestre
- **ModuleID(requerido)** ID del modulo

## Student
Se usa para almacenar a los alumnos
```go
type Student struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	DNI         string    `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName    string    `json:"full_name" gorm:"type:varchar(250)"`
	Phone       string    `json:"phone" gorm:"type:varchar(32)"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	BirthDate   time.Time `json:"birth_date"`
	BirthPlace  string    `json:"birth_place"`
	Country     string    `json:"country"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	Region      string    `json:"region"`
	MarketStall string    `json:"market_stall"`
	CivilStatus string    `json:"civil_status"`
	IsWork      string    `json:"is_work"` // si || no

	UserID          uint `json:"user_id"`
	StudentStatusID uint `json:"student_status_id"`

	Reviews []Review `json:"reviews, omitempty"`
	User    User     `json:"user"`
}
```

### Campos
- **DNI(requerido)** Campo para almacenar el DNI del alumno
- **FullName(requerido)** Campo para almacenar el nombre completo del alumno se recomiendo usar el siguiente formato `APELLIDO APELLIDO, Nombres`
- **Phone(requerido)** Campo para almacenar el telefono o el numero de calular del alumno
- **Gender(requerido)** Campo para almacenar el sexo del alumno `F` o `M` segun lo que coresponda
- **Address(opcinal)** Campo para almacenar la direccion de recidencia del alumno
- **BirthDate(requerido)** Campo para almacenar la fecha de nacimiento del alumno
- **BirthPlace(opcinal)**  Campo para almacenar el lugar de nacimiento del alumno
- **Country(opcinal)**  Campo para almacenar el pais de origen del alumno
- **District(opcinal)**  Campo para almacenar el distrito de nacimiento del alumno
- **Province(opcinal)**  Campo para almacenar la provincia de nacimiento del alumno
- **Region(opcinal)**  Campo para almacenar la región de nacimiento del alumno
- **MarketStall(opcinal)** Campo para almacenar en la que trabaja actualmente el alumno en cuestión
- **CivilStatus(opcinal)** Campo para almacenar el estado civil del alumno
- **IsWork(opcinal)** Campo para almacenar si el alumno trabaja o no
    - `si`
    - `no`
- **UserID(requerido)** ID del usuario a la que esta asignado el alumno
- **StudentStatusID(requerido)** ID del estado actual del alumno


## StudentHistory
Se usa para almacenar el historial del alumno cada que se hace alguna modificacion, o acciones en el sistema
como por ejemplo:
- Inscribir a un poceros de admisión
- Rendir la evaluación
- Revisiones de las practicas modulares
- Etc.

```go
type StudentHistory struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Description string    `json:"description"`
	StudentID   uint      `json:"student_id"`
	UserID      uint      `json:"user_id"`
	Type        uint      `json:"type"`
	Date        time.Time `json:"date"`
}
```
### Campos
- **Description(requerido)** Campo para almacenar una descripción corta del nuevo historial del alumno
- **StudentID(requerido)** ID del estudiante al que se le esta creando el nuevo historial
- **UserID(requerido)** ID del usuario que esta creando el nuevo historial
- **Type(requerido)**  ID tipo deaccion que se realizo
    - 1 : Cuando hubo uan accion de crear
    - 2 : Cuando hubo una accion de actualizar
    - 3 : Cuando hubo una accion de imprimir
- **Date(requerido)** Campo para almacenar la fecha en la que se creo el nuevo historial


## StudentProgram
Se usa para relacionar los alumnos con los programas de estudios para 
hacerle el seguimiento de forma adecuada

esta tabla se creo con el propósito de que un alumno puede estudiar mas de un programa de estudios en 
la misma institución deforma que con esta tabla se podrá guardar de todo su historial durante la estadía 
en la que estuvo en cada programa de estudios.

```go
type StudentProgram struct {
	ID            uint `json:"id" gorm:"primary_key"`
	StudentID     uint `json:"student_id"`
	ProgramID     uint `json:"program_id"`
	ByDefault     bool `json:"by_default"`
	YearAdmission uint `json:"year_admission"`
	YearPromotion uint `json:"year_promotion"`
}
```

### Campos
- **StudentID(requerido)** ID del alumno al que se va relacionar aun programa de estudios
- **ProgramID(requerido)** ID del programa de estudios con lo que se relacionara el alumno
- **ByDefault(requerido)** Campo que define cual es la relación principal del estudiante con el programa
- **YearAdmission(opcional)** Campo para almacenar el año de admisión del alumno
- **YearAdmission(opcional)** Campo para almacenar el año de promoción o egreso del alumno


## StudentStatus
Se usa para almacenar los diferentes estados que puede tener el alumno

```go
type StudentStatus struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
```

### Campos
- **Name()** Campo usado para almacenar nombre del estado del alumno

### Datos
los todos de esta tabla se insertan de forma dinámica al migrar los modelos a la base de datos

1. **No asignado**   : Es cuando un alumno fue insertado de forma dinámica desde los archivos de Excel
2. **Postulante**    : Cunado un alumno se inscribe en un proceso de admisión
3. **Exonerado**     : Cunado un alumno es Exonerado en el proceso de admisión
4. **Trasladado**    : Cunado un alumno se traslado desde otras institución
5. **Rechazado**     : Cuando un alumno no pudo postular a un porgrma que se inscribio en el porceso de admision
6. **Aprobado**      : Cunado un alumno ingreso a la institución mediante el proceso de admisión
7. **Prematriculado**    : Prematriculado
8. **Matriculado**       : Matriculado
9. **Expulsado**         : Expulsado
10. **Egresado**         : Egresados


## Teacher
Se usa para almacenar los datos de los profesores
```go
type Teacher struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	DNI       string `json:"dni" gorm:"type:varchar(15); not null; unique"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`

	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`

	Address string `json:"address"`
	Phone   string `json:"phone"`
	WorkConditions string    `json:"work_conditions"`
	EducationLevel string    `json:"education_level"`
	AdmissionDate  time.Time `json:"admission_date"`
	RetirementDate time.Time `json:"retirement_date"`
	Specialty string `json:"specialty"`

	UserID uint `json:"user_id"`

	Type      string `json:"type" gorm:"-"`
	ProgramID uint   `json:"program_id" gorm:"-"`

	TeacherPrograms []TeacherProgram `json:"teacher_programs"`
}
```

### Campos
- **DNI(requerido)** Campo para almacenar el DNI del profesor
- **LastName(requerido)** Campo para almacenar los apellidos del profesor
- **FirstName(requerido)** Campo para almacenar los nombre profesor
- **BirthDate(requerido)** Campo para almacenar la fecha de nacimiento del profesor
- **Gender(requerido)** Campo para almacenar el sexo del profesor `F` o `M` según lo que corresponda
- **Address(requerido)** Campo para almacenar la direccion de recidencia del profesor
- **Phone(requerido)** Campo para almacenar el telefono o el numero de calular del profesor
- **WorkConditions(requerido)** Campo para almacenar la condición laboral profesor
- **EducationLevel(requerido)** Campo para almacenar nivel de educación del profesor
- **AdmissionDate(requerido)** Campo para almacenar fecha de ingreso al instituto del profesor
- **RetirementDate(requerido)** Campo para almacenar la fecha de retiro del instituto del profesor
- **Specialty(requerido)** Campo para almacenar la especialidad profesor
- **UserID(requerido)** ID del usuario a la que esta asignado el profesor
- **Type(euxiliar)** Si es un profesor transversal o de carrera
- **ProgramID(euxiliar)** ID del programa a la que se esta registrando
- **TeacherPrograms(requerido)** Referencia a todo los programas de estudios en la que esta el profesor


## TeacherAction
Se usa para registrar toda las acciones que realiza el profesor en el sistema como por ejemplo 
registrar las notas del estudiante hacer supervisiones y entre otros.

```go
type TeacherAction struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Action      string `json:"action"`
	Description string `json:"description"`

	TeacherID uint `json:"teacher_id"`
}
```

### Campos
- **Action(requerido)** Campo para registrar la acción del profesor como
    - `create`  : crear un nuevo registro
    - `update`  : actualizar un nuevo registro
    - `delete`  : eliminar un nuevo registro
    - `print`  : imprimir un nuevo registro
- **Description(opcional)** Campo para registrar la descripción corta de la acción
- **TeacherID(requerido)** ID del profesor quien esta realizando la acción


## TeacherProgram
Se usa para relacionar los profesores con los programas de estudios para 

ya que un profesor puede estar asignado en diferentes programas de estudios
- como profesores de carrera
- y los profesores transversales

```go
type TeacherProgram struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	TeacherID uint   `json:"teacher_id"`
	ProgramID uint   `json:"program_id"`
	ByDefault bool   `json:"by_default"`
	Type      string `json:"type"` //  // 
}
```

### Campos
- **TeacherID(requerido)** ID del profesor que esta relacionado al programa de estudios
- **ProgramID(requerido)** ID del programa de estudios a la que esta relacionado el profesor
- **ByDefault(requerido)** Asignar el programa por defecto a la que pertenece el profesor
- **Type(requerido)** 
    - `cross`   : transversales
    - `career`  : de carrera


## Unity
Se usa para registrar todas los unidades didácticas que se dictan en un programa de estudios 
por semestre

```go
type Unity struct {
	ID     uint    `json:"id" gorm:"primary_key"`
	Name   string  `json:"name" gorm:"type:varchar(128); not null"`
	Credit float32 `json:"credit" gorm:"not null"`
	Hours  uint    `json:"hours"  gorm:"not null"`
	State  bool    `json:"state" gorm:"default:'true'"`

	ModuleID   uint `json:"module_id"`
	SemesterID uint `json:"semester_id"`
}
```

### Campos
- **Name()** Campo para registrar el nombre de la unidad
- **Credit()** Campo para registrar el numero de créditos que equivale de la unidad
- **Hours()** Campo para registrar el numero de horas de avance de la unidad
- **State()** Campo para registrar el estado de la unidad
- **ModuleID()** ID del modulo a la que pertenece esta unidad
- **SemesterID()** ID del semestre a la que pertenece este unidad


