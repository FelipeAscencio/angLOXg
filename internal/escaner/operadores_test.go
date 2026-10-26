package escaner

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Casos de prueba para tokens de un solo carácter.
var casosUnCaracter = []casoDeEscaneo{
	{
		nombre: "agrupadores",
		fuente: "(){}",
		esperado: []tokenEsperado{
			{token.LEFT_PAREN, "(", nil, 1},
			{token.RIGHT_PAREN, ")", nil, 1},
			{token.LEFT_BRACE, "{", nil, 1},
			{token.RIGHT_BRACE, "}", nil, 1},
		},
	},
	{
		nombre: "aritméticos",
		fuente: "+ - * / %",
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 1},
			{token.MINUS, "-", nil, 1},
			{token.STAR, "*", nil, 1},
			{token.SLASH, "/", nil, 1},
			{token.PERCENT, "%", nil, 1},
		},
	},
	{
		nombre: "separadores",
		fuente: ",;",
		esperado: []tokenEsperado{
			{token.COMMA, ",", nil, 1},
			{token.SEMICOLON, ";", nil, 1},
		},
	},
	{
		nombre:   "carácter inesperado",
		fuente:   "@",
		errores:  []string{"carácter inesperado"},
		esperado: []tokenEsperado{},
	},
	{
		nombre: "corchetes",
		fuente: "[]",
		esperado: []tokenEsperado{
			{token.LEFT_BRACKET, "[", nil, 1},
			{token.RIGHT_BRACKET, "]", nil, 1},
		},
	},
	{
		nombre: "indexado de un array",
		fuente: "a[0]",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "a", nil, 1},
			{token.LEFT_BRACKET, "[", nil, 1},
			{token.NUMBER, "0", float64(0), 1},
			{token.RIGHT_BRACKET, "]", nil, 1},
		},
	},
	{
		nombre: "punto",
		fuente: ".",
		esperado: []tokenEsperado{
			{token.DOT, ".", nil, 1},
		},
	},
	{
		nombre: "acceso a una propiedad",
		fuente: "p.x",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "p", nil, 1},
			{token.DOT, ".", nil, 1},
			{token.IDENTIFIER, "x", nil, 1},
		},
	},
	{
		nombre:  "los errores se acumulan y el escaneo sigue",
		fuente:  "@+#",
		errores: []string{"carácter inesperado", "carácter inesperado"},
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 1},
		},
	},
	{
		nombre: "Operador ternario",
		fuente: "? :",
		esperado: []tokenEsperado{
			{token.QUESTION, "?", nil, 1},
			{token.COLON, ":", nil, 1},
		},
	},
}

// Casos de prueba para tokens de uno o dos caracteres.
var casosDosCaracteres = []casoDeEscaneo{
	{
		nombre: "versión de un carácter",
		fuente: "! = < >",
		esperado: []tokenEsperado{
			{token.BANG, "!", nil, 1},
			{token.EQUAL, "=", nil, 1},
			{token.LESS, "<", nil, 1},
			{token.GREATER, ">", nil, 1},
		},
	},
	{
		nombre: "versión de dos caracteres",
		fuente: "!= == <= >=",
		esperado: []tokenEsperado{
			{token.BANG_EQUAL, "!=", nil, 1},
			{token.EQUAL_EQUAL, "==", nil, 1},
			{token.LESS_EQUAL, "<=", nil, 1},
			{token.GREATER_EQUAL, ">=", nil, 1},
		},
	},
	{
		nombre: "se consume el lexema más largo posible",
		fuente: "!==",
		esperado: []tokenEsperado{
			{token.BANG_EQUAL, "!=", nil, 1},
			{token.EQUAL, "=", nil, 1},
		},
	},
	{
		nombre: "sin espacios entre operadores",
		fuente: "<=<",
		esperado: []tokenEsperado{
			{token.LESS_EQUAL, "<=", nil, 1},
			{token.LESS, "<", nil, 1},
		},
	},
}

// Casos de prueba para espacios en blanco y comentarios.
var casosComentarios = []casoDeEscaneo{
	{
		nombre:   "comentario hasta el fin de línea",
		fuente:   "// esto no es código",
		esperado: []tokenEsperado{},
	},
	{
		nombre: "el comentario no se come la línea siguiente",
		fuente: "// comentario\n+",
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 2},
		},
	},
	{
		nombre: "comentario al final de una línea con código",
		fuente: "+ // comentario",
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 1},
		},
	},
	{
		nombre: "la barra sola sigue siendo división",
		fuente: "/ /",
		esperado: []tokenEsperado{
			{token.SLASH, "/", nil, 1},
			{token.SLASH, "/", nil, 1},
		},
	},
	{
		nombre: "los espacios en blanco se descartan",
		fuente: " \t\r+",
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 1},
		},
	},
	{
		nombre: "el salto de línea incrementa el contador",
		fuente: "+\n-\n*",
		esperado: []tokenEsperado{
			{token.PLUS, "+", nil, 1},
			{token.MINUS, "-", nil, 2},
			{token.STAR, "*", nil, 3},
		},
	},
}

func TestTokensDeUnCaracter(t *testing.T) {
	correrCasos(t, casosUnCaracter)
}

func TestTokensDeUnoODosCaracteres(t *testing.T) {
	correrCasos(t, casosDosCaracteres)
}

func TestEspaciosYComentarios(t *testing.T) {
	correrCasos(t, casosComentarios)
}
