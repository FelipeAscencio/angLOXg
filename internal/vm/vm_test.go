package vm

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/compilador"
	"github.com/FelipeAscencio/angLOXg/internal/escaner"
)

// casoDeEjecucion es una entrada de una tabla de pruebas: un fuente y lo que la
type casoDeEjecucion struct {
	nombre string
	fuente string

	esperado []string

	exito bool
}

// ejecutar compila un fuente y lo corre en la máquina virtual, devolviendo las
func ejecutar(t *testing.T, fuente string) (lineas []string, exito bool) {
	t.Helper()

	tokens, erroresLexicos := escaner.Nuevo(fuente).EscanearTokens()
	if len(erroresLexicos) > 0 {
		t.Fatalf("el fuente de la prueba no escanea: %v", erroresLexicos)
	}

	var salida bytes.Buffer

	defer func() {
		if recuperado := recover(); recuperado != nil {
			t.Fatalf("la máquina virtual entró en pánico con %q: %v", fuente, recuperado)
		}
	}()

	exito = Nueva(&salida).Ejecutar(compilador.Compilar(tokens))

	texto := salida.String()
	if texto == "" {
		return []string{}, exito
	}

	return strings.Split(strings.TrimSuffix(texto, "\n"), "\n"), exito
}

// correrCasos ejecuta una tabla de pruebas.
func correrCasos(t *testing.T, casos []casoDeEjecucion) {
	t.Helper()

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenidas, exito := ejecutar(t, caso.fuente)

			if exito != caso.exito {
				t.Errorf("la ejecución devolvió %v y se esperaba %v\nsalida: %v",
					exito, caso.exito, obtenidas)
			}

			if len(obtenidas) != len(caso.esperado) {
				t.Fatalf("se obtuvieron %d líneas y se esperaban %d\nobtenidas: %v",
					len(obtenidas), len(caso.esperado), obtenidas)
			}

			for i, esperada := range caso.esperado {
				if obtenidas[i] != esperada {
					t.Errorf("línea %d: se obtuvo %q y se esperaba %q", i, obtenidas[i], esperada)
				}
			}
		})
	}
}

var casosAritmetica = []casoDeEjecucion{
	{nombre: "suma", fuente: "print 1 + 2;", esperado: []string{"3"}, exito: true},
	{nombre: "resta", fuente: "print 5 - 3;", esperado: []string{"2"}, exito: true},
	{nombre: "producto", fuente: "print 4 * 3;", esperado: []string{"12"}, exito: true},
	{nombre: "división", fuente: "print 10 / 4;", esperado: []string{"2.5"}, exito: true},
	{nombre: "módulo", fuente: "print 10 % 3;", esperado: []string{"1"}, exito: true},
	{nombre: "negación", fuente: "print -5;", esperado: []string{"-5"}, exito: true},
	{
		nombre:   "precedencia",
		fuente:   "print 1 + 2 * 3;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "los paréntesis le ganan a la precedencia",
		fuente:   "print (1 + 2) * 3;",
		esperado: []string{"9"},
		exito:    true,
	},
	{
		nombre:   "la resta asocia a la izquierda",
		fuente:   "print 1 - 2 - 3;",
		esperado: []string{"-4"},
		exito:    true,
	},
	{
		nombre:   "un entero se imprime sin parte decimal",
		fuente:   "print 3.0;",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "módulo por cero",
		fuente:   "print 10 % 0;",
		esperado: []string{"Error en ejecución: Módulo por cero."},
		exito:    false,
	},
}

var casosCadenas = []casoDeEjecucion{
	{nombre: "cadena literal", fuente: `print "hola";`, esperado: []string{"hola"}, exito: true},
	{
		nombre:   "concatenación",
		fuente:   `print "ho" + "la";`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "cadena vacía",
		fuente:   `print "";`,
		esperado: []string{""},
		exito:    true,
	},
	{
		nombre:   "caracteres multibyte",
		fuente:   `print "ñandú";`,
		esperado: []string{"ñandú"},
		exito:    true,
	},
	{
		nombre:   "sumar un número a una cadena no se puede",
		fuente:   `print 1 + "a";`,
		esperado: []string{"Error en ejecución: Los operandos deben ser dos números o dos cadenas."},
		exito:    false,
	},
}

var casosComparaciones = []casoDeEjecucion{
	{nombre: "menor", fuente: "print 1 < 2;", esperado: []string{"true"}, exito: true},
	{nombre: "mayor", fuente: "print 1 > 2;", esperado: []string{"false"}, exito: true},
	{nombre: "menor o igual", fuente: "print 2 <= 2;", esperado: []string{"true"}, exito: true},
	{nombre: "mayor o igual", fuente: "print 2 >= 3;", esperado: []string{"false"}, exito: true},
	{nombre: "igualdad", fuente: "print 1 == 1;", esperado: []string{"true"}, exito: true},
	{nombre: "desigualdad", fuente: "print 1 != 2;", esperado: []string{"true"}, exito: true},
	{
		nombre:   "igualdad de cadenas",
		fuente:   `print "a" == "a";`,
		esperado: []string{"true"},
		exito:    true,
	},
	{
		nombre:   "tipos distintos nunca son iguales",
		fuente:   `print 1 == "1";`,
		esperado: []string{"false"},
		exito:    true,
	},
	{nombre: "nil es igual a nil", fuente: "print nil == nil;", esperado: []string{"true"}, exito: true},
}

// En Lox sólo nil y false son falsos. Todo lo demás, incluidos el 0 y la cadena
var casosVeracidad = []casoDeEjecucion{
	{nombre: "nil es falso", fuente: "print !nil;", esperado: []string{"true"}, exito: true},
	{nombre: "false es falso", fuente: "print !false;", esperado: []string{"true"}, exito: true},
	{nombre: "true es verdadero", fuente: "print !true;", esperado: []string{"false"}, exito: true},
	{
		nombre:   "el cero es verdadero",
		fuente:   "print !0;",
		esperado: []string{"false"},
		exito:    true,
	},
	{
		nombre:   "la cadena vacía es verdadera",
		fuente:   `print !"";`,
		esperado: []string{"false"},
		exito:    true,
	},
	{
		nombre:   "el cero entra en un condicional",
		fuente:   `if (0) print "entra"; else print "no entra";`,
		esperado: []string{"entra"},
		exito:    true,
	},
	{nombre: "doble negación", fuente: "print !!nil;", esperado: []string{"false"}, exito: true},
}

func TestAritmetica(t *testing.T) {
	correrCasos(t, casosAritmetica)
}

func TestCadenas(t *testing.T) {
	correrCasos(t, casosCadenas)
}

func TestComparaciones(t *testing.T) {
	correrCasos(t, casosComparaciones)
}

func TestVeracidad(t *testing.T) {
	correrCasos(t, casosVeracidad)
}
