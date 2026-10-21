package compilador

import (
	"fmt"
	"os"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Precedencias ordenadas de menor a mayor.
type Precedencia int

// Constantes de precedencia para el Pratt Parser.
const (
	PrecNone Precedencia = iota
	PrecAssignment
	PrecOr
	PrecAnd
	PrecEquality
	PrecComparison
	PrecTerm
	PrecFactor
	PrecUnary
	PrecCall
	PrecPrimary
)

// La ReglaParseo asocia funciones de parseo a un tipo de token.
type ReglaParseo struct {
	prefijo     func(*Compilador)
	infijo      func(*Compilador)
	precedencia Precedencia
}

// El Compilador orquesta la traducción de tokens a bytecode (Single-Pass).
type Compilador struct {
	tokens []token.Token
	actual int
	chunk  *Chunk
	reglas map[token.TipoDeToken]ReglaParseo
}

// Toma una lista de tokens y devuelve un Chunk listo para la VM.
func Compilar(tokens []token.Token) *Chunk {
	comp := &Compilador{
		tokens: tokens,
		actual: 0,
		chunk:  NuevoChunk(),
	}

	comp.iniciarReglas()
	comp.parsearPrecedencia(PrecAssignment)
	comp.emitirByte(uint8(OpReturn))
	return comp.chunk
}

// Configura la tabla del Pratt Parser.
func (c *Compilador) iniciarReglas() {
	c.reglas = map[token.TipoDeToken]ReglaParseo{
		token.LEFT_PAREN: {prefijo: agrupamiento, infijo: nil, precedencia: PrecNone},
		token.MINUS:      {prefijo: unario, infijo: binario, precedencia: PrecTerm},
		token.PLUS:       {prefijo: nil, infijo: binario, precedencia: PrecTerm},
		token.SLASH:      {prefijo: nil, infijo: binario, precedencia: PrecFactor},
		token.STAR:       {prefijo: nil, infijo: binario, precedencia: PrecFactor},
		token.NUMBER:     {prefijo: numero, infijo: nil, precedencia: PrecNone},
	}
}

// Parsea expresiones según la precedencia y emite bytecode correspondiente.
func (c *Compilador) parsearPrecedencia(prec Precedencia) {
	c.avanzar()
	reglaPrefijo := c.obtenerRegla(c.tokenAnterior().Tipo).prefijo
	if reglaPrefijo == nil {
		fmt.Fprintf(os.Stderr, "Se esperaba una expresión en la línea %d.\n", c.tokenAnterior().Linea)
		return
	}

	reglaPrefijo(c)
	for prec <= c.obtenerRegla(c.tokenActual().Tipo).precedencia {
		c.avanzar()
		reglaInfijo := c.obtenerRegla(c.tokenAnterior().Tipo).infijo
		reglaInfijo(c)
	}
}

// ===================
// Funciones de parseo
// ===================

func numero(c *Compilador) {
	valor := c.tokenAnterior().Literal.(float64)
	c.emitirConstante(valor)
}

func agrupamiento(c *Compilador) {
	c.parsearPrecedencia(PrecAssignment)
	c.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de la expresión.")
}

func unario(c *Compilador) {
	tipoOperador := c.tokenAnterior().Tipo

	// Compilamos el operando
	c.parsearPrecedencia(PrecUnary)

	// Emitimos la instrucción correspondiente
	switch tipoOperador {
	case token.MINUS:
		c.emitirByte(uint8(OpNegate))
	}
}

func binario(c *Compilador) {
	tipoOperador := c.tokenAnterior().Tipo
	regla := c.obtenerRegla(tipoOperador)

	// Compilamos el operando derecho con un nivel más de precedencia
	c.parsearPrecedencia(regla.precedencia + 1)

	switch tipoOperador {
	case token.PLUS:
		c.emitirByte(uint8(OpAdd))
	case token.MINUS:
		c.emitirByte(uint8(OpSubtract))
	case token.STAR:
		c.emitirByte(uint8(OpMultiply))
	case token.SLASH:
		c.emitirByte(uint8(OpDivide))
	}
}

// ====================
// Navegación y emisión
// ====================

func (c *Compilador) avanzar() {
	c.actual++
}

func (c *Compilador) tokenActual() token.Token {
	return c.tokens[c.actual]
}

func (c *Compilador) tokenAnterior() token.Token {
	return c.tokens[c.actual-1]
}

func (c *Compilador) consumir(tipo token.TipoDeToken, mensaje string) {
	if c.tokenActual().Tipo == tipo {
		c.avanzar()
		return
	}
	
	fmt.Fprintf(os.Stderr, "[Línea %d] Error: %s\n", c.tokenActual().Linea, mensaje)
}

func (c *Compilador) emitirByte(b uint8) {
	c.chunk.Escribir(b, c.tokenAnterior().Linea)
}

func (c *Compilador) emitirConstante(valor any) {
	indice := c.chunk.AgregarConstante(valor)
	c.emitirByte(uint8(OpConstant))
	c.emitirByte(uint8(indice))
}

func (c *Compilador) obtenerRegla(tipo token.TipoDeToken) ReglaParseo {
	if regla, existe := c.reglas[tipo]; existe {
		return regla
	}

	return ReglaParseo{nil, nil, PrecNone}
}
