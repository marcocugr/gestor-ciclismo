package dominio_test

import (
	"errors"
	"testing"

	"github.com/marcocugr/gestor-ciclismo/dominio"
)

// TestPuntosValidosYNegativos verifica las invariantes sobre el objeto valor Puntos (issue #7, #13).
func TestPuntosValidosYNegativos(t *testing.T) {
	// Puntos válidos: mayor o igual que cero
	pts0, err := dominio.NuevosPuntos(0)
	if err != nil {
		t.Fatalf("se esperaba que 0 puntos fuera válido, se obtuvo error: %v", err)
	}
	if pts0.Valor() != 0 {
		t.Errorf("se esperaba valor 0, se obtuvo %d", pts0.Valor())
	}

	pts129, err := dominio.NuevosPuntos(129)
	if err != nil {
		t.Fatalf("se esperaba que 129 puntos fuera válido, se obtuvo error: %v", err)
	}
	if pts129.Valor() != 129 {
		t.Errorf("se esperaba valor 129, se obtuvo %d", pts129.Valor())
	}

	// Puntos inválidos: negativos
	_, errNeg := dominio.NuevosPuntos(-1)
	if !errors.Is(errNeg, dominio.ErrPuntosNegativos) {
		t.Errorf("se esperaba ErrPuntosNegativos para valor -1, se obtuvo: %v", errNeg)
	}
}

// TestDiferenciaEntreNoParticiparY0Puntos verifica la distinción crítica del ranking (issue #7).
func TestDiferenciaEntreNoParticiparY0Puntos(t *testing.T) {
	noPart := dominio.ResultadoNoParticipo()
	if noPart.Participo() {
		t.Errorf("se esperaba que Participo() fuera false para celda en blanco")
	}
	if _, ok := noPart.Puntos(); ok {
		t.Errorf("se esperaba que no tuviera puntos asociados")
	}
	if noPart.EsCeroPuntos() {
		t.Errorf("no participar no debe considerarse como 0 puntos")
	}

	pts0, err := dominio.NuevosPuntos(0)
	if err != nil {
		t.Fatalf("error creando puntos 0: %v", err)
	}
	part0 := dominio.ResultadoParticipo(pts0)
	if !part0.Participo() {
		t.Errorf("se esperaba que Participo() fuera true para quien compitió")
	}
	pts, ok := part0.Puntos()
	if !ok || pts.Valor() != 0 {
		t.Errorf("se esperaba tener 0 puntos asociados")
	}
	if !part0.EsCeroPuntos() {
		t.Errorf("se esperaba que EsCeroPuntos() fuera true")
	}

	if noPart.Igual(part0) {
		t.Errorf("no participar y participar con 0 puntos no pueden ser iguales")
	}
}

// TestNormalizacionEIdentidadNombre verifica la limpieza y equivalencia de nombres (issue #9, #13).
func TestNormalizacionEIdentidadNombre(t *testing.T) {
	// Nombres que deben normalizarse a la misma identidad
	n1, err1 := dominio.NuevoNombre("CRAUSE , PEDRI (SUB-23)(MADRID)")
	if err1 != nil {
		t.Fatalf("error creando nombre 1: %v", err1)
	}

	n2, err2 := dominio.NuevoNombre("Crause,  Pedri")
	if err2 != nil {
		t.Fatalf("error creando nombre 2: %v", err2)
	}

	if !n1.Igual(n2) {
		t.Errorf("se esperaba que n1 y n2 fueran equivalentes tras normalización. n1='%s', n2='%s'", n1.Valor(), n2.Valor())
	}

	// Normalización de tildes (LÓPEZ = LOPEZ)
	n3, err3 := dominio.NuevoNombre("LÓPEZ, ÁLVARO")
	if err3 != nil {
		t.Fatalf("error creando nombre con tildes: %v", err3)
	}
	n4, err4 := dominio.NuevoNombre("Lopez, Alvaro")
	if err4 != nil {
		t.Fatalf("error creando nombre sin tildes: %v", err4)
	}
	if !n3.Igual(n4) {
		t.Errorf("se esperaba que nombres con y sin tildes fueran iguales tras normalizar")
	}

	// Nombres inválidos: vacío o sin apellidos
	_, errVacio := dominio.NuevoNombre("   ")
	if !errors.Is(errVacio, dominio.ErrNombreVacio) {
		t.Errorf("se esperaba ErrNombreVacio, se obtuvo: %v", errVacio)
	}

	_, errSoloNombre := dominio.NuevoNombre("Pedri")
	if !errors.Is(errSoloNombre, dominio.ErrNombreSinApellidos) {
		t.Errorf("se esperaba ErrNombreSinApellidos para un único nombre sin apellidos, se obtuvo: %v", errSoloNombre)
	}
}

