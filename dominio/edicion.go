package dominio

import (
	"fmt"
	"strings"
)

// Edicion representa la celebración de una carrera en una temporada concreta (issue #8, #10, #15).
// Ocupa una columna en la tabla del ranking federativo y sobre ella se registran los resultados.
// Conserva el nombre oficial publicado ese año, ya que sirve para que Alfonso identifique
// qué ediciones corresponden a cada carrera.
//
// Los resultados se indexan por el objeto valor Nombre para preservar el tipo del dominio.
// La ausencia de un Nombre en resultados representa que el corredor no participó en la edición,
// equivalente a una celda en blanco del ranking. Un resultado registrado con Puntos{0}
// representa que participó y obtuvo cero puntos.
type Edicion struct {
	carreraID     string
	temporada     Temporada
	nombreOficial string
	resultados    map[Nombre]Resultado
}

// NuevaEdicion instancia una nueva edición asociada a una carrera, una temporada y un nombre oficial.
// Valida que el identificador de carrera y el nombre oficial no estén vacíos, y que la temporada sea
// un objeto valor válido (año positivo y que ya haya comenzado en el sistema) (issue #8, #13).
func NuevaEdicion(carreraID string, temporada Temporada, nombreOficial string) (*Edicion, error) {
	cID := strings.TrimSpace(carreraID)
	if cID == "" {
		return nil, ErrorEdicionInvalida{Motivo: "el identificador de carrera no puede estar vacío"}
	}
	nOficial := strings.TrimSpace(nombreOficial)
	if nOficial == "" {
		return nil, ErrorEdicionInvalida{Motivo: "el nombre oficial de la edición no puede estar vacío"}
	}
	if !temporada.EsValida() {
		return nil, ErrorEdicionInvalida{
			Motivo: fmt.Sprintf("la temporada es inválida o no ha comenzado (año %d)", temporada.Anio()),
		}
	}
	return &Edicion{
		carreraID:     cID,
		temporada:     temporada,
		nombreOficial: nOficial,
		resultados:    make(map[Nombre]Resultado),
	}, nil
}

// CarreraID devuelve el identificador de la carrera a la que pertenece esta edición.
func (e *Edicion) CarreraID() string {
	return e.carreraID
}

// Temporada devuelve la temporada en la que se disputó la edición.
func (e *Edicion) Temporada() Temporada {
	return e.temporada
}

// NombreOficial devuelve el nombre con el que se publicó la prueba en el ranking ese año.
func (e *Edicion) NombreOficial() string {
	return e.nombreOficial
}

// RegistrarResultado almacena el resultado de un corredor en esta edición (issue #13).
// Devuelve ErrorResultadoDuplicado si ya existe un resultado previo registrado para ese corredor.
func (e *Edicion) RegistrarResultado(corredor *Corredor, resultado Resultado) error {
	if corredor == nil {
		return ErrParametroNulo
	}
	nombre := corredor.Nombre()
	if _, existe := e.resultados[nombre]; existe {
		return ErrorResultadoDuplicado{
			NombreCorredor: nombre.Valor(),
			NombreEdicion:  e.nombreOficial,
		}
	}
	e.resultados[nombre] = resultado
	return nil
}

// ResultadoDe obtiene el resultado de un corredor participante en esta edición.
// Si devuelve false, el corredor no participó y su celda aparece vacía en el ranking.
func (e *Edicion) ResultadoDe(corredor *Corredor) (Resultado, bool) {
	if corredor == nil {
		return Resultado{}, false
	}
	res, existe := e.resultados[corredor.Nombre()]
	return res, existe
}

// Resultados devuelve una copia del mapa de resultados de la edición, conservando el tipo de dominio Nombre.
func (e *Edicion) Resultados() map[Nombre]Resultado {
	copia := make(map[Nombre]Resultado, len(e.resultados))
	for k, v := range e.resultados {
		copia[k] = v
	}
	return copia
}
