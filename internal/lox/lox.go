package lox

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/escaner"
	"github.com/FelipeAscencio/angLOXg/internal/interprete"
	"github.com/FelipeAscencio/angLOXg/internal/parser"
	"github.com/FelipeAscencio/angLOXg/internal/semantica"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Ejecuta código Lox completo y escribe la salida o errores en el escritor provisto.
func Ejecutar(fuente string, salida io.Writer) bool {
	sentencias, ok := parsearFuente(fuente, salida)
	if !ok {
		return false
	}

	intp := interprete.NuevoInterprete(salida)
	analizador := semantica.Nuevo(intp)

	if err := analizador.Resolver(sentencias); err != nil {
		fmt.Fprintln(salida, err.Error())
		return false
	}

	if err := intp.Interpretar(sentencias); err != nil {
		fmt.Fprintln(salida, err.Error())
		return false
	}

	return true
}

// Ejecuta las fases de escaneo, parseo y análisis semántico sin llegar a interpretar.
func Resolver(fuente string, salida io.Writer) bool {
	sentencias, ok := parsearFuente(fuente, salida)
	if !ok {
		return false
	}

	intp := interprete.NuevoInterprete(salida)
	analizador := semantica.Nuevo(intp)
	if err := analizador.Resolver(sentencias); err != nil {
		fmt.Fprintln(salida, err.Error())
		return false
	}

	fmt.Fprintln(salida, "Análisis semántico completado sin errores. El scope es correcto.")
	return true
}

// Escanea y parsea código Lox, imprimiendo el árbol sintáctico (AST) resultante.
func Parsear(fuente string, salida io.Writer) bool {
	sentencias, ok := parsearFuente(fuente, salida)
	if !ok {
		return false
	}

	imprimirArbol(sentencias, salida)

	return true
}

// Escanea código Lox e imprime los tokens generados por la fase léxica.
func Escanear(fuente string, salida io.Writer) bool {
	tokens, ok := escanearFuente(fuente, salida)
	if !ok {
		return false
	}

	imprimirTokens(tokens, salida)

	return true
}

// ==============================
// Funciones auxiliares internas.
// ==============================

func escanearFuente(fuente string, salida io.Writer) ([]token.Token, bool) {
	tokens, errores := escaner.Nuevo(fuente).EscanearTokens()
	for _, err := range errores {
		fmt.Fprintln(salida, err)
	}

	return tokens, len(errores) == 0
}

func imprimirTokens(tokens []token.Token, salida io.Writer) {
	for _, tok := range tokens {
		fmt.Fprintln(salida, tok)
	}
}

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

func imprimirArbol(sentencias []sintaxis.Stmt, salida io.Writer) {
	for _, sentencia := range sentencias {
		fmt.Fprintln(salida, sintaxis.RepresentarSentencia(sentencia))
	}
}
