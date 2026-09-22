package escaner

import (
	"fmt"
	"strconv"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

const byteFinal = byte(0)

// Realiza el análisis léxico del código fuente.
type Escaner struct {
	fuente  []byte
	tokens  []token.Token
	errores []error
	inicio  int
	actual  int
	linea   int
}

// Crea una nueva instancia del escáner.
func Nuevo(fuente string) *Escaner {
	return &Escaner{
		fuente: []byte(fuente),
		linea:  1,
	}
}

// Recorre todo el fuente y retorna la lista de tokens y errores léxicos.
func (e *Escaner) EscanearTokens() ([]token.Token, []error) {
	for !e.esFinal() {
		e.inicio = e.actual
		e.escanearToken()
	}

	e.inicio = e.actual
	e.agregarToken(token.EOF, nil)

	return e.tokens, e.errores
}

// Reconoce el siguiente token a partir de la posición actual.
func (e *Escaner) escanearToken() {
	c := e.avanzar()

	switch c {
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

	// Tokens de uno o dos caracteres.
	case '!':
		e.agregarTokenSiCoincide('=', token.BANG_EQUAL, token.BANG)
	case '=':
		e.agregarTokenSiCoincide('=', token.EQUAL_EQUAL, token.EQUAL)
	case '<':
		e.agregarTokenSiCoincide('=', token.LESS_EQUAL, token.LESS)
	case '>':
		e.agregarTokenSiCoincide('=', token.GREATER_EQUAL, token.GREATER)

	// Comentarios o división.
	case '/':
		if !e.coincide('/') {
			e.agregarToken(token.SLASH, nil)
			break
		}

		for e.mirar() != '\n' && !e.esFinal() {
			e.avanzar()
		}

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

func (e *Escaner) agregarTokenSiCoincide(esperado byte, conCoincidencia, sinCoincidencia token.TipoDeToken) {
	if e.coincide(esperado) {
		e.agregarToken(conCoincidencia, nil)
		return
	}

	e.agregarToken(sinCoincidencia, nil)
}

// Consume un literal de cadena multilínea.
func (e *Escaner) escanearCadena() {
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

	e.avanzar()

	literal := string(e.fuente[e.inicio+1 : e.actual-1])
	e.agregarTokenEnLinea(token.STRING, literal, lineaInicial)
}

// Consume un literal numérico decimal (float64).
func (e *Escaner) escanearNumero() {
	for esDigito(e.mirar()) {
		e.avanzar()
	}

	if e.mirar() == '.' && esDigito(e.mirarSiguiente()) {
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

// Consume un identificador o palabra reservada.
func (e *Escaner) escanearIdentificador() {
	for esAlfanumerico(e.mirar()) {
		e.avanzar()
	}

	e.agregarToken(token.ResolverIdentificador(e.lexema()), nil)
}

// ========================
// Núcleo y manejo interno.
// ========================

func (e *Escaner) lexema() string {
	return string(e.fuente[e.inicio:e.actual])
}

func (e *Escaner) agregarToken(tipo token.TipoDeToken, literal any) {
	e.agregarTokenEnLinea(tipo, literal, e.linea)
}

func (e *Escaner) agregarTokenEnLinea(tipo token.TipoDeToken, literal any, linea int) {
	e.tokens = append(e.tokens, token.Nuevo(tipo, e.lexema(), literal, linea))
}

func (e *Escaner) registrarError(formato string, args ...any) {
	e.registrarErrorEnLinea(e.linea, formato, args...)
}

func (e *Escaner) registrarErrorEnLinea(linea int, formato string, args ...any) {
	e.errores = append(e.errores, &Error{
		Linea:   linea,
		Mensaje: fmt.Sprintf(formato, args...),
	})
}

// ======================
// Auxiliares de lectura.
// ======================

func (e *Escaner) esFinal() bool {
	return e.actual >= len(e.fuente)
}

func (e *Escaner) mirar() byte {
	if e.esFinal() {
		return byteFinal
	}
	return e.fuente[e.actual]
}

func (e *Escaner) mirarSiguiente() byte {
	if e.actual+1 >= len(e.fuente) {
		return byteFinal
	}
	return e.fuente[e.actual+1]
}

func (e *Escaner) avanzar() byte {
	c := e.mirar()
	e.actual++
	return c
}

func (e *Escaner) coincide(esperado byte) bool {
	if e.mirar() != esperado {
		return false
	}
	e.avanzar()
	return true
}

func esDigito(c byte) bool {
	return c >= '0' && c <= '9'
}

func esLetra(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func esAlfanumerico(c byte) bool {
	return esLetra(c) || esDigito(c)
}
