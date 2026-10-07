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
* **Resultado:** Desempeño de un corredor en una edición concreta. Modela dos situaciones mutuamente excluyentes: *no participó* (celda en blanco) o *participó con N puntos* (incluyendo explícitamente $0$ puntos).
* **Temporada:** Año natural de competición deportiva (p. ej., 2024, 2025, 2026).
* **Histórico:** Conjunto de temporadas que Alfonso consulta para valorar el rendimiento de un corredor en una carrera: la temporada en curso más las dos anteriores (3 temporadas en total).
* **Total:** Suma acumulada de puntos obtenidos por un corredor en una temporada o conjunto de carreras. Es un valor calculado, no un dato almacenado.

---

## 4. Problemas identificados en los issues #7–#14

El análisis detallado de HU001 reveló una cascada de problemas conceptuales resueltos secuencialmente:

* **Issue #7 (Diferencia entre no participar y obtener 0 puntos):** En el ranking federativo, una celda con `0` significa que el ciclista tomó la salida y no sumó puntos; una celda en blanco indica ausencia de participación. Tratar ambos casos como cero falsearía la evaluación deportiva de Alfonso.
* **Issue #8 (Carrera vs. Prueba vs. Edición):** La HU habla de "carreras", pero el ranking puntúa "pruebas" o etapas individuales. Se acordó que cada columna es una edición/prueba, y la carrera es la entidad continua en el tiempo.
* **Issue #9 (Identificación de corredores):** El ranking no provee DNI ni identificador numérico federativo; únicamente incluye cadenas de texto como `"CRAUSE , PEDRI (SUB-23)(MADRID)"`. Se acordó normalizar el nombre (eliminando acentos, espacios redundantes y categorías) y acotar el modelo a los corredores del equipo Extremadura Pebetero.
* **Issue #10 (Relación entre carrera y ediciones de distintos años):** El nombre de una carrera cambia entre temporadas (edición romana, año, patrocinador). No se puede deducir la igualdad solo por texto automáticamente en este milestone; una carrera tiene identidad propia y agrupa sus ediciones anuales, con la regla estricta de como máximo una edición por temporada.
* **Issue #11 (Datos almacenados vs. calculados):** El total de puntos y la posición en el ranking mostrados en la web son sumas agregadas. Almacenarlos crearía redundancia y riesgo de incoherencia. Solo se almacenan resultados individuales; los totales se calculan.
* **Issue #12 (Definición del histórico y temporadas):** Se clarificó que el histórico necesario para Alfonso abarca la temporada actual y las dos anteriores (3 temporadas), prevaleciendo esta decisión sobre menciones preliminares de 5 años.
* **Issue #13 (Datos inválidos, invariantes y errores):** Se determinaron los datos intolerables: puntos negativos, nombres vacíos o sin apellidos, carreras sin identificador, temporadas imposibles, resultados duplicados para el mismo corredor en una edición y más de una edición de la misma carrera en el mismo año.
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
3. El agregado raíz (`HistoricoPuntuaciones`) centraliza las consultas requeridas por HU001 y encapsula el cálculo de totales.

---

## 6. Entidades identificadas

* **`Corredor`:**
  * Representa al ciclista del equipo Extremadura Pebetero.
  * Su identidad reside en su `Nombre` normalizado (objeto valor).
  * Persiste a través de temporadas con independencia de cambios de categoría o club.
* **`Carrera`:**
  * Representa la competición a lo largo de los años.
  * Identidad propia invariante mediante un `id` (cadena única).
  * Posee un nombre descriptivo y gestiona su conjunto de ediciones anuales.
* **`Edicion`:**
  * Representa la disputa concreta de una carrera en una temporada.
  * Posee identidad contextual asociada a la tupla `(carreraID, temporada)`.
  * Conserva el `nombreOficial` tal como fue publicado en el ranking federativo ese año.
  * Mantiene el registro de resultados de los corredores participantes.
* **`HistoricoPuntuaciones`:**
  * Agregado raíz que coordina el catálogo de corredores, carreras y ediciones para resolver las consultas del histórico de HU001.

---

## 7. Objetos valor identificados

* **`Nombre`:**
  * Representa el nombre normalizado del corredor.
  * Inmutable, sin identidad independiente; dos nombres iguales representan a la misma persona.
  * Se construye limpiando acentos, mayúsculas/minúsculas, espacios y sufijos de categoría/región.
* **`Puntos`:**
  * Representa una puntuación deportiva no negativa ($N \ge 0$).
  * Inmutable. Encapsula la validación de que los puntos no pueden ser negativos.
* **`Resultado`:**
  * Modela el resultado de un corredor en una prueba.
  * Inmutable. Distingue entre `ResultadoNoParticipo()` y `ResultadoParticipo(puntos)`.
