package parser

import "testing"

// Casos de prueba para verificar la sincronización y recuperación de errores sintácticos.
var casosSincronizacion = []casoDePrograma{
	{
		nombre: "una sentencia rota no se lleva a las sanas",
		fuente: "var = 1;\nprint 2;",
		errores: []string{
			"se esperaba el nombre de la variable",
		},
		esperado: []string{"(print 2)"},
	},
	{
		nombre: "varios errores en una sola pasada",
		fuente: "var = 1;\nprint 2;\nvar = 3;\nprint 4;",
		errores: []string{
			"se esperaba el nombre de la variable",
			"se esperaba el nombre de la variable",
		},
		esperado: []string{"(print 2)", "(print 4)"},
	},
	{
		nombre: "cada error se reporta en su línea",
		fuente: "print 1;\n\nvar = 2;",
		errores: []string{
			"[línea 3] Error de sintaxis: se esperaba el nombre de la variable",
		},
		esperado: []string{"(print 1)"},
	},
	{
		nombre: "un error dentro de un bloque pierde el bloque, no su contenido",
		fuente: "{ var = 1; print 2; }\nprint 3;",
		errores: []string{
			"se esperaba el nombre de la variable",
			"se esperaba una expresión",
		},
		esperado: []string{"(print 2)", "(print 3)"},
	},
	{
		nombre:   "se sincroniza en la palabra que abre una sentencia",
		fuente:   "var x = ;\nwhile (a) print 1;",
		errores:  []string{"se esperaba una expresión"},
		esperado: []string{"(while a (print 1))"},
	},
	{
		nombre:   "un error al final del archivo no cuelga el parseo",
		fuente:   "print 1;\nprint",
		errores:  []string{"se esperaba una expresión"},
		esperado: []string{"(print 1)"},
	},
	{
		nombre:   "solo errores",
		fuente:   "var = 1;\nvar = 2;",
		errores:  []string{"se esperaba el nombre de la variable", "se esperaba el nombre de la variable"},
		esperado: []string{},
	},
}

func TestSincronizacionDeErrores(t *testing.T) {
	correrCasosDePrograma(t, casosSincronizacion)
}
