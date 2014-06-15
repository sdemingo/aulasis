<!---
Properties:
@>status:close
-->

# Empezando con Aulasis

Aulasis te permitirá gestionar el flujo de trabajo habitual de un aula
exportando tareas previamente editadas por ti a los alumnos a través de sus
navegadores, y recibir las soluciones de estos enviadas a ti a través de un
sencillo formulario.

Aulasis mantiene la información almacenada a dos niveles en el sistema de
ficheros. En el mismo directorio en el que se encuentra el fichero ejecutable
encontrarás un directorio llamado `courses/`. En el se encuentra toda la
información editable para tus alumnos. Esta información se encuentra recogida a
dos niveles. En un primer nivel de directorios encontramos los **cursos**, bajo
estos encontramos otro nivel de directorios correspondiente a las **actividades**.

```
  courses/
     |--- curso1/
           |--- actividad1/
                 |--- info.md
                 |--- imagen.jpg
           |--- actividad2/
                 |---  info.md
     |--- curso2/
               ....
     |--- meta.xml
```

## Creando o editando un curso

Par crear un curso simplemente hemos de crear un directorio bajo el directorio
`courses/` con un nombre que no exista. Tras esto, simplemente hemos de
registrar su nombre y descripción el el fichero `courses/meta.xml`. Creamos una
nueva etiqueta `course` justo bajo la etiqueta de cierre de la anterior. Esta
etiqueta a su vez contendrá las siguientes etiquetas:

   - **name** : Indica el nombre del curso que aparece en la página principal
   - **path** : Indica el nombre del directorio que acabas de crear
   - **description** : Indica la descripción del curso que aparece en la página principal


## Creando o editando tareas

Para crear una nueva tarea simplemente crearemos un nuevo directorio
dentro del curso en el que queramos situarla. El único fichero
obligatorio para cada tarea es el fichero `info.md`. Este fichero
contiene la descripción de la tarea en formato
[markdown](http://es.wikipedia.org/wiki/Markdown), un sencillo
mecanismo de edición rápida de documentos muy sencillo de
utilizar. 


Dentro del directorio de la actividad podemos meter todos los ficheros
que quieres servir junto con esta: imágenes, ficheros con código,
archivos PDF, etc. Para enlazar los recursos estáticos servidos dentro
de cada carpeta de actividad hemos de tener en cuenta que la ruta de
estos ha de comenzar por `/courses/dirCurso/dirTarea`. Si por ejemplo
quisiéramos visualizar una imagen contenida dentro del directorio de
esta actividad escribiríamos:

```
![Alt text](/courses/inicio/empezando/imagen.jpg)
```

 Los ficheros enviados por los alumnos a través del formulario de entrega
(mostrado solo en tareas abiertas) son almacenados bajo el mismo directorio
donde se ha definido la actividad, en un subdirectorio llamado `submitted`.


## Propiedades de las tareas

Junto con el contenido, en el fichero `info.md` podemos
establecer ciertas propiedades usando el prefijo `@>` seguido del
nombre de la propiedad que queramos declarar.

Una de las propiedades válidas para las tareas es el estado. El valor
de la propiedad estado de una actividad o tarea puede estar este entre
los siguientes:

   - **open**: Una tarea abierta es pública y ofrece un formulario en
       la parte inferior para enviar archivos al sistema.
   - **closed**: Una tarea cerrada es pública pero no permite el envío
       de ficheros.
   - **hide**: Una tarea oculta no muestra su contenido ni se lista al
       cargar el curso.

Para definir, por ejemplo, el estado de una tarea como abierta (open)
hemos de incluir la siguiente linea en nuestro fichero `info.md`:

```
 <!---
 Properties:
 @>status:open
  -->
```

De igual manera que con la anterior, podemos definir la propiedad de
registro de entregas o `logfile`. Esta propiedad activa la creación de
un fichero donde aulasis irá registrando la entrega de la tarea,
apuntando la hora y la ip desde la cual se ha recibido la
información. Para activar esta propiedad incluiremos en un entorno
"properties" como el anterior la línea:

```
 @>logfile:yes
```