* **`Temporada`:**
  * Representa un año natural de competición.
  * Inmutable. Valida que el año pertenezca a un rango cronológico válido.

---

## 8. Relaciones entre entidades y objetos valor

* **`Carrera` $1 \longleftrightarrow N$ `Edicion`:** Una carrera puede tener varias ediciones en distintas temporadas, pero **como máximo una edición por temporada**.
* **`Edicion` $N \longleftrightarrow 1$ `Temporada`:** Cada edición pertenece a exactamente una temporada.
* **`Edicion` $1 \longleftrightarrow N$ `Resultado`:** Una edición almacena los resultados de los corredores. Cada corredor tiene **como máximo un resultado** por edición.
* **`Resultado` $N \longleftrightarrow 1$ `Corredor`:** Cada resultado está asociado al corredor correspondiente mediante su objeto valor `Nombre`.
* **`HistoricoPuntuaciones` $1 \longleftrightarrow N$ `Carrera` / `Corredor`:** El agregado organiza y consulta el histórico relacionando corredores, carreras y temporadas.

---

## 9. Datos almacenados

El modelo almacena exclusivamente los datos atómicos necesarios:

* Identificador y nombre canónico de cada `Carrera`.
* Nombre normalizado de cada `Corredor`.
* Temporada y nombre oficial publicado de cada `Edicion`.
* `Resultado` (participación y puntos) de cada corredor en cada edición.

---

## 10. Valores calculados

En estricta conformidad con el issue #11:

* **Total de puntos por temporada:** No se persiste en el modelo; se calcula dinámicamente sumando los puntos de las ediciones en las que el corredor participó durante esa temporada (`CalcularTotalPuntos`).
* **Total de puntos en el histórico:** Se calcula bajo demanda a partir de los resultados en las temporadas consultadas.
* **Posición en el ranking:** No forma parte del modelo; es un cálculo derivado que no pertenece al alcance de HU001.

---

## 11. Reglas e invariantes

El modelo garantiza en tiempo de ejecución las siguientes invariantes del dominio:

1. **Unicidad de edición por temporada:** Una carrera no puede tener más de una edición en la misma temporada (`ErrEdicionDuplicadaEnTemporada`).
2. **Unicidad de resultado por corredor en edición:** Un corredor no puede registrar más de un resultado en una misma edición (`ErrResultadoDuplicado`).
3. **Puntos no negativos:** Todo valor de puntos debe ser un entero mayor o igual que cero (`ErrPuntosNegativos`).
4. **Completitud del nombre:** Un corredor debe contar obligatoriamente con nombre y apellidos tras la normalización (`ErrNombreSinApellidos`, `ErrNombreVacio`).
5. **Completitud de entidades:** No se permiten carreras sin ID ni nombre, ni ediciones sin nombre oficial o desvinculadas de carrera/temporada.
6. **Estados válidos no erróneos:**
   * Que un corredor no participe en una carrera es un estado normal (`ResultadoNoParticipo`).
   * Que un corredor participe y obtenga 0 puntos es un estado normal (`ResultadoParticipo(0)`).
   * Que un corredor no tenga puntos en todo el histórico es admisible.
   * Que una carrera sea nueva y cuente únicamente con su primera edición es admisible.

---

## 12. Datos inválidos y errores

Se han definido tipos de error explícitos en el paquete `dominio`:

* `ErrPuntosNegativos`: Intentar crear puntos con valor $< 0$.
* `ErrNombreVacio`: Cadena de nombre vacía o con solo espacios.
* `ErrNombreSinApellidos`: Nombre que no contiene al menos dos componentes tras normalizar.
* `ErrTemporadaInvalida`: Año fuera de rango cronológico válido.
* `ErrCarreraSinID` / `ErrCarreraSinNombre`: Carrera sin identificador o sin nombre.
* `ErrEdicionSinNombre`: Edición sin denominación oficial.
* `ErrEdicionDuplicadaEnTemporada`: Conflicto de dos ediciones para la misma carrera en el mismo año.
* `ErrResultadoDuplicado`: Conflicto al registrar dos veces a un corredor en una edición.
* `ErrCorredorNulo` / `ErrEdicionNula`: Referencias nulas pasadas a operaciones del dominio.

---

## 13. Decisiones sobre inmutabilidad, identidad y composición

* **Inmutabilidad de Objetos Valor:** En Go, los campos de `Puntos`, `Temporada`, `Nombre` y `Resultado` no son exportados (minúsculas). No existen métodos *setter*; solo constructores validadores y métodos *getter* de sólo lectura. Los objetos valor se pasan por copia, evitando modificaciones colaterales.
* **Identidad de Entidades:** Las entidades (`Corredor`, `Carrera`, `Edicion`) se gestionan mediante punteros y comparan su igualdad a través de sus identificadores (`id` en Carrera, `Nombre` en Corredor).
* **Composición sobre herencia:** Go utiliza composición directa. Las relaciones se modelan mediante referencias y mapas internos no exportados, encapsulando el acceso a través de métodos del agregado.

