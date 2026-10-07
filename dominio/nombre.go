package dominio

import (
	"regexp"
	"strings"
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
//  3. Colapsa espacios redundantes y recorta extremos.
//  4. Valida estrictamente la presencia de nombre y apellidos: si hay coma, debe existir texto
//     válido a ambos lados tras limpiar espacios; si no hay coma, debe haber al menos dos palabras.
//  5. Normaliza comas con el formato "APELLIDOS, NOMBRE" sin espacios residuales.
//  6. Convierte todo a mayúsculas para comparación uniforme.
func NuevoNombre(raw string) (Nombre, error) {
	// 1. Eliminar paréntesis y su contenido
	limpio := parentesisRegex.ReplaceAllString(raw, " ")

	// 2. Normalizar acentos y diacríticos
	limpio = quitarAcentos(limpio)

	// 3. Colapsar espacios y recortar
	limpio = espaciosRegex.ReplaceAllString(limpio, " ")
	limpio = strings.TrimSpace(limpio)

	if limpio == "" {
		return Nombre{}, ErrorNombreInvalido{
			Motivo: "el nombre no puede estar vacío",
			Texto:  raw,
		}
	}

	// 4. Validar que incluya tanto nombre como apellidos
	if strings.Contains(limpio, ",") {
		partes := strings.Split(limpio, ",")
		if len(partes) != 2 {
			return Nombre{}, ErrorNombreInvalido{
				Motivo: "formato con coma inválido, debe contener apellidos y nombre separados por una única coma",
				Texto:  raw,
			}
		}
		lado1 := strings.TrimSpace(espaciosRegex.ReplaceAllString(partes[0], " "))
		lado2 := strings.TrimSpace(espaciosRegex.ReplaceAllString(partes[1], " "))
		if lado1 == "" || lado2 == "" {
			return Nombre{}, ErrorNombreInvalido{
				Motivo: "el corredor debe incluir tanto apellidos como nombre a ambos lados de la coma",
				Texto:  raw,
			}
		}
		limpio = lado1 + ", " + lado2
	} else {
		palabras := strings.Fields(limpio)
		if len(palabras) < 2 {
			return Nombre{}, ErrorNombreInvalido{
				Motivo: "el corredor debe incluir al menos nombre y apellidos",
				Texto:  raw,
			}
		}
		limpio = strings.Join(palabras, " ")
	}

	// 5. Convertir a mayúsculas
	limpio = strings.ToUpper(limpio)

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
