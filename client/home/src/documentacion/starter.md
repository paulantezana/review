---
title: "Documentación"
date: "2017-08-10"
---

## Introduccion
🚀 SIntitucional es un sistema desarrollado para IESTP públicas que pretende automatizar 
la mayoría de los trabajos que se realizan en una IESTP iniciando desde la admisión de 
los alumnos pasando por el proceso de exámenes de admisión, matriculación y haciendo un 
seguimiento durante todo el trayecto que lleva el estudiante en la IESTP en su formación 
profesional, además permite hacer un seguimiento a los egresados de dicha IESTP.

El sistema se ha dividido en 7 grandes módulos cada una de ella con un propósito 
especifico interconectados entre si con un servidor centralizado en 
[Microsoft Azure](https://azure.microsoft.com/es-es) – mas adelante en [DigitalOcean](https://www.digitalocean.com/)

## Modularización.
En javascript el patron modular emula el concepto de clases, de manera que somos capaces de incluir métodos públicos/privados y propiedades dentro de un único objeto, protegiendo las datos particulares del ámbito global, lo que ayuda a evitar la colisión de nombres de funciones y variables ya definidas a lo largo de nuestro proyecto, o API’s de terceros, a continuación unos conceptos previos para poder entender mejor el patrón modular.

## Objeto literal
EL patron modular se basa en parte en los objetos literales por ende es importante entenderlo.
Un objeto literal es descrito como cero o más pares nombre/valor, separados por comas entre llaves.
Los nombres dentro del objeto pueden ser cadenas o identificadores que son seguidas por 2 puntos, dichos objetos también pueden contener otros objetos y funciones.

```javascript
let objetoLiteral = {
    /* los objetos literales pueden contener propiedades y métodos */
    saludo : "soy un objeto literal",
    miFuncion : function(){
      // código
    }
};
/* accediendo a una propiedad de nuestro objeto literal persona */
objetoLiteral.saludo
```
Un ejemplo de un modulo usando un objeto literal.

```javascript
var persona = {
    /* definiendo propiedades */
    nombre : "adan",
    edad   : 33,
    /* método simple */
    comer  : function(){
        console.log(this.nombre + " esta comiendo.");
    }
};
/* accediendo al método comer de nuestro objeto literal persona */
persona.comer();
```

## Módulo
Un módulo es una unidad independiente funcional que forma parte de la estructura de una aplicación.
Podemos usar funciones y closures(cierres) para crear módulos.

```javascript
let modulo = (function(){
    //- - -
});
```
Un ejemplo más completo:

```javascript
var automovil = (function(colorDeAuto){
    var color = colorDeAuto;
    return{
        avanzar : function(){
            console.log("el auto "+ color +" esta avanzando");
        },

        retroceder : function(){
            console.log("el auto "+ color +" esta retrocediendo");
        }
    }
})("azul");
/* accediendo los metodos retroceder y avanzar de nuestro módulo */
automovil.retroceder();
automovil.avanzar();
```

[**Seguir Leendo**](https://frontendlabs.io/2643--patron-modular-en-javascript)
