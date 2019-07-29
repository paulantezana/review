---
title: "Instalar Servidor"
date: "2019-27-02"
---

## Usuario
Crear un nuevo usuario
```bash
$ adduser paul
```

Añadir el usuario al grupo sudo
```bash
$ gpasswd -a paul sudo
```
Crear carpeta .ssh
```bash
$ mkdir /home/paul/.ssh
$ chmod 700 /home/mynewuser/.ssh
```

Crear archivo de claves autorizadas
```bash
$ touch /home/paul/.ssh/authorized_keys
```

A partir de ahí, iniciará sesión como usuario y creará su clave SSH. Por lo general, uso una tecla más pesada con más rondas de KDF, aunque puede demorar el inicio de sesión de unos segundos a minutos, dependiendo de la cantidad de rondas de KDF que use.

Por ejemplo, para generar una clave RSA, usaría:

```bash
$ ssh-keygen -a 1000 -b 4096 -C "" -E sha256 -o -t rsa
```
Para una clave ED25519, usaría:
```bash
$ ssh-keygen -a 1000 -C "" -E sha256 -o -t ed25519
```

`-a-` KDF Rounds (función de derivación de clave).

`-b-` Tamaño de bit (se aplica a RSA, pero no a ED25519).

`-C-` Establece que el comentario de la clave esté en blanco.

`-e-` Establece el hash de clave utilizado (sha256 es el predeterminado).

`-o-` Utiliza el nuevo formato OpenSSH para las claves.

`-t-` Especifica El tipo de clave (RSA / ED25519).


Con 1,000 rondas de KDF, la clave tarda unos segundos en generarse cuando se usa una frase de contraseña, y también demorará unos segundos en iniciar sesión. El uso de KDF genera una clave más segura, aunque hay que tener cuidado ya que establecerlo demasiado alto definitivamente causará retrasos graves al intentar iniciar sesión (es decir, 20,000 rondas tomarán un promedio de 2-4 minutos para generar y lo mismo para iniciar sesión ).

Una vez que se genere su clave pública / privada, coloque la clave pública en:

```bash
/home/paul/.ssh/authorized_keys
```

