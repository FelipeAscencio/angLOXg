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
	prefijo     func(*Compilador, bool)
	infijo      func(*Compilador, bool)
	precedencia Precedencia
}

// Estructura para simular el rastro de las variables locales en la pila[cite: 39].
type Local struct {
	nombre      string
	profundidad int
}

// El Compilador orquesta la traducción de tokens a bytecode (Single-Pass).
type Compilador struct {
	tokens []token.Token
	actual int
	chunk  *Chunk
	reglas map[token.TipoDeToken]ReglaParseo
	locales          []Local
	profundidadScope int
}

// Toma una lista de tokens y devuelve un Chunk listo para la VM.
func Compilar(tokens []token.Token) *Chunk {
	comp := &Compilador{
		tokens:           tokens,
		actual:           0,
		chunk:            NuevoChunk(),
		locales:          make([]Local, 0),
		profundidadScope: 0,
	}

	comp.iniciarReglas()
	for !comp.esFin() {
		comp.declaracion()
	}

	comp.emitirByte(uint8(OpReturn))
	return comp.chunk
}

// Configura la tabla del Pratt Parser.
func (c *Compilador) iniciarReglas() {
	c.reglas = map[token.TipoDeToken]ReglaParseo{
		token.LEFT_PAREN:    {prefijo: agrupamiento, infijo: nil, precedencia: PrecNone},
		token.MINUS:         {prefijo: unario, infijo: binario, precedencia: PrecTerm},
		token.PLUS:          {prefijo: nil, infijo: binario, precedencia: PrecTerm},
		token.SLASH:         {prefijo: nil, infijo: binario, precedencia: PrecFactor},
		token.STAR:          {prefijo: nil, infijo: binario, precedencia: PrecFactor},
		token.NUMBER:        {prefijo: numero, infijo: nil, precedencia: PrecNone},
		token.STRING:        {prefijo: cadena, infijo: nil, precedencia: PrecNone},
		token.FALSE:         {prefijo: literal, infijo: nil, precedencia: PrecNone},
		token.TRUE:          {prefijo: literal, infijo: nil, precedencia: PrecNone},
		token.NIL:           {prefijo: literal, infijo: nil, precedencia: PrecNone},
		token.BANG:          {prefijo: unario, infijo: nil, precedencia: PrecNone},
		token.BANG_EQUAL:    {prefijo: nil, infijo: binario, precedencia: PrecEquality},
		token.EQUAL_EQUAL:   {prefijo: nil, infijo: binario, precedencia: PrecEquality},
		token.GREATER:       {prefijo: nil, infijo: binario, precedencia: PrecComparison},
		token.GREATER_EQUAL: {prefijo: nil, infijo: binario, precedencia: PrecComparison},
		token.LESS:          {prefijo: nil, infijo: binario, precedencia: PrecComparison},
		token.LESS_EQUAL:    {prefijo: nil, infijo: binario, precedencia: PrecComparison},
		token.IDENTIFIER:    {prefijo: variable, infijo: nil, precedencia: PrecNone},
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

	puedeAsignar := prec <= PrecAssignment
	reglaPrefijo(c, puedeAsignar)
	for prec <= c.obtenerRegla(c.tokenActual().Tipo).precedencia {
		c.avanzar()
		reglaInfijo := c.obtenerRegla(c.tokenAnterior().Tipo).infijo
		reglaInfijo(c, puedeAsignar)
	}

	if puedeAsignar && c.coincide(token.EQUAL) {
		fmt.Fprintf(os.Stderr, "[Línea %d] Error: Objetivo de asignación no válido.\n", c.tokenAnterior().Linea)
	}
}

// ===================
// Funciones de parseo
// ===================

func variable(c *Compilador, puedeAsignar bool) {
	nombre := c.tokenAnterior().Lexema
	arg := c.resolverLocal(nombre)
	var opGet, opSet uint8
	if arg != -1 {
		// La variable existe en el scope local.
		opGet = uint8(OpGetLocal)
		opSet = uint8(OpSetLocal)
	} else {
		// Asumimos que es global.
		arg = c.chunk.AgregarConstante(nombre)
		opGet = uint8(OpGetGlobal)
		opSet = uint8(OpSetGlobal)
	}

	if puedeAsignar && c.coincide(token.EQUAL) {
		c.expresion()
		c.emitirByte(opSet)
		c.emitirByte(uint8(arg))
	} else {
		c.emitirByte(opGet)
		c.emitirByte(uint8(arg))
	}
}

// Simula la búsqueda en la pila recorriendo el arreglo de locales hacia atrás.
func (c *Compilador) resolverLocal(nombre string) int {
	for i := len(c.locales) - 1; i >= 0; i-- {
		if c.locales[i].nombre == nombre {
			return i
		}
	}

	return -1
}

func numero(c *Compilador, puedeAsignar bool) {
	valor := c.tokenAnterior().Literal.(float64)
	c.emitirConstante(valor)
}

func cadena(c *Compilador, puedeAsignar bool) {
	valor := c.tokenAnterior().Literal.(string)
	c.emitirConstante(valor)
}

func literal(c *Compilador, puedeAsignar bool) {
	switch c.tokenAnterior().Tipo {
	case token.FALSE:
		c.emitirByte(uint8(OpFalse))
	case token.TRUE:
		c.emitirByte(uint8(OpTrue))
	case token.NIL:
		c.emitirByte(uint8(OpNil))
	}
}

func agrupamiento(c *Compilador, puedeAsignar bool) {
	c.parsearPrecedencia(PrecAssignment)
	c.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de la expresión.")
}

func unario(c *Compilador, puedeAsignar bool) {
	tipoOperador := c.tokenAnterior().Tipo
	c.parsearPrecedencia(PrecUnary)
	switch tipoOperador {
	case token.MINUS:
		c.emitirByte(uint8(OpNegate))
	case token.BANG:
		c.emitirByte(uint8(OpNot))
	}
}

func binario(c *Compilador, puedeAsignar bool) {
	tipoOperador := c.tokenAnterior().Tipo
	regla := c.obtenerRegla(tipoOperador)
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
	case token.BANG_EQUAL:
		c.emitirByte(uint8(OpEqual))
		c.emitirByte(uint8(OpNot))
	case token.EQUAL_EQUAL:
		c.emitirByte(uint8(OpEqual))
	case token.GREATER:
		c.emitirByte(uint8(OpGreater))
	case token.GREATER_EQUAL:
		c.emitirByte(uint8(OpLess))
		c.emitirByte(uint8(OpNot))
	case token.LESS:
		c.emitirByte(uint8(OpLess))
	case token.LESS_EQUAL:
		c.emitirByte(uint8(OpGreater))
		c.emitirByte(uint8(OpNot))
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

// =================================
// Parseo de estructuras (Gramática)
// =================================

func (c *Compilador) declaracion() {
	if c.coincide(token.VAR) {
		c.declaracionVariable()
	} else {
		c.sentencia()
	}
}

func (c *Compilador) declaracionVariable() {
	c.consumir(token.IDENTIFIER, "Se esperaba el nombre de la variable.")
	nombreToken := c.tokenAnterior()
	if c.coincide(token.EQUAL) {
		c.expresion()
	} else {
		c.emitirByte(uint8(OpNil))
	}

	c.consumir(token.SEMICOLON, "Se esperaba ';' después de la declaración de la variable.")
	if c.profundidadScope > 0 {
		// Es una variable local.
		c.locales = append(c.locales, Local{
			nombre:      nombreToken.Lexema,
			profundidad: c.profundidadScope,
		})
	} else {
		// Es una variable global.
		indiceGlobal := c.chunk.AgregarConstante(nombreToken.Lexema)
		c.emitirByte(uint8(OpDefineGlobal))
		c.emitirByte(uint8(indiceGlobal))
	}
}

func (c *Compilador) sentencia() {
	if c.coincide(token.PRINT) {
		c.sentenciaImprimir()
	} else if c.coincide(token.LEFT_BRACE) {
		c.empezarScope()
		c.bloque()
		c.terminarScope()
	} else {
		c.sentenciaExpresion()
	}
}

func (c *Compilador) sentenciaImprimir() {
	c.expresion()
	c.consumir(token.SEMICOLON, "Se esperaba ';' después del valor.")
	c.emitirByte(uint8(OpPrint))
}

func (c *Compilador) sentenciaExpresion() {
	c.expresion()
	c.consumir(token.SEMICOLON, "Se esperaba ';' después de la expresión.")
	c.emitirByte(uint8(OpPop))
}

func (c *Compilador) expresion() {
	c.parsearPrecedencia(PrecAssignment)
}

func (c *Compilador) bloque() {
	for !c.esFin() && c.tokenActual().Tipo != token.RIGHT_BRACE {
		c.declaracion()
	}

	c.consumir(token.RIGHT_BRACE, "Se esperaba '}' al final del bloque.")
}

func (c *Compilador) empezarScope() {
	c.profundidadScope++
}

func (c *Compilador) terminarScope() {
	c.profundidadScope--
	for len(c.locales) > 0 && c.locales[len(c.locales)-1].profundidad > c.profundidadScope {
		c.emitirByte(uint8(OpPop))
		c.locales = c.locales[:len(c.locales)-1]
	}
}

func (c *Compilador) coincide(tipo token.TipoDeToken) bool {
	if c.esFin() {
		return false
	}

	if c.tokenActual().Tipo == tipo {
		c.avanzar()
		return true
	}

	return false
}

func (c *Compilador) esFin() bool {
	return c.tokenActual().Tipo == token.EOF
}
