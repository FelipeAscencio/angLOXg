package parser

import (
	"fmt"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Representa un error de sintaxis asociado a una línea del código fuente.
type Error struct {
	Linea   int
	Mensaje string
}

// Implementa la interfaz error de Go.
func (e *Error) Error() string {
	return fmt.Sprintf("[línea %d] Error de sintaxis: %s", e.Linea, e.Mensaje)
}

// Descarta tokens hasta un punto seguro para retomar el análisis sintáctico.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Tipo == token.SEMICOLON {
			return
		}

		switch p.peek().Tipo {
		case token.FUN, token.VAR, token.FOR, token.IF,
			token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
