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
	capturada   bool
}

// Representa un upvalue resuelto en tiempo de compilación.
type UpvalueLocal struct {
	indice  uint8
	esLocal bool
}

// El Compilador orquesta la traducción de tokens a bytecode (Single-Pass).
type Compilador struct {
	padre            *Compilador
	tokens           []token.Token
	actual           int
	funcion          *Funcion
	locales          []Local
	upvalues         []UpvalueLocal
	profundidadScope int
	reglas           map[token.TipoDeToken]ReglaParseo
	esInicializador  bool
}

// Toma una lista de tokens y devuelve un Chunk listo para la VM.
func Compilar(tokens []token.Token) *Funcion {
	comp := &Compilador{
		padre:            nil,
		tokens:           tokens,
		actual:           0,
		funcion:          NuevaFuncion(),
		locales:          make([]Local, 0),
		upvalues:         make([]UpvalueLocal, 0),
		profundidadScope: 0,
	}

	comp.locales = append(comp.locales, Local{profundidad: 0, nombre: "", capturada: false})
	comp.iniciarReglas()
	for !comp.esFin() {
		comp.declaracion()
	}

	comp.emitirByte(uint8(OpReturn))
	return comp.funcion
}

// Configura la tabla del Pratt Parser.
func (c *Compilador) iniciarReglas() {
	c.reglas = map[token.TipoDeToken]ReglaParseo{
		token.LEFT_PAREN:    {prefijo: agrupamiento, infijo: llamada, precedencia: PrecCall},
		token.MINUS:         {prefijo: unario, infijo: binario, precedencia: PrecTerm},
		token.PLUS:          {prefijo: nil, infijo: binario, precedencia: PrecTerm},
		token.SLASH:         {prefijo: nil, infijo: binario, precedencia: PrecFactor},
		token.STAR:          {prefijo: nil, infijo: binario, precedencia: PrecFactor},
		token.PERCENT:       {prefijo: nil, infijo: binario, precedencia: PrecFactor},
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
		token.AND:           {prefijo: nil, infijo: and, precedencia: PrecAnd},
		token.OR:            {prefijo: nil, infijo: or, precedencia: PrecOr},
		token.FUN:           {prefijo: nil, infijo: nil, precedencia: PrecNone},
		token.DOT:           {prefijo: nil, infijo: punto, precedencia: PrecCall},
		token.THIS:          {prefijo: este, infijo: nil, precedencia: PrecNone},
		token.SUPER:         {prefijo: superior, infijo: nil, precedencia: PrecNone},
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
	} else if arg = c.resolverUpvalue(nombre); arg != -1 {
		// La variable es un upvalue (capturada de un scope padre).
		opGet = uint8(OpGetUpvalue)
		opSet = uint8(OpSetUpvalue)
	} else {
		// Asumimos que es global.
		arg = c.funcion.Chunk.AgregarConstante(nombre)
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

func (c *Compilador) resolverLocal(nombre string) int {
	for i := len(c.locales) - 1; i >= 0; i-- {
		if c.locales[i].nombre == nombre {
			return i
		}
	}

	return -1
}

func (c *Compilador) resolverUpvalue(nombre string) int {
	if c.padre == nil {
		return -1
	}

	local := c.padre.resolverLocal(nombre)
	if local != -1 {
		c.padre.locales[local].capturada = true
		return c.agregarUpvalue(uint8(local), true)
	}

	upvalue := c.padre.resolverUpvalue(nombre)
	if upvalue != -1 {
		return c.agregarUpvalue(uint8(upvalue), false)
	}

	return -1
}

func (c *Compilador) agregarUpvalue(indice uint8, esLocal bool) int {
	for i, upv := range c.upvalues {
		if upv.indice == indice && upv.esLocal == esLocal {
			return i
		}
	}

	c.upvalues = append(c.upvalues, UpvalueLocal{indice: indice, esLocal: esLocal})
	c.funcion.UpvalueCount = len(c.upvalues)
	return len(c.upvalues) - 1
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
	case token.PERCENT:
		c.emitirByte(uint8(OpModulo))
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

func and(c *Compilador, puedeAsignar bool) {
	saltoFinal := c.emitirSalto(uint8(OpJumpIfFalse))
	c.emitirByte(uint8(OpPop))
	c.parsearPrecedencia(PrecAnd)
	c.parchearSalto(saltoFinal)
}

func or(c *Compilador, puedeAsignar bool) {
	saltoElse := c.emitirSalto(uint8(OpJumpIfFalse))
	saltoFinal := c.emitirSalto(uint8(OpJump))
	c.parchearSalto(saltoElse)
	c.emitirByte(uint8(OpPop))
	c.parsearPrecedencia(PrecOr)
	c.parchearSalto(saltoFinal)
}

// ====================
// Navegación y emisión
// ====================

func (c *Compilador) avanzar() {
	c.actual++
}

func (c *Compilador) tokenActual() token.Token {
	if c.actual >= len(c.tokens) {
		return token.Nuevo(token.EOF, "", nil, 0)
	}

	return c.tokens[c.actual]
}

func (c *Compilador) tokenAnterior() token.Token {
	if c.actual == 0 {
		return c.tokenActual()
	}

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
	c.funcion.Chunk.Escribir(b, c.tokenAnterior().Linea)
}

func (c *Compilador) emitirSalto(instruccion uint8) int {
	c.emitirByte(instruccion)
	c.emitirByte(0xff)
	c.emitirByte(0xff)
	return len(c.funcion.Chunk.Instrucciones) - 2
}

func (c *Compilador) parchearSalto(offset int) {
	salto := len(c.funcion.Chunk.Instrucciones) - offset - 2
	if salto > 65535 {
		fmt.Fprintf(os.Stderr, "El código a saltar es demasiado grande.\n")
	}

	c.funcion.Chunk.Instrucciones[offset] = uint8((salto >> 8) & 0xff)
	c.funcion.Chunk.Instrucciones[offset+1] = uint8(salto & 0xff)
}

func (c *Compilador) emitirLoop(inicio int) {
	c.emitirByte(uint8(OpLoop))
	salto := len(c.funcion.Chunk.Instrucciones) - inicio + 2
	if salto > 65535 {
		fmt.Fprintf(os.Stderr, "El cuerpo del bucle es demasiado grande.\n")
	}

	c.emitirByte(uint8((salto >> 8) & 0xff))
	c.emitirByte(uint8(salto & 0xff))
}

func (c *Compilador) emitirConstante(valor any) {
	indice := c.funcion.Chunk.AgregarConstante(valor)
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
	if c.coincide(token.CLASS) {
		c.declaracionClase()
	} else if c.coincide(token.FUN) {
		c.declaracionFuncion()
	} else if c.coincide(token.VAR) {
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
		indiceGlobal := c.funcion.Chunk.AgregarConstante(nombreToken.Lexema)
		c.emitirByte(uint8(OpDefineGlobal))
		c.emitirByte(uint8(indiceGlobal))
	}
}

func (c *Compilador) declaracionFuncion() {
	c.consumir(token.IDENTIFIER, "Se esperaba el nombre de la función.")
	nombreToken := c.tokenAnterior()
	indiceGlobal := c.funcion.Chunk.AgregarConstante(nombreToken.Lexema)

	esLocal := c.profundidadScope > 0
	if esLocal {
		c.locales = append(c.locales, Local{
			nombre:      nombreToken.Lexema,
			profundidad: c.profundidadScope,
		})
	}

	compiladorFuncion := &Compilador{
		padre:            c,
		tokens:           c.tokens,
		actual:           c.actual,
		funcion:          NuevaFuncion(),
		locales:          make([]Local, 0),
		upvalues:         make([]UpvalueLocal, 0),
		profundidadScope: 0,
		reglas:           c.reglas,
	}

	compiladorFuncion.funcion.Nombre = nombreToken.Lexema
	compiladorFuncion.locales = append(compiladorFuncion.locales, Local{profundidad: 0, nombre: "", capturada: false})
	compiladorFuncion.consumir(token.LEFT_PAREN, "Se esperaba '(' después del nombre de la función.")
	if compiladorFuncion.tokenActual().Tipo != token.RIGHT_PAREN {
		for {
			compiladorFuncion.funcion.Aridad++
			if compiladorFuncion.funcion.Aridad > 255 {
				fmt.Fprintf(os.Stderr, "[Línea %d] Error: No se pueden tener más de 255 parámetros.\n", compiladorFuncion.tokenActual().Linea)
			}

			compiladorFuncion.consumir(token.IDENTIFIER, "Se esperaba el nombre del parámetro.")
			paramNombre := compiladorFuncion.tokenAnterior().Lexema
			compiladorFuncion.locales = append(compiladorFuncion.locales, Local{
				nombre:      paramNombre,
				profundidad: 1,
				capturada:   false,
			})

			if !compiladorFuncion.coincide(token.COMMA) {
				break
			}
		}
	}

	compiladorFuncion.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de los parámetros.")
	compiladorFuncion.consumir(token.LEFT_BRACE, "Se esperaba '{' antes del cuerpo de la función.")
	compiladorFuncion.profundidadScope = 1
	compiladorFuncion.bloque()
	compiladorFuncion.emitirRetornoImplicito()
	c.actual = compiladorFuncion.actual
	indiceFuncion := c.funcion.Chunk.AgregarConstante(compiladorFuncion.funcion)
	c.emitirByte(uint8(OpClosure))
	c.emitirByte(uint8(indiceFuncion))
	for _, upv := range compiladorFuncion.upvalues {
		if upv.esLocal {
			c.emitirByte(1)
		} else {
			c.emitirByte(0)
		}
		c.emitirByte(upv.indice)
	}

	if esLocal {
		return
	}

	c.emitirByte(uint8(OpDefineGlobal))
	c.emitirByte(uint8(indiceGlobal))
}

func (c *Compilador) sentencia() {
	if c.coincide(token.PRINT) {
		c.sentenciaImprimir()
	} else if c.coincide(token.IF) {
		c.sentenciaIf()
	} else if c.coincide(token.WHILE) {
		c.sentenciaWhile()
	} else if c.coincide(token.FOR) {
		c.sentenciaFor()
	} else if c.coincide(token.LEFT_BRACE) {
		c.empezarScope()
		c.bloque()
		c.terminarScope()
	} else if c.coincide(token.RETURN) {
		c.sentenciaReturn()
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

func (c *Compilador) sentenciaIf() {
	c.consumir(token.LEFT_PAREN, "Se esperaba '(' después del 'if'.")
	c.expresion()
	c.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de la condición.")
	saltoIf := c.emitirSalto(uint8(OpJumpIfFalse))
	c.emitirByte(uint8(OpPop))
	c.sentencia()
	saltoElse := c.emitirSalto(uint8(OpJump))
	c.parchearSalto(saltoIf)
	c.emitirByte(uint8(OpPop))
	if c.coincide(token.ELSE) {
		c.sentencia()
	}

	c.parchearSalto(saltoElse)
}

func (c *Compilador) sentenciaWhile() {
	inicioLoop := len(c.funcion.Chunk.Instrucciones)
	c.consumir(token.LEFT_PAREN, "Se esperaba '(' después del 'while'.")
	c.expresion()
	c.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de la condición.")
	saltoSalida := c.emitirSalto(uint8(OpJumpIfFalse))
	c.emitirByte(uint8(OpPop))
	c.sentencia()
	c.emitirLoop(inicioLoop)
	c.parchearSalto(saltoSalida)
	c.emitirByte(uint8(OpPop))
}

func (c *Compilador) sentenciaFor() {
	c.empezarScope()
	c.consumir(token.LEFT_PAREN, "Se esperaba '(' después del 'for'.")
	if c.coincide(token.SEMICOLON) {
		// No hacemos nada.
	} else if c.coincide(token.VAR) {
		c.declaracionVariable()
	} else {
		c.sentenciaExpresion()
	}

	inicioLoop := len(c.funcion.Chunk.Instrucciones)
	saltoSalida := -1
	if !c.coincide(token.SEMICOLON) {
		c.expresion()
		c.consumir(token.SEMICOLON, "Se esperaba ';' después de la condición del for.")

		saltoSalida = c.emitirSalto(uint8(OpJumpIfFalse))
		c.emitirByte(uint8(OpPop))
	}

	if !c.coincide(token.RIGHT_PAREN) {
		saltoCuerpo := c.emitirSalto(uint8(OpJump))
		inicioIncremento := len(c.funcion.Chunk.Instrucciones)
		c.expresion()
		c.emitirByte(uint8(OpPop))
		c.consumir(token.RIGHT_PAREN, "Se esperaba ')' después del for.")
		c.emitirLoop(inicioLoop)
		inicioLoop = inicioIncremento
		c.parchearSalto(saltoCuerpo)
	}

	c.sentencia()
	c.emitirLoop(inicioLoop)
	if saltoSalida != -1 {
		c.parchearSalto(saltoSalida)
		c.emitirByte(uint8(OpPop))
	}

	c.terminarScope()
}

func (c *Compilador) sentenciaReturn() {
	if c.coincide(token.SEMICOLON) {
		// Un return pelado dentro de un inicializador devuelve la instancia.
		if c.esInicializador {
			c.emitirByte(uint8(OpGetLocal))
			c.emitirByte(0)
		} else {
			c.emitirByte(uint8(OpNil))
		}

		c.emitirByte(uint8(OpReturn))

		return
	}

	if c.esInicializador {
		fmt.Fprintf(os.Stderr, "[Línea %d] Error: No se puede devolver un valor desde un inicializador.\n", c.tokenActual().Linea)
	}

	c.expresion()
	c.consumir(token.SEMICOLON, "Se esperaba ';' después del valor de retorno.")
	c.emitirByte(uint8(OpReturn))
}

func llamada(c *Compilador, puedeAsignar bool) {
	cantidadArgumentos := 0
	if c.tokenActual().Tipo != token.RIGHT_PAREN {
		for {
			c.expresion()
			if cantidadArgumentos == 255 {
				fmt.Fprintf(os.Stderr, "[Línea %d] Error: No se pueden tener más de 255 argumentos.\n", c.tokenActual().Linea)
			}

			cantidadArgumentos++
			if !c.coincide(token.COMMA) {
				break
			}
		}
	}

	c.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de los argumentos.")
	c.emitirByte(uint8(OpCall))
	c.emitirByte(uint8(cantidadArgumentos))
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
		if c.locales[len(c.locales)-1].capturada {
			c.emitirByte(uint8(OpCloseUpvalue))
		} else {
			c.emitirByte(uint8(OpPop))
		}

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

// Reconoce "class Nombre { metodos }" y emite el bytecode que arma la clase.
//
// La clase queda en la pila mientras se le cuelgan los métodos, y recién al
// final se la define como variable.
func (c *Compilador) declaracionClase() {
	c.consumir(token.IDENTIFIER, "Se esperaba el nombre de la clase.")
	nombreToken := c.tokenAnterior()
	indiceNombre := c.funcion.Chunk.AgregarConstante(nombreToken.Lexema)

	// La clase se crea y se define como variable antes de los métodos, para
	// que su nombre ya se pueda leer mientras se la arma.
	c.emitirByte(uint8(OpClass))
	c.emitirByte(uint8(indiceNombre))
	c.definirClase(nombreToken)

	tieneSuperclase := c.coincide(token.LESS)
	if tieneSuperclase {
		c.consumir(token.IDENTIFIER, "Se esperaba el nombre de la superclase.")
		nombreSuper := c.tokenAnterior()

		if nombreSuper.Lexema == nombreToken.Lexema {
			fmt.Fprintf(os.Stderr, "[Línea %d] Error: Una clase no puede heredar de sí misma.\n", nombreSuper.Linea)
		}

		// La superclase queda en la pila como una local llamada "super". Los
		// métodos la capturan como upvalue, así cada "super" apunta al nivel
		// donde fue escrito y no a la clase de la instancia.
		c.empezarScope()
		c.emitirLecturaDeVariable(nombreSuper.Lexema)
		c.locales = append(c.locales, Local{nombre: "super", profundidad: c.profundidadScope})

		c.emitirLecturaDeVariable(nombreToken.Lexema)
		c.emitirByte(uint8(OpInherit))
	}

	c.consumir(token.LEFT_BRACE, "Se esperaba '{' antes del cuerpo de la clase.")

	// Los métodos se cuelgan de la clase, así que hay que traerla a la pila.
	c.emitirLecturaDeVariable(nombreToken.Lexema)
	for c.tokenActual().Tipo != token.RIGHT_BRACE && !c.esFin() {
		c.metodo()
	}

	c.consumir(token.RIGHT_BRACE, "Se esperaba '}' después del cuerpo de la clase.")
	c.emitirByte(uint8(OpPop))

	if tieneSuperclase {
		c.terminarScope()
	}
}

// Define la clase que quedó en la pila como variable global o local.
func (c *Compilador) definirClase(nombreToken token.Token) {
	if c.profundidadScope > 0 {
		c.locales = append(c.locales, Local{
			nombre:      nombreToken.Lexema,
			profundidad: c.profundidadScope,
		})

		return
	}

	indiceGlobal := c.funcion.Chunk.AgregarConstante(nombreToken.Lexema)
	c.emitirByte(uint8(OpDefineGlobal))
	c.emitirByte(uint8(indiceGlobal))
}

// Emite la lectura de una variable, resolviéndola como local, upvalue o global.
func (c *Compilador) emitirLecturaDeVariable(nombre string) {
	arg := c.resolverLocal(nombre)
	operacion := uint8(OpGetLocal)

	if arg == -1 {
		if arg = c.resolverUpvalue(nombre); arg != -1 {
			operacion = uint8(OpGetUpvalue)
		} else {
			arg = c.funcion.Chunk.AgregarConstante(nombre)
			operacion = uint8(OpGetGlobal)
		}
	}

	c.emitirByte(operacion)
	c.emitirByte(uint8(arg))
}

// Compila un método y lo cuelga de la clase que está en la pila.
func (c *Compilador) metodo() {
	c.consumir(token.IDENTIFIER, "Se esperaba el nombre del método.")
	nombreToken := c.tokenAnterior()
	indiceNombre := c.funcion.Chunk.AgregarConstante(nombreToken.Lexema)

	c.compilarCuerpoDeMetodo(nombreToken.Lexema)

	c.emitirByte(uint8(OpMethod))
	c.emitirByte(uint8(indiceNombre))
}

// Compila el cuerpo de un método en su propio chunk y emite el OpClosure.
//
// El slot cero del marco se llama "this": es donde la máquina deja el receptor
// antes de entrar, así que "this" se resuelve como cualquier otra local.
func (c *Compilador) compilarCuerpoDeMetodo(nombre string) {
	compiladorMetodo := &Compilador{
		padre:            c,
		tokens:           c.tokens,
		actual:           c.actual,
		funcion:          NuevaFuncion(),
		locales:          make([]Local, 0),
		upvalues:         make([]UpvalueLocal, 0),
		profundidadScope: 0,
		reglas:           c.reglas,
		esInicializador:  nombre == NombreDelInicializador,
	}

	compiladorMetodo.funcion.Nombre = nombre
	compiladorMetodo.locales = append(compiladorMetodo.locales, Local{
		profundidad: 0,
		nombre:      "this",
		capturada:   false,
	})

	compiladorMetodo.consumir(token.LEFT_PAREN, "Se esperaba '(' después del nombre del método.")
	if compiladorMetodo.tokenActual().Tipo != token.RIGHT_PAREN {
		for {
			compiladorMetodo.funcion.Aridad++
			compiladorMetodo.consumir(token.IDENTIFIER, "Se esperaba el nombre del parámetro.")
			compiladorMetodo.locales = append(compiladorMetodo.locales, Local{
				nombre:      compiladorMetodo.tokenAnterior().Lexema,
				profundidad: 1,
				capturada:   false,
			})

			if !compiladorMetodo.coincide(token.COMMA) {
				break
			}
		}
	}

	compiladorMetodo.consumir(token.RIGHT_PAREN, "Se esperaba ')' después de los parámetros.")
	compiladorMetodo.consumir(token.LEFT_BRACE, "Se esperaba '{' antes del cuerpo del método.")
	compiladorMetodo.profundidadScope = 1
	compiladorMetodo.bloque()
	compiladorMetodo.emitirRetornoImplicito()

	c.actual = compiladorMetodo.actual
	indiceFuncion := c.funcion.Chunk.AgregarConstante(compiladorMetodo.funcion)
	c.emitirByte(uint8(OpClosure))
	c.emitirByte(uint8(indiceFuncion))
	for _, upv := range compiladorMetodo.upvalues {
		if upv.esLocal {
			c.emitirByte(1)
		} else {
			c.emitirByte(0)
		}

		c.emitirByte(upv.indice)
	}
}

// Emite el retorno que corre cuando el cuerpo termina sin pasar por un return.
//
// Un inicializador devuelve la instancia, que vive en el slot cero; cualquier
// otro cuerpo devuelve nil.
func (c *Compilador) emitirRetornoImplicito() {
	if c.esInicializador {
		c.emitirByte(uint8(OpGetLocal))
		c.emitirByte(0)
	} else {
		c.emitirByte(uint8(OpNil))
	}

	c.emitirByte(uint8(OpReturn))
}

// Reconoce el acceso a una propiedad con punto, de lectura o de escritura.
func punto(c *Compilador, puedeAsignar bool) {
	c.consumir(token.IDENTIFIER, "Se esperaba el nombre de la propiedad después del '.'.")
	indiceNombre := c.funcion.Chunk.AgregarConstante(c.tokenAnterior().Lexema)

	if puedeAsignar && c.coincide(token.EQUAL) {
		c.expresion()
		c.emitirByte(uint8(OpSetProperty))
	} else {
		c.emitirByte(uint8(OpGetProperty))
	}

	c.emitirByte(uint8(indiceNombre))
}

// Reconoce "this", que se resuelve como la local del slot cero del método.
func este(c *Compilador, puedeAsignar bool) {
	variable(c, false)
}

// Reconoce "super.metodo", que busca el método en la superclase capturada.
func superior(c *Compilador, puedeAsignar bool) {
	c.consumir(token.DOT, "Se esperaba '.' después de 'super'.")
	c.consumir(token.IDENTIFIER, "Se esperaba el nombre del método de la superclase.")
	indiceNombre := c.funcion.Chunk.AgregarConstante(c.tokenAnterior().Lexema)

	// Primero el receptor y después la superclase: la máquina los saca en ese
	// orden para armar el método ligado.
	c.emitirLecturaDeVariable("this")
	c.emitirLecturaDeVariable("super")

	c.emitirByte(uint8(OpGetSuper))
	c.emitirByte(uint8(indiceNombre))
}
