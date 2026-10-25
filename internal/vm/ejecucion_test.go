package vm

import "testing"

var casosVariables = []casoDeEjecucion{
	{
		nombre:   "global con valor inicial",
		fuente:   "var x = 1; print x;",
		esperado: []string{"1"},
		exito:    true,
	},
	{
		nombre:   "asignación a una global",
		fuente:   "var x = 1; x = 2; print x;",
		esperado: []string{"2"},
		exito:    true,
	},
	{
		nombre:   "la asignación devuelve el valor asignado",
		fuente:   "var x = 1; print x = 5;",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre:   "local dentro de un bloque",
		fuente:   "{ var y = 2; print y; }",
		esperado: []string{"2"},
		exito:    true,
	},
	{
		nombre:   "una local tapa a la global del mismo nombre",
		fuente:   "var a = 1; { var a = 2; print a; } print a;",
		esperado: []string{"2", "1"},
		exito:    true,
	},
	{
		nombre:   "bloques anidados",
		fuente:   "{ var a = 1; { var b = 2; print a + b; } }",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "leer una variable que no existe",
		fuente:   "print x;",
		esperado: []string{"Error en ejecución: Variable no definida 'x'."},
		exito:    false,
	},
	{
		nombre:   "asignar a una variable que no existe",
		fuente:   "x = 1;",
		esperado: []string{"Error en ejecución: Variable no definida 'x'."},
		exito:    false,
	},
}

var casosControlDeFlujo = []casoDeEjecucion{
	{
		nombre:   "condicional por la rama verdadera",
		fuente:   `if (true) print "si"; else print "no";`,
		esperado: []string{"si"},
		exito:    true,
	},
	{
		nombre:   "condicional por la rama falsa",
		fuente:   `if (false) print "si"; else print "no";`,
		esperado: []string{"no"},
		exito:    true,
	},
	{
		nombre:   "condicional sin rama falsa que no entra",
		fuente:   `if (false) print "si"; print "fin";`,
		esperado: []string{"fin"},
		exito:    true,
	},
	{
		nombre:   "el else se pega al if más cercano",
		fuente:   `if (true) if (false) print 1; else print 2;`,
		esperado: []string{"2"},
		exito:    true,
	},
	{
		nombre:   "bucle mientras",
		fuente:   "var i = 0; while (i < 3) { print i; i = i + 1; }",
		esperado: []string{"0", "1", "2"},
		exito:    true,
	},
	{
		nombre:   "bucle mientras que no entra nunca",
		fuente:   "while (false) print 1;",
		esperado: []string{},
		exito:    true,
	},
	{
		nombre:   "bucle for",
		fuente:   "for (var i = 0; i < 3; i = i + 1) print i;",
		esperado: []string{"0", "1", "2"},
		exito:    true,
	},
	{
		nombre:   "bucle for con la variable de afuera",
		fuente:   "var i = 0; for (; i < 2; i = i + 1) print i;",
		esperado: []string{"0", "1"},
		exito:    true,
	},
	{
		nombre:   "la conjunción no evalúa el lado derecho si el izquierdo es falso",
		fuente:   "print false and 1;",
		esperado: []string{"false"},
		exito:    true,
	},
	{
		nombre:   "la conjunción sigue si el izquierdo es verdadero",
		fuente:   "print true and 1;",
		esperado: []string{"1"},
		exito:    true,
	},
	{
		nombre:   "la disyunción no evalúa el lado derecho si el izquierdo es verdadero",
		fuente:   "print 1 or false;",
		esperado: []string{"1"},
		exito:    true,
	},
	{
		nombre:   "la disyunción sigue si el izquierdo es falso",
		fuente:   "print false or 2;",
		esperado: []string{"2"},
		exito:    true,
	},
	{
		nombre:   "el cortocircuito evita el error del lado derecho",
		fuente:   "print false and sinDefinir;",
		esperado: []string{"false"},
		exito:    true,
	},
}

