package parser

import (
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Reconoce declaraciones permitidas en cuerpos de programas o bloques.
func (p *Parser) declaration() sintaxis.Stmt {
	switch {
	case p.match(token.VAR):
		return p.varDeclaration()

	case p.match(token.FUN):
		return p.funDeclaration()

	default:
		return p.statement()
	}
}

// Reconoce declaraciones de variables ("var x;" o "var x = expresion;").
func (p *Parser) varDeclaration() sintaxis.Stmt {
	nombre := p.consume(token.IDENTIFIER, "se esperaba el nombre de la variable")

	var inicializador sintaxis.Expr
	if p.match(token.EQUAL) {
		inicializador = p.expression()
	}

	p.consume(token.SEMICOLON, "se esperaba ';' después de la declaración")

	return &sintaxis.Var{Name: nombre, Initializer: inicializador}
}

// Reconoce y arma sentencias de acuerdo con su tipo.
func (p *Parser) statement() sintaxis.Stmt {
	switch {
	case p.match(token.PRINT):
		return p.printStatement()

	case p.match(token.IF):
		return p.ifStatement()

	case p.match(token.WHILE):
		return p.whileStatement()

	case p.match(token.FOR):
		return p.forStatement()

	case p.match(token.RETURN):
		return p.returnStatement()

	case p.match(token.LEFT_BRACE):
		return &sintaxis.Block{Statements: p.block()}

	default:
		return p.expressionStatement()
	}
}

// Reconoce sentencias de impresión ("print expresion;").
func (p *Parser) printStatement() sintaxis.Stmt {
	valor := p.expression()
	p.consume(token.SEMICOLON, "se esperaba ';' después del valor a imprimir")

	return &sintaxis.Print{Value: valor}
}

// Reconoce expresiones sueltas terminadas en punto y coma.
func (p *Parser) expressionStatement() sintaxis.Stmt {
	expression := p.expression()
	p.consume(token.SEMICOLON, "se esperaba ';' después de la expresión")

	return &sintaxis.ExpressionStmt{Expression: expression}
}

// Reconoce declaraciones dentro de un bloque delimitado por llaves.
func (p *Parser) block() []sintaxis.Stmt {
	sentencias := []sintaxis.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		sentencias = append(sentencias, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "se esperaba '}' para cerrar el block")

	return sentencias
}

// Reconoce estructuras condicionales ("if").
func (p *Parser) ifStatement() sintaxis.Stmt {
	p.consume(token.LEFT_PAREN, "se esperaba '(' después de 'if'")
	condicion := p.expression()
	p.consume(token.RIGHT_PAREN, "se esperaba ')' después de la condición")

	entonces := p.statement()

	var sino sintaxis.Stmt
	if p.match(token.ELSE) {
		sino = p.statement()
	}

	return &sintaxis.If{Condition: condicion, Then: entonces, Else: sino}
}

// Reconoce bucles iterativos ("while").
func (p *Parser) whileStatement() sintaxis.Stmt {
	p.consume(token.LEFT_PAREN, "se esperaba '(' después de 'while'")
	condicion := p.expression()
	p.consume(token.RIGHT_PAREN, "se esperaba ')' después de la condición")

	return &sintaxis.While{Condition: condicion, Body: p.statement()}
}

// Reconoce bucles for y los traduce internamente a construcciones while.
func (p *Parser) forStatement() sintaxis.Stmt {
	p.consume(token.LEFT_PAREN, "se esperaba '(' después de 'for'")

	inicializador := p.forInitializer()

	var condicion sintaxis.Expr
	if !p.check(token.SEMICOLON) {
		condicion = p.expression()
	}

	p.consume(token.SEMICOLON, "se esperaba ';' después de la condición")

	var incremento sintaxis.Expr
	if !p.check(token.RIGHT_PAREN) {
		incremento = p.expression()
	}

	p.consume(token.RIGHT_PAREN, "se esperaba ')' para cerrar el 'for'")

	cuerpo := p.statement()

	if incremento != nil {
		cuerpo = &sintaxis.Block{Statements: []sintaxis.Stmt{
			cuerpo,
			&sintaxis.ExpressionStmt{Expression: incremento},
		}}
	}

	if condicion == nil {
		condicion = &sintaxis.Literal{Value: true}
	}

	var bucle sintaxis.Stmt = &sintaxis.While{Condition: condicion, Body: cuerpo}

	if inicializador != nil {
		bucle = &sintaxis.Block{Statements: []sintaxis.Stmt{inicializador, bucle}}
	}

	return bucle
}

// Reconoce el inicializador de un bucle for (declaración, expresión o vacío).
func (p *Parser) forInitializer() sintaxis.Stmt {
	switch {
	case p.match(token.SEMICOLON):
		return nil

	case p.match(token.VAR):
		return p.varDeclaration()

	default:
		return p.expressionStatement()
	}
}

// Reconoce declaraciones de funciones ("fun nombre(...) { ... }").
func (p *Parser) funDeclaration() sintaxis.Stmt {
	nombre := p.consume(token.IDENTIFIER, "se esperaba el nombre de la función")
	p.consume(token.LEFT_PAREN, "se esperaba '(' después del nombre de la función")

	parameters := p.parameters()

	p.consume(token.LEFT_BRACE, "se esperaba '{' para abrir el cuerpo de la función")

	return &sintaxis.Function{
		Name:   nombre,
		Params: parameters,
		Body:   p.block(),
	}
}

// Reconoce la lista de parámetros formales de una función.
func (p *Parser) parameters() []token.Token {
	nombres := []token.Token{}

	if !p.check(token.RIGHT_PAREN) {
		for {
			nombres = append(nombres, p.consume(token.IDENTIFIER, "se esperaba el nombre del parámetro"))

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	p.consume(token.RIGHT_PAREN, "se esperaba ')' después de los parámetros")

	return nombres
}

// Reconoce sentencias de retorno ("return").
func (p *Parser) returnStatement() sintaxis.Stmt {
	palabraClave := p.previous()

	var valor sintaxis.Expr
	if !p.check(token.SEMICOLON) {
		valor = p.expression()
	}

	p.consume(token.SEMICOLON, "se esperaba ';' después del retorno")

	return &sintaxis.Return{Keyword: palabraClave, Value: valor}
}
