package dominio

// RegistroHistorico representa la consulta de HU001 para una temporada concreta.
// TieneEdicion indica si la carrera se disputó esa temporada. Si TieneEdicion es
// true y TieneResultado es false, el corredor no participó y su celda estaba vacía
// en el ranking. Si TieneResultado es true, Resultado contiene los puntos obtenidos,
// incluidos cero puntos.
type RegistroHistorico struct {
	Temporada      Temporada
	NombreEdicion  string
	TieneEdicion   bool
	TieneResultado bool
	Resultado      Resultado
}

// HistoricoPuntuaciones actúa como agregado raíz del dominio para coordinar el histórico
// de carreras, ediciones y corredores según los requerimientos de HU001 e issue #15.
// Utiliza el objeto valor Nombre como clave para preservar los tipos de dominio.
type HistoricoPuntuaciones struct {
	corredores map[Nombre]*Corredor
	carreras   map[string]*Carrera
}

// NuevoHistoricoPuntuaciones inicializa una nueva instancia del histórico de puntuaciones.
func NuevoHistoricoPuntuaciones() *HistoricoPuntuaciones {
	return &HistoricoPuntuaciones{
		corredores: make(map[Nombre]*Corredor),
		carreras:   make(map[string]*Carrera),
	}
}

// RegistrarCorredor añade un corredor al histórico garantizando su unicidad por nombre normalizado.
// Si el corredor ya existe en el histórico, devuelve ErrorCorredorYaRegistrado preservando los datos previos.
func (h *HistoricoPuntuaciones) RegistrarCorredor(corredor *Corredor) error {
	if corredor == nil {
		return ErrParametroNulo
	}
	nombre := corredor.Nombre()
	if _, existe := h.corredores[nombre]; existe {
		return ErrorCorredorYaRegistrado{NombreCorredor: nombre.Valor()}
	}
	h.corredores[nombre] = corredor
	return nil
}

// RegistrarCarrera añade una carrera al histórico por su identificador.
// Si la carrera ya existe en el histórico, devuelve ErrorCarreraYaRegistrada preservando los datos previos.
func (h *HistoricoPuntuaciones) RegistrarCarrera(carrera *Carrera) error {
	if carrera == nil {
		return ErrParametroNulo
	}
	id := carrera.ID()
	if _, existe := h.carreras[id]; existe {
		return ErrorCarreraYaRegistrada{CarreraID: id}
	}
	h.carreras[id] = carrera
	return nil
}

// ObtenerCorredor busca un corredor por su objeto valor Nombre.
func (h *HistoricoPuntuaciones) ObtenerCorredor(nombre Nombre) (*Corredor, bool) {
	c, existe := h.corredores[nombre]
	return c, existe
}

// ObtenerCarrera busca una carrera por su identificador único.
func (h *HistoricoPuntuaciones) ObtenerCarrera(id string) (*Carrera, bool) {
	c, existe := h.carreras[id]
	return c, existe
}

// ConsultarHistoricoCorredorEnCarrera resuelve la necesidad central de HU001:
// permite a Alfonso consultar los puntos obtenidos por un corredor en una carrera concreta
// a lo largo de un conjunto de temporadas (típicamente los 2 años anteriores más el actual, issue #12).
func (h *HistoricoPuntuaciones) ConsultarHistoricoCorredorEnCarrera(
	corredor *Corredor,
	carreraID string,
	temporadas []Temporada,
) ([]RegistroHistorico, error) {
	if corredor == nil {
		return nil, ErrParametroNulo
	}
	carrera, existe := h.carreras[carreraID]
	if !existe {
		return nil, ErrorCarreraInvalida{Motivo: "la carrera solicitada no existe en el histórico"}
	}

	registros := make([]RegistroHistorico, 0, len(temporadas))
	for _, temp := range temporadas {
		edicion, tieneEd := carrera.EdicionEnTemporada(temp)
		if !tieneEd {
			// Carrera sin edición en esa temporada (por ejemplo, carrera nueva, issue #15)
			registros = append(registros, RegistroHistorico{
				Temporada:      temp,
				NombreEdicion:  "",
				TieneEdicion:   false,
				TieneResultado: false,
				Resultado:      Resultado{},
			})
			continue
		}

		res, tieneRes := edicion.ResultadoDe(corredor)
		if !tieneRes {
			// El corredor no participó en esta edición: celda vacía en el ranking.
			registros = append(registros, RegistroHistorico{
				Temporada:      temp,
				NombreEdicion:  edicion.NombreOficial(),
				TieneEdicion:   true,
				TieneResultado: false,
				Resultado:      Resultado{},
			})
		} else {
			// El corredor tiene un resultado registrado y participó en la edición.
			registros = append(registros, RegistroHistorico{
				Temporada:      temp,
				NombreEdicion:  edicion.NombreOficial(),
				TieneEdicion:   true,
				TieneResultado: true,
				Resultado:      res,
			})
		}
	}

	return registros, nil
}

// CalcularTotalPuntos calcula la suma total de puntos obtenidos por un corredor en una temporada (issue #11).
// Esta información no se almacena en el modelo de datos, sino que se calcula bajo demanda sumando
// los puntos de todas las ediciones disputadas en esa temporada en las que participó el corredor.
func (h *HistoricoPuntuaciones) CalcularTotalPuntos(corredor *Corredor, temporada Temporada) int {
	if corredor == nil {
		return 0
	}
	total := 0
	for _, carrera := range h.carreras {
		if edicion, existe := carrera.EdicionEnTemporada(temporada); existe {
			if res, ok := edicion.ResultadoDe(corredor); ok {
				total += res.Puntos().Valor()
			}
		}
	}
	return total
}
