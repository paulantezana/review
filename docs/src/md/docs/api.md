# API
Debido a que es un sistema centralizado los datos que viajan desde es servidor a los clientes es mediante un servicio web (API SERVICE)
El api servicie esta desarrollado con el lenguaje de programación [GoLang]. usando [JWT]. para la encriptación de los datos que viajan desde el servidor al cliente y viceversa.

## EJemplo de uso
En este ejemplo mostraremos para el proceso de Login y una consulta al webservce de la lista de los alumnos

Para poder conectar al api servicie debe hacer usando la siguiente url
````bash
    http://13.68.218.250:1323/api/v1
````
### Login
Todos los datos en el servidor están protegidos con proceso de LOGIN que este mismo tiene una vigencia de 24 horas, pasado este tiempo deberá renovar sus credenciales para seguir consumiendo los datos del servidor

Para loguearse en el servidor puede hacer desde cualquier tipo de cliente o con cualquier lenguaje de programación.

Lo parámetros que se deben enviar al servidor para loguearse son las siguientes

* URI
    ````bash
    http://13.68.218.250:1323/api/v1/public/user/login
    ````
* Metodo:
    ````bash
    POST
    ````
* Header: 
    ````bash
    'content-type: application/json'
    ````
* Body
    ```json
    {
        "user_name": "sa",
        "password": "sa"
    }
    ```

### Respuesta
Al enviar todos los datos de manera correcta el servidor le responde con todos los datos del usuario entre ellos podemos destacar. token (credenciales para conectarse con una vigencia de 24 hora), roles (todos los roles que se le asigno a este usuario que se loguero), success (si es true significa que los datos de logueo son correctos, pero si es falso, el logue es incorrecto en este caso el servidor no enviara más datos)
````json
{
    "message": "Bienvenido al sistema sa",
    "success": true,
    "data": {
        "user": {
            "id": 1,
            "user_name": "sa",
            "password": "",
            "key": "",
            "state": true,
            "avatar": "",
            "email": "paul.antezana.2@gmail.com",
            "role_id": 1,
            "old_password": "",
            "students": null,
            "teachers": null,
            "reviews": null,
            "coos": null
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJ1c2VyX25hbWUiOiJzYSIsInBhc3N3b3JkIjoiIiwia2V5IjoiIiwic3RhdGUiOnRydWUsImF2YXRhciI6IiIsImVtYWlsIjoicGF1bC5hbnRlemFuYS4yQGdtYWlsLmNvbSIsInJvbGVfaWQiOjEsIm9sZF9wYXNzd29yZCI6IiIsInN0dWRlbnRzIjpudWxsLCJ0ZWFjaGVycyI6bnVsbCwicmV2aWV3cyI6bnVsbCwiY29vcyI6bnVsbH0sImV4cCI6MTU0OTc4OTA4NCwiaXNzIjoicGF1bGFudGV6YW5hIn0.4d6bArYFHNf4bIAcXH9pKuIQhcuRC0v_j3XBEQIj_pg",
        "licenses": {
            "programs": [],
            "subsidiaries": [
                {
                    "id": 2,
                    "name": "COMBAPATA"
                },
                {
                    "id": 3,
                    "name": "PITUMARCA"
                },
                {
                    "id": 1,
                    "name": "SICUANI"
                }
            ]
        }
    }
}
````


[GoLang]:(https://golang.org/)
[JWT]:(https://jwt.io/)