package dominio

import "strings"

// Carrera representa una competición deportiva continua a lo largo del tiempo (issue #8, #10, #15).
// Es una entidad que mantiene su identidad independiente entre temporadas, aunque su denominación
// oficial, patrocinadores o número de edición varíen en las publicaciones anuales del ranking.
type Carrera struct {
	id        string
	nombre    string
	ediciones map[int]*Edicion // Mapeo de temporada (año) a su edición correspondiente
}

// NuevaCarrera instancia una nueva Carrera con su identificador propio e invariable y su nombre canónico.
func NuevaCarrera(id, nombre string) (*Carrera, error) {
	idLimpio := strings.TrimSpace(id)
	if idLimpio == "" {
		return nil, ErrorCarreraInvalida{Motivo: "el identificador de la carrera no puede estar vacío"}
	}
	nombreLimpio := strings.TrimSpace(nombre)
	if nombreLimpio == "" {
		return nil, ErrorCarreraInvalida{Motivo: "el nombre de la carrera no puede estar vacío"}
	}
	return &Carrera{
		id:        idLimpio,
		nombre:    nombreLimpio,
		ediciones: make(map[int]*Edicion),
	}, nil
}

// ID devuelve el identificador único de la carrera.
func (c *Carrera) ID() string {
	return c.id
}

// Nombre devuelve el nombre descriptivo de la carrera.
func (c *Carrera) Nombre() string {
	return c.nombre
}

// AgregarEdicion vincula una edición a la carrera garantizando dos invariantes del dominio:
//  1. La edición debe pertenecer obligatoriamente a esta carrera receptora (coherencia de ID).
//  2. Una carrera no puede tener más de una edición en la misma temporada (issue #10, #13).
func (c *Carrera) AgregarEdicion(edicion *Edicion) error {
	if edicion == nil {
		return ErrParametroNulo
	}
	if edicion.CarreraID() != c.id {
		return ErrorEdicionCarreraIncompatible{
			CarreraReceptoraID: c.id,
			EdicionCarreraID:   edicion.CarreraID(),
			NombreEdicion:      edicion.NombreOficial(),
		}
	}
	anio := edicion.Temporada().Anio()
	if _, existe := c.ediciones[anio]; existe {
		return ErrorEdicionDuplicadaEnTemporada{
			CarreraID: c.id,
			Anio:      anio,
		}
	}
	c.ediciones[anio] = edicion
	return nil
}

// EdicionEnTemporada devuelve la edición celebrada en la temporada dada, si existe.
func (c *Carrera) EdicionEnTemporada(temporada Temporada) (*Edicion, bool) {
	ed, existe := c.ediciones[temporada.Anio()]
	return ed, existe
}

// Ediciones devuelve una copia de las ediciones vinculadas a la carrera.
func (c *Carrera) Ediciones() []*Edicion {
	lista := make([]*Edicion, 0, len(c.ediciones))
	for _, ed := range c.ediciones {
		lista = append(lista, ed)
	}
	return lista
}
