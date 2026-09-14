package escaner

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Identificadores: arrancan con letra ASCII o guion bajo, y siguen con letras,
// dígitos o guiones bajos.
var casosIdentificadores = []casoDeEscaneo{
	{
		nombre: "identificador simple",
		fuente: "contador",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "contador", nil, 1},
		},
	},
	{
		nombre: "puede arrancar con guion bajo",
		fuente: "_privado",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "_privado", nil, 1},
		},
	},
	{
		nombre: "puede contener dígitos y guiones bajos",
		fuente: "version_2",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "version_2", nil, 1},
		},
	},
	{
		nombre: "no puede arrancar con un dígito",
		// "1abc" no es un identificador inválido: son dos tokens distintos.
		fuente: "1abc",
		esperado: []tokenEsperado{
			{token.NUMBER, "1", float64(1), 1},
			{token.IDENTIFIER, "abc", nil, 1},
		},
	},
	{
		nombre: "varios identificadores separados",
		fuente: "a b_1 _c",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "a", nil, 1},
			{token.IDENTIFIER, "b_1", nil, 1},
			{token.IDENTIFIER, "_c", nil, 1},
		},
	},
	{
		nombre: "class, this y super son identificadores comunes",
		fuente: "class this super",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "class", nil, 1},
			{token.IDENTIFIER, "this", nil, 1},
			{token.IDENTIFIER, "super", nil, 1},
		},
	},
}

// Palabras reservadas: se reconocen recién después de consumir el lexema
// completo, así que nunca se parte un identificador al medio.
var casosPalabrasClave = []casoDeEscaneo{
	{
		nombre: "todas las palabras reservadas",
		fuente: "and else false fun for if nil or print return true var while",
		esperado: []tokenEsperado{
			{token.AND, "and", nil, 1},
			{token.ELSE, "else", nil, 1},
			// true, false y nil son palabras clave, no literales: no llevan
			// valor resuelto.
			{token.FALSE, "false", nil, 1},
			{token.FUN, "fun", nil, 1},
			{token.FOR, "for", nil, 1},
			{token.IF, "if", nil, 1},
			{token.NIL, "nil", nil, 1},
			{token.OR, "or", nil, 1},
			{token.PRINT, "print", nil, 1},
			{token.RETURN, "return", nil, 1},
			{token.TRUE, "true", nil, 1},
			{token.VAR, "var", nil, 1},
			{token.WHILE, "while", nil, 1},
		},
	},
	{
		nombre: "una palabra clave que es prefijo de un identificador",
		fuente: "orden printer variable",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "orden", nil, 1},
			{token.IDENTIFIER, "printer", nil, 1},
			{token.IDENTIFIER, "variable", nil, 1},
		},
	},
	{
		nombre: "las palabras clave distinguen mayúsculas",
		fuente: "Var IF nIl",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "Var", nil, 1},
			{token.IDENTIFIER, "IF", nil, 1},
			{token.IDENTIFIER, "nIl", nil, 1},
		},
	},
	{
		nombre: "una palabra clave pegada a un guion bajo es un identificador",
		fuente: "var_",
		esperado: []tokenEsperado{
			{token.IDENTIFIER, "var_", nil, 1},
		},
	},
}

// Programas completos, con todas las piezas del escáner trabajando juntas.
var casosPrograma = []casoDeEscaneo{
	{
		nombre: "declaración con condicional",
		fuente: "var x = 10;\nif (x >= 5) print \"grande\";",
		esperado: []tokenEsperado{
			{token.VAR, "var", nil, 1},
			{token.IDENTIFIER, "x", nil, 1},
			{token.EQUAL, "=", nil, 1},
			{token.NUMBER, "10", float64(10), 1},
			{token.SEMICOLON, ";", nil, 1},
			{token.IF, "if", nil, 2},
			{token.LEFT_PAREN, "(", nil, 2},
			{token.IDENTIFIER, "x", nil, 2},
			{token.GREATER_EQUAL, ">=", nil, 2},
			{token.NUMBER, "5", float64(5), 2},
			{token.RIGHT_PAREN, ")", nil, 2},
			{token.PRINT, "print", nil, 2},
			{token.STRING, `"grande"`, "grande", 2},
			{token.SEMICOLON, ";", nil, 2},
		},
	},
	{
		nombre: "función con módulo y comentario",
		fuente: "// devuelve el resto\nfun resto(a, b) { return a % b; }",
		esperado: []tokenEsperado{
			{token.FUN, "fun", nil, 2},
			{token.IDENTIFIER, "resto", nil, 2},
			{token.LEFT_PAREN, "(", nil, 2},
			{token.IDENTIFIER, "a", nil, 2},
			{token.COMMA, ",", nil, 2},
			{token.IDENTIFIER, "b", nil, 2},
			{token.RIGHT_PAREN, ")", nil, 2},
			{token.LEFT_BRACE, "{", nil, 2},
			{token.RETURN, "return", nil, 2},
			{token.IDENTIFIER, "a", nil, 2},
			{token.PERCENT, "%", nil, 2},
			{token.IDENTIFIER, "b", nil, 2},
			{token.SEMICOLON, ";", nil, 2},
			{token.RIGHT_BRACE, "}", nil, 2},
		},
	},
	{
		nombre:   "fuente vacío",
		fuente:   "",
		esperado: []tokenEsperado{},
	},
	{
		nombre: "condiciones lógicas",
		fuente: "while (a != nil and !b or c == true) { }",
		esperado: []tokenEsperado{
			{token.WHILE, "while", nil, 1},
			{token.LEFT_PAREN, "(", nil, 1},
			{token.IDENTIFIER, "a", nil, 1},
			{token.BANG_EQUAL, "!=", nil, 1},
			{token.NIL, "nil", nil, 1},
			{token.AND, "and", nil, 1},
			{token.BANG, "!", nil, 1},
			{token.IDENTIFIER, "b", nil, 1},
			{token.OR, "or", nil, 1},
			{token.IDENTIFIER, "c", nil, 1},
			{token.EQUAL_EQUAL, "==", nil, 1},
			{token.TRUE, "true", nil, 1},
			{token.RIGHT_PAREN, ")", nil, 1},
			{token.LEFT_BRACE, "{", nil, 1},
			{token.RIGHT_BRACE, "}", nil, 1},
		},
	},
	{
		nombre: "bucle con división y else",
		fuente: "for (;;) { x = x / 2; } else false",
		esperado: []tokenEsperado{
			{token.FOR, "for", nil, 1},
			{token.LEFT_PAREN, "(", nil, 1},
			{token.SEMICOLON, ";", nil, 1},
			{token.SEMICOLON, ";", nil, 1},
			{token.RIGHT_PAREN, ")", nil, 1},
			{token.LEFT_BRACE, "{", nil, 1},
			{token.IDENTIFIER, "x", nil, 1},
			{token.EQUAL, "=", nil, 1},
			{token.IDENTIFIER, "x", nil, 1},
			{token.SLASH, "/", nil, 1},
			{token.NUMBER, "2", float64(2), 1},
			{token.SEMICOLON, ";", nil, 1},
			{token.RIGHT_BRACE, "}", nil, 1},
			{token.ELSE, "else", nil, 1},
			{token.FALSE, "false", nil, 1},
		},
	},
}

func TestIdentificadores(t *testing.T) {
	correrCasos(t, casosIdentificadores)
}

func TestPalabrasClave(t *testing.T) {
	correrCasos(t, casosPalabrasClave)
}

func TestProgramasCompletos(t *testing.T) {
	correrCasos(t, casosPrograma)
}
