package token

// TipoDeToken identifica la categoría léxica de un token.
type TipoDeToken int

// Categorías léxicas que reconoce angLOXg.
//
// El orden de estas constantes debe coincidir con el de "nombresDeTipos".
const (
	// Tokens de un solo carácter.
	LEFT_PAREN TipoDeToken = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	MINUS
	PLUS
	SEMICOLON
	STAR
	PERCENT

	// El "/" es un token de un solo carácter, salvo cuando viene duplicado
	// ("//"): ahí arranca un comentario y el escáner descarta el resto de la
	// línea.
	SLASH

	// Tokens de uno o dos caracteres.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literales.
	IDENTIFIER
	STRING
	NUMBER

	// Palabras clave.
	AND
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	TRUE
	VAR
	WHILE

	// Fin de archivo.
	EOF
)

// Nombre legible de cada categoría, indexado por el valor del "TipoDeToken".
var nombresDeTipos = [...]string{
	LEFT_PAREN:    "LEFT_PAREN",
	RIGHT_PAREN:   "RIGHT_PAREN",
	LEFT_BRACE:    "LEFT_BRACE",
	RIGHT_BRACE:   "RIGHT_BRACE",
	COMMA:         "COMMA",
	MINUS:         "MINUS",
	PLUS:          "PLUS",
	SEMICOLON:     "SEMICOLON",
	STAR:          "STAR",
	PERCENT:       "PERCENT",
	SLASH:         "SLASH",
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	IDENTIFIER:    "IDENTIFIER",
	STRING:        "STRING",
	NUMBER:        "NUMBER",
	AND:           "AND",
	ELSE:          "ELSE",
	FALSE:         "FALSE",
	FUN:           "FUN",
	FOR:           "FOR",
	IF:            "IF",
	NIL:           "NIL",
	OR:            "OR",
	PRINT:         "PRINT",
	RETURN:        "RETURN",
	TRUE:          "TRUE",
	VAR:           "VAR",
	WHILE:         "WHILE",
	EOF:           "EOF",
}

// Chequeo en tiempo de compilación de que no falte ni sobre ningún nombre en
// "nombresDeTipos". Si las dos listas se desincronizan, el paquete no compila.
const _ = uint(len(nombresDeTipos) - (int(EOF) + 1))
const _ = uint((int(EOF) + 1) - len(nombresDeTipos))

// Devuelve el nombre legible de la categoría léxica para la impresión por formato.
func (t TipoDeToken) String() string {
	if t < 0 || int(t) >= len(nombresDeTipos) {
		return "UNKNOWN"
	}

	return nombresDeTipos[t]
}
