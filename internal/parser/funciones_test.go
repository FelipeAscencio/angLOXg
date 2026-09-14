package parser

import "testing"

var casosLlamadas = []casoDePrograma{
	{
		nombre:   "llamada sin argumentos",
		fuente:   "f();",
		esperado: []string{"(expr (call f))"},
	},
	{
		nombre:   "llamada con un argumento",
		fuente:   "f(1);",
		esperado: []string{"(expr (call f 1))"},
	},
	{
		nombre:   "llamada con varios argumentos",
		fuente:   "f(1, 2, 3);",
		esperado: []string{"(expr (call f 1 2 3))"},
	},
	{
		nombre:   "los argumentos pueden ser expresiones",
		fuente:   "f(1 + 2, g(3));",
		esperado: []string{"(expr (call f (+ 1 2) (call g 3)))"},
	},
	{
		nombre: "llamadas encadenadas",
		fuente: "f()();",
		// Lo que se llama la segunda vez es el resultado de la primera.
		esperado: []string{"(expr (call (call f)))"},
	},
	{
		nombre:   "la llamada agrupa antes que el unario",
		fuente:   "-f(1);",
		esperado: []string{"(expr (- (call f 1)))"},
	},
	{
		nombre:   "la llamada agrupa antes que el producto",
		fuente:   "f(1) * 2;",
		esperado: []string{"(expr (* (call f 1) 2))"},
	},
	{
		nombre:  "llamada sin cerrar",
		fuente:  "f(1;",
		errores: []string{"se esperaba ')'"},
	},
}

var casosDeclaracionFuncion = []casoDePrograma{
	{
		nombre:   "función sin parámetros",
		fuente:   "fun saludar() { print 1; }",
		esperado: []string{"(fun saludar () (print 1))"},
	},
	{
		nombre:   "función con un parámetro",
		fuente:   "fun doble(n) { return n * 2; }",
		esperado: []string{"(fun doble (n) (return (* n 2)))"},
	},
	{
		nombre:   "función con varios parámetros",
		fuente:   "fun sumar(a, b) { return a + b; }",
		esperado: []string{"(fun sumar (a b) (return (+ a b)))"},
	},
	{
		nombre:   "función con cuerpo vacío",
		fuente:   "fun nada() {}",
		esperado: []string{"(fun nada ())"},
	},
	{
		nombre:   "función con varias sentencias",
		fuente:   "fun f() { var x = 1; print x; }",
		esperado: []string{"(fun f () (var x 1) (print x))"},
	},
	{
		nombre:   "retorno sin valor",
		fuente:   "fun f() { return; }",
		esperado: []string{"(fun f () (return))"},
	},
	{
		nombre: "función recursiva",
		fuente: "fun fact(n) { if (n <= 1) return 1; return n * fact(n - 1); }",
		esperado: []string{
			"(fun fact (n) (if (<= n 1) (return 1)) (return (* n (call fact (- n 1)))))",
		},
	},
	{
		nombre:  "función sin nombre",
		fuente:  "fun () {}",
		errores: []string{"se esperaba el nombre de la función"},
	},
	{
		nombre:  "función sin paréntesis",
		fuente:  "fun f {}",
		errores: []string{"se esperaba '(' después del nombre de la función"},
	},
	{
		nombre:  "parámetro que no es un nombre",
		fuente:  "fun f(1) {}",
		errores: []string{"se esperaba el nombre del parámetro"},
	},
	{
		nombre:  "función sin cuerpo",
		fuente:  "fun f();",
		errores: []string{"se esperaba '{' para abrir el cuerpo de la función"},
	},
}

func TestLlamadas(t *testing.T) {
	correrCasosDePrograma(t, casosLlamadas)
}

func TestDeclaracionDeFunciones(t *testing.T) {
	correrCasosDePrograma(t, casosDeclaracionFuncion)
}
