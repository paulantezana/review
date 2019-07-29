---
title: "Estructura de archivos del servidor"
date: "2017-08-10"
---

## Carpetas:
* 📁 **config** : Paquete que contiene todas las configuraciones del sistema como la base de datos, versiones, etc.
    - 📄 **config.go** : Se encarga de leer el archivo `config.json` y mapearlo en una estructura que pueda ser utilizado desde cualquier parte de la aplicación
    - 📄 **database.go** : Contiene la unica funcion `func GetConnection() *gorm.DB` Esta es una función que se encarga de conectar a la base de datos y devuelve un puntero de la conexión que podrá ser usado desde cualquier parte del sistema.
    Puede cambiar fácilmente el motor de la base de datos ya sea `MySql`,  `MongoDB`  y entre otros ya que en este sistema se esta haciendo uso de un orm denominado [gorm](http://gorm.io)
    Si desea cambiar el nombre de la base de datos, usuario, contraseña puede hacerlo desde al archivo `config.json` que se encuentra en la carpeta raíz de este sistema.
    - 📄 **email.go** : Contiene la unica funcion `func SendEmail(to string, subject string, tem string) error` En este archivo se encuentra la configuración para enviar correos electrónicos mediante una cuenta de Gmail. 
    Cada ves que se quiera enviar un correo electrónico deberá llmarar esta funcion
* 📁 **controller** : **CORE** En este paquete se encuentra una de las partes más importantes del sistema pues se encarga de controlar los datos, realizar consultas a la base de datos y devolver en formato json, xml, etc. al usuario que realizo la petición
Pues el controlador esta estrechamente relacionada con las rutas del API SERVICE de este sistema
    - 📄 ... Ingrese a cada modulo del sistema para ver los detalles de cada archivo
* 📁 **endpoint** : Es en este paquete es donde se encuentra los canales de comunicación del la API y Websocket
    - 📄 **api.go** : Es aquí donde reside todas las rutas del api, tenemos tres tipos de rutas
    - 📄 **ws.go** : Es aquí donde reside todas las rutas de los Websockets del sistema como del chat, comentarios, notificaciones, etc. que estas comunicaciones son en tiempo real.
* 📁 **migration** : Paquete encargado de las migraciones de la base de datos
    - 📄 **migration.go** Básicamente se encarga de crear las tablas en la base de datos usando los modelos que se encuentran en la carpeta models además ingresas los datos iniciales del sistema como el usuario, configuracion
* 📁 **models** : Paquete donde se encuentras los modelos de la base de datos que se podrán migrar usando el maquete `migration`
    - 📄 ... Ingrese a cada modulo del sistema para ver los detalles de cada archivo
* 📁 **static** : En esta carpeta se almacenan todos los archivos estáticos que el cliente necesita como logo, fotos de perfil, pdfs, etc.
    - 📁 **apps** : Esta carpeta se usa para almacenar los logos de los diferentes sistemas que se despliegan usan este api servicie.
    - 📁 **books** : Esta carpeta se usa para almacenar todos los archivos que genera el sistema de biblioteca como los pdfs, portas ye entre otros
    - 📁 **chat** : Esta carpeta se usa para guardar todos los medios que genera el sistema de chat
    - 📁 **profile** : Esta carpeta se usa para almacenar las fotos de perfil cada usuario
    - 📄 **nationalEmblem.jpg** : Escudo nacional del Perú
    - 📄 **ministry.jpg** : Encabezado del ministerio formato largo
    - 📄 **ministrySmall.jpg** : Encabezado del ministerio formato corto
    - 📄 **logo.jpg** : Logo por defecto de la institución
    - 📄 **book.jpg** : Portada por defecto de los libros del sistema de biblioteca
    - 📄 **data-set.min.js** : Archivo de JavaScript que los clientes necesitan para generar gráficos estadísticos.
* 📁 **temp** : Esta carpeta se usa para almacenar todos los archivos temporales que el sistema genera de forma dinámica.
Es recomendable eliminar todo el contenido de esta carpeta para liberar espacio en la memoria
**No eliminar la carpeta solo el contenido**
* 📁 **templates** : En este paquete se encuentra las plantillas de HTML, Excel, etc. Que son usadas por el sistema y el usuario final
    - 📄 **email.html** : Esta plantilla se usa para enviar los correos electrónicos de recuperación de contraseña de un usuario.
    - 📄 **templateCompany.xlsx** : Esta plantilla que facilitar al usuario para subir de forma masiva las empresas desde un archivo Excel.
    - 📄 **templateCourseStudent.xlsx** : Esta plantilla que facilitar al usuario para subir de forma masiva los alumnos de los cursos del sistema de certificación desde un archivo Excel.
    - 📄 **templateStudent.xlsx** : Esta plantilla que facilitar al usuario para subir de forma masiva los alumnos de un programa de estudios en específico es útil para los coordinadores de un programa de estudios.
    - 📄 **templateStudentSA.xlsx** : Esta plantilla que facilitar al usuario para subir de forma masiva los alumnos de todos los programas de estudios.
    - 📄 **templateTeacher.xlsx** : Esta plantilla que facilitar al usuario para subir de forma masiva los profesores de un programa de estudios en específico es útil para los coordinadores de un programa de estudios.
    - 📄 **templateTeacherSubsidiary.xlsx** : Esta plantilla que facilitar al usuario para subir de forma masiva los profesores de todos los programas de estudios.
* 📁 **utilities** : En este paquete se encuentran todos los utilitarios que el sistema usa en diferentes procesos como para recibir y enviar datos, paginación, etc.
    - 📄 **counter.go** : Función para hacer conteos
    - 📄 **notice.go** : Estructura que sirve para enviar las notificaciones
    - 📄 **request.go** : Estructura que sirve como una plantilla para recibir datos desde el cliente y también calcula la paginación para cada consulta.
    - 📄 **response.go** : Estructura que sirve como una plantilla para enviar datos al cliente.
    - 📄 **token.go** : Permite firmar los claves JWT de un usuario con vigencia de 24 horas usando el método HS256

## Archivos:
* 📄 **.editorconfig**: Contiene la definición de la configuración para mantener la codificación estándar entre diferentes editores e IDEs, considera que en algunos editores tendrás que instalar un plugin adicional para que funcione, consulta el sitio [editorconfig.org](http://editorconfig.org/) para saber si tu editor o ide lo soporta nativamente o requiere algún plugin.
* 📄 **.gitignore**: Indica que archivos y directorios ignorará Git al momento de sincronizar el proyecto, la configuración que se propone ha sido generada en el sitio [gitignore.io](https://www.gitignore.io/) y es esta: [osx,node,macos,linux,windows,visualstudiocode](https://www.gitignore.io/api/osx,node,macos,linux,windows,visualstudiocode) siéntete libre de modificarla a tus necesidades.
* 📄 **main.go** : **Es el archivo que inicia todo el sistema.**
A ejecutar este archivo se ejecutará todo el sistema y también si se Desa compilar el sistema para diferentes plataformas como Windows, Linux y Mac se debe hacer referencia a este archivo
Para ejecutar en tu sistema operativo actual puedes ejecutar con el siguiente comando
```go
go run main.go
```
Compilar para linux desde windows
```bash
GOOS=linux GOARCH=amd64 go build -o instituteL main.go
```
Compilar para mac desde windows
```bash
GOOS=darwin GOARCH=amd64 go build -o instituteM main.go
```
Si aun no esta familiarizado con el lenguaje de programación GO puedes ingresar al [siguiente enlace](https://golang.org/doc/code.html)