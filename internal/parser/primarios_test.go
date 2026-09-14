package parser

import "testing"

// Expresiones primarias: lo que no se descompone en nada más chico.
var casosPrimarios = []casoDeParseo{
	{
		nombre:   "número entero",
		fuente:   "123",
		esperado: "123",
	},
	{
		nombre:   "número decimal",
		fuente:   "45.67",
		esperado: "45.67",
	},
	{
		nombre:   "cadena",
		fuente:   `"hola"`,
		esperado: `"hola"`,
	},
	{
		nombre:   "verdadero",
		fuente:   "true",
		esperado: "true",
	},
	{
		nombre:   "falso",
		fuente:   "false",
		esperado: "false",
	},
	{
		nombre:   "nulo",
		fuente:   "nil",
		esperado: "nil",
	},
	{
		nombre:   "variable",
		fuente:   "contador",
		esperado: "contador",
	},
	{
		nombre:   "agrupación",
		fuente:   "(1)",
		esperado: "(group 1)",
	},
	{
		nombre:   "agrupación anidada",
		fuente:   "((1))",
		esperado: "(group (group 1))",
	},
	{
		nombre:  "paréntesis sin cerrar",
		fuente:  "(1",
		errores: []string{"se esperaba ')'"},
	},
	{
		nombre:  "no hay expresión",
		fuente:  ";",
		errores: []string{"se esperaba una expresión"},
	},
}

func TestExpresionesPrimarias(t *testing.T) {
	correrCasosDeExpresion(t, casosPrimarios)
}
