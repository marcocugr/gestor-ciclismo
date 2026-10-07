package dominio

// RegistroHistorico representa el dato de consulta para una temporada concreta en una carrera (HU001).
type RegistroHistorico struct {
	Temporada     Temporada
	NombreEdicion string
	Resultado     Resultado
	TieneEdicion  bool
}

// HistoricoPuntuaciones actúa como agregado raíz del dominio para coordinar el histórico
// de carreras, ediciones y corredores según los requerimientos de HU001 e issue #15.
type HistoricoPuntuaciones struct {
	corredores map[string]*Corredor // Indexado por Nombre.Valor()
	carreras   map[string]*Carrera  // Indexado por Carrera.ID()
}

// NuevoHistoricoPuntuaciones inicializa una nueva instancia del histórico de puntuaciones.
func NuevoHistoricoPuntuaciones() *HistoricoPuntuaciones {
	return &HistoricoPuntuaciones{
		corredores: make(map[string]*Corredor),
		carreras:   make(map[string]*Carrera),
	}
}

// RegistrarCorredor añade un corredor al histórico garantizando su unicidad por nombre normalizado.
func (h *HistoricoPuntuaciones) RegistrarCorredor(corredor *Corredor) error {
	if corredor == nil {
		return ErrCorredorNulo
	}
	h.corredores[corredor.Nombre().Valor()] = corredor
	return nil
}

// RegistrarCarrera añade una carrera al histórico por su identificador.
func (h *HistoricoPuntuaciones) RegistrarCarrera(carrera *Carrera) error {
	if carrera == nil {
		return ErrCarreraSinID
	}
	h.carreras[carrera.ID()] = carrera
	return nil
}

// ObtenerCorredor busca un corredor por su objeto valor Nombre.
func (h *HistoricoPuntuaciones) ObtenerCorredor(nombre Nombre) (*Corredor, bool) {
	c, existe := h.corredores[nombre.Valor()]
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
		return nil, ErrCorredorNulo
	}
	carrera, existe := h.carreras[carreraID]
	if !existe {
		return nil, ErrCarreraSinID
	}

	registros := make([]RegistroHistorico, 0, len(temporadas))
	for _, temp := range temporadas {
		edicion, tieneEd := carrera.EdicionEnTemporada(temp)
		if !tieneEd {
			// Carrera sin edición en esa temporada (por ejemplo, carrera nueva, issue #15)
			registros = append(registros, RegistroHistorico{
				Temporada:     temp,
				NombreEdicion: "",
				Resultado:     ResultadoNoParticipo(),
				TieneEdicion:  false,
			})
			continue
		}

		res, tieneRes := edicion.ResultadoDe(corredor)
		if !tieneRes {
			// El corredor no aparece en esa edición (no participó / celda en blanco, issue #7)
			registros = append(registros, RegistroHistorico{
				Temporada:     temp,
				NombreEdicion: edicion.NombreOficial(),
				Resultado:     ResultadoNoParticipo(),
				TieneEdicion:  true,
			})
		} else {
			registros = append(registros, RegistroHistorico{
				Temporada:     temp,
				NombreEdicion: edicion.NombreOficial(),
				Resultado:     res,
				TieneEdicion:  true,
			})
		}
	}

	return registros, nil
}

// CalcularTotalPuntos calcula la suma total de puntos obtenidos por un corredor en una temporada (issue #11).
// Esta información no se almacena en el modelo de datos, sino que se calcula bajo demanda sumando
// los puntos de todas las ediciones disputadas en esa temporada.
func (h *HistoricoPuntuaciones) CalcularTotalPuntos(corredor *Corredor, temporada Temporada) int {
	if corredor == nil {
		return 0
	}
	total := 0
	for _, carrera := range h.carreras {
		if edicion, existe := carrera.EdicionEnTemporada(temporada); existe {
			if res, ok := edicion.ResultadoDe(corredor); ok && res.Participo() {
				if pts, tienePts := res.Puntos(); tienePts {
					total += pts.Valor()
				}
			}
		}
	}
	return total
}
