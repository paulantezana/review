---
title: "Modelos de la base de datos del sistema de seguimiento de egresados."
date: "2019-27-02"
---
## Poll
Se usa para almacenar las encuestas con diferentes configuraciones cada uno con un propósito especifico.
```go
type Poll struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Message     string    `json:"message"`
	Weather     bool      `json:"weather"` //definite / undefined
	State       bool      `json:"state" gorm:"default:'true'"`

	ProgramID uint `json:"program_id"`

	Questions []Question `json:"questions, omitempty"`
}
```

### Campos
- **Name(requerido)** : Campo para registrar el nombre de la encuesta con la que se usara como título en las interfaces del cuestionario.
- **Description(requerido)** : Campo para registrar la descripción de la encuesta o también puede ingresar las instrucción o tutorial para que los alumnos puedan orientarse a la hora de responder las preguntas que se les planteo.
- **StartDate(opcional)** : Campo para registrar la fecha de inicio en la que se apertura el cuestionario para los alumnos.
- **EndDate(opcional)** : Campo para registrar la fecha en que se cerrar la encuesta por lo tanto pasando esta fecha los alumnos ya no podrán enviar sus respuestas.
- **Message(opcional)** : Campo para registra el mensaje de agradecimiento o respuesta cuando un estudiante ha resultado y enviar el cuestionario correspondiente.
- **Weather(opcional)** : Campo para definir el tiempo
	+ `false`	: sin definir tiempo
	+ `false`	: tiempo definido
	+ Valor por defecto `false`
- **State(opcional)** : Estado de la encuesta
- **ProgramID(requerido)** : ID del programa de estudios a la que pertenece o está dirigido esta encuesta.

## Question
Se usa para almacenar todas las preguntas de un cuestionario.

```go
type Question struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Position uint   `json:"position"`
	State    bool   `json:"state" gorm:"default:'true'"`

	TypeQuestionID uint `json:"type_question_id"`
	PollID         uint `json:"poll_id"`

	MultipleQuestions []MultipleQuestion `json:"multiple_questions"`
	AnswerDetails     []AnswerDetail     `json:"answer_details"`
}
```

### Campos
- **Name(requerido)** : Campo para registrar la pregunta en cuestión que se desea plantear al alumno
- **Required(opcional)** : Campo para especificar si una pregunta es requerida 
	+ `false`	: opcional
	+ `true`	: requerido
	+ Valor por defecto `false`
- **Position(requerido)** : Campo par indicar el orden de las preguntas con este campo podrá ordenar las preguntas según su necesidad.
- **State(opcional)** : Campo para especificar el estado de la pregunta
- **TypeQuestionID(requerido)** : ID de tipo de pregunta que se le está planteando al alumno.
- **PollID(requerido)** : ID del cuestionario a la que pertenece esta pregunta.
- **MultipleQuestions(opcional)** : Cuando se elije el tipo de pregunta las posibles repuestas pueden cambiar como es el caso para marcar una alternativa, puede a ver mas de una alternativa o 10, 20, etc. alternativas ha responder pude ser para responder una sola o varias alternativas según el tipo de pregunta que se escogió anteriormente.
- **AnswerDetails(opcional)** : Referencia a la lista de respuesta que tiene la pregunta por cada alumno y en cada programa de estudios.


## TypeQuestion
Se usa para almacenar el tipo de preguntas que se puede plantear a un alumno según cada tipo de pregunta se almacenara las respuestas y el análisis para las estadísticas de cada una.

```go
type TypeQuestion struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
```
### Campos
- **Name(requerido)** : Campo para registrar el nombre del tipo de pregunta

### Datos
Debido a que la interfaz, el comportamiento, análisis, almacenamiento de ca tipo de pregunta es diferente el sistema por defecto soporta los siguientes tipos de preguntas.

- **Name(requerido)** : 
	1. Respuesta breve: Respuesta corta de tipo texto que tiene menor a 5 palabras aproximadamente.
	2. Párrafo: Repuesta de tipo párrafo en la que el análisis de los datos que realizara manualmente por el usuario.
	3. Opción múltiple: Respuesta de opción múltiple para elegir una opción en especifica como marcar verdadero o false o preguntas de ese estilo.
	4. Casillas de verificación: Respuesta de casillas para responder varias alternativas que se planteó.

## MultipleQuestion
Se usa para almacenar las preguntas de tipo múltiple siempre y cuando el usuario haya elegido este tipo de pregunta.

```go
type MultipleQuestion struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Label string `json:"label"`
	State bool   `json:"state" gorm:"default:'true'"`

	QuestionID uint `json:"question_id"`
}
```

### Campos
- **Label(requerido)** : Campo para registrar el nombre la posible respuesta.
- **State(opcional)** : Estado de la posible respuesta.
- **QuestionID(requerido)** : ID de la pregunta a la que pertenece esta posible respuesta.


## Answer
Se usa para almacenar las encuestas respondidas por los alumnos de los diferentes programas de estudio.
```go
type Answer struct {
	ID    uint `json:"id" gorm:"primary_key"`
	State bool `json:"state" gorm:"default:'true'"`

	StudentID uint `json:"student_id"`
	PollID    uint `json:"poll_id"`

	AnswerDetails []AnswerDetail `json:"answer_details, omitempty"`
}
```

### Campos
- **State(opcional)** : Campo para especificar el estado en la que se encuentra la respuesta.
- **StudentID(requerido)** : ID del estudiante quien resolvió la encuesta
- **PollID(requerido)** : ID de la encuesta a la que pertenece estas respuestas.
- **AnswerDetails(opcional)** : Referencia a los detalles de las respuestas

## AnswerDetail
Se usa para almacenar todas las respuestas de cada pregunta sea del tipo que sea esta es la tabla mas importante a la hora de evaluar y analizar las respuestas de cada alumno por pregunta.
```go
type AnswerDetail struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Answer string `json:"answer"`
	State  bool   `json:"state" gorm:"default:'true'"`

	QuestionID uint `json:"question_id"`
	AnswerID   uint `json:"answer_id"`
}
```
### Campos
- **Answer(requerido)** : Campo para almacenar la respuesta de cualquier tipo
- **State(requerido)** : Campo para indicar el estado de la respuesta
- **QuestionID(requerido)** : ID de la pregunta a la que corresponde la respuesta
- **AnswerID(requerido)** : ID del cuestionario a la que pertenece esta respuesta