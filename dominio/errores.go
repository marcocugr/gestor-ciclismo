package dominio

import (
	"errors"
	"fmt"
)

// Errores de programación y precondiciones de API.
// No representan reglas de negocio del dominio de ciclismo de Alfonso,
// sino validaciones técnicas de invocación ante referencias nulas.
var (
	// ErrParametroNulo indica que se ha pasado una referencia nula a una función o método.
	ErrParametroNulo = errors.New("parámetro nulo no permitido")
)

// ErrorPuntosNegativos describe el intento de instanciar una puntuación con un valor menor que cero (issue #7, #13).
type ErrorPuntosNegativos struct {
	Valor int
}

func (e ErrorPuntosNegativos) Error() string {
	return fmt.Sprintf("los puntos no pueden ser negativos, se recibió %d", e.Valor)
}

// ErrorNombreInvalido describe el fallo de validación de un nombre de corredor (issue #9, #13).
type ErrorNombreInvalido struct {
	Motivo string
	Texto  string
}

func (e ErrorNombreInvalido) Error() string {
	return fmt.Sprintf("nombre de corredor inválido (%s): %q", e.Motivo, e.Texto)
}

// ErrorTemporadaInvalida describe un año de temporada no positivo (issue #12, #13).
type ErrorTemporadaInvalida struct {
	Anio int
}

func (e ErrorTemporadaInvalida) Error() string {
	return fmt.Sprintf("el año de la temporada es inválido: %d (debe ser un año positivo)", e.Anio)
}

// ErrorTemporadaNoComenzada describe una temporada futura que todavía no ha comenzado en el sistema (issue #13).
type ErrorTemporadaNoComenzada struct {
	Anio       int
	AnioActual int
}

func (e ErrorTemporadaNoComenzada) Error() string {
	return fmt.Sprintf("la temporada %d todavía no ha comenzado (año actual del sistema: %d)", e.Anio, e.AnioActual)
}

// ErrorCarreraInvalida describe que una carrera no cumple las condiciones mínimas para existir (issue #10, #13).
type ErrorCarreraInvalida struct {
	Motivo string
}

func (e ErrorCarreraInvalida) Error() string {
	return fmt.Sprintf("carrera inválida: %s", e.Motivo)
}

// ErrorEdicionInvalida describe que una edición carece de datos mínimos válidos (issue #8, #13).
type ErrorEdicionInvalida struct {
	Motivo string
}

func (e ErrorEdicionInvalida) Error() string {
	return fmt.Sprintf("edición inválida: %s", e.Motivo)
}

// ErrorEdicionDuplicadaEnTemporada indica que ya existe una edición registrada para esa carrera en la temporada (issue #10, #13).
type ErrorEdicionDuplicadaEnTemporada struct {
	CarreraID string
	Anio      int
}

func (e ErrorEdicionDuplicadaEnTemporada) Error() string {
	return fmt.Sprintf("la carrera %q ya tiene una edición registrada para la temporada %d", e.CarreraID, e.Anio)
}

// ErrorEdicionCarreraIncompatible indica que se intenta asociar una edición a una carrera distinta a la suya.
type ErrorEdicionCarreraIncompatible struct {
	CarreraReceptoraID string
	EdicionCarreraID   string
	NombreEdicion      string
}

func (e ErrorEdicionCarreraIncompatible) Error() string {
	return fmt.Sprintf("la edición %q con carrera %q no puede agregarse a la carrera receptora %q", e.NombreEdicion, e.EdicionCarreraID, e.CarreraReceptoraID)
}

// ErrorResultadoDuplicado indica que un corredor ya posee un resultado en la edición indicada (issue #13).
type ErrorResultadoDuplicado struct {
	NombreCorredor string
	NombreEdicion  string
}

func (e ErrorResultadoDuplicado) Error() string {
	return fmt.Sprintf("el corredor %q ya tiene un resultado registrado en la edición %q", e.NombreCorredor, e.NombreEdicion)
}

// ErrorCorredorYaRegistrado indica que el corredor ya existe en el histórico de puntuaciones.
type ErrorCorredorYaRegistrado struct {
	NombreCorredor string
}

func (e ErrorCorredorYaRegistrado) Error() string {
	return fmt.Sprintf("el corredor %q ya se encuentra registrado en el histórico", e.NombreCorredor)
}

// ErrorCarreraYaRegistrada indica que la carrera ya existe en el histórico de puntuaciones.
type ErrorCarreraYaRegistrada struct {
	CarreraID string
}

func (e ErrorCarreraYaRegistrada) Error() string {
	return fmt.Sprintf("la carrera %q ya se encuentra registrada en el histórico", e.CarreraID)
}
