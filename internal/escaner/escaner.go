package escaner

import (
	"fmt"
	"strconv"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Centinela que devuelven "mirar" y "mirarSiguiente" cuando ya no queda fuente
// por leer.
const byteFinal = byte(0)

// Escaner recorre el código fuente carácter por carácter y lo convierte en una
// lista de tokens.
//
// El lexema que se está capturando en cada momento es el que queda entre
// "inicio" y "actual" dentro de "fuente".
type Escaner struct {
	fuente  []byte        // Fuente crudo, sin significado todavía.
	tokens  []token.Token // Tokens ya reconocidos.
	errores []error       // Fallas léxicas acumuladas.

	inicio int // Byte donde arranca el lexema que se está leyendo.
	actual int // Byte donde está parado el cursor.
	linea  int // Línea donde está parado el cursor.
}

// Nuevo construye un escáner posicionado al principio del fuente.
func Nuevo(fuente string) *Escaner {
	return &Escaner{
		fuente: []byte(fuente),
		linea:  1,
	}
}

// EscanearTokens recorre todo el fuente y devuelve los tokens encontrados junto
// con las fallas léxicas detectadas.
//
// El escaneo continúa después de cada error, así que se reportan todos en una
// sola pasada. Si la lista de errores no viene vacía, los tokens no son
// confiables.
func (e *Escaner) EscanearTokens() ([]token.Token, []error) {
	for !e.esFinal() {
		// Arranca un lexema nuevo.
		e.inicio = e.actual
		e.escanearToken()
	}

	// La lista siempre se cierra con un EOF.
	e.inicio = e.actual
	e.agregarToken(token.EOF, nil)

	return e.tokens, e.errores
}

// escanearToken reconoce un único token a partir de la posición actual.
func (e *Escaner) escanearToken() {
	c := e.avanzar()

	switch c {
	// Los espacios en blanco no generan tokens.
	case ' ', '\r', '\t':
	case '\n':
		e.linea++

	// Tokens de un solo carácter.
	case '(':
		e.agregarToken(token.LEFT_PAREN, nil)
	case ')':
		e.agregarToken(token.RIGHT_PAREN, nil)
	case '{':
		e.agregarToken(token.LEFT_BRACE, nil)
	case '}':
		e.agregarToken(token.RIGHT_BRACE, nil)
	case ',':
		e.agregarToken(token.COMMA, nil)
	case '-':
		e.agregarToken(token.MINUS, nil)
	case '+':
		e.agregarToken(token.PLUS, nil)
	case ';':
		e.agregarToken(token.SEMICOLON, nil)
	case '*':
		e.agregarToken(token.STAR, nil)
	case '%':
		e.agregarToken(token.PERCENT, nil)

	// Tokens de uno o dos caracteres: se consume el lexema más largo posible,
	// así que ante "!=" nunca se devuelve un BANG suelto.
	case '!':
		e.agregarTokenSiCoincide('=', token.BANG_EQUAL, token.BANG)
	case '=':
		e.agregarTokenSiCoincide('=', token.EQUAL_EQUAL, token.EQUAL)
	case '<':
		e.agregarTokenSiCoincide('=', token.LESS_EQUAL, token.LESS)
	case '>':
		e.agregarTokenSiCoincide('=', token.GREATER_EQUAL, token.GREATER)

	// La barra es división, salvo que venga duplicada: ahí arranca un
	// comentario y se descarta el resto de la línea.
	case '/':
		if !e.coincide('/') {
			e.agregarToken(token.SLASH, nil)
			break
		}

		// El '\n' no se consume acá: lo procesa la vuelta siguiente del loop,
		// que es la que lleva la cuenta de las líneas.
		for e.mirar() != '\n' && !e.esFinal() {
			e.avanzar()
		}

	// Literales de cadena.
	case '"':
		e.escanearCadena()

	default:
		switch {
		case esDigito(c):
			e.escanearNumero()
		case esLetra(c):
			e.escanearIdentificador()
		default:
			e.registrarError("carácter inesperado %q", c)
		}
	}
}

// agregarTokenSiCoincide agrega "conCoincidencia" si el byte actual es el
// esperado (y lo consume), o "sinCoincidencia" en caso contrario.
func (e *Escaner) agregarTokenSiCoincide(esperado byte, conCoincidencia, sinCoincidencia token.TipoDeToken) {
	if e.coincide(esperado) {
		e.agregarToken(conCoincidencia, nil)
		return
	}

	e.agregarToken(sinCoincidencia, nil)
}

// escanearCadena consume un literal de cadena. Se llama con la comilla de
// apertura ya consumida.
//
// No hay secuencias de escape, así que la primera comilla doble cierra el
// literal. Un salto de línea no lo termina: las cadenas son multilínea.
func (e *Escaner) escanearCadena() {
	// Se guarda la línea de apertura, que es la que se reporta al final.
	lineaInicial := e.linea

	for !e.esFinal() && e.mirar() != '"' {
		if e.mirar() == '\n' {
			e.linea++
		}

		e.avanzar()
	}

	if e.esFinal() {
		e.registrarErrorEnLinea(lineaInicial, "cadena sin cerrar")
		return
	}

	// Consumimos la comilla de cierre.
	e.avanzar()

	// El literal es el contenido, ya sin las comillas.
	literal := string(e.fuente[e.inicio+1 : e.actual-1])
	e.agregarTokenEnLinea(token.STRING, literal, lineaInicial)
}

// escanearNumero consume un literal numérico. Se llama con el primer dígito ya
// consumido.
//
// Tanto "42" como "4.2" resuelven a un float64. No se admite notación
// científica ni prefijos de base.
func (e *Escaner) escanearNumero() {
	// Parte entera.
	for esDigito(e.mirar()) {
		e.avanzar()
	}

	// Parte decimal, sólo si atrás del punto viene otro dígito. Ese chequeo es
	// el que deja afuera a "123." (el punto queda suelto) y, como después no
	// se vuelve a mirar si hay otro punto, también a "1.2.3".
	if e.mirar() == '.' && esDigito(e.mirarSiguiente()) {
		// Consumimos el punto.
		e.avanzar()

		for esDigito(e.mirar()) {
			e.avanzar()
		}
	}

	lexema := e.lexema()

	literal, err := strconv.ParseFloat(lexema, 64)
	if err != nil {
		e.registrarError("número inválido %q", lexema)
		return
	}

	e.agregarToken(token.NUMBER, literal)
}

// escanearIdentificador consume un identificador o una palabra reservada. Se
// llama con el primer carácter ya consumido.
//
// Primero se consume el lexema más largo posible y recién después se pregunta
// si es una palabra clave, así "orden" nunca se parte en un OR seguido de un
// identificador "den".
func (e *Escaner) escanearIdentificador() {
	for esAlfanumerico(e.mirar()) {
		e.avanzar()
	}

	e.agregarToken(token.ResolverIdentificador(e.lexema()), nil)
}

// ---------- Núcleo ---------- //

// lexema devuelve el lexema que se está capturando.
func (e *Escaner) lexema() string {
	return string(e.fuente[e.inicio:e.actual])
}

// agregarToken agrega a la lista un token con el lexema actual, en la línea
// actual.
func (e *Escaner) agregarToken(tipo token.TipoDeToken, literal any) {
	e.agregarTokenEnLinea(tipo, literal, e.linea)
}

// agregarTokenEnLinea agrega un token atribuyéndolo a una línea puntual, para
// los lexemas que abarcan varias líneas y dejaron el cursor más adelante.
func (e *Escaner) agregarTokenEnLinea(tipo token.TipoDeToken, literal any, linea int) {
	e.tokens = append(e.tokens, token.Nuevo(tipo, e.lexema(), literal, linea))
}

// registrarError anota una falla léxica en la línea actual y sigue adelante.
func (e *Escaner) registrarError(formato string, args ...any) {
	e.registrarErrorEnLinea(e.linea, formato, args...)
}

// registrarErrorEnLinea anota una falla léxica atribuyéndola a una línea
// puntual.
func (e *Escaner) registrarErrorEnLinea(linea int, formato string, args ...any) {
	e.errores = append(e.errores, &Error{
		Linea:   linea,
		Mensaje: fmt.Sprintf(formato, args...),
	})
}

// ---------- Auxiliares ---------- //

// esFinal indica si ya se consumió todo el fuente.
func (e *Escaner) esFinal() bool {
	return e.actual >= len(e.fuente)
}

// mirar devuelve el byte actual sin consumirlo.
func (e *Escaner) mirar() byte {
	if e.esFinal() {
		return byteFinal
	}

	return e.fuente[e.actual]
}

// mirarSiguiente devuelve el byte siguiente al actual sin consumir nada.
//
// Lo usa el escaneo de números: al encontrar un punto hay que confirmar que
// atrás venga un dígito antes de darlo por parte del literal.
func (e *Escaner) mirarSiguiente() byte {
	if e.actual+1 >= len(e.fuente) {
		return byteFinal
	}

	return e.fuente[e.actual+1]
}

// avanzar consume el byte actual y lo devuelve.
func (e *Escaner) avanzar() byte {
	c := e.mirar()
	e.actual++

	return c
}

// coincide consume el byte actual sólo si es el esperado, e informa si lo hizo.
// Es la combinación de "mirar" y "avanzar" que usan los tokens de dos
// caracteres.
func (e *Escaner) coincide(esperado byte) bool {
	if e.mirar() != esperado {
		return false
	}

	e.avanzar()

	return true
}

// esDigito indica si el byte es un dígito decimal ASCII.
func esDigito(c byte) bool {
	return c >= '0' && c <= '9'
}

// esLetra indica si el byte puede abrir un identificador: letras ASCII y el
// guion bajo.
func esLetra(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// esAlfanumerico indica si el byte puede continuar un identificador.
func esAlfanumerico(c byte) bool {
	return esLetra(c) || esDigito(c)
}
