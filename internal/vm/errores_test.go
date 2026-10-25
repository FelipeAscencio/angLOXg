package vm

import "testing"

// Operandos del tipo equivocado. Lo importante de estos casos es que la máquina
var casosOperandosInvalidos = []casoDeEjecucion{
	{
		nombre:   "resta con una cadena",
		fuente:   `print 1 - "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "producto con una cadena",
		fuente:   `print 2 * "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "división con una cadena",
		fuente:   `print 2 / "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "módulo con una cadena",
		fuente:   `print 10 % "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "comparación con una cadena",
		fuente:   `print 1 < "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "comparación de mayor con una cadena",
		fuente:   `print 1 > "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "negar una cadena",
		fuente:   `print -"a";`,
		esperado: []string{"Error en ejecución: El operando debe ser un número."},
		exito:    false,
	},
	{
		nombre:   "negar el nulo",
		fuente:   "print -nil;",
		esperado: []string{"Error en ejecución: El operando debe ser un número."},
		exito:    false,
	},
	{
		nombre:   "operandos nulos",
		fuente:   "print nil - nil;",
		esperado: []string{"Error en ejecución: Los operandos deben ser números."},
		exito:    false,
	},
	{
		nombre:   "división por cero",
		fuente:   "print 1 / 0;",
		esperado: []string{"Error en ejecución: División por cero."},
		exito:    false,
	},
}

// El nulo de Lox se escribe "nil", no con el formato por omisión de Go.
var casosNulo = []casoDeEjecucion{
	{
		nombre:   "imprimir el nulo",
		fuente:   "print nil;",
		esperado: []string{"nil"},
		exito:    true,
	},
	{
		nombre:   "una variable sin valor inicial arranca en nulo",
		fuente:   "var x; print x;",
		esperado: []string{"nil"},
		exito:    true,
	},
	{
		nombre:   "una función sin return devuelve nulo",
		fuente:   "fun f() {} print f();",
		esperado: []string{"nil"},
		exito:    true,
	},
	{
		nombre:   "un return pelado devuelve nulo",
		fuente:   "fun f() { return; } print f();",
		esperado: []string{"nil"},
		exito:    true,
	},
}

// El ámbito de una función anidada no se escapa a la función que la contiene.
var casosAmbitoDeFunciones = []casoDeEjecucion{
	{
		nombre: "una función anidada no queda visible desde afuera",
		fuente: `
			fun externa() {
				fun interna() { return 42; }
				return interna();
			}
			print externa();
			print interna();
		`,
		esperado: []string{"42", "Error en ejecución: Variable no definida 'interna'."},
		exito:    false,
	},
	{
		nombre: "una función anidada puede ser recursiva",
		fuente: `
			fun externa() {
				fun fact(n) { if (n <= 1) return 1; return n * fact(n - 1); }
				return fact(5);
			}
			print externa();
		`,
		esperado: []string{"120"},
		exito:    true,
	},
	{
		nombre: "dos funciones pueden tener una anidada del mismo nombre",
		fuente: `
			fun a() { fun ayuda() { return 1; } return ayuda(); }
			fun b() { fun ayuda() { return 2; } return ayuda(); }
			print a();
			print b();
		`,
		esperado: []string{"1", "2"},
		exito:    true,
	},
}

func TestOperandosInvalidos(t *testing.T) {
	correrCasos(t, casosOperandosInvalidos)
}

func TestRepresentacionDelNulo(t *testing.T) {
	correrCasos(t, casosNulo)
}

func TestAmbitoDeFuncionesAnidadas(t *testing.T) {
	correrCasos(t, casosAmbitoDeFunciones)
}