---

## 14. Estructura de paquetes, módulos y archivos

El código se organiza siguiendo las mejores prácticas idiomáticas de Go:

```
gestor-ciclismo/
├── go.mod                     # Módulo Go con versión y toolchain fijadas
├── iv.yaml                    # Metadatos del Objetivo 2 (lenguaje y ruta de la entidad principal)
├── README.md                  # Descripción del proyecto y guía de entorno
├── docs/                      # Documentación del análisis de dominio
│   └── modelado-hu001.md      # Este documento de análisis consolidado
└── dominio/                   # Paquete con el modelo de dominio puro
    ├── errores.go             # Errores del dominio
    ├── puntos.go              # Objeto valor Puntos
    ├── temporada.go           # Objeto valor Temporada
    ├── nombre.go              # Objeto valor Nombre (con normalización)
    ├── resultado.go           # Objeto valor Resultado (participación vs puntos)
    ├── corredor.go            # Entidad Corredor
    ├── carrera.go             # Entidad Carrera
    ├── edicion.go             # Entidad Edicion
    ├── historico.go           # Agregado HistoricoPuntuaciones y consultas de HU001
    └── dominio_test.go        # Comprobaciones de compilación e invariantes
```

---

## 15. Correspondencia entre el análisis y el código

| Concepto del Dominio | Tipo en Go | Archivo | Justificación |
|---|---|---|---|
| Corredor | `type Corredor struct` | `dominio/corredor.go` | Entidad con ciclo de vida e identidad por nombre |
| Nombre de corredor | `type Nombre struct` | `dominio/nombre.go` | Objeto valor normalizado (ignora tildes, mayúsculas, etc.) |
| Carrera | `type Carrera struct` | `dominio/carrera.go` | Entidad continua con identificador único |
| Edición / Prueba | `type Edicion struct` | `dominio/edicion.go` | Entidad de una carrera en una temporada con nombre oficial |
| Temporada | `type Temporada struct` | `dominio/temporada.go` | Objeto valor que representa un año |
| Puntos | `type Puntos struct` | `dominio/puntos.go` | Objeto valor para enteros $\ge 0$ |
| Resultado | `type Resultado struct` | `dominio/resultado.go` | Objeto valor que distingue "no participó" de $0$ puntos |
| Histórico | `type HistoricoPuntuaciones struct` | `dominio/historico.go` | Agregado raíz para coordinar consultas y calcular totales |
| Invariantes y Errores | `var Err... = errors.New(...)` | `dominio/errores.go` | Errores de dominio tipados para datos inválidos |

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

---

## 17. Justificación del lenguaje Go

Tras el debate en el issue #5, se seleccionó **Go** por las siguientes razones:

1. **Tipado estático y compilación:** Permite detectar errores de tipos en tiempo de compilación y definir modelos estrictos.
2. **Distinción explícita de ausencia de datos:** A diferencia del valor cero por defecto de un `int`, el modelo implementa el tipo `Resultado` para diferenciar de forma infalible la ausencia de participación de una puntuación de cero puntos.
3. **Gestión determinista de versiones:** Go 1.21+ gestiona toolchains y módulos en espacio de usuario a través de `go.mod` sin necesidad de gestores externos complejos ni permisos de administrador.
4. **Inmutabilidad y encapsulación:** La visibilidad por paquetes de Go (identificadores con inicial minúscula) garantiza que los campos de los objetos valor no puedan ser modificados externamente.

---

## 18. Herramientas de desarrollo y calidad

* **Go Toolchain:** Gestión nativa de dependencias mediante `go.mod` fijando `go 1.24.4` y `toolchain go1.24.4`.
* **Linter de código Go:** `golangci-lint` (issue #16) para verificar buenas prácticas, formateo, variables no utilizadas y análisis estático.
* **Integración Continua para YAML:** Workflow de GitHub Actions (`.github/workflows/lint-yaml.yml`) (issue #17) para validar automáticamente la sintaxis de `iv.yaml` y otros ficheros YAML del repositorio.

---

## 19. Cómo se ejecutan las comprobaciones disponibles

1. **Compilar el proyecto:**
   ```bash
   go build ./...
   ```
2. **Ejecutar los tests unitarios del modelo:**
   ```bash
   go test -v ./...
   ```
3. **Ejecutar el linter estático:**
   ```bash
   golangci-lint run ./...
   ```
4. **Validar la sintaxis de los ficheros YAML:**
   ```bash
   python3 -c "import yaml; yaml.safe_load(open('iv.yaml'))"
   ```
