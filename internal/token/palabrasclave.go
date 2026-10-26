package token

// Mapa de palabras reservadas del lenguaje Lox.
var palabrasClave = map[string]TipoDeToken{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// Resuelve un lexema alfanumérico para determinar si corresponde a una palabra clave reservada o si se trata de un identificador general.
func ResolverIdentificador(lexema string) TipoDeToken {
	if tipo, ok := palabrasClave[lexema]; ok {
		return tipo
	}

	return IDENTIFIER
}
