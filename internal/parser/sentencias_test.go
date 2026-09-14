package parser

import "testing"

var casosSentencias = []casoDePrograma{
	{
		nombre:   "impresión",
		fuente:   "print 1;",
		esperado: []string{"(print 1)"},
	},
	{
		nombre:   "impresión de una expresión compuesta",
		fuente:   "print 1 + 2 * 3;",
		esperado: []string{"(print (+ 1 (* 2 3)))"},
	},
	{
		nombre:   "sentencia de expresión",
		fuente:   "1 + 2;",
		esperado: []string{"(expr (+ 1 2))"},
	},
	{
		nombre:   "declaración con valor inicial",
		fuente:   "var x = 5;",
		esperado: []string{"(var x 5)"},
	},
	{
		nombre:   "declaración sin valor inicial",
		fuente:   "var x;",
		esperado: []string{"(var x)"},
	},
	{
		nombre:   "varias sentencias",
		fuente:   "print 1; print 2; var x = 3;",
		esperado: []string{"(print 1)", "(print 2)", "(var x 3)"},
	},
	{
		nombre:   "programa vacío",
		fuente:   "",
		esperado: []string{},
	},
	{
		nombre:   "solo un comentario",
		fuente:   "// nada que hacer",
		esperado: []string{},
	},
	{
		nombre:   "bloque vacío",
		fuente:   "{}",
		esperado: []string{"(block)"},
	},
	{
		nombre:   "bloque con sentencias",
		fuente:   "{ print 1; print 2; }",
		esperado: []string{"(block (print 1) (print 2))"},
	},
	{
		nombre:   "bloques anidados",
		fuente:   "{ { var x = 1; } }",
		esperado: []string{"(block (block (var x 1)))"},
	},
	{
		nombre:  "falta el punto y coma",
		fuente:  "print 1",
		errores: []string{"se esperaba ';'"},
	},
	{
		nombre:  "falta el punto y coma de una declaración",
		fuente:  "var x = 1",
		errores: []string{"se esperaba ';'"},
	},
	{
		nombre:  "declaración sin nombre",
		fuente:  "var = 1;",
		errores: []string{"se esperaba el nombre de la variable"},
	},
	{
		nombre:  "bloque sin cerrar",
		fuente:  "{ print 1;",
		errores: []string{"se esperaba '}'"},
	},
}

func TestSentencias(t *testing.T) {
	correrCasosDePrograma(t, casosSentencias)
}
