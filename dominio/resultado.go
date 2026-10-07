package dominio

// Resultado representa los puntos obtenidos por un corredor que participó en una edición concreta (issue #7, #15).
// Es un objeto valor inmutable. La ausencia de un resultado para un Nombre en
// una Edicion representa que el corredor no participó, equivalente a una celda
// en blanco en el ranking. Un resultado con Puntos de valor 0 representa que
// participó y no obtuvo puntos.
type Resultado struct {
	puntos Puntos
}

// ResultadoParticipo crea un resultado para un corredor que participó en la edición.
func ResultadoParticipo(puntos Puntos) Resultado {
	return Resultado{
		puntos: puntos,
	}
}

// Puntos devuelve la puntuación obtenida por el corredor participante.
func (r Resultado) Puntos() Puntos {
	return r.puntos
}

// EsCeroPuntos indica si el corredor participó y obtuvo cero puntos.
func (r Resultado) EsCeroPuntos() bool {
	return r.puntos.Valor() == 0
}

// Igual compara dos resultados por su puntuación.
func (r Resultado) Igual(otro Resultado) bool {
	return r.puntos.Igual(otro.puntos)
}
