package escaner

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Casos de prueba para literales de cadena del escáner.
var casosCadenas = []casoDeEscaneo{
	{
		nombre: "cadena simple",
		fuente: `"hola"`,
		esperado: []tokenEsperado{
			{token.STRING, `"hola"`, "hola", 1},
		},
	},
	{
		nombre: "cadena vacía",
		fuente: `""`,
		esperado: []tokenEsperado{
			{token.STRING, `""`, "", 1},
		},
	},
	{
		nombre: "el contenido no se interpreta",
		fuente: `"// esto no es un comentario, y + no es un operador"`,
		esperado: []tokenEsperado{
			{
				token.STRING,
				`"// esto no es un comentario, y + no es un operador"`,
				"// esto no es un comentario, y + no es un operador",
				1,
			},
		},
	},
	{
		nombre: "los caracteres multibyte pasan intactos",
		fuente: `"ñandú en Belén"`,
		esperado: []tokenEsperado{
			{token.STRING, `"ñandú en Belén"`, "ñandú en Belén", 1},
		},
	},
	{
		nombre: "cadena multilínea",
		fuente: "\"primera\nsegunda\"",
		esperado: []tokenEsperado{
			{token.STRING, "\"primera\nsegunda\"", "primera\nsegunda", 1},
		},
	},
	{
		nombre: "la cadena multilínea mueve el contador de líneas",
		fuente: "\"primera\nsegunda\"+",
		esperado: []tokenEsperado{
			{token.STRING, "\"primera\nsegunda\"", "primera\nsegunda", 1},
			{token.PLUS, "+", nil, 2},
		},
	},
	{
		nombre: "dos cadenas seguidas",
		fuente: `"a" "b"`,
		esperado: []tokenEsperado{
			{token.STRING, `"a"`, "a", 1},
			{token.STRING, `"b"`, "b", 1},
		},
	},
	{
		nombre:   "cadena sin cerrar",
		fuente:   `"hola`,
		errores:  []string{"cadena sin cerrar"},
		esperado: []tokenEsperado{},
	},
	{
		nombre:  "la cadena sin cerrar se reporta en la línea donde abre",
		fuente:  "+\n\"hola\nmundo",
		errores: []string{"[línea 2] Error de escaneo: cadena sin cerrar"},
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 1},
		},
	},
	{
		nombre:   "las comillas simples no delimitan cadenas",
		fuente:   `'`,
		errores:  []string{"carácter inesperado"},
		esperado: []tokenEsperado{},
	},
}

func TestLiteralesDeCadena(t *testing.T) {
	correrCasos(t, casosCadenas)
}
