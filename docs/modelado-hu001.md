# Análisis y Modelado del Dominio: HU001

Este documento recoge el análisis del dominio y el diseño del modelo técnico correspondiente al **Objetivo 2** y al **Milestone 0**, aplicando la metodología de **Domain Driven Design (DDD)** sobre la historia de usuario **HU001** del proyecto *Gestor Ciclismo*.

---

## 1. Historia de usuario analizada

La historia de usuario sobre la que se fundamenta este trabajo es **[HU001] Histórico de puntuación disperso en varias páginas** (issue #2):

> **Como** director deportivo (Alfonso Rodríguez, del equipo Sub-23 Extremadura Pebetero),  
> **quiero** consultar en un único lugar el histórico de puntuación de cada corredor por carrera de años anteriores (actualmente disperso en distintas URLs anuales de la RFEC),  
> **para** valorar el rendimiento real de los corredores en carreras específicas antes de tomar decisiones de convocatoria.

Alfonso lleva más de 30 años dirigiendo equipos y viaja constantemente durante la temporada. Actualmente debe consultar de forma manual tablas anuales en la web de la Real Federación Española de Ciclismo (RFEC), donde cada año tiene una URL independiente con decenas de columnas de pruebas ($P_1, P_2, \dots$), cruzando datos entre hoteles, aeropuertos y coches.

---

## 2. Objetivo del Milestone 0

El **Milestone 0** (*Modelado del problema*) define un **Producto Mínimamente Viable (PMV) interno**. Su propósito es:

1. Establecer mediante **Domain Driven Design (DDD)** un **lenguaje ubicuo** común entre el cliente (Alfonso / Manu) y el desarrollador (Marco).
2. Modelar en código las estructuras de datos, tipos, entidades, objetos valor, relaciones e invariantes necesarios para HU001.
3. Dejar preparada una base sólida, tipada y validada para que en el siguiente milestone (Milestone 1) se pueda implementar la lógica de negocio sin ambigüedades.

En este milestone **no** se implementa scraping web, interfaces gráficas, persistencia en base de datos ni algoritmos de convocatoria automática. La historia de usuario HU001 permanece abierta porque este hito cubre exclusivamente el modelado del dominio.

---

## 3. Lenguaje ubicuo construido

A partir de las conversaciones con el cliente en los issues #2, #8 y #14, se ha consolidado el siguiente vocabulario del dominio:

* **Carrera:** Competición deportiva tal como la entiende Alfonso, que mantiene su identidad a lo largo de los años aunque varíe su denominación oficial, patrocinador o número de edición (p. ej., "Clásica de Valladolid").
* **Edición:** Celebración de una carrera en una temporada concreta. Corresponde a una columna del ranking oficial federativo y es donde se registran los resultados de los corredores. Conserva el nombre oficial publicado ese año.
* **Prueba:** Sinónimo de edición en el contexto del ranking federativo (cada columna $P_i$ puntuable).
* **Corredor:** Deportista del equipo Extremadura Pebetero cuya trayectoria y puntuaciones se registran, identificado por su nombre y apellidos.
* **Nombre:** Objeto valor identificador del corredor, normalizado para ser independiente de mayúsculas/minúsculas, tildes, espacios redundantes y categorías entre paréntesis.
* **Puntos:** Número entero mayor o igual que cero ($N \ge 0$) que obtiene un corredor en una edición. No tiene límite superior.
* **Resultado:** Desempeño de un corredor en una edición concreta. Una celda vacía del ranking se representa mediante la ausencia del Nombre del corredor en los resultados de esa edición. Un resultado existente con Puntos de valor 0 representa que el corredor participó y no obtuvo puntos.
* **Temporada:** Año natural de competición deportiva (p. ej., 2024, 2025, 2026). Debe ser un año positivo y que ya haya comenzado en el sistema.
* **Histórico:** Conjunto de temporadas que Alfonso consulta para valorar el rendimiento de un corredor en una carrera: la temporada en curso más las dos anteriores (3 temporadas en total).

---

## 4. Problemas identificados en los issues #7–#14

Los issues #7 a #14 corresponden a debates y aclaraciones previas sobre el dominio del problema. No constituyen tareas de código independientes que se cierren de forma aislada, sino una cascada de análisis cuyos acuerdos se consolidan en la modelización conjunta:

* **Issue #7 (Diferencia entre no participar y obtener 0 puntos):** En el ranking federativo, una celda con `0` significa que el ciclista tomó la salida y no sumó puntos; una celda en blanco indica ausencia de participación. Tratar ambos casos como cero falsearía la evaluación deportiva de Alfonso.
* **Issue #8 (Carrera vs. Prueba vs. Edición):** La HU habla de "carreras", pero el ranking puntúa "pruebas" o etapas individuales. Se acordó que cada columna es una edición/prueba, y la carrera es la entidad continua en el tiempo.
* **Issue #9 (Identificación de corredores):** El ranking no provee DNI ni identificador numérico federativo; únicamente incluye cadenas de texto como `"CRAUSE , PEDRI (SUB-23)(MADRID)"`. Se acordó normalizar el nombre (eliminando acentos, espacios redundantes y categorías) y acotar el modelo a los corredores del equipo Extremadura Pebetero.
* **Issue #10 (Relación entre carrera y ediciones de distintos años):** El nombre de una carrera cambia entre temporadas (edición romana, año, patrocinador). No se puede deducir la igualdad solo por texto automáticamente en este milestone; una carrera tiene identidad propia y agrupa sus ediciones anuales, con la regla estricta de como máximo una edición por temporada.
* **Issue #11 (Datos almacenados vs. calculados):** El total de puntos y la posición en el ranking mostrados en la web son sumas agregadas. Almacenarlos crearía redundancia y riesgo de incoherencia. Solo se almacenan resultados individuales; los totales se calculan.
* **Issue #12 (Definición del histórico y temporadas):** Se clarificó que el histórico necesario para Alfonso abarca la temporada actual y las dos anteriores (3 temporadas), prevaleciendo esta decisión sobre menciones preliminares de 5 años.
* **Issue #13 (Datos inválidos, invariantes y errores):** Se determinaron los datos intolerables: puntos negativos, nombres vacíos o sin apellidos (o con comas sin texto válido a ambos lados), carreras sin identificador, temporadas que aún no han comenzado, resultados duplicados para el mismo corredor en una edición y más de una edición de la misma carrera en el mismo año.
* **Issue #14 (Construcción del lenguaje ubicuo):** Alineación explícita del vocabulario entre cliente y desarrollador para evitar ambigüedades técnicas.

---

## 5. Cómo se consolidan esos problemas en el issue #15

El **issue #15** no introduce una funcionalidad nueva, sino que transforma la cascada de decisiones de los issues #7 a #14 en una arquitectura de dominio coherente bajo los principios de Domain Driven Design (DDD):

```
       #7 (Ausencia vs 0)  ──┐
       #8 (Carrera/Edición) ─┤
       #9 (Identidad ciclista)─┤
      #10 (Ediciones anuales)─┼──> #15 (Modelo Consolidado DDD)
      #11 (Almacenado/Calc)  ─┤       ├── Entidades (Corredor, Carrera, Edicion, Historico)
      #12 (Temporadas)       ─┤       ├── Objetos Valor (Nombre, Puntos, Resultado, Temporada)
      #13 (Invariantes)      ─┤       ├── Invariantes y Validaciones
      #14 (Lenguaje Ubicuo)  ─┘       └── Estructura de paquetes Go
```

El modelo unifica estas decisiones:
1. Las entidades poseen identidad propia y ciclo de vida (`Corredor`, `Carrera`, `Edicion`).
2. Los valores inmutables se aíslan en objetos valor (`Nombre`, `Puntos`, `Resultado`, `Temporada`).
3. El agregado raíz (`HistoricoPuntuaciones`) centraliza las consultas requeridas por HU001, previene duplicados y encapsula el cálculo de totales.

---

## 6. Entidades identificadas

* **`Corredor`:**
  * Representa al ciclista del equipo Extremadura Pebetero.
  * Su identidad reside en su `Nombre` normalizado (objeto valor).
  * Persiste a través de temporadas con independencia de cambios de categoría o club.
* **`Carrera`:**
  * Representa la competición a lo largo de los años.
  * Identidad propia invariante mediante un `id` (cadena única).
  * Posee un nombre descriptivo y gestiona su conjunto de ediciones anuales, garantizando que cada edición pertenezca a esta carrera y que no existan dos ediciones en el mismo año.
* **`Edicion`:**
  * Representa la disputa concreta de una carrera en una temporada.
  * Posee identidad contextual asociada a la tupla `(carreraID, temporada)`.
  * Conserva el `nombreOficial` tal como fue publicado en el ranking federativo ese año.
  * Mantiene el registro de resultados indexado por el objeto valor `Nombre` para conservar el tipo de dominio.
* **`HistoricoPuntuaciones`:**
  * Agregado raíz que coordina el catálogo de corredores, carreras y ediciones para resolver las consultas del histórico de HU001, rechazando registros duplicados de corredores y carreras para no sobrescribir datos previos.

---

## 7. Objetos valor identificados

* **`Nombre`:**
  * Representa el nombre normalizado del corredor.
  * Inmutable, sin identidad independiente; dos nombres iguales representan a la misma persona.
  * Se construye limpiando acentos, mayúsculas/minúsculas, espacios y sufijos de categoría/región.
  * Valida estrictamente que contenga nombre y apellidos; si incluye coma, exige texto no vacío a ambos lados de la misma.
* **`Puntos`:**
  * Representa una puntuación deportiva no negativa ($N \ge 0$).
  * Inmutable. Encapsula la validación de que los puntos no pueden ser negativos mediante el tipo de error `ErrorPuntosNegativos`.
* **`Resultado`:**
  * Modela el resultado de un corredor en una prueba.
  * Inmutable. Una celda vacía del ranking se representa mediante la ausencia del Nombre del corredor en los resultados de esa edición. Un resultado existente con Puntos de valor 0 representa que el corredor participó y no obtuvo puntos.
* **`Temporada`:**
  * Representa un año natural de competición.
  * Inmutable. Valida que el año sea estrictamente positivo y que no corresponda a una temporada futura que todavía no haya comenzado en el sistema (`ErrorTemporadaNoComenzada`).

---

## 8. Relaciones entre entidades y objetos valor

* **`Carrera` $1 \longleftrightarrow N$ `Edicion`:** Una carrera puede tener varias ediciones en distintas temporadas, pero **como máximo una edición por temporada**. Toda edición agregada debe pertenecer obligatoriamente a esa carrera (`CarreraID() == c.id`).
* **`Edicion` $N \longleftrightarrow 1$ `Temporada`:** Cada edición pertenece a exactamente una temporada válida y ya comenzada.
* **`Edicion` $1 \longleftrightarrow N$ `Resultado`:** Una edición almacena los resultados indexados por el objeto valor `Nombre`. Cada corredor tiene **como máximo un resultado** por edición.
* **`Resultado` $N \longleftrightarrow 1$ `Corredor`:** Cada resultado está asociado al corredor correspondiente mediante su objeto valor `Nombre`.
* **`HistoricoPuntuaciones` $1 \longleftrightarrow N$ `Carrera` / `Corredor`:** El agregado organiza y consulta el histórico relacionando corredores, carreras y temporadas sin permitir sobrescrituras duplicadas.

---

## 9. Datos almacenados

El modelo almacena exclusivamente los datos atómicos necesarios:

* Identificador y nombre canónico de cada `Carrera`.
* Nombre normalizado de cada `Corredor`.
* Temporada y nombre oficial publicado de cada `Edicion`.
* `Resultado` y puntos de cada corredor que participó en cada edición.

---

## 10. Valores calculados

En estricta conformidad con el issue #11:

* **Total de puntos por temporada:** No se persiste en el modelo; se calcula dinámicamente sumando los puntos de las ediciones en las que el corredor participó durante esa temporada (`CalcularTotalPuntos`).
* **Total de puntos en el histórico:** Se calcula bajo demanda a partir de los resultados en las temporadas consultadas.
* **Posición en el ranking:** No forma parte del modelo; es un cálculo derivado que no pertenece al alcance de HU001.

---

## 11. Reglas e invariantes

El modelo garantiza en tiempo de ejecución las siguientes invariantes del dominio:

1. **Unicidad de edición por temporada:** Una carrera no puede tener más de una edición en la misma temporada (`ErrorEdicionDuplicadaEnTemporada`).
2. **Coherencia de carrera en edición:** Una edición solo puede agregarse a la carrera cuyo identificador coincida con el `carreraID` de la edición (`ErrorEdicionCarreraIncompatible`).
3. **Unicidad de resultado por corredor en edición:** Un corredor no puede registrar más de un resultado en una misma edición (`ErrorResultadoDuplicado`).
4. **Puntos no negativos:** Todo valor de puntos debe ser un entero mayor o igual que cero (`ErrorPuntosNegativos`).
5. **Completitud y validez del nombre:** Un corredor debe contar obligatoriamente con nombre y apellidos tras la normalización. Si se usa coma, debe haber texto a ambos lados (`ErrorNombreInvalido`).
6. **Temporada válida y comenzada:** La temporada debe ser positiva y no puede pertenecer a un año futuro no comenzado (`ErrorTemporadaNoComenzada`). La edición rechaza temporadas no válidas.
7. **Unicidad en el catálogo del histórico:** No se permite registrar por duplicado un corredor o una carrera en el agregado `HistoricoPuntuaciones`, preservando los datos previos (`ErrorCorredorYaRegistrado`, `ErrorCarreraYaRegistrada`).
8. **Representación de ausencia:** Una celda vacía del ranking se representa mediante la ausencia del Nombre del corredor en los resultados de esa edición. Un resultado existente con Puntos de valor 0 representa que el corredor participó y no obtuvo puntos.

---

## 12. Datos inválidos y errores

Se han definido tipos de error estructurados y expresivos que proporcionan contexto detallado:

* `ErrorPuntosNegativos`: Incluye el valor numérico inválido recibido.
* `ErrorNombreInvalido`: Especifica el motivo concreto (vacío, sin apellidos, formato de coma incompleto) y la cadena original.
* `ErrorTemporadaInvalida`: Detalla el año no positivo recibido.
* `ErrorTemporadaNoComenzada`: Indica el año solicitado y el año actual del sistema.
* `ErrorCarreraInvalida` / `ErrorEdicionInvalida`: Detalla el motivo de la invalidez.
* `ErrorEdicionDuplicadaEnTemporada`: Especifica el ID de carrera y el año en conflicto.
* `ErrorEdicionCarreraIncompatible`: Detalla el ID de la carrera receptora, el ID de la carrera de la edición y el nombre oficial de la edición.
* `ErrorResultadoDuplicado`: Incluye el nombre del corredor y el nombre de la edición en conflicto.
* `ErrorCorredorYaRegistrado` / `ErrorCarreraYaRegistrada`: Informa del elemento duplicado en el histórico.

*Diferenciación técnica:* Las referencias nulas pasadas a funciones o métodos se tratan como precondiciones técnicas mediante `ErrParametroNulo`, distinguiéndolas claramente de las reglas de negocio del dominio de Alfonso.

---

## 13. Decisiones sobre inmutabilidad, identidad y composición

* **Inmutabilidad de Objetos Valor:** En Go, los campos de `Puntos`, `Temporada`, `Nombre` y `Resultado` no son exportados (minúsculas). No existen métodos *setter*; solo constructores validadores y métodos *getter* de sólo lectura. Los objetos valor se pasan por copia, evitando modificaciones colaterales.
* **Preservación de tipos en mapas:** Los resultados y corredores se indexan mediante el objeto valor `Nombre` (`map[Nombre]Resultado`), conservando el tipo de dominio en lugar de degradarlo a `string`.
* **Identidad de Entidades:** Las entidades (`Corredor`, `Carrera`, `Edicion`) se gestionan mediante punteros y comparan su igualdad a través de sus identificadores (`id` en Carrera, `Nombre` en Corredor).
* **Composición sobre herencia:** Go utiliza composición directa. Las relaciones se modelan mediante referencias y mapas internos, encapsulando el acceso a través de métodos del agregado.

---

## 14. Estructura de paquetes, módulos y archivos

El código se organiza siguiendo las convenciones idiomáticas de Go:

```
gestor-ciclismo/
├── go.mod                     # Módulo Go github.com/manuusnchz/gestor-ciclismo fijado a 1.24.4
├── iv.yaml                    # Metadatos del Objetivo 2 (lenguaje y ruta de la entidad principal)
├── README.md                  # Descripción del proyecto y guía de entorno
├── docs/                      # Documentación del análisis de dominio
│   └── modelado-hu001.md      # Este documento de análisis consolidado
└── dominio/                   # Paquete con el modelo de dominio puro
    ├── errores.go             # Errores del dominio contextuales y precondiciones
    ├── puntos.go              # Objeto valor Puntos
    ├── temporada.go           # Objeto valor Temporada (validación de año comenzado)
    ├── nombre.go              # Objeto valor Nombre (normalización y validación estricta)
    ├── resultado.go           # Objeto valor Resultado (representación canónica de participación)
    ├── corredor.go            # Entidad Corredor
    ├── carrera.go             # Entidad Carrera (validación de ediciones compatibles)
    ├── edicion.go             # Entidad Edicion (indexada por Nombre, temporada validada)
    └── historico.go           # Agregado HistoricoPuntuaciones y consultas de HU001
```

---

## 15. Correspondencia entre el análisis y el código

| Concepto del Dominio | Tipo en Go | Archivo | Justificación |
|---|---|---|---|
| Corredor | `type Corredor struct` | `dominio/corredor.go` | Entidad con ciclo de vida e identidad por nombre |
| Nombre de corredor | `type Nombre struct` | `dominio/nombre.go` | Objeto valor normalizado (ignora tildes, mayúsculas, comas completas) |
| Carrera | `type Carrera struct` | `dominio/carrera.go` | Entidad continua con identificador único y control de ID de edición |
| Edición / Prueba | `type Edicion struct` | `dominio/edicion.go` | Entidad de una carrera en una temporada con resultados indexados por `Nombre` |
| Temporada | `type Temporada struct` | `dominio/temporada.go` | Objeto valor que representa un año comenzado |
| Puntos | `type Puntos struct` | `dominio/puntos.go` | Objeto valor para enteros $\ge 0$ |
| Resultado | `type Resultado struct` | `dominio/resultado.go` | Objeto valor inmutable para los puntos de un corredor participante; la ausencia del Nombre representa una celda vacía |
| Histórico | `type HistoricoPuntuaciones struct` | `dominio/historico.go` | Agregado raíz para coordinar consultas sin sobrescribir duplicados |
| Invariantes y Errores | `type Error... struct` | `dominio/errores.go` | Errores estructurados con contexto y valores implicados |

---

## 16. Elementos que quedan fuera del alcance de este milestone

Conforme a las decisiones de HU001 y el Milestone 0, quedan explícitamente excluidos:

* Algoritmo de selección automática de los 7 corredores (pertenece a HU002).
* Ponderación por carreras similares o perfiles de etapa.
* Comparación automática o cruce heurístico de nombres de carreras entre temporadas (Milestone 1).
* Gestión de bajas, lesiones y descansos semanales.
* Extracción automática de datos de la web de la RFEC mediante scraping o TOR.
* Representación y clasificación general agregada de vueltas por etapas.
* Persistencia en base de datos, APIs REST o interfaces de usuario.
* Baterías de tests de comportamiento y pruebas de integración (pertenecientes al Objetivo 4).

---

## 17. Justificación del lenguaje Go

Tras el debate en el issue #5, se seleccionó **Go** por las siguientes razones:

1. **Tipado estático y compilación:** Permite detectar errores de tipos en tiempo de compilación y definir modelos estrictos.
2. **Distinción explícita de ausencia de datos:** La ausencia del Nombre del corredor en los resultados representa una celda vacía del ranking, mientras que un resultado existente con Puntos de valor 0 representa participación sin puntos.
3. **Gestión determinista de dependencias:** Go gestiona dependencias en espacio de usuario a través de `go.mod` sin necesidad de gestores externos complejos ni permisos de administrador.
4. **Inmutabilidad y encapsulación:** La visibilidad por paquetes de Go (identificadores con inicial minúscula) garantiza que los campos de los objetos valor no puedan ser modificados externamente.

---

## 18. Herramientas de desarrollo y calidad

* **Go Compilador y Módulos (#6):** El módulo declara en `go.mod` el requisito mínimo `go 1.24.4`, que coincide con la versión disponible en el entorno de desarrollo. No se utiliza una directiva `toolchain`; la gestión de dependencias se realiza mediante los módulos nativos de Go en espacio de usuario.
* **Linter de código Go (#16):** Configurado `golangci-lint` con `.golangci.yml` (linters: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `unconvert`).
* **Integración Continua para YAML (#17):** Workflow de GitHub Actions (`.github/workflows/lint-yaml.yml`) para validar automáticamente la sintaxis de `iv.yaml` y otros ficheros YAML del repositorio mediante `yq`.

---

## 19. Cómo se ejecutan las comprobaciones disponibles

1. **Compilar el paquete de dominio:**
   ```bash
   go build ./...
   ```
2. **Ejecutar el análisis estático (linter):**
   ```bash
   golangci-lint run ./...
   ```
3. **Validar la sintaxis de los ficheros YAML:**
   ```bash
   python3 -c "import yaml; yaml.safe_load(open('iv.yaml'))"
   ```
