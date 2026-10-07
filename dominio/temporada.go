package dominio

import "time"

// Temporada representa un año de competición deportiva en el dominio ciclista (issue #12, #15).
// Es un objeto valor inmutable y sin identidad propia, determinado únicamente por su año.
type Temporada struct {
	anio int
}

// NuevaTemporada crea una nueva Temporada tras validar que el año sea estrictamente positivo
// y que la temporada no corresponda a un año futuro que todavía no ha comenzado en el sistema (issue #13).
func NuevaTemporada(anio int) (Temporada, error) {
	if anio <= 0 {
		return Temporada{}, ErrorTemporadaInvalida{Anio: anio}
	}
	anioActual := time.Now().Year()
	if anio > anioActual {
		return Temporada{}, ErrorTemporadaNoComenzada{Anio: anio, AnioActual: anioActual}
	}
	return Temporada{anio: anio}, nil
}

// EsValida indica si la temporada tiene un año positivo y que ya ha comenzado en el sistema.
func (t Temporada) EsValida() bool {
	return t.anio > 0 && t.anio <= time.Now().Year()
}

// Anio devuelve el año numérico correspondiente a la temporada.
func (t Temporada) Anio() int {
	return t.anio
}

// Igual compara dos temporadas por su año.
func (t Temporada) Igual(otra Temporada) bool {
	return t.anio == otra.anio
}
