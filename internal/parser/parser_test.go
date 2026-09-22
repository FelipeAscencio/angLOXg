package parser

import (
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/escaner"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
)

// Representa una entrada de prueba para expresiones individuales.
type casoDeParseo struct {
	nombre   string
	fuente   string
	esperado string
	errores  []string
}

// Escanea y parsea una expresión suelta, retornando su representación y errores.
func representarExpresion(t *testing.T, fuente string) (string, []error) {
	t.Helper()

	tokens, erroresLexicos := escaner.Nuevo(fuente).EscanearTokens()
	if len(erroresLexicos) > 0 {
		t.Fatalf("el fuente de la prueba no escanea: %v", erroresLexicos)
	}

	expresion, errores := Nuevo(tokens).ParsearExpresion()
	if expresion == nil {
		return "", errores
	}

	return sintaxis.Representar(expresion), errores
}

// Ejecuta una tabla de pruebas para expresiones.
func correrCasosDeExpresion(t *testing.T, casos []casoDeParseo) {
	t.Helper()

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenido, errores := representarExpresion(t, caso.fuente)
			verificarErrores(t, errores, caso.errores)

			if obtenido != caso.esperado {
				t.Errorf("se obtuvo %q y se esperaba %q", obtenido, caso.esperado)
			}
		})
	}
}

// Verifica que los errores de sintaxis coincidan con los esperados.
func verificarErrores(t *testing.T, obtenidos []error, esperados []string) {
	t.Helper()

	if len(obtenidos) != len(esperados) {
		t.Fatalf("se obtuvieron %d errores y se esperaban %d\nobtenidos: %v",
			len(obtenidos), len(esperados), obtenidos)
	}

	for i, esperado := range esperados {
		if !strings.Contains(obtenidos[i].Error(), esperado) {
			t.Errorf("error %d: %q, se esperaba que contuviera %q",
				i, obtenidos[i].Error(), esperado)
		}
	}
}

// Representa una tabla de prueba para un programa completo con múltiples sentencias.
type casoDePrograma struct {
	nombre   string
	fuente   string
	esperado []string
	errores  []string
}

// Escanea y parsea un programa entero, retornando las sentencias representadas y sus errores.
func representarPrograma(t *testing.T, fuente string) ([]string, []error) {
	t.Helper()

	tokens, erroresLexicos := escaner.Nuevo(fuente).EscanearTokens()
	if len(erroresLexicos) > 0 {
		t.Fatalf("el fuente de la prueba no escanea: %v", erroresLexicos)
	}

	sentencias, errores := Nuevo(tokens).Parsear()

	representadas := make([]string, 0, len(sentencias))
	for _, sentencia := range sentencias {
		representadas = append(representadas, sintaxis.RepresentarSentencia(sentencia))
	}

	return representadas, errores
}

// Ejecuta una tabla de pruebas para programas completos.
func correrCasosDePrograma(t *testing.T, casos []casoDePrograma) {
	t.Helper()

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenidas, errores := representarPrograma(t, caso.fuente)
			verificarErrores(t, errores, caso.errores)

			if len(obtenidas) != len(caso.esperado) {
				t.Fatalf("se obtuvieron %d sentencias y se esperaban %d\nobtenidas: %v",
					len(obtenidas), len(caso.esperado), obtenidas)
			}

			for i, esperada := range caso.esperado {
				if obtenidas[i] != esperada {
					t.Errorf("sentencia %d: se obtuvo %q y se esperaba %q",
						i, obtenidas[i], esperada)
				}
			}
		})
	}
}
