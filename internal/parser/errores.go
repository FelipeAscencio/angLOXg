package parser

import (
	"fmt"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Error representa una falla de sintaxis, atada a la línea donde se detectó.
type Error struct {
	Linea   int
	Mensaje string
}

// Error arma el mensaje de la falla.
//
// Se llama así, y no "Mensaje", porque es el método que Go exige para que el
// tipo sirva como error.
func (e *Error) Error() string {
	return fmt.Sprintf("[línea %d] Error de sintaxis: %s", e.Linea, e.Mensaje)
}

// synchronize descarta tokens hasta llegar a un punto donde tenga sentido
// retomar el parseo.
//
// Después de un error, el cursor queda en medio de una construcción rota y todo
// lo que siga sería ruido. Saltar hasta el próximo punto y coma, o hasta la
// próxima palabra que abra una sentencia, evita esa catarata de errores falsos.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		// El punto y coma cierra una sentencia: lo que viene después arranca
		// limpio.
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
