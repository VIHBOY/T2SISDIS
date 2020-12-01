# Tarea 2 Sistemas Dsitribuidos


## Comenzando 🚀

_Estas instrucciones te permitirán obtener una copia del proyecto en funcionamiento en tus máquinas virtuales para propósitos de desarrollo y pruebas._

### Pre-requisitos 📋

_Que cosas necesitas para instalar el software y como instalarlas_

```
Las maquinas virtuales ya tienen todo lo necesario para realizar pruebas.
```

### Pasos a Seguir 🔧

_Se debe seguir estos pasos para poder ejecutar la tarea_

```
1. En la Maquina dist28, se debe ejecutar make

Lo que le otorgara a la maquina el rol de NameNode

2. En la Maquina dist27, se debe ejecutar make

Lo que le otorgara a la maquina el rol de DataNode3

3. En la Maquina dist26, se debe ejecutar make

Lo que le otorgara a la maquina el rol de Datanode2

4. En la Maquina dist27, se debe ejecutar make

Lo que le otorgara a la maquina el rol de DataNode1, esta posee tanto cliente uploader como downloader.
```
Existe un  archivo llamado clienteseparado.go el cual posee los clientes uploader y downloader para probar solicitudes al mismo tiempo

## Consideraciones NameNode ⚙️

En todo momento se debe mantener encendido el servidor.

## Consideraciones Datanode3 ⚙️

La forma en que esta implementado el aspecto de Datanode3 permite que esta este apagada y siga funcionando el procedimiento de la tarea.

## Consideraciones Datanode2 ⚙️

La forma en que esta implementado el aspecto de Datanode2 permite que esta este apagada y siga funcionando el procedimiento de la tarea.

## Consideraciones Datanode1 ⚙️

En todo momento se debe mantener encendido el servidor, dado que esta posee los clientes de uploader y downloader.


Aclaraciones.


Existe un menú en Datanode1 el cual pregunta por si desea subir o bajar un archivo, cabe recalcar que como se trata de libros, estos deben estar en formato pdf y no tener guion bajos en su nombre, el archivo debe estar en la maquina 25 para ser subido, y este será bajado con el nombre con un 1 añadido en su titulo.

Como se indico antes, existe un archivo llamado clienteseparado.go en la maquina 25 el cual se puede utilizar para probar conexiones simultaneas, cabe decir que el datanode1 no puede caerse (el datanode 2 y el 3 si) es una de las pocas indicaciones que damos.



## Construido con

[VS CODE] - Editor de texto

## Autores

**Joaquin Concha** - 201773569-4 *VIHBOY*
**Renato Bassi** - 201773521-k *bassisi*
