package lox_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// casoDeClases es una entrada de una tabla de pruebas de punta a punta: un
// fuente y lo que el intérprete tiene que escribir al ejecutarlo.
type casoDeClases struct {
	nombre   string
	fuente   string
	esperado []string
	exito    bool
}

// ejecutarFuente corre un programa completo por el tree-walk y devuelve las
// líneas que escribió.
func ejecutarFuente(fuente string) ([]string, bool) {
	var salida bytes.Buffer
	exito := lox.Ejecutar(fuente, &salida)

	texto := salida.String()
	if texto == "" {
		return []string{}, exito
	}

	return strings.Split(strings.TrimSuffix(texto, "\n"), "\n"), exito
}

// correrCasosDeClases ejecuta una tabla de pruebas. Un error de ejecución
// ocupa dos renglones: el mensaje y la línea del fuente.
func correrCasosDeClases(t *testing.T, casos []casoDeClases) {
	t.Helper()

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenidas, exito := ejecutarFuente(caso.fuente)

			if exito != caso.exito {
				t.Errorf("la ejecución devolvió %v y se esperaba %v\nsalida: %v",
					exito, caso.exito, obtenidas)
			}

			if len(obtenidas) != len(caso.esperado) {
				t.Fatalf("se obtuvieron %d líneas y se esperaban %d\nobtenidas: %v",
					len(obtenidas), len(caso.esperado), obtenidas)
			}

			for i, esperada := range caso.esperado {
				if !strings.Contains(obtenidas[i], esperada) {
					t.Errorf("línea %d: se obtuvo %q y se esperaba que contuviera %q",
						i, obtenidas[i], esperada)
				}
			}
		})
	}
}

var casosDeclaracionEInstancia = []casoDeClases{
	{
		nombre:   "una clase se imprime con su nombre",
		fuente:   "class Punto {} print Punto;",
		esperado: []string{"Punto"},
		exito:    true,
	},
	{
		nombre:   "llamar a la clase devuelve una instancia",
		fuente:   "class Punto {} print Punto();",
		esperado: []string{"instancia de Punto"},
		exito:    true,
	},
	{
		nombre:   "la instancia se puede guardar en una variable",
		fuente:   "class Punto {} var p = Punto(); print p;",
		esperado: []string{"instancia de Punto"},
		exito:    true,
	},
	{
		nombre:   "una clase declarada dentro de un bloque",
		fuente:   "{ class Punto {} print Punto(); }",
		esperado: []string{"instancia de Punto"},
		exito:    true,
	},
	{
		nombre: "dos clases distintas",
		fuente: "class A {} class B {} print A(); print B();",
		esperado: []string{
			"instancia de A",
			"instancia de B",
		},
		exito: true,
	},
	{
		nombre: "la clase no acepta argumentos todavía",
		fuente: "class Punto {} Punto(1);",
		// El constructor "init" viene en una tanda siguiente.
		esperado: []string{"Se esperaban 0 argumentos pero se obtuvieron 1.", "[línea 1]"},
		exito:    false,
	},
}

var casosPropiedades = []casoDeClases{
	{
		nombre:   "escribir y leer una propiedad",
		fuente:   "class P {} var p = P(); p.x = 3; print p.x;",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "varias propiedades",
		fuente:   "class P {} var p = P(); p.x = 3; p.y = 4; print p.x + p.y;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre: "una propiedad guarda cualquier tipo",
		fuente: `class P {} var p = P();
			p.texto = "hola"; p.flag = true; p.nada = nil;
			print p.texto; print p.flag; print p.nada;`,
		esperado: []string{"hola", "true", "nil"},
		exito:    true,
	},
	{
		nombre: "dos instancias no comparten propiedades",
		fuente: "class P {} var a = P(); var b = P(); a.x = 1; b.x = 2; print a.x; print b.x;",
		// Cada instancia tiene su propia tabla de propiedades.
		esperado: []string{"1", "2"},
		exito:    true,
	},
	{
		nombre:   "una propiedad se puede reescribir",
		fuente:   "class P {} var p = P(); p.x = 1; p.x = 2; print p.x;",
		esperado: []string{"2"},
		exito:    true,
	},
	{
		nombre:   "una propiedad se puede leer y escribir en la misma expresión",
		fuente:   "class P {} var p = P(); p.x = 1; p.x = p.x + 10; print p.x;",
		esperado: []string{"11"},
		exito:    true,
	},
	{
		nombre:   "la escritura devuelve el valor asignado",
		fuente:   "class P {} var p = P(); print p.x = 5;",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre: "una propiedad puede guardar otra instancia",
		fuente: `class P {} var a = P(); var b = P();
			a.hijo = b; a.hijo.valor = 42;
			print a.hijo.valor;`,
		esperado: []string{"42"},
		exito:    true,
	},
	{
		nombre: "una propiedad puede guardar una función",
		fuente: `class P {} fun saludar() { return "hola"; }
			var p = P(); p.f = saludar; print p.f();`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "las propiedades se pueden usar en un bucle",
		fuente:   "class P {} var p = P(); p.n = 0; while (p.n < 3) { p.n = p.n + 1; } print p.n;",
		esperado: []string{"3"},
		exito:    true,
	},
}

var casosErroresDePropiedades = []casoDeClases{
	{
		nombre:   "leer una propiedad que no existe",
		fuente:   "class P {} var p = P(); print p.z;",
		esperado: []string{"Propiedad no definida 'z'.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "leer una propiedad de un número",
		fuente:   "print (1).x;",
		esperado: []string{"Solo las instancias tienen propiedades.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "leer una propiedad de una cadena",
		fuente:   `print "hola".largo;`,
		esperado: []string{"Solo las instancias tienen propiedades.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "escribir una propiedad en un número",
		fuente:   "var a = 1; a.x = 1;",
		esperado: []string{"Solo las instancias tienen propiedades.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "leer una propiedad del nulo",
		fuente:   "var a = nil; print a.x;",
		esperado: []string{"Solo las instancias tienen propiedades.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "leer una propiedad de una clase",
		fuente:   "class P {} print P.x;",
		esperado: []string{"Solo las instancias tienen propiedades.", "[línea 1]"},
		exito:    false,
	},
}

func TestDeclaracionEInstanciaDeClases(t *testing.T) {
	correrCasosDeClases(t, casosDeclaracionEInstancia)
}

func TestPropiedadesDeInstancias(t *testing.T) {
	correrCasosDeClases(t, casosPropiedades)
}

func TestErroresDePropiedades(t *testing.T) {
	correrCasosDeClases(t, casosErroresDePropiedades)
}

var casosMetodos = []casoDeClases{
	{
		nombre:   "llamada a un método sin argumentos",
		fuente:   `class A { saludar() { print "hola"; } } A().saludar();`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "método con parámetros y retorno",
		fuente:   "class A { sumar(x, y) { return x + y; } } print A().sumar(2, 3);",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre:   "un método se puede guardar en una variable",
		fuente:   `class A { saludar() { return "hola"; } } var m = A().saludar; print m();`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "un método se imprime como función",
		fuente:   "class A { m() {} } print A().m;",
		esperado: []string{"<fn m>"},
		exito:    true,
	},
	{
		nombre:   "varios métodos en la misma clase",
		fuente:   "class A { uno() { return 1; } dos() { return 2; } } var a = A(); print a.uno(); print a.dos();",
		esperado: []string{"1", "2"},
		exito:    true,
	},
	{
		nombre: "una propiedad tapa a un método del mismo nombre",
		fuente: `class A { m() { return "metodo"; } } var a = A(); a.m = "propiedad"; print a.m;`,
		// Las propiedades se buscan antes que los métodos.
		esperado: []string{"propiedad"},
		exito:    true,
	},
	{
		nombre:   "llamar a un método que no existe",
		fuente:   "class A {} A().noExiste();",
		esperado: []string{"Propiedad no definida 'noExiste'.", "[línea 1]"},
		exito:    false,
	},
}

var casosThis = []casoDeClases{
	{
		nombre:   "this devuelve la propia instancia",
		fuente:   "class A { yo() { return this; } } var a = A(); print a.yo() == a;",
		esperado: []string{"true"},
		exito:    true,
	},
	{
		nombre:   "un método escribe una propiedad con this",
		fuente:   "class A { poner() { this.x = 7; } } var a = A(); a.poner(); print a.x;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "un método lee una propiedad con this",
		fuente:   "class A { ver() { return this.x; } } var a = A(); a.x = 3; print a.ver();",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "un método llama a otro método con this",
		fuente:   "class A { doble() { return this.valor() * 2; } valor() { return 21; } } print A().doble();",
		esperado: []string{"42"},
		exito:    true,
	},
	{
		nombre: "this queda ligado a la instancia correcta",
		fuente: `class A { nombre() { return this.n; } }
			var a = A(); var b = A();
			a.n = "primera"; b.n = "segunda";
			print a.nombre(); print b.nombre();`,
		esperado: []string{"primera", "segunda"},
		exito:    true,
	},
	{
		nombre: "un método extraído se lleva su instancia",
		fuente: `class A { nombre() { return this.n; } }
			var a = A(); a.n = "primera";
			var m = a.nombre;
			var b = A(); b.n = "segunda";
			print m();`,
		// El método quedó ligado a "a" cuando se lo sacó, así que sigue viendo
		// esa instancia aunque después exista otra.
		esperado: []string{"primera"},
		exito:    true,
	},
	{
		nombre:   "un método acumula estado en la instancia",
		fuente:   "class C { subir() { this.n = this.n + 1; return this.n; } } var c = C(); c.n = 0; print c.subir(); print c.subir();",
		esperado: []string{"1", "2"},
		exito:    true,
	},
	{
		nombre:   "this fuera de una clase es un error",
		fuente:   "print this;",
		esperado: []string{"No se puede usar 'this' fuera de una clase"},
		exito:    false,
	},
	{
		nombre:   "this dentro de una función suelta es un error",
		fuente:   "fun suelta() { return this; }",
		esperado: []string{"No se puede usar 'this' fuera de una clase"},
		exito:    false,
	},
}

func TestMetodos(t *testing.T) {
	correrCasosDeClases(t, casosMetodos)
}

func TestReferenciasAThis(t *testing.T) {
	correrCasosDeClases(t, casosThis)
}
