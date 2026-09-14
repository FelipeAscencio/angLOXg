package parser

import (
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/escaner"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
)

// casoDeParseo es una entrada de una tabla de pruebas: un fuente y la
// representación del árbol que tiene que salir de parsearlo.
type casoDeParseo struct {
	nombre string
	fuente string

	// Representación esperada del árbol. Vacía si se espera que falle.
	esperado string

	// Subcadenas que tienen que aparecer en los errores, una por error y en
	// orden. Si es nil, el parseo tiene que terminar sin fallas.
	errores []string
}

// representarExpresion escanea y parsea una expresión suelta, y devuelve su
// representación junto con los errores de sintaxis.
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

// correrCasosDeExpresion ejecuta una tabla de pruebas de expresiones.
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

// verificarErrores chequea que las fallas de sintaxis sean las esperadas.
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

// casoDePrograma es una tabla de pruebas para un programa entero, que puede
// tener varias sentencias.
type casoDePrograma struct {
	nombre   string
	fuente   string
	esperado []string
	errores  []string
}

// representarPrograma escanea y parsea un programa entero, y devuelve la
// representación de cada sentencia junto con los errores de sintaxis.
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

// correrCasosDePrograma ejecuta una tabla de pruebas de programas.
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
