package lox_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Verifica la cobertura de errores de ejecución, casos límite y comportamientos especiales de Lox.
func TestCoberturaDeErroresYExtremos(t *testing.T) {
	casos := []struct {
		nombre string
		fuente string
		falla  bool
		error  string
	}{
		// 1. Operadores Unarios.
		{"Negación aritmética inválida", `print -"hola";`, true, "operando"},
		
		// 2. Operadores Binarios.
		{"Resta inválida", `print 5 - "hola";`, true, "operandos"},
		{"Multiplicación inválida", `print 5 * false;`, true, "operandos"},
		{"División inválida", `print true / 2;`, true, "operandos"},
		{"Suma mixta", `print 2 + "hola";`, true, "dos números o dos cadenas"},
		{"Menor qué inválido", `print 5 < "5";`, true, "operandos"},
		{"Mayor o igual inválido", `print 5 >= false;`, true, "operandos"},
		
		// 3. Casos por Cero.
		{"División por cero", `print 10 / 0;`, true, "cero"},
		{"Módulo por cero", `print 10 % 0;`, true, "cero"},
		
		// 4. Llamadas a Funciones.
		{"Llamar a no-invocable", `var a = 5; a();`, true, "Solo se pueden llamar"},
		{"Faltan argumentos", `fun f(a, b) {} f(1);`, true, "argumentos"},
		{"Sobran argumentos", `fun f(a) {} f(1, 2, 3);`, true, "argumentos"},
		
		// 5. Variables e Igualdad.
		{"Igualdad tipos distintos", `print 1 == "1";`, false, ""},
		{"Igualdad de nulos", `print nil == nil;`, false, ""},
		{"Cortocircuito AND", `print false and error_no_existe;`, false, ""},
		{"Cortocircuito OR", `print true or error_no_existe;`, false, ""},
		{"Impresión de nil", `var nulo; print nulo;`, false, ""},
		{"Evaluación de Truthiness", `if ("cadena") print "es verdad"; if (nil) print "no";`, false, ""},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			var salida bytes.Buffer
			ok := lox.Ejecutar(c.fuente, &salida)
			if c.falla && ok {
				t.Errorf("Se esperaba un error pero el script pasó: %s", c.fuente)
			}
			
			if !c.falla && !ok {
				t.Errorf("El script falló inesperadamente: %s\nSalida: %s", c.fuente, salida.String())
			}
			
			if c.falla && c.error != "" {
				if !strings.Contains(salida.String(), c.error) {
					t.Errorf("Se esperaba que el error contuviera '%s', pero la salida fue: %s", c.error, salida.String())
				}
			}
		})
	}
}
