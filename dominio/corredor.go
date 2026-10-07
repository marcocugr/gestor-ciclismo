package dominio

// Corredor representa a un ciclista cuya trayectoria y resultados se registran (issue #9, #15).
// Es una entidad cuya identidad proviene de su Nombre normalizado (objeto valor) y persiste
// a lo largo de las distintas temporadas, con independencia de que en la fuente federativa
// varíe su club o categoría deportiva.
type Corredor struct {
	nombre Nombre
}

// NuevoCorredor instancia un nuevo Corredor a partir de su objeto valor Nombre.
func NuevoCorredor(nombre Nombre) (*Corredor, error) {
	if nombre.Valor() == "" {
		return nil, ErrNombreVacio
	}
	return &Corredor{
		nombre: nombre,
	}, nil
}

// Nombre devuelve el objeto valor Nombre que identifica al corredor.
func (c *Corredor) Nombre() Nombre {
	return c.nombre
}

// MismoCorredor comprueba si dos entidades Corredor representan a la misma persona
// comparando la igualdad de sus nombres normalizados.
func (c *Corredor) MismoCorredor(otro *Corredor) bool {
	if c == nil || otro == nil {
		return false
	}
	return c.nombre.Igual(otro.nombre)
}
