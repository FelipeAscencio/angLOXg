package token

// Palabras reservadas del lenguaje. Cualquier lexema alfanumérico que no esté
// en esta tabla es un identificador.
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

// ResolverIdentificador toma un lexema alfanumérico ya escaneado y devuelve la
// palabra clave que le corresponde, o IDENTIFIER si no es una.
func ResolverIdentificador(lexema string) TipoDeToken {
	if tipo, ok := palabrasClave[lexema]; ok {
		return tipo
	}

	return IDENTIFIER
}
