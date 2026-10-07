package dominio

import "strings"

// Edicion representa la celebración de una carrera en una temporada concreta (issue #8, #10, #15).
// Ocupa una columna en la tabla del ranking federativo y sobre ella se registran los resultados.
// Conserva el nombre oficial publicado ese año, ya que sirve para que Alfonso identifique
// qué ediciones corresponden a cada carrera.
type Edicion struct {
	carreraID     string
	temporada     Temporada
	nombreOficial string
	resultados    map[string]Resultado // Indexado por Nombre.Valor()
}

// NuevaEdicion instancia una nueva edición asociada a una carrera, una temporada y un nombre oficial.
func NuevaEdicion(carreraID string, temporada Temporada, nombreOficial string) (*Edicion, error) {
	cID := strings.TrimSpace(carreraID)
	if cID == "" {
		return nil, ErrCarreraSinID
	}
	nOficial := strings.TrimSpace(nombreOficial)
	if nOficial == "" {
		return nil, ErrEdicionSinNombre
	}
	return &Edicion{
		carreraID:     cID,
		temporada:     temporada,
		nombreOficial: nOficial,
		resultados:    make(map[string]Resultado),
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
// Devuelve ErrResultadoDuplicado si ya existe un resultado previo registrado para ese corredor.
func (e *Edicion) RegistrarResultado(corredor *Corredor, resultado Resultado) error {
	if corredor == nil {
		return ErrCorredorNulo
	}
	clave := corredor.Nombre().Valor()
	if _, existe := e.resultados[clave]; existe {
		return ErrResultadoDuplicado
	}
	e.resultados[clave] = resultado
	return nil
}

// ResultadoDe obtiene el resultado de un corredor en esta edición si participó o tiene registro.
func (e *Edicion) ResultadoDe(corredor *Corredor) (Resultado, bool) {
	if corredor == nil {
		return Resultado{}, false
	}
	res, existe := e.resultados[corredor.Nombre().Valor()]
	return res, existe
}

// Resultados devuelve una copia del mapa de resultados de la edición.
func (e *Edicion) Resultados() map[string]Resultado {
	copia := make(map[string]Resultado, len(e.resultados))
	for k, v := range e.resultados {
		copia[k] = v
	}
	return copia
}
