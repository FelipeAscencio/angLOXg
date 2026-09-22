package token

// Mapa de palabras reservadas del lenguaje Lox. Cualquier lexema alfanumérico
// que no figure en este conjunto se interpreta por defecto como un identificador.
var palabrasClave = map[string]TipoDeToken{
	"and":    AND,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// Resuelve un lexema alfanumérico para determinar si corresponde a una palabra clave
// reservada o si se trata de un identificador general.
func ResolverIdentificador(lexema string) TipoDeToken {
	if tipo, ok := palabrasClave[lexema]; ok {
		return tipo
	}

	return IDENTIFIER
}