var casosFunciones = []casoDeEjecucion{
	{
		nombre:   "función sin parámetros",
		fuente:   `fun saludar() { print "hola"; } saludar();`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "función con parámetros y retorno",
		fuente:   "fun sumar(a, b) { return a + b; } print sumar(2, 3);",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre:   "una función sin return devuelve nil",
		fuente:   "fun nada() {} print nada() == nil;",
		esperado: []string{"true"},
		exito:    true,
	},
	{
		nombre:   "un return corta la función",
		fuente:   `fun f() { return 1; print "nunca"; } print f();`,
		esperado: []string{"1"},
		exito:    true,
	},
	{
		nombre:   "una función se imprime con su nombre",
		fuente:   "fun f() {} print f;",
		esperado: []string{"<fn f>"},
		exito:    true,
	},
	{
		nombre:   "recursión",
		fuente:   "fun fact(n) { if (n <= 1) return 1; return n * fact(n - 1); } print fact(5);",
		esperado: []string{"120"},
		exito:    true,
	},
	{
		nombre:   "una función llama a otra",
		fuente:   "fun doble(n) { return n * 2; } fun cuadruple(n) { return doble(doble(n)); } print cuadruple(3);",
		esperado: []string{"12"},
		exito:    true,
	},
	{
		nombre:   "llamar a algo que no es una función",
		fuente:   "var x = 1; x();",
		esperado: []string{"Error en ejecución: Solo se pueden llamar funciones y clases."},
		exito:    false,
	},
	{
		nombre:   "faltan argumentos",
		fuente:   "fun f(a) {} f();",
		esperado: []string{"Error en ejecución: Se esperaban 1 argumentos pero se obtuvieron 0."},
		exito:    false,
	},
	{
		nombre:   "sobran argumentos",
		fuente:   "fun f(a) {} f(1, 2);",
		esperado: []string{"Error en ejecución: Se esperaban 1 argumentos pero se obtuvieron 2."},
		exito:    false,
	},
}

// Los closures son la parte más delicada de la máquina: una función que se
var casosClosures = []casoDeEjecucion{
	{
		nombre:   "una función anidada lee una variable de la de afuera",
		fuente:   "fun externa() { var a = 1; fun interna() { return a; } return interna(); } print externa();",
		esperado: []string{"1"},
		exito:    true,
	},
	{
		nombre: "la variable capturada sobrevive a la función que la creó",
		fuente: `
			fun crearContador() {
				var cuenta = 0;
				fun incrementar() { cuenta = cuenta + 1; return cuenta; }
				return incrementar;
			}
			var contar = crearContador();
			print contar();
			print contar();
			print contar();
		`,
		esperado: []string{"1", "2", "3"},
		exito:    true,
	},
	{
		nombre: "dos closures del mismo generador no comparten estado",
		fuente: `
			fun crearContador() {
				var cuenta = 0;
				fun incrementar() { cuenta = cuenta + 1; return cuenta; }
				return incrementar;
			}
			var a = crearContador();
			var b = crearContador();
			print a();
			print a();
			print b();
		`,
		esperado: []string{"1", "2", "1"},
		exito:    true,
	},
	{
		nombre: "el closure captura el valor del parámetro",
		fuente: `
			fun sumador(n) {
				fun sumar(x) { return x + n; }
				return sumar;
			}
			var sumarDiez = sumador(10);
			print sumarDiez(5);
		`,
		esperado: []string{"15"},
		exito:    true,
	},
}

func TestVariables(t *testing.T) {
	correrCasos(t, casosVariables)
}

func TestControlDeFlujo(t *testing.T) {
	correrCasos(t, casosControlDeFlujo)
}

func TestFunciones(t *testing.T) {
	correrCasos(t, casosFunciones)
}

func TestClosures(t *testing.T) {
	correrCasos(t, casosClosures)
}