// TestInvarianteEdicionUnicaPorTemporada verifica que una carrera no admita dos ediciones el mismo año (issue #10, #13).
func TestInvarianteEdicionUnicaPorTemporada(t *testing.T) {
	carrera, err := dominio.NuevaCarrera("clasica-valladolid", "Clásica de Valladolid")
	if err != nil {
		t.Fatalf("error creando carrera: %v", err)
	}

	temp2025, _ := dominio.NuevaTemporada(2025)
	ed1, _ := dominio.NuevaEdicion(carrera.ID(), temp2025, "IV Clásica de Valladolid")
	ed2, _ := dominio.NuevaEdicion(carrera.ID(), temp2025, "IV Clásica de Valladolid Repetida")

	if err := carrera.AgregarEdicion(ed1); err != nil {
		t.Fatalf("no se esperaba error al agregar primera edición: %v", err)
	}

	errDuplicado := carrera.AgregarEdicion(ed2)
	if !errors.Is(errDuplicado, dominio.ErrEdicionDuplicadaEnTemporada) {
		t.Errorf("se esperaba ErrEdicionDuplicadaEnTemporada al registrar dos ediciones en 2025, se obtuvo: %v", errDuplicado)
	}
}

// TestInvarianteResultadoUnicoPorCorredor verifica que no haya resultados duplicados en una edición (issue #13).
func TestInvarianteResultadoUnicoPorCorredor(t *testing.T) {
	temp, _ := dominio.NuevaTemporada(2026)
	ed, _ := dominio.NuevaEdicion("vuelta-valencia-etapa1", temp, "1ª Etapa Vuelta a Valencia")

	nombre, _ := dominio.NuevoNombre("GARCIA, JUAN")
	corredor, _ := dominio.NuevoCorredor(nombre)

	pts, _ := dominio.NuevosPuntos(15)
	res1 := dominio.ResultadoParticipo(pts)
	res2 := dominio.ResultadoNoParticipo()

	if err := ed.RegistrarResultado(corredor, res1); err != nil {
		t.Fatalf("no se esperaba error registrando primer resultado: %v", err)
	}

	errDuplicado := ed.RegistrarResultado(corredor, res2)
	if !errors.Is(errDuplicado, dominio.ErrResultadoDuplicado) {
		t.Errorf("se esperaba ErrResultadoDuplicado, se obtuvo: %v", errDuplicado)
	}
}

// TestCalculoTotalPuntosHistorico verifica que el total se calcula dinámicamente y no se persiste (issue #11).
func TestCalculoTotalPuntosHistorico(t *testing.T) {
	historico := dominio.NuevoHistoricoPuntuaciones()

	nom, _ := dominio.NuevoNombre("CRAUSE, PEDRI")
	corredor, _ := dominio.NuevoCorredor(nom)
	_ = historico.RegistrarCorredor(corredor)

	c1, _ := dominio.NuevaCarrera("carrera-1", "Carrera 1")
	c2, _ := dominio.NuevaCarrera("carrera-2", "Carrera 2")

	temp2025, _ := dominio.NuevaTemporada(2025)
	ed1, _ := dominio.NuevaEdicion(c1.ID(), temp2025, "Carrera 1 2025")
	ed2, _ := dominio.NuevaEdicion(c2.ID(), temp2025, "Carrera 2 2025")

	pts10, _ := dominio.NuevosPuntos(10)
	pts25, _ := dominio.NuevosPuntos(25)

	_ = ed1.RegistrarResultado(corredor, dominio.ResultadoParticipo(pts10))
	_ = ed2.RegistrarResultado(corredor, dominio.ResultadoParticipo(pts25))

	_ = c1.AgregarEdicion(ed1)
	_ = c2.AgregarEdicion(ed2)

	_ = historico.RegistrarCarrera(c1)
	_ = historico.RegistrarCarrera(c2)

	total2025 := historico.CalcularTotalPuntos(corredor, temp2025)
	if total2025 != 35 {
		t.Errorf("se esperaba total calculado de 35 puntos, se obtuvo %d", total2025)
	}

	temp2024, _ := dominio.NuevaTemporada(2024)
	total2024 := historico.CalcularTotalPuntos(corredor, temp2024)
	if total2024 != 0 {
		t.Errorf("se esperaba total calculado de 0 puntos para temporada sin carreras, se obtuvo %d", total2024)
	}
}

