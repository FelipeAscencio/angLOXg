package parser

import (
	"fmt"

	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Recorre la lista de tokens producida por el escáner para construir el árbol sintáctico (AST).
type Parser struct {
	tokens  []token.Token
	actual  int
	errores []error
}

// Crea una nueva instancia del parser posicionada en el primer token.
func Nuevo(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Procesa todos los tokens y retorna las sentencias del programa y los errores encontrados.
func (p *Parser) Parsear() ([]sintaxis.Stmt, []error) {
	sentencias := []sintaxis.Stmt{}

	for !p.isAtEnd() {
		if statement := p.declarationWithRecovery(); statement != nil {
			sentencias = append(sentencias, statement)
		}
	}

	return sentencias, p.errores
}

// Parsea una declaración y maneja la recuperación de errores mediante pánico y sincronización.
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

// Parsea una única expresión suelta y retorna su árbol junto con posibles errores.
func (p *Parser) ParsearExpresion() (arbol sintaxis.Expr, errores []error) {
	defer func() {
		if recuperado := recover(); recuperado != nil {
			p.recordPanic(recuperado)
			arbol, errores = nil, p.errores
		}
	}()

	arbol = p.expression()

	return arbol, p.errores
}

// =======================
// Reglas de la gramática.
// =======================

func (p *Parser) expression() sintaxis.Expr {
	return p.assignment()
}

// Reconoce elementos primarios (literales, variables, agrupaciones).
func (p *Parser) primary() sintaxis.Expr {
	switch {
	case p.match(token.FALSE):
		return &sintaxis.Literal{Value: false}

	case p.match(token.TRUE):
		return &sintaxis.Literal{Value: true}

	case p.match(token.NIL):
		return &sintaxis.Literal{Value: nil}

	case p.match(token.NUMBER, token.STRING):
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

// ==================
// Manejo de errores.
// ==================

func (p *Parser) fail(donde token.Token, mensaje string) {
	panic(&Error{Linea: donde.Linea, Mensaje: mensaje})
}

func (p *Parser) recordPanic(recuperado any) {
	falla, esDeSintaxis := recuperado.(*Error)
	if !esDeSintaxis {
		panic(recuperado)
	}

	p.errores = append(p.errores, falla)
}

// ===================================
// Auxiliares de lectura y navegación.
// ===================================

func (p *Parser) isAtEnd() bool {
	return p.peek().Tipo == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.actual]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.actual-1]
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.actual++
	}

	return p.previous()
}

func (p *Parser) check(tipos ...token.TipoDeToken) bool {
	if p.isAtEnd() && !contains(tipos, token.EOF) {
		return false
	}

	return contains(tipos, p.peek().Tipo)
}

func (p *Parser) match(tipos ...token.TipoDeToken) bool {
	if !p.check(tipos...) {
		return false
	}

	p.advance()

	return true
}

func (p *Parser) consume(tipo token.TipoDeToken, mensaje string) token.Token {
	if p.check(tipo) {
		return p.advance()
	}

	p.fail(p.peek(), mensaje)

	panic(fmt.Sprintf("consume siguió después de fail en %v", p.peek()))
}

func contains(tipos []token.TipoDeToken, buscado token.TipoDeToken) bool {
	for _, tipo := range tipos {
		if tipo == buscado {
			return true
		}
	}

	return false
}
