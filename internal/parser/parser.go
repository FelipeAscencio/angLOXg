package parser

import (
	"fmt"

	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Parser recorre la lista de tokens que produjo el escáner y arma con ella
// un árbol de sintaxis.
//
// El recorrido es un descenso recursivo: hay una función por cada regla de la
// gramática, y las reglas se llaman entre sí siguiendo la precedencia de los
// operadores, de la más floja a la más fuerte.
type Parser struct {
	tokens  []token.Token
	actual  int     // Token en el que está parado el cursor.
	errores []error // Fallas de sintaxis acumuladas.
}

// Nuevo construye un parser posicionado en el primer token.
func Nuevo(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parsear recorre todos los tokens y devuelve las sentencias del programa
// junto con las fallas de sintaxis encontradas.
//
// Una declaración que falla no corta el parseo: se descarta y se sigue con la
// siguiente, para poder reportar varios errores en una sola pasada. Si la lista
// de errores no viene vacía, el árbol está incompleto.
func (p *Parser) Parsear() ([]sintaxis.Stmt, []error) {
	sentencias := []sintaxis.Stmt{}

	for !p.isAtEnd() {
		if statement := p.declarationWithRecovery(); statement != nil {
			sentencias = append(sentencias, statement)
		}
	}

	return sentencias, p.errores
}

// declarationWithRecovery parsea una declaración y, si falla, anota el
// error y deja el cursor en un punto desde donde se pueda seguir.
func (p *Parser) declarationWithRecovery() (statement sintaxis.Stmt) {
	defer func() {
		if recuperado := recover(); recuperado != nil {
			p.recordPanic(recuperado)
			p.synchronize()
			statement = nil
		}
	}()

	return p.declaration()
}

// ParsearExpresion parsea una única expresión y devuelve su árbol junto con
// las fallas de sintaxis encontradas.
//
// Si el parseo falla, el árbol devuelto es nil.
func (p *Parser) ParsearExpresion() (arbol sintaxis.Expr, errores []error) {
	// Un error de sintaxis viaja como pánico hasta acá, que es el borde del
	// paquete. Nunca escapa hacia afuera.
	defer func() {
		if recuperado := recover(); recuperado != nil {
			p.recordPanic(recuperado)
			arbol, errores = nil, p.errores
		}
	}()

	arbol = p.expression()

	return arbol, p.errores
}

// ---------- Reglas de la gramática ---------- //

// expression es la regla de arranque de toda expresión.
func (p *Parser) expression() sintaxis.Expr {
	return p.assignment()
}

// primary reconoce lo que no se descompone en nada más chico: los literales,
// las variables y las expresiones entre paréntesis.
func (p *Parser) primary() sintaxis.Expr {
	switch {
	case p.match(token.FALSE):
		return &sintaxis.Literal{Value: false}

	case p.match(token.TRUE):
		return &sintaxis.Literal{Value: true}

	case p.match(token.NIL):
		return &sintaxis.Literal{Value: nil}

	case p.match(token.NUMBER, token.STRING):
		// El escáner ya dejó el valor resuelto en el token.
		return &sintaxis.Literal{Value: p.previous().Literal}

	case p.match(token.IDENTIFIER):
		return &sintaxis.Variable{Name: p.previous()}

	case p.match(token.LEFT_PAREN):
		interior := p.expression()
		p.consume(token.RIGHT_PAREN, "se esperaba ')' para cerrar la agrupación")

		return &sintaxis.Grouping{Expression: interior}

	default:
		p.fail(p.peek(), "se esperaba una expresión")

		return nil
	}
}

// ---------- Manejo de errores ---------- //

// fail corta el parseo con una falla de sintaxis en el token indicado.
//
// El pánico lo atrapa el borde del paquete; sirve para volver de un descenso
// recursivo de varios niveles sin arrastrar un "error" por cada regla.
func (p *Parser) fail(donde token.Token, mensaje string) {
	panic(&Error{Linea: donde.Linea, Mensaje: mensaje})
}

// recordPanic anota la falla que venía viajando como pánico. Si no es una
// falla de sintaxis, se vuelve a lanzar: es un error de programación nuestro.
func (p *Parser) recordPanic(recuperado any) {
	falla, esDeSintaxis := recuperado.(*Error)
	if !esDeSintaxis {
		panic(recuperado)
	}

	p.errores = append(p.errores, falla)
}

// ---------- Auxiliares ---------- //

// isAtEnd indica si el cursor llegó al token de fin de archivo.
func (p *Parser) isAtEnd() bool {
	return p.peek().Tipo == token.EOF
}

// peek devuelve el token en el que está parado el cursor, sin consumirlo.
func (p *Parser) peek() token.Token {
	return p.tokens[p.actual]
}

// previous devuelve el último token consumido.
func (p *Parser) previous() token.Token {
	return p.tokens[p.actual-1]
}

// advance consume el token actual y lo devuelve.
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.actual++
	}

	return p.previous()
}

// check indica si el token actual es de alguno de los tipos dados, sin
// consumirlo.
func (p *Parser) check(tipos ...token.TipoDeToken) bool {
	if p.isAtEnd() && !contains(tipos, token.EOF) {
		return false
	}

	return contains(tipos, p.peek().Tipo)
}

// match consume el token actual si es de alguno de los tipos dados, e
// informa si lo hizo.
func (p *Parser) match(tipos ...token.TipoDeToken) bool {
	if !p.check(tipos...) {
		return false
	}

	p.advance()

	return true
}

// consume exige que el token actual sea del tipo dado y lo consume. Si no lo
// es, corta el parseo con el mensaje indicado.
func (p *Parser) consume(tipo token.TipoDeToken, mensaje string) token.Token {
	if p.check(tipo) {
		return p.advance()
	}

	p.fail(p.peek(), mensaje)

	// Inalcanzable: "fail" siempre entra en pánico.
	panic(fmt.Sprintf("consume siguió después de fail en %v", p.peek()))
}

// contains indica si un tipo está en una lista de tipos.
func contains(tipos []token.TipoDeToken, buscado token.TipoDeToken) bool {
	for _, tipo := range tipos {
		if tipo == buscado {
			return true
		}
	}

	return false
}
