package dominio

// Puntos representa la cantidad de puntos obtenida por un corredor en una prueba o edición (issue #7, #15).
// Es un objeto valor inmutable: su estado no se puede alterar tras la creación.
// El valor debe ser un número entero mayor o igual que cero. Un valor de 0 indica que el
// corredor participó pero no sumó puntos.
type Puntos struct {
	valor int
}

// NuevosPuntos crea una nueva instancia del objeto valor Puntos validando que sea >= 0.
// Si el valor recibido es negativo, devuelve ErrorPuntosNegativos con el valor recibido (issue #7, #13).
func NuevosPuntos(valor int) (Puntos, error) {
	if valor < 0 {
		return Puntos{}, ErrorPuntosNegativos{Valor: valor}
	}
	return Puntos{valor: valor}, nil
}

// Valor devuelve el valor numérico entero de los puntos.
func (p Puntos) Valor() int {
	return p.valor
}

// Igual compara dos objetos valor Puntos por igualdad de valor.
func (p Puntos) Igual(otro Puntos) bool {
	return p.valor == otro.valor
}
