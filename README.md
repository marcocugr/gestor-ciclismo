# Gestor Ciclismo

*Problema:*
Un director deportivo de un equipo de ciclismo debe seleccionar la plantilla de 7 corredores de entre los 20 que suele tener un equipo, de todas las carreras que el equipo tiene durante toda la temporada, y esto es un engorro ya que debe mirar la puntuacion en ranking que ha tenido cada corredor en esa carrera u otras similares otros años para saber cual será el que más puntos le sume al equipo en cada carrera para así ganar el ranking nacional a final de año.

*Para quien va dirigido:*
Para todos los directivos de equipos ciclistas de carretera, sea de la categoría que sea, quitarles horas y horas de busqueda de información de sus corredores a la hora de elegir la plantilla. En mi caso personal, usaré el calendario y los datos (puntos del ranking) de todos los corredores del equipo Sub23 de ciclismo de carretera Extremadura Pebetero, el cual es mi actual equipo con el que tengo una gran cercanía, tanto con deportistas, ya que son mis compañeros, como con directivos, y así será mas fácil de que me faciliten algún dato que falte, o me comenten sus preocupaciones y mejoras que creen que podría hacer en el sistema.

*Que datos existen ya:*
En la página web de la Real Federación Española de ciclismo se encuentra tanto todo el calendario nacional para el año en curso, como las puntuaciones en ranking, de 2 años anteriores de cada corredor en cada carrera más el año actual.
Los datos del ranking de años pasados se obtendrá de la URL : https://rfec.com/index.php/es/smartweb/seccion/clasificacioncircuito/rfec/carretera/ELITE-SUB23/2026/26RANKELITESUB , de donde podríamos extraerlos mediante scrapping , en formato csv para así poderlos procesar, y el calendario de carreras se buscará en https://rfec.com/index.php/es/smartweb/seccion/seccion/rfec//calendario_carretera , aunque como mi idea es enfocarlo a un equipo completo, cada equipo tiene su propio calendario cerrado, hecho en pdf a principio del año desde el cual basaríamos nuestro producto, pero siempre tenemos el calendario total en la web por si necesitamos algo más de información. Para hacer el scrapping de estos datos, debemos hacerlo de una manera algo peculiar, debido a que desde cualquier navegador normal, está prohibido por el administrador de la web entrar en modo desarrollador, por lo que después de muchas pruebas, usando VPNs, distintos navegadores, etc, he llegado a la solución de que mediante TOR, no hay problema en entrar al modo desarrollador de la web por lo que lo haré desde ahí.

*Criterios a seguir:*
Mi idea, no es hacer algo simple como tan solo sumar puntos y convocar a los que tengan más puntos totales, sino tener en cuenta, a parte de este primer criterio que he dicho, factores como en que carreras ha sumado cada corredor los puntos, y darle un peso mayou a aquellos que lo hayan hecho en la misma carrera para la que se va a hacer la convocatoria. También se tendrá en cuenta el descanso de cada corredor, ya que cada X carreras habrá que dar a ciertos corredores descanso para que no se lesionen. También tendré en cuenta que algún corredor puede estar lesionado durante la temporada por lo que mientras dure su lesión, no se le podrá convocar.

*¿Es un problema real del que tengo conocimiento personal o interés?:*
Si, soy ciclista de carretera y todos mis directores deportivos me han expresado su preocupación por este problema y lo necesaria que sería una plataforma de este estilo.

*¿Requiere para su solución el despliegue de una aplicación en la nube?:*
Para empezar, un director deportivo está viajando continuamente y está muy poco en su casa por lo que no deben de depender de una base de datos local, para así poder hacerlo en todo momento y desde cualquier lado. Por otro lado, tanto el calendario como el ranking de puntos se va actualizando semanalmente por lo que no sería viable tenerlo en local

*¿La solución requiere lógica de negocio (procesamiento/cálculo) y NO es solo almacenamiento/búsqueda (CRUD)?:*
Si, ya que requiere procesar todos los datos tanto de los puntos de cada corredor como de las carreras por lo que necesita de esos cálculos y procesamiento

