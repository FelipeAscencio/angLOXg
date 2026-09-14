package escaner

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Literales numéricos: todos, con o sin parte decimal, resuelven a un float64.
var casosNumeros = []casoDeEscaneo{
	{
		nombre: "entero",
		fuente: "123",
		esperado: []tokenEsperado{
			// El literal es float64, no int.
			{token.NUMBER, "123", float64(123), 1},
		},
	},
	{
		nombre: "cero",
		fuente: "0",
		esperado: []tokenEsperado{
			{token.NUMBER, "0", float64(0), 1},
		},
	},
	{
		nombre: "decimal",
		fuente: "3.14159",
		esperado: []tokenEsperado{
			{token.NUMBER, "3.14159", 3.14159, 1},
		},
	},
	{
		nombre: "decimal con parte entera cero",
		fuente: "0.5",
		esperado: []tokenEsperado{
			{token.NUMBER, "0.5", 0.5, 1},
		},
	},
	{
		nombre: "el signo menos no es parte del literal",
		// El menos unario lo resuelve el parser, no el escáner.
		fuente: "-7",
		esperado: []tokenEsperado{
			{token.MINUS, "-", nil, 1},
			{token.NUMBER, "7", float64(7), 1},
		},
	},
	{
		nombre: "números pegados a operadores",
		fuente: "1+2",
		esperado: []tokenEsperado{
			{token.NUMBER, "1", float64(1), 1},
			{token.PLUS, "+", nil, 1},
			{token.NUMBER, "2", float64(2), 1},
		},
	},
	{
		nombre: "el módulo entre dos números",
		fuente: "10 % 3",
		esperado: []tokenEsperado{
			{token.NUMBER, "10", float64(10), 1},
			{token.PERCENT, "%", nil, 1},
			{token.NUMBER, "3", float64(3), 1},
		},
	},
	{
		nombre: "los números llevan la cuenta de líneas",
		fuente: "1\n2",
		esperado: []tokenEsperado{
			{token.NUMBER, "1", float64(1), 1},
			{token.NUMBER, "2", float64(2), 2},
		},
	},
	{
		nombre: "un número no puede terminar en punto",
		// El punto no se traga: se cierra el número y el punto suelto, que no
		// es un token del lenguaje, queda como error.
		fuente:  "123.",
		errores: []string{"carácter inesperado"},
		esperado: []tokenEsperado{
			{token.NUMBER, "123", float64(123), 1},
		},
	},
	{
		nombre:  "un número no puede empezar con punto",
		fuente:  ".5",
		errores: []string{"carácter inesperado"},
		esperado: []tokenEsperado{
			{token.NUMBER, "5", float64(5), 1},
		},
	},
	{
		nombre: "un número no puede tener dos puntos",
		// Se consume la primera parte decimal y nada más.
		fuente:  "1.2.3",
		errores: []string{"carácter inesperado"},
		esperado: []tokenEsperado{
			{token.NUMBER, "1.2", 1.2, 1},
			{token.NUMBER, "3", float64(3), 1},
		},
	},
}

func TestLiteralesNumericos(t *testing.T) {
	correrCasos(t, casosNumeros)
}
