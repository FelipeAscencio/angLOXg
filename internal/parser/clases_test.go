package parser

import "testing"

// Declaración de clases. En esta etapa el cuerpo tiene que estar vacío: los
// métodos vienen después.
var casosDeclaracionClase = []casoDePrograma{
	{
		nombre:   "clase con cuerpo vacío",
		fuente:   "class Punto {}",
		esperado: []string{"(class Punto)"},
	},
	{
		nombre:   "varias clases",
		fuente:   "class A {} class B {}",
		esperado: []string{"(class A)", "(class B)"},
	},
	{
		nombre:   "la clase se puede declarar dentro de un bloque",
		fuente:   "{ class A {} }",
		esperado: []string{"(block (class A))"},
	},
	{
		nombre:  "clase sin nombre",
		fuente:  "class {}",
		errores: []string{"se esperaba el nombre de la clase"},
	},
	{
		nombre:  "clase sin llave de apertura",
		fuente:  "class A }",
		errores: []string{"se esperaba '{'"},
	},
	{
		nombre:  "clase sin llave de cierre",
		fuente:  "class A {",
		errores: []string{"se esperaba '}'"},
	},
	{
		nombre:  "los métodos todavía no se admiten",
		fuente:  "class A { metodo() {} }",
		errores: []string{"el cuerpo de la clase tiene que estar vacío"},
	},
}

// Acceso a propiedades con punto.
var casosPropiedades = []casoDePrograma{
	{
		nombre:   "lectura de una propiedad",
		fuente:   "p.x;",
		esperado: []string{"(expr (get p x))"},
	},
	{
		nombre:   "lecturas encadenadas",
		fuente:   "a.b.c;",
		esperado: []string{"(expr (get (get a b) c))"},
	},
	{
		nombre:   "lectura sobre el resultado de una llamada",
		fuente:   "f().x;",
		esperado: []string{"(expr (get (call f) x))"},
	},
	{
		nombre:   "llamada sobre una propiedad",
		fuente:   "p.f();",
		esperado: []string{"(expr (call (get p f)))"},
	},
	{
		nombre:   "escritura de una propiedad",
		fuente:   "p.x = 3;",
		esperado: []string{"(expr (set p x 3))"},
	},
	{
		nombre: "escritura sobre una propiedad anidada",
		fuente: "a.b.c = 1;",
		// Sólo el último punto es una escritura; los de antes son lecturas.
		esperado: []string{"(expr (set (get a b) c 1))"},
	},
	{
		nombre:   "la escritura asocia a la derecha",
		fuente:   "p.x = p.y = 1;",
		esperado: []string{"(expr (set p x (set p y 1)))"},
	},
	{
		nombre:   "el valor escrito puede ser una expresión",
		fuente:   "p.x = 1 + 2;",
		esperado: []string{"(expr (set p x (+ 1 2)))"},
	},
	{
		nombre:   "instanciación y uso",
		fuente:   "var p = Punto(); p.x = 3; print p.x;",
		esperado: []string{"(var p (call Punto))", "(expr (set p x 3))", "(print (get p x))"},
	},
	{
		nombre:  "el punto tiene que ir seguido de un nombre",
		fuente:  "p.1;",
		errores: []string{"se esperaba el nombre de la propiedad"},
	},
	{
		nombre:  "el punto tiene que ir seguido de algo",
		fuente:  "p.;",
		errores: []string{"se esperaba el nombre de la propiedad"},
	},
}

func TestDeclaracionDeClases(t *testing.T) {
	correrCasosDePrograma(t, casosDeclaracionClase)
}

func TestAccesoAPropiedades(t *testing.T) {
	correrCasosDePrograma(t, casosPropiedades)
}
