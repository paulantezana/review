---
title: "Logica de negocio"
date: "2019-27-02"
---

## Procesos de admisión
En este proceso pues nos encontramos con un típico CRUD de crear, leer, modificar, eliminar ya algunas funciones adicionales.
* **GetAdmissionSettings** : Función que permite realizar una consulta a todos los registros existentes en la base de datos
* **GetAdmissionSettingByID** : Función que permite realizar una consulta detallada sobre un registro en específico pasando como parámetro el ID del registro que se desea consultar.
* **CreateAdmissionSetting** : Función que permite crear un nuevo registro en la base de datos pasando
* **UpdateAdmissionSetting** : Función que permite actualizar los datos de registro en especifico en la base de datos
* **DeleteAdmissionSetting** : Función que permite eliminar un registro en específico
* **ShowInWebAdmissionSetting** : 🔥 Función que permite indicar cual de los procesos de admisión se mostrara por defecto en la web

## Admisión
Es una de las partes mas importantes del proceso de admisión pues en este es donde se realiza el gran porcentaje de cálculos que se requieren para registrar, actualizar, reportar de un proceso de admisión
* **GetAdmissionsPaginate** : Permite realizar una paginación de todos los alumnos que se registraron en un proceso de admisión, filtrando los datos por cada proceso de admisión que se apertura.
* **GetAdmissionsByID** : Permite realizar una consultar todos los detalles de un alumno que esta en el proceso de admisión.
* **GetAdmissionsPaginateExam** : Permite realizar una paginación de todos los alumnos que se registraron en un proceso de admisión, filtrando los datos por cada proceso de admisión que se apertura y que no este anulado de esta manera se podra ingresar las notas del examen de cada alumno.
* **UpdateStudentAdmission** : Función que permite actualizar los datos del estudiante cuando este ya existe en la base de datos caso contrario creara un nuevo registro del estudiante.
    - Si existe: actualizara los datos del alumno que se pasó como parámetro.
    - No existe: Creara un nuevo alumno con todos los datos correspondientes y además creara una nueva cuenta de usuario en el sistema con los siguientes valores: usuario = `DNIST`, contraseña = `DNIST`. y obviamente con el rol estudiante para hacer la lógica de las restricciones en el sistema.
    - El estado del alumno cambia a ` StudentStatusID = 1` que corresponde aun alumno no asignado, y también se crea un nuevo historial indicando que los datos del alumno has sido modificados o creados en un proceso de admisión.