// TestConsultaHistoricoHU001 verifica la consulta de histórico de un corredor por carrera (HU001, issue #12).
func TestConsultaHistoricoHU001(t *testing.T) {
	historico := dominio.NuevoHistoricoPuntuaciones()

	nom, _ := dominio.NuevoNombre("RODRIGUEZ, MARCOS")
	corredor, _ := dominio.NuevoCorredor(nom)
	_ = historico.RegistrarCorredor(corredor)

	carrera, _ := dominio.NuevaCarrera("clasica-valladolid", "Clásica de Valladolid")

	temp2024, _ := dominio.NuevaTemporada(2024)
	temp2025, _ := dominio.NuevaTemporada(2025)
	temp2026, _ := dominio.NuevaTemporada(2026)

	// En 2024 participó y obtuvo 0 puntos
	ed2024, _ := dominio.NuevaEdicion(carrera.ID(), temp2024, "III Clásica de Valladolid")
	pts0, _ := dominio.NuevosPuntos(0)
	_ = ed2024.RegistrarResultado(corredor, dominio.ResultadoParticipo(pts0))
	_ = carrera.AgregarEdicion(ed2024)

	// En 2025 participó y obtuvo 45 puntos
	ed2025, _ := dominio.NuevaEdicion(carrera.ID(), temp2025, "IV Clásica de Valladolid")
	pts45, _ := dominio.NuevosPuntos(45)
	_ = ed2025.RegistrarResultado(corredor, dominio.ResultadoParticipo(pts45))
	_ = carrera.AgregarEdicion(ed2025)

	// En 2026 no participó en esa carrera (no hay resultado para él en ed2026)
	ed2026, _ := dominio.NuevaEdicion(carrera.ID(), temp2026, "V Clásica de Valladolid")
	_ = carrera.AgregarEdicion(ed2026)

	_ = historico.RegistrarCarrera(carrera)

	temporadasHistorico := []dominio.Temporada{temp2024, temp2025, temp2026}
	registros, err := historico.ConsultarHistoricoCorredorEnCarrera(corredor, carrera.ID(), temporadasHistorico)
	if err != nil {
		t.Fatalf("error consultando histórico: %v", err)
	}

	if len(registros) != 3 {
		t.Fatalf("se esperaban 3 registros de temporadas, se obtuvieron %d", len(registros))
	}

	// 2024: participó con 0 puntos
	if !registros[0].Resultado.Participo() || !registros[0].Resultado.EsCeroPuntos() {
		t.Errorf("en 2024 debía constar que participó y obtuvo 0 puntos")
	}

	// 2025: participó con 45 puntos
	if !registros[1].Resultado.Participo() {
		t.Errorf("en 2025 debía constar participación")
	}
	if p, ok := registros[1].Resultado.Puntos(); !ok || p.Valor() != 45 {
		t.Errorf("en 2025 se esperaban 45 puntos, se obtuvo: %v", p.Valor())
	}

	// 2026: no participó (celda en blanco)
	if registros[2].Resultado.Participo() {
		t.Errorf("en 2026 debía constar que no participó")
	}
}
