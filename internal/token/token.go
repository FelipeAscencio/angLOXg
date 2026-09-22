package token

import "fmt"

// Representa la unidad mínima con significado léxico producida por el escáner.
type Token struct {
	Tipo    TipoDeToken // Categoría léxica del token.
	Lexema  string      // Caracteres crudos tal como aparecen en el código fuente.
	Literal any         // Valor ya resuelto del literal, o nil si no es uno.
	Linea   int         // Línea del código fuente para informes de error precisos.
}

// Construye un nuevo token.
func Nuevo(tipo TipoDeToken, lexema string, literal any, linea int) Token {
	return Token{
		Tipo:    tipo,
		Lexema:  lexema,
		Literal: literal,
		Linea:   linea,
	}
}

// Arma una representación legible del token.
func (t Token) String() string {
	switch t.Tipo {
	case IDENTIFIER:
		return fmt.Sprintf("%s<%s>", t.Tipo, t.Lexema)
	case STRING, NUMBER:
		return fmt.Sprintf("%s<%v>", t.Tipo, t.Literal)
	default:
		return t.Tipo.String()
	}
}
