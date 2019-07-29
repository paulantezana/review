---
title: "Modelo de la base de datos de sistema de mensajería(Chat)."
date: "2019-27-02"
---

## Group
Se usa para almacenar los grupos del chat
```go
type Group struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Date     time.Time `json:"date"`
	IsActive bool      `json:"is_active" gorm:"default:'true'"`

	UserGroups []UserGroup `json:"user_groups"`
}
```
### Campos
- **Name(requerido)** Campo para registrar el nombre del grupo
- **Avatar(opcional)** Campo para registrar la `URL` del avatar del grupo
- **Date(requerido)** Campo para registrar la fecha en la que se creo el grupo de chat
- **IsActive(opcional)** Campo para indicar si el grupo está activo o no
	+ `false`	: Grupo desactivado.
	+ `true`	: Grupo activo
	+ Valor por defecto `true`
- **UserGroups(opcional)** Referencia a los integrantes del grupo de chat


## UserGroup
Se usa para almacenar todos los integrantes de un grupo.

```go
type UserGroup struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Date     time.Time `json:"date"`
	IsActive bool      `json:"is_active" gorm:"default:'true'"`
	IsAdmin  bool      `json:"is_admin"`

	UserID  uint `json:"user_id"`
	GroupID uint `json:"group_id"`
}
```

### Campos
- **Date(requerido)** : Campo para registrar la fecha en la que el usuario se integró al grupo.
- **IsActive(requerido)** : Campo para indicar si el usuario este activo en el grupo.
- **IsAdmin(opcional)** : Campo para indicar el admin del grupo el cual tiene mas privilegios que el resto de los integrantes del grupo.
- **UserID(requerido)** : ID del usuario quien se está integrando a grupo.
- **GroupID(requerido)** : ID del grupo a la que perteneces este grupo de integrantes.

## Message
Se usa para almacenar todos los mensajes que un usuario enviar a cualquier otro usuario ya sea una sola o un grupo.

```go
type Message struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	Subject          string    `json:"subject"`
	Body             string    `json:"body"`
	BodyType         uint8     `json:"body_type"` // 0 = plain string || 1 == file
	FilePath         string    `json:"file_path"`
	Date             time.Time `json:"date"`
	ExpiryDate       uint      `json:"expiry_date"`
	IsReminder       bool      `json:"is_reminder"`
	NextReminderDate time.Time `json:"next_reminder_date"`

	CreatorID           uint `json:"creator_id"`
	ReminderFrequencyID uint `json:"reminder_frequency_id"`
}
```

### Campos
- **Subject(opcinal)** : Campo para registrar el asunto del mensaje que se envió.
- **Body(requerido)** : Campo para registrar el contenido del mensaje que se envió.
- **BodyType(opcinal)** : Campo para registrar el tipo de mensaje.
	+ `0` : mensaje de tipo texto o parrafo comun
	+ `1` : de tipo multimedi
	+ Valor por defecto `0`
- **FilePath(opcinal)** : Campo para registrar la ruta de un archivo siempre i cuando cunado el tipo de archivo sea de tipo `1`.
- **Date(requerido)** : Campo para registrar en la fecha en la que se creó el mensaje
- **ExpiryDate(opcinal)** : Campo para registrar fecha de expiración de un recordatorio.
- **IsReminder(opcinal)** : Campo para rindicar si el mensaje es un recordatorio.
- **NextReminderDate(opcinal)** : Campo para indicar la siguiente fecha en la que el recordatorio lanzara una nueva notificación del mensaje.
- **CreatorID(requerido)** : ID del usuario que esta enviando el nuevo mensaje.
- **ReminderFrequencyID(opcinal)** : ID de frecuentas preestablecidas en el sistema.

## MessageRecipient
Es la tabla más importante en este modelo de datos. Todo el modelo de datos gira en torno a esta tabla solamente. Uno de los principales objetivos detrás de la creación de esta tabla es mantener la asignación entre los mensajes y sus destinatarios. Por lo tanto, la columna `RecipientID` de esta tabla significa los identificadores de los destinatarios, y esta columna se refiere a la columna identificada de la tabla de usuarios.
```go
type MessageRecipient struct {
	ID     uint `json:"id" gorm:"primary_key"`
	IsRead bool `json:"is_read"`

	RecipientID      uint `json:"recipient_id"`
	RecipientGroupID uint `json:"recipient_group_id"`
	MessageID        uint `json:"message_id"`
}
```
### Campos
- **IsRead(opcinal)** : Campo que indica si un usuario ya leo el mensaje que le enviaron otros usuarios.
- **RecipientID(opcinal)** : ID del destinatario en este caso es el ID del usuario
- **RecipientGroupID(opcinal)** : ID del grupo a la que se está enviando el mensaje
- **MessageID(requerido)** : ID del mensaje a la que perteneces.