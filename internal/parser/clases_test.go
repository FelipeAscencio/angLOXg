package parser

import "testing"

// Declaración de clases.
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
		nombre:   "clase con un método",
		fuente:   "class A { saludar() { print 1; } }",
		esperado: []string{"(class A (fun saludar () (print 1)))"},
	},
	{
		nombre:   "clase con varios métodos",
		fuente:   "class A { uno() {} dos() {} }",
		esperado: []string{"(class A (fun uno ()) (fun dos ()))"},
	},
	{
		nombre:   "método con parámetros",
		fuente:   "class A { sumar(x, y) { return x + y; } }",
		esperado: []string{"(class A (fun sumar (x y) (return (+ x y))))"},
	},
	{
		nombre:  "un método no lleva la palabra fun",
		fuente:  "class A { fun metodo() {} }",
		errores: []string{"se esperaba el nombre del método"},
	},
	{
		nombre:  "método sin paréntesis",
		fuente:  "class A { metodo {} }",
		errores: []string{"se esperaba '('"},
	},
}

var casosThis = []casoDePrograma{
	{
		nombre:   "this suelto en un método",
		fuente:   "class A { yo() { return this; } }",
		esperado: []string{"(class A (fun yo () (return this)))"},
	},
	{
		nombre:   "lectura de una propiedad de this",
		fuente:   "class A { ver() { print this.x; } }",
		esperado: []string{"(class A (fun ver () (print (get this x))))"},
	},
	{
		nombre:   "escritura de una propiedad de this",
		fuente:   "class A { poner() { this.x = 1; } }",
		esperado: []string{"(class A (fun poner () (expr (set this x 1))))"},
	},
	{
		nombre:   "llamada a un método de this",
		fuente:   "class A { a() { this.b(); } b() {} }",
		esperado: []string{"(class A (fun a () (expr (call (get this b)))) (fun b ()))"},
	},
	{
		nombre: "this no se puede asignar",
		fuente: "class A { m() { this = 1; } }",
		errores: []string{
			"destino de asignación inválido",
			"se esperaba una expresión",
		},
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

func TestReferenciasAThis(t *testing.T) {
	correrCasosDePrograma(t, casosThis)
}

var casosHerencia = []casoDePrograma{
	{
		nombre:   "clase con superclase",
		fuente:   "class Perro < Animal {}",
		esperado: []string{"(class Perro <Animal)"},
	},
	{
		nombre:   "clase con superclase y métodos",
		fuente:   "class Perro < Animal { ladrar() { print 1; } }",
		esperado: []string{"(class Perro <Animal (fun ladrar () (print 1)))"},
	},
	{
		nombre:   "llamada a un método de la superclase",
		fuente:   "class Perro < Animal { hablar() { super.hablar(); } }",
		esperado: []string{"(class Perro <Animal (fun hablar () (expr (call (super hablar)))))"},
	},
	{
		nombre:   "un método de la superclase se puede guardar",
		fuente:   "class Perro < Animal { m() { var f = super.hablar; } }",
		esperado: []string{"(class Perro <Animal (fun m () (var f (super hablar))))"},
	},
	{
		nombre:  "la superclase tiene que ser un nombre",
		fuente:  "class Perro < 1 {}",
		errores: []string{"se esperaba el nombre de la superclase"},
	},
	{
		nombre: "super tiene que ir seguido de un punto",
		fuente: "class P < A { m() { super; } }",
		errores: []string{
			"se esperaba '.' después de 'super'",
			"se esperaba una expresión",
		},
	},
	{
		nombre: "super tiene que nombrar un método",
		fuente: "class P < A { m() { super.1; } }",
		errores: []string{
			"se esperaba el nombre del método de la superclase",
			"se esperaba una expresión",
		},
	},
}

func TestHerencia(t *testing.T) {
	correrCasosDePrograma(t, casosHerencia)
}
