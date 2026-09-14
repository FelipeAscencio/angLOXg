package token

import "fmt"

// Token es la unidad mínima con significado que produce el escáner.
//
// Los únicos tipos dinámicos que puede contener "Literal" son:
//
//	float64 -> para los tokens NUMBER.
//	string  -> para los tokens STRING, ya sin las comillas.
//	nil     -> para todo el resto, incluidos TRUE, FALSE y NIL, donde el propio
//	           "Tipo" alcanza para saber a qué valor resuelven.
type Token struct {
	Tipo    TipoDeToken // Categoría léxica del token.
	Lexema  string      // Los caracteres crudos, tal cual aparecen en el fuente.
	Literal any         // Valor ya resuelto del literal, o nil si no es uno.
	Linea   int         // Línea del fuente, para reportar errores más precisos.
}

// Nuevo construye un token.
func Nuevo(tipo TipoDeToken, lexema string, literal any, linea int) Token {
	return Token{
		Tipo:    tipo,
		Lexema:  lexema,
		Literal: literal,
		Linea:   linea,
	}
}

// String arma una representación legible del token: el nombre de la categoría,
// acompañado del nombre del identificador o del valor del literal si los tiene.
//
// Se llama así, y no "Texto", porque es el método que Go espera para poder
// imprimir el valor con el paquete "fmt".
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
