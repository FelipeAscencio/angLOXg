package lox

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/escaner"
	"github.com/FelipeAscencio/angLOXg/internal/interprete"
	"github.com/FelipeAscencio/angLOXg/internal/parser"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Ejecutar corre código Lox y escribe la salida en "salida". Devuelve false si
// el código tenía errores.
//
// Toda la salida, incluidos los mensajes de error, va a "salida" y no a la
// salida de errores del proceso.
func Ejecutar(fuente string, salida io.Writer) bool {
	sentencias, ok := parsearFuente(fuente, salida)
	if !ok {
		return false
	}

	intp := interprete.NuevoInterprete(salida)
	if err := intp.Interpretar(sentencias); err != nil {
		fmt.Fprintln(salida, err.Error())
		return false
	}

	return true
}

// Resolver escanea, parsea, y realiza el análisis semántico sin ejecutar.
// Es lo que corre el CLI con "--resolucion".
func Resolver(fuente string, salida io.Writer) bool {
	_, ok := parsearFuente(fuente, salida)
	if !ok {
		return false
	}

	// TODO: acá va el análisis semántico (resolución de variables locales).
	fmt.Fprintln(salida, "El análisis semántico está en construcción.")
	return true
}

// Parsear escanea y parsea código Lox, y vuelca el árbol de sintaxis en
// "salida", sin ejecutarlo. Es lo que corre el CLI con "--arbol".
func Parsear(fuente string, salida io.Writer) bool {
	sentencias, ok := parsearFuente(fuente, salida)
	if !ok {
		return false
	}

	imprimirArbol(sentencias, salida)

	return true
}

// Escanear escanea código Lox y vuelca los tokens en "salida", sin ejecutarlo.
// Es lo que corre el CLI con "--escaneo".
func Escanear(fuente string, salida io.Writer) bool {
	tokens, ok := escanearFuente(fuente, salida)
	if !ok {
		return false
	}

	imprimirTokens(tokens, salida)

	return true
}

// escanearFuente corre el escáner y reporta las fallas léxicas que encontró.
func escanearFuente(fuente string, salida io.Writer) ([]token.Token, bool) {
	tokens, errores := escaner.Nuevo(fuente).EscanearTokens()
	for _, err := range errores {
		fmt.Fprintln(salida, err)
	}

	// Con errores léxicos los tokens no son confiables y el proceso se corta.
	return tokens, len(errores) == 0
}

// imprimirTokens vuelca la lista de tokens, uno por línea.
func imprimirTokens(tokens []token.Token, salida io.Writer) {
	for _, tok := range tokens {
		fmt.Fprintln(salida, tok)
	}
}

// parsearFuente encadena el escáner con el parser y reporta las fallas de
// cualquiera de los dos.
func parsearFuente(fuente string, salida io.Writer) ([]sintaxis.Stmt, bool) {
	tokens, ok := escanearFuente(fuente, salida)
	if !ok {
		return nil, false
	}

	sentencias, errores := parser.Nuevo(tokens).Parsear()
	for _, err := range errores {
		fmt.Fprintln(salida, err)
	}

	return sentencias, len(errores) == 0
}

// imprimirArbol vuelca el árbol de cada sentencia, una por línea.
func imprimirArbol(sentencias []sintaxis.Stmt, salida io.Writer) {
	for _, sentencia := range sentencias {
		fmt.Fprintln(salida, sintaxis.RepresentarSentencia(sentencia))
	}
}
