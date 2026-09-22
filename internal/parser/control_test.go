package parser

import "testing"

// Casos de prueba para operaciones de asignación.
var casosAsignacion = []casoDePrograma{
	{
		nombre:   "asignación simple",
		fuente:   "a = 1;",
		esperado: []string{"(expr (= a 1))"},
	},
	{
		nombre:   "asignación de una expresión",
		fuente:   "a = 1 + 2;",
		esperado: []string{"(expr (= a (+ 1 2)))"},
	},
	{
		nombre:   "la asignación asocia a la derecha",
		fuente:   "a = b = 1;",
		esperado: []string{"(expr (= a (= b 1)))"},
	},
	{
		nombre:   "destino inválido",
		fuente:   "1 = 2;",
		errores:  []string{"destino de asignación inválido"},
	},
	{
		nombre:   "destino inválido con una expresión",
		fuente:   "a + b = 1;",
		errores:  []string{"destino de asignación inválido"},
	},
}

// Casos de prueba para operadores lógicos.
var casosLogicos = []casoDePrograma{
	{
		nombre:   "disyunción",
		fuente:   "a or b;",
		esperado: []string{"(expr (or a b))"},
	},
	{
		nombre:   "conjunción",
		fuente:   "a and b;",
		esperado: []string{"(expr (and a b))"},
	},
	{
		nombre:   "la conjunción agrupa antes que la disyunción",
		fuente:   "a or b and c;",
		esperado: []string{"(expr (or a (and b c)))"},
	},
	{
		nombre:   "la igualdad agrupa antes que la conjunción",
		fuente:   "a and b == c;",
		esperado: []string{"(expr (and a (== b c)))"},
	},
	{
		nombre:   "los lógicos asocian a la izquierda",
		fuente:   "a or b or c;",
		esperado: []string{"(expr (or (or a b) c))"},
	},
	{
		nombre:   "la asignación agrupa después de los lógicos",
		fuente:   "a = b or c;",
		esperado: []string{"(expr (= a (or b c)))"},
	},
}

// Casos de prueba para estructuras condicionales (if/else).
var casosCondicionales = []casoDePrograma{
	{
		nombre:   "condicional sin rama falsa",
		fuente:   "if (a) print 1;",
		esperado: []string{"(if a (print 1))"},
	},
	{
		nombre:   "condicional con rama falsa",
		fuente:   "if (a) print 1; else print 2;",
		esperado: []string{"(if a (print 1) (print 2))"},
	},
	{
		nombre:   "condicional con bloques",
		fuente:   "if (a) { print 1; } else { print 2; }",
		esperado: []string{"(if a (block (print 1)) (block (print 2)))"},
	},
	{
		nombre:   "el else se pega al if más cercano",
		fuente:   "if (a) if (b) print 1; else print 2;",
		esperado: []string{"(if a (if b (print 1) (print 2)))"},
	},
	{
		nombre:   "condicional sin paréntesis de apertura",
		fuente:   "if a) print 1;",
		errores:  []string{"se esperaba '(' después de 'if'"},
		esperado: []string{"(print 1)"},
	},
	{
		nombre:   "condicional sin paréntesis de cierre",
		fuente:   "if (a print 1;",
		errores:  []string{"se esperaba ')'"},
	},
}

// Casos de prueba para bucles (while y for).
var casosBucles = []casoDePrograma{
	{
		nombre:   "bucle mientras",
		fuente:   "while (a) print 1;",
		esperado: []string{"(while a (print 1))"},
	},
	{
		nombre:   "bucle mientras con bloque",
		fuente:   "while (a) { print 1; }",
		esperado: []string{"(while a (block (print 1)))"},
	},
	{
		nombre:   "bucle mientras sin paréntesis",
		fuente:   "while a) print 1;",
		errores:  []string{"se esperaba '(' después de 'while'"},
		esperado: []string{"(print 1)"},
	},
	{
		nombre:   "bucle for sin ninguna de sus tres partes",
		fuente:   "for (;;) print 1;",
		esperado: []string{"(while true (print 1))"},
	},
	{
		nombre:   "bucle for completo",
		fuente:   "for (var i = 0; i < 10; i = i + 1) print i;",
		esperado: []string{"(block (var i 0) (while (< i 10) (block (print i) (expr (= i (+ i 1))))))"},
	},
	{
		nombre:   "bucle for con inicializador de expresión",
		fuente:   "for (i = 0; i < 3;) print i;",
		esperado: []string{"(block (expr (= i 0)) (while (< i 3) (print i)))"},
	},
	{
		nombre:   "bucle for sin inicializador",
		fuente:   "for (; i < 3; i = i + 1) print i;",
		esperado: []string{"(while (< i 3) (block (print i) (expr (= i (+ i 1)))))"},
	},
	{
		nombre: "bucle for sin paréntesis",
		fuente: "for ;; print 1;",
		errores: []string{
			"se esperaba '(' después de 'for'",
			"se esperaba una expresión",
		},
		esperado: []string{"(print 1)"},
	},
}

func TestAsignacion(t *testing.T) {
	correrCasosDePrograma(t, casosAsignacion)
}

func TestOperadoresLogicos(t *testing.T) {
	correrCasosDePrograma(t, casosLogicos)
}

func TestCondicionales(t *testing.T) {
	correrCasosDePrograma(t, casosCondicionales)
}

func TestBucles(t *testing.T) {
	correrCasosDePrograma(t, casosBucles)
}
