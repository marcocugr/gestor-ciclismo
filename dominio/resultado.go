package dominio

// Resultado representa el resultado obtenido por un corredor en una edición concreta (issue #7, #15).
// Es un objeto valor inmutable que modela explícitamente dos estados mutuamente excluyentes:
//  1. No participó (equivalente a celda en blanco en el ranking). No contiene puntos.
//  2. Participó con N puntos (donde N >= 0). Incluye explícitamente el caso de participar y obtener 0 puntos.
type Resultado struct {
	participo bool
	puntos    Puntos
}

// ResultadoNoParticipo crea un objeto valor que representa la ausencia de participación (celda en blanco).
func ResultadoNoParticipo() Resultado {
	return Resultado{
		participo: false,
		puntos:    Puntos{},
	}
}

// ResultadoParticipo crea un objeto valor que representa la participación del corredor con una puntuación obtenida.
func ResultadoParticipo(puntos Puntos) Resultado {
	return Resultado{
		participo: true,
		puntos:    puntos,
	}
}

// Participo indica si el corredor efectivamente tomó la salida y participó en la prueba.
func (r Resultado) Participo() bool {
	return r.participo
}

// Puntos devuelve los puntos del resultado y un booleano indicando si son válidos (es decir, si participó).
// Si el corredor no participó, devuelve Puntos{} vacío y false.
func (r Resultado) Puntos() (Puntos, bool) {
	if !r.participo {
		return Puntos{}, false
	}
	return r.puntos, true
}

// EsCeroPuntos indica si el corredor participó en la prueba y obtuvo exactamente cero puntos.
func (r Resultado) EsCeroPuntos() bool {
	return r.participo && r.puntos.Valor() == 0
}

// Igual compara dos resultados por su estado y puntuación.
func (r Resultado) Igual(otro Resultado) bool {
	if r.participo != otro.participo {
		return false
	}
	if !r.participo {
		return true
	}
	return r.puntos.Igual(otro.puntos)
}
