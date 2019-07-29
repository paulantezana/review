---
title: "Tema de WordPress para IESTP"
date: "2019-27-02"
---

sninstitute: Es el nombre del tema de [WordPress] el que se desarrolló según a las necesidades de una IESTP.

## Src
Es el directorio donde tendremos los archivos del proyecto en fase de desarrollo, y se estructura de la siguiente manera:
* 📁 **scss** : Contiene los archivos **.scss** que compilarán a archivos CSS.
    + 📁 **components** : Contiene los archivos partials de la presentación **SCSS** o CSS de los componentes.
* 📁 **scripts** : Contiene los archivos JS que serán compilados con Babel y unificados con Browserify.
    + 📁 **components** : Contiene los archivos de la programación JS de los componentes.
    + 📁 **helpers** :  Contiene los archivos de la programación JS de códigos auxiliares que no sean componentes, como la conexión a una API, funciones para formatear o filtrar contenido, etc.
    + 📄 **admin.js** : Es el archivo JS del que permite modificar el admin de [WordPress].
    + 📄 **app.js** : Es el archivo principal JS del tema y del sitio web, en el que se podrá importar los componentes que se requieran de la carpeta components, helpers o de las dependencias que se tenga en node_modules.
    + 📄 **editor.js** : Es el archivo JS del editor de [WordPress] que viene integrada por defecto, en el que se podrá importar los componentes que se requieran de la carpeta components, helpers o de las dependencias que se tenga en node_modules.
    los JS de este archivo son para mejorar la experiencia de usuario para que el editor sea igual el formato del contenido que se mostrara en el sitio web y además agregar funciones adicionales al editor.
    + 📄 **login.js** : Es el archivo JS que permite cambiar las características del **Login** de sitio [WordPress]
* 📁 **raw** : Contiene las imágenes del proyecto sin optimizar.

## Assets
Es el directorio donde tendremos la versión para indexar a los scripts a [WordPress]
## Page Templates
Las plantillas de página se utilizan para cambiar el aspecto de una página.

Las plantillas de página muestran el contenido dinámico del sitio en una página, por ejemplo, publicaciones, actualizaciones de noticias, eventos de calendario, archivos multimedia, etc. Puede decidir que desea que su página de inicio se vea de una manera específica, que es muy diferente a otras partes de su sitio.

Este tema cuenta con 4 plantillas
* 📄 **page-full.php** : *Ancho completo sin sidebar*

┌────────────┐
│            │
│            │
│            │
│            │
└────────────┘

* 📄 **page-full-width.php** : *Página sin sidebar*
* 📄 **page-no-title.php** : *Ancho completo sin sidebar sin titulo*
* 📄 **page-sidebar-left.php** : *Página con sidebar a la izquierda*
```bash
```

## Inc

[WordPress]:https://wordpress.org/