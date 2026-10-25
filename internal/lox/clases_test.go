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