## Servidor
Para poder instalar es necesario tener los permisos del autor del código fuente del sistema que puede descargar desde [github]( https://github.com/paulantezana/review).

El sistema está desarrollado en el lenguaje de programación GO que es un lenguaje de programación concurrente y compilado inspirado en la sintaxis de C. Ha sido desarrollado por Google. 

De esta manera puede compilar en binarios para cualquier sistema operativo ya sea Windows, Linux, Mac.

## 1 compilar
* Para linux desde windows
```bash
GOOS=linux GOARCH=amd64 go build -o instituteL main.go
```

* Para mac desde windows
```bash
GOOS=darwin GOARCH=amd64 go build -o instituteM main.go
```
en este caso conpilaremos para linux ya que tengo configurado un servidor Ubuntu 

posteriormente subiremos los archivos al servidor

## 2 Subir al servidor
Tenemos varios métodos para subir archivos al servidor como [filezilla](https://filezilla-project.org/) en este caso usaremos comando scp que está disponible tanto en Linux como en Mac. Si estás usando Windows puede descargarse [GitBash](https://git-scm.com/downloads) para usar este comando.

Subir archivo
comando: scp miApp usuario@ip directorio
```bash
$ scp -i KeyPairSednaServer.pem instituteL ubuntu@13.68.218.250:/home/ubuntu/instituteapp/institute_server
instituteL                              100%   14MB  47.8KB/s   04:59
```

Subir el resto de carpetas y los archivos necesario para el funcionamiento del sistema
```bash
$ scp -i KeyPairSednaServer.pem -rp static/ ubuntu@13.68.218.250:/home/ubuntu/instituteapp/institute_server
$ scp -i KeyPairSednaServer.pem -rp temp/ ubuntu@13.68.218.250:/home/ubuntu/instituteapp/institute_server
$ scp -i KeyPairSednaServer.pem -rp templates/ ubuntu@13.68.218.250:/home/ubuntu/instituteapp/institute_server
$ scp -i KeyPairSednaServer.pem config.json  ubuntu@13.68.218.250:/home/ubuntu/instituteapp/institute_server
```

## 3 Cambiar los permisos
es recomendable cambiar los permisos a 775 de los siguientes archivos
```bash
$ chmod 775 instituteL
$ chmod 775 templates/
$ chmod 775 static/
$ chmod 775 temp/
```

## 4 Base de datos
### Instalacion
El sistema usa una base de datos Postgres por defecto, puede cambiar el motor de la base de datos desde el codigo fuente

Debido a que es nuestra primera vez utilizando apt en esta sesión, debemos refrescar nuestro índice de paquetes locales. Podemos instalar el paquete Postgres
```bash
$ sudo apt-get update
$ sudo apt-get install postgresql
```

Editar el archivo de configuración pg_hba.conf
```bash
$ sudo vim /etc/postgresql/10/main/pg_hba.conf
```
Actualice la parte inferior del archivo, debe quedar de la siguiente forma
```vim
# Database administrative login by Unix domain socket
local   all             postgres                                trust

# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     peer
# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
# IPv6 local connections:
host    all             all             ::1/128                 md5
```

Edite el siguiente archivo - postgresql.conf
```
$ sudo vim /etc/postgresql/10/main/postgresql.conf
```

Descomenta y edita las siguientes líneas -
```vim
listen_addresses = 'localhost'
```

Reinicial postgres
```bash
$ sudo /etc/init.d/postgresql restart
```

Ejecuta el siguiente comando para ingresar a la base de datos
```bash
$ psql -U postgres
```

### Crear un Nuevo Rol
Por defecto, Postgress utiliza un concepto llamado "roles" que maneja identificación y autorización. Estos son, de algún modo, similares a los estilos de cuentas en Unix, pero Postgres no distingue entre usuarios y grupos y en su lugar prefiere ser más flexible con el término "rol"

Al concluir la instalación Postgres está listo para utilizar la identificación ident, lo que significa que asocia los roles de Postgres con una cuenta de sistema Unix/Linux. Si el rol existe en Postres, un nombre de usuario Unix/Linux con el mismo nombre podrá identificarse como ese rol.

para realizar esta operación primero debe ingresar a la base de datos con el usuario postgres
```sql
CREATE ROLE institute_user LOGIN PASSWORD 'newright789'
```

### Crear la base de datos
De forma predeterminada, otra suposición que hace el sistema de autenticación de Postgres es que habrá una base de datos con el mismo nombre que el rol que se utiliza para iniciar la sesión, a la que el rol tiene acceso.

en este caso usaremo el usuario institute_user que creamos anteriormente.
```sql
CREATE DATABASE institute OWNER institute_user
```

## Test
ejecute el binario que subio al servidor
```bash
$ ./instituteL
(D:/golang/src/github.com/paulantezana/review/migration/migration.go:15)
[2019-03-03 19:21:53]  [3.90ms]  CREATE TABLE "message_recipients" ("id" serial,"created_at" timestamp with time zone,"updated_at" timestamp with time zone,"is_read" boolean,"recipient_id" integer,"recipient_group_id" integer,"message_id" integer , PRIMARY KEY ("id"))
[0 rows affected or returned ]

(D:/golang/src/github.com/paulantezana/review/migration/migration.go:15)
[2019-03-03 19:21:53]  [6.40ms]  CREATE TABLE "reminder_frequencies" ("id" serial,"created_at" timestamp with time zone,"updated_at" timestamp with time zone,"name" text,"frequency" numeric,"is_active" boolean , PRIMARY KEY ("id"))
[0 rows affected or returned ]

(D:/golang/src/github.com/paulantezana/review/migration/migration.go:15)
[2019-03-03 19:21:53]  [6.15ms]  CREATE TABLE "sessions" ("id" serial,"created_at" timestamp with time zone,"updated_at" timestamp with time zone,"ip_address" text,"user_id" integer,"last_activity" timestamp with time zone , PRIMARY KEY ("id"))
[0 rows affected or returned ]

(D:/golang/src/github.com/paulantezana/review/migration/migration.go:15)
[2019-03-03 19:21:53]  [3.92ms]  CREATE TABLE "user_groups" ("id" serial,"created_at" timestamp with time zone,"updated_at" timestamp with time zone,"date" timestamp with time zone,"is_active" boolean DEFAULT 'true',"is_admin" boolean,"user_id" integer,"group_id" integer , PRIMARY KEY ("id"))
[0 rows affected or returned ]

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.3.dev
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1323
```
como se puede apreciar la aplicacion corre de manera correcta

## 5 Crear servicio
crear un nuevo archivo de servicio
```bash
sudo vim /etc/systemd/system/institute.service
```
configuracion
```vim
[Unit]
Description="API service de GO para el sistema institucional"

[Service]
ExecStart=/home/ubuntu/instituteapp/institute_server/instituteL
WorkingDirectory=/home/ubuntu/instituteapp/institute_server
User=ubuntu
Restart=always

[Install]
WantedBy=multi-user.target
```

activar el servicio
```bash
sudo systemctl enable institute.service
```

iniciar el servicio
```bash
sudo systemctl start institute.service
```

verificar el servicio
```bash
sudo systemctl status institute.service
```

## 6 Proxy
## 7 Letsencrypt