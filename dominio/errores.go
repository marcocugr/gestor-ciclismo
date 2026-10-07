package dominio

import "errors"

// Errores de dominio derivados de las reglas e invariantes analizadas en HU001 e issues #7-#15.
var (
	// ErrPuntosNegativos indica que se ha intentado instanciar una puntuación menor que cero (issue #7, #13).
	ErrPuntosNegativos = errors.New("los puntos no pueden ser negativos")

	// ErrNombreVacio indica que la cadena del nombre no contiene caracteres válidos (issue #9, #13).
	ErrNombreVacio = errors.New("el nombre del corredor no puede estar vacío")

	// ErrNombreSinApellidos indica que el nombre no incluye tanto nombre como apellidos tras normalizar (issue #9, #13).
	ErrNombreSinApellidos = errors.New("el corredor debe incluir nombre y apellidos")

	// ErrTemporadaInvalida indica un año de temporada no admisible en el dominio (issue #12, #13).
	ErrTemporadaInvalida = errors.New("el año de la temporada es inválido")

	// ErrCarreraSinID indica que la carrera carece de identificador único (issue #10, #13).
	ErrCarreraSinID = errors.New("el identificador de la carrera no puede estar vacío")

	// ErrCarreraSinNombre indica que la carrera no tiene un nombre descriptivo asignado (issue #10, #13).
	ErrCarreraSinNombre = errors.New("el nombre de la carrera no puede estar vacío")

	// ErrEdicionSinNombre indica que la edición carece de nombre oficial publicado en el ranking (issue #8, #13).
	ErrEdicionSinNombre = errors.New("el nombre oficial de la edición no puede estar vacío")

	// ErrEdicionDuplicadaEnTemporada indica que ya existe una edición registrada para esa carrera en la temporada (issue #10, #13).
	ErrEdicionDuplicadaEnTemporada = errors.New("una carrera no puede tener más de una edición en la misma temporada")

	// ErrResultadoDuplicado indica que ya se ha registrado un resultado para ese corredor en la edición (issue #13).
	ErrResultadoDuplicado = errors.New("el corredor ya tiene un resultado registrado en esta edición")

	// ErrCorredorNulo indica que se ha recibido una referencia nula a un corredor (issue #13).
	ErrCorredorNulo = errors.New("el corredor no puede ser nulo")

	// ErrEdicionNula indica que se ha recibido una referencia nula a una edición (issue #13).
	ErrEdicionNula = errors.New("la edición no puede ser nula")
)
