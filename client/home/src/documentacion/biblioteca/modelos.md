---
title: "Modelos del sistema de biblioteca."
date: "2019-27-02"
---

## Category
Se usa para almacenar las categorías de los libros.
Use el `ParentID` para poder registrar categorías de infinitos niveles
```go
type Category struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	ParentID uint   `json:"parent_id"`
	State    bool   `json:"state" gorm:"default:'true'"`

	Books []Book `json:"books, omitempty"`
}
```
### Campos
- **Name(requerido)** : Campo para registrar el nombre de la categoría.
- **ParentID(requerido)** : ID de la categoría padre para crear categorías de niveles infinitos.
- **State(requerido)** : Campo para indicar el estado de la categoría.
- **Books(requerido)** : Referencia a los libros que pertenecen a la categoría.

## Book
Se usa para registrar los libros por categorías
```go
type Book struct {
	ID               uint       `json:"id" gorm:"primary_key"`
	Name             string     `json:"name"`
	ShortDescription string     `json:"short_description"`
	LongDescription  string     `json:"long_description"`
	Author           string     `json:"author"`
	Editorial        string     `json:"editorial"`
	YearEdition      uint       `json:"year_edition"`
	Version          uint       `json:"version"`
	EnableDownload   bool       `json:"enable_download"`
	Avatar           string     `json:"avatar"`
	Pdf              string     `json:"pdf"`
	State            bool       `json:"state" gorm:"default:'true'"`
	Views            uint32     `json:"views"`

	CategoryID uint `json:"category_id"`
	UserID     uint `json:"user_id"`

	Comments []Comment `json:"comments, omitempty"`
	Readings []Reading `json:"readings, omitempty"`
	Likes    []Like    `json:"likes, omitempty"`
}
```

### Campos
- **Name(requerido)** Campo para registrar el nombre o titulo
- **ShortDescription(opcional)** Campo para registrar una descripción corta
- **LongDescription(opcional)** Campo para registrar una descripción larga
- **Author(opcional)** Campo para registrar el autor
- **Editorial(opcional)** Campo para registrar el editorial
- **YearEdition(opcional)** Campo para registrar la edición del año
- **Version(opcional)** Campo para registrar la versión
- **EnableDownload(opcional)** Campo para habilitar la descarga
- **Avatar(opcional)** Campo para registrar la URL de la foto de la portada
- **Pdf(opcional)** Campo para registrar la URL del PDF que se subió
- **State(opcional)** Estado
- **Views(opcional)** Campo para registrar el número de visitas que tiene el libro
- **CategoryID(opcional)** ID de la categoría a la que pertenece el libro
- **UserID(opcional)** ID del usuario quien registro el libro
- **Comments(opcional)** Referencia a la tabla `Comment` de todos los comentarios acerca del libro
- **Readings(opcional)** Referencias a las lecturas de los libros
- **Likes(opcional)** Referencia likes que tiene el libro


## Reading
Se usa para registrar toda las lecturas que realiza un usuario acerca de un libro en especifico
```go
type Reading struct {
	ID   uint      `json:"id" gorm:"primary_key"`
	Date time.Time `json:"date"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
```

### Campos
- **Date(requerido)** Campo para registrar la fecha en la que se realizo la lectura del libro
- **UserID(requerido)** ID del usuario quien realizo la lectura
- **BookID(requerido)** ID del libro a la que se izo la lectura

## Comment
Se usa para registrar los cometarios use el `ParentID` para responde un comentario existente

```go
type Comment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ParentID  uint      `json:"parent_id"` // Parent comment ID
	Votes     uint32    `json:"votes"`
	HasVote   int8      `json:"has_vote" gorm:"-"` // if current user has vote
	Body      string    `json:"body"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`

	User     []User    `json:"user, omitempty"`
	Children []Comment `json:"children, omitempty"`
}
```

### Campos
- **ParentID(requerido)** ID del comentario padre que se usa para responder un comentario existente
- **Votes(requerido)** Campo para registrar el número de votos del comentario
- **HasVote(requerido)** Campo para registrar si el usuario que está tratando de votar ya voto
- **Body(requerido)** Campo para registrar el contenido del comentario
- **UserID(requerido)** ID del usuario que hizo el comentario
- **BookID(requerido)** ID del libro en la que se está haciendo el comentario
- **User(referencia)** Referencia a todos los usuarios que asieron comentarios en el libro
- **Children(referencia)** Referencia a las respuestas de los comentarios


## Like
Se usa para registras los likes o las votaciones que tiene un libro
```go
type Like struct {
	ID    uint  `json:"id" gorm:"primary_key"`
	Stars uint8 `json:"stars"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
```

### Campos
- **Stars(requerido)** Campo para registrar el número de likes o votos acerca de un libro
- **UserID(requerido)** ID del usuario que realizo el voto
- **BookID(requerido)** ID del libro a la que se está realizando la votación


## Vote
Se usa para registrar las votaciones de un comentario que se izo en un libro o una respuesta aun comentario existente
```go
type Vote struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CommentID uint `json:"comment_id" gorm:"not null"`
	UserID    uint `json:"user_id" gorm:"not null"`
	Value     bool `json:"value" gorm:"not null"`
}
```
### Campos
- **CommentID(requerido)** ID del comentario a la que se está realizando la votación
- **UserID(requerido)** ID del usuario quien está realizando la votación
- **Value(requerido)** El valor de la votación
