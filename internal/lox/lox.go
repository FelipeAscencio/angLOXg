package lox

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/escaner"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Ejecutar corre código Lox y escribe la salida en "salida". Devuelve false si
// el código tenía errores.
//
// Toda la salida, incluidos los mensajes de error, va a "salida" y no a la
// salida de errores del proceso.
func Ejecutar(fuente string, salida io.Writer) bool {
	tokens, ok := escanearFuente(fuente, salida)
	if !ok {
		return false
	}

	// TODO: acá van el parser y el intérprete. Hasta que existan, el modo de
	// ejecución vuelca los tokens igual que el modo de escaneo.
	imprimirTokens(tokens, salida)

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
