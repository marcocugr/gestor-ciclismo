package dominio

// Temporada representa un año de competición deportiva en el dominio ciclista (issue #12, #15).
// Es un objeto valor inmutable y sin identidad propia, determinado únicamente por su año.
type Temporada struct {
	anio int
}

// NuevaTemporada crea una nueva Temporada tras validar que el año sea válido.
// Se rechazan años menores o iguales a 1900 o futuros lejanos como datos inválidos del dominio (issue #13).
func NuevaTemporada(anio int) (Temporada, error) {
	if anio < 1900 || anio > 2100 {
		return Temporada{}, ErrTemporadaInvalida
	}
	return Temporada{anio: anio}, nil
}

// Anio devuelve el año numérico correspondiente a la temporada.
func (t Temporada) Anio() int {
	return t.anio
}

// Igual compara dos temporadas por su año.
func (t Temporada) Igual(otra Temporada) bool {
	return t.anio == otra.anio
}