*¿He incluido la licencia libre y el fichero .gitignore?:*
Si.

*¿He subido al repo y enlazado correctamente la foto de la tarjeta del juego de rol en el README.md?:*
Si

*¿He enviado el URL del Pull Request (PR) y no el de la rama ni el de un commit?:*
Si.

*¿El título de mi PR en el repo de la asignatura incluye [IV-26-27]?:*
Si.

*¿He usado el formato de versión semántica v0.0.z en proyectos/objetivo-0.md?:*
Si.

**Fotografías del juego de rol**

**Fotografía de la tarjeta de cliente**
![Fotografía de la tarjeta de cliente](imgs/cliente.jpeg)

**Fotografía de la tarjeta de desarrollador**
![Fotografía de la tarjeta de desarrollador](imgs/desarrollador.jpeg)

**Fotografía de la tarjeta de validación**
![Fotografía de la tarjeta de validación](imgs/validacion.jpeg)

**Configuración del repositorio**

Las evidencias y capturas de la configuración del entorno (claves SSH e identidad de Git) se encuentran aquí:

* [Ver configuración del repositorio](docs/configuracion.md)


## Estado del proyecto

* **Objetivo 0**: completado (problema identificado, repositorio configurado).
* **Objetivo 1**: planificación del proyecto — jornada de usuario, historias de
  usuario y milestones iniciales definidos.
* **Objetivo 2**: modelado del problema en Go, lenguaje ubicuo de HU001 y configuración del entorno de desarrollo.

## Entorno de desarrollo y gestión de versiones de Go

El proyecto utiliza **Go** sobre **Linux** como plataforma de desarrollo acordada (issues #5 y #6).

### Instalación de Go en espacio de usuario

Para trabajar en el proyecto sin modificar el sistema ni requerir permisos de administrador (`root`):

1. **Instalación inicial:** Descargar el binario oficial de Go desde [golang.org/dl](https://golang.org/dl/) y descomprimirlo en el directorio de usuario (por ejemplo, `~/.local/go` o `~/go`):
   ```bash
   tar -C ~/.local -xzf go1.24.4.linux-amd64.tar.gz
   export PATH=$HOME/.local/go/bin:$HOME/go/bin:$PATH
   ```

2. **Gestión y cambio de versiones en espacio de usuario:** Si se requiere una versión diferente de Go, se puede instalar y descargar directamente en espacio de usuario sin tocar el sistema:
   ```bash
   go install golang.org/dl/go1.24.4@latest
   go1.24.4 download
   ```

3. **Fuente única de verdad (`go.mod`):**
   Las líneas `go` y `toolchain` en `go.mod` definen de forma determinista la versión del lenguaje y de la toolchain para todos los colaboradores. Desde Go 1.21, el comando `go` detecta estas directivas y descarga/ejecuta automáticamente la toolchain exacta indicada.

4. **Gestión de dependencias:**
   Las dependencias del proyecto se gestionan mediante los módulos de Go (`go.mod` y `go.sum`), con versiones fijas y sumas de verificación criptográficas. Se descargan en `$GOPATH/pkg/mod` (en espacio de usuario).

5. **Compilación y verificación:**
   Desde un clon limpio con Go configurado en el usuario:
   ```bash
   go build ./...
   go test -v ./...
   ```

### Análisis estático de código (Linter)

El proyecto utiliza **`golangci-lint`** como herramienta de análisis estático acordada en el issue #16 para validar buenas prácticas y detectar posibles errores:

1. **Instalación en espacio de usuario (sin permisos de administrador):**
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
   ```

2. **Ejecución del linter:**
   ```bash
   golangci-lint run ./...
   ```
   La configuración se encuentra definida en `.golangci.yml`.

## Documentación adicional

* [Jornada de usuario](docs/user-journeys.md)
* [Historias de usuario](docs/historias-de-usuario.md)
* [Milestones](docs/milestones.md)
* [Referencia a persona](docs/personas.md)
* [Análisis y modelado de HU001](docs/modelado-hu001.md)
