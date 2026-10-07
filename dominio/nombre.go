package dominio

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// parentesisRegex elimina cualquier contenido entre paréntesis, como categorías (SUB-23) o regiones (MADRID).
	parentesisRegex = regexp.MustCompile(`\([^)]*\)`)
	// espaciosRegex colapsa múltiples espacios en blanco consecutivos en uno solo.
	espaciosRegex = regexp.MustCompile(`\s+`)
)

// Nombre representa el nombre completo normalizado de un corredor ciclista (issue #9, #15).
// Es un objeto valor inmutable que constituye la base de identidad del corredor entre temporadas.
type Nombre struct {
	valor string
}

// NuevoNombre normaliza y valida un nombre sin formato procedente del ranking oficial.
// Realiza los siguientes pasos de normalización (acordados en issues #9 y #13):
//  1. Elimina texto entre paréntesis (categoría y comunidad autónoma).
//  2. Elimina acentos y tildes (p. ej. LÓPEZ -> LOPEZ).
//  3. Corrige espacios alrededor de comas (p. ej. "CRAUSE , PEDRI" -> "CRAUSE, PEDRI").
//  4. Colapsa espacios redundantes y recorta extremos.
//  5. Convierte todo a mayúsculas para comparación uniforme.
//  6. Valida que no esté vacío y contenga tanto nombre como apellidos.
func NuevoNombre(raw string) (Nombre, error) {
	// 1. Eliminar paréntesis y su contenido
	limpio := parentesisRegex.ReplaceAllString(raw, " ")

	// 2. Normalizar acentos y diacríticos
	limpio = quitarAcentos(limpio)

	// 3. Normalizar comas y espacios alrededor
	limpio = strings.ReplaceAll(limpio, " ,", ",")
	limpio = strings.ReplaceAll(limpio, ", ", ",")

	// 4. Colapsar espacios y recortar
	limpio = espaciosRegex.ReplaceAllString(limpio, " ")
	limpio = strings.TrimSpace(limpio)

	// 5. Convertir a mayúsculas
	limpio = strings.ToUpper(limpio)

	if limpio == "" {
		return Nombre{}, ErrNombreVacio
	}

	// 6. Validar que incluya nombre y apellidos
	if !tieneNombreYApellidos(limpio) {
		return Nombre{}, ErrNombreSinApellidos
	}

	return Nombre{valor: limpio}, nil
}

// Valor devuelve la cadena normalizada del nombre.
func (n Nombre) Valor() string {
	return n.valor
}

// Igual compara dos nombres normalizados por su valor exacto.
func (n Nombre) Igual(otro Nombre) bool {
	return n.valor == otro.valor
}

// quitarAcentos transforma vocales con tilde y diacríticos a sus equivalentes básicos.
func quitarAcentos(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case 'Á', 'À', 'Â', 'Ä':
			b.WriteRune('A')
		case 'á', 'à', 'â', 'ä':
			b.WriteRune('a')
		case 'É', 'È', 'Ê', 'Ë':
			b.WriteRune('E')
		case 'é', 'è', 'ê', 'ë':
			b.WriteRune('e')
		case 'Í', 'Ì', 'Î', 'Ï':
			b.WriteRune('I')
		case 'í', 'ì', 'î', 'ï':
			b.WriteRune('i')
		case 'Ó', 'Ò', 'Ô', 'Ö':
			b.WriteRune('O')
		case 'ó', 'ò', 'ô', 'ö':
			b.WriteRune('o')
		case 'Ú', 'Ù', 'Û', 'Ü':
			b.WriteRune('U')
		case 'ú', 'ù', 'û', 'ü':
			b.WriteRune('u')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// tieneNombreYApellidos comprueba si la cadena normalizada cuenta con al menos dos componentes.
func tieneNombreYApellidos(s string) bool {
	if strings.Contains(s, ",") {
		partes := strings.Split(s, ",")
		if len(partes) >= 2 && len(strings.TrimSpace(partes[0])) > 0 && len(strings.TrimSpace(partes[1])) > 0 {
			return true
		}
	}
	palabras := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == ','
	})
	return len(palabras) >= 2
}
