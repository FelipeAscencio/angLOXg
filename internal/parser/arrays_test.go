package parser

import "testing"

var casosArrays = []casoDePrograma{
	{
		nombre:   "array vacío",
		fuente:   "[];",
		esperado: []string{"(expr (array))"},
	},
	{
		nombre:   "array con un elemento",
		fuente:   "[1];",
		esperado: []string{"(expr (array 1))"},
	},
	{
		nombre:   "array con varios elementos",
		fuente:   "[1, 2, 3];",
		esperado: []string{"(expr (array 1 2 3))"},
	},
	{
		nombre:   "los elementos pueden ser expresiones",
		fuente:   "[1 + 2, f(), a];",
		esperado: []string{"(expr (array (+ 1 2) (call f) a))"},
	},
	{
		nombre:   "los elementos pueden ser de tipos distintos",
		fuente:   `[1, "dos", true, nil];`,
		esperado: []string{`(expr (array 1 "dos" true nil))`},
	},
	{
		nombre:   "arrays anidados",
		fuente:   "[[1, 2], [3]];",
		esperado: []string{"(expr (array (array 1 2) (array 3)))"},
	},
	{
		nombre:   "array guardado en una variable",
		fuente:   "var a = [1, 2];",
		esperado: []string{"(var a (array 1 2))"},
	},
	{
		nombre:  "array sin cerrar",
		fuente:  "[1, 2;",
		errores: []string{"se esperaba ']'"},
	},
	{
		nombre:  "coma sin elemento atrás",
		fuente:  "[1, ];",
		errores: []string{"se esperaba una expresión"},
	},
}

var casosIndexado = []casoDePrograma{
	{
		nombre:   "lectura por posición",
		fuente:   "a[0];",
		esperado: []string{"(expr (index a 0))"},
	},
	{
		nombre:   "el índice puede ser una expresión",
		fuente:   "a[i + 1];",
		esperado: []string{"(expr (index a (+ i 1)))"},
	},
	{
		nombre:   "lectura sobre un literal",
		fuente:   "[1, 2][0];",
		esperado: []string{"(expr (index (array 1 2) 0))"},
	},
	{
		nombre: "lecturas encadenadas",
		fuente: "m[0][1];",
		esperado: []string{"(expr (index (index m 0) 1))"},
	},
	{
		nombre:   "lectura sobre el resultado de una llamada",
		fuente:   "f()[0];",
		esperado: []string{"(expr (index (call f) 0))"},
	},
	{
		nombre:   "llamada sobre un elemento",
		fuente:   "a[0]();",
		esperado: []string{"(expr (call (index a 0)))"},
	},
	{
		nombre:   "propiedad de un elemento",
		fuente:   "a[0].x;",
		esperado: []string{"(expr (get (index a 0) x))"},
	},
	{
		nombre:   "elemento de una propiedad",
		fuente:   "p.lista[0];",
		esperado: []string{"(expr (index (get p lista) 0))"},
	},
	{
		nombre: "el indexado liga más fuerte que la aritmética",
		fuente: "print a[0] + 1;",
		esperado: []string{"(print (+ (index a 0) 1))"},
	},
	{
		nombre:  "índice sin cerrar",
		fuente:  "a[0;",
		errores: []string{"se esperaba ']'"},
	},
	{
		nombre:  "índice vacío",
		fuente:  "a[];",
		errores: []string{"se esperaba una expresión"},
	},
	{
		nombre: "no se puede asignar a una posición",
		fuente: "a[0] = 1;",
		errores: []string{"destino de asignación inválido"},
	},
}

func TestLiteralesDeArray(t *testing.T) {
	correrCasosDePrograma(t, casosArrays)
}

func TestIndexadoDeArrays(t *testing.T) {
	correrCasosDePrograma(t, casosIndexado)
}
