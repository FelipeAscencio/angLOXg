package parser

import (
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// declaration reconoce todo lo que puede aparecer en el cuerpo de un programa o
// de un bloque.
//
// Se separa de "statement" porque no todos los lugares admiten una declaración:
// "if (a) var x = 1;" no tiene sentido, ya que la variable moriría al instante.
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

// varDeclaration reconoce "var x;" y "var x = expresion;".
func (p *Parser) varDeclaration() sintaxis.Stmt {
	nombre := p.consume(token.IDENTIFIER, "se esperaba el nombre de la variable")

	// El valor inicial es opcional: sin él, la variable arranca en nil.
	var inicializador sintaxis.Expr
	if p.match(token.EQUAL) {
		inicializador = p.expression()
	}

	p.consume(token.SEMICOLON, "se esperaba ';' después de la declaración")

	return &sintaxis.Var{Name: nombre, Initializer: inicializador}
}

// statement reconoce cuál de las formas de sentencia viene y la arma.
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

// printStatement reconoce "print expresion;".
func (p *Parser) printStatement() sintaxis.Stmt {
	valor := p.expression()
	p.consume(token.SEMICOLON, "se esperaba ';' después del valor a imprimir")

	return &sintaxis.Print{Value: valor}
}

// expressionStatement reconoce una expresión suelta terminada en punto y coma.
func (p *Parser) expressionStatement() sintaxis.Stmt {
	expression := p.expression()
	p.consume(token.SEMICOLON, "se esperaba ';' después de la expresión")

	return &sintaxis.ExpressionStmt{Expression: expression}
}

// block reconoce las declaraciones de adentro de un par de llaves. Se llama
// con la llave de apertura ya consumida.
func (p *Parser) block() []sintaxis.Stmt {
	sentencias := []sintaxis.Stmt{}

	// El chequeo de fin de archivo evita quedarse dando vueltas para siempre
	// cuando la llave de cierre nunca llega.
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		sentencias = append(sentencias, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "se esperaba '}' para cerrar el block")

	return sentencias
}

// ifStatement reconoce "if (condicion) sentencia" con un "else" opcional.
func (p *Parser) ifStatement() sintaxis.Stmt {
	p.consume(token.LEFT_PAREN, "se esperaba '(' después de 'if'")
	condicion := p.expression()
	p.consume(token.RIGHT_PAREN, "se esperaba ')' después de la condición")

	entonces := p.statement()

	// El "else" se engancha acá mismo, en la llamada más profunda que esté
	// abierta. Por eso en "if (a) if (b) x; else y;" el else queda con el if
	// interno, que es la lectura habitual en los lenguajes con llaves.
	var sino sintaxis.Stmt
	if p.match(token.ELSE) {
		sino = p.statement()
	}

	return &sintaxis.If{Condition: condicion, Then: entonces, Else: sino}
}

// whileStatement reconoce "while (condicion) sentencia".
func (p *Parser) whileStatement() sintaxis.Stmt {
	p.consume(token.LEFT_PAREN, "se esperaba '(' después de 'while'")
	condicion := p.expression()
	p.consume(token.RIGHT_PAREN, "se esperaba ')' después de la condición")

	return &sintaxis.While{Condition: condicion, Body: p.statement()}
}

// forStatement reconoce "for (inicializador; condicion; incremento) sentencia"
// y lo traduce a un bucle mientras.
//
// No hay un nodo propio para el for: un for hace exactamente lo mismo que un
// mientras con el inicializador delante y el incremento al final del cuerpo, y
// armarlo así deja un nodo menos que atender en cada recorrido del árbol.
func (p *Parser) forStatement() sintaxis.Stmt {
	p.consume(token.LEFT_PAREN, "se esperaba '(' después de 'for'")

	// Las tres partes son opcionales.
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

	// El incremento corre al final de cada vuelta.
	if incremento != nil {
		cuerpo = &sintaxis.Block{Statements: []sintaxis.Stmt{
			cuerpo,
			&sintaxis.ExpressionStmt{Expression: incremento},
		}}
	}

	// Sin condición, el bucle es infinito.
	if condicion == nil {
		condicion = &sintaxis.Literal{Value: true}
	}

	var bucle sintaxis.Stmt = &sintaxis.While{Condition: condicion, Body: cuerpo}

	// El inicializador corre una sola vez, antes del bucle. El bloque que lo
	// envuelve le da su propio ámbito, para que la variable del for no se
	// escape al código de alrededor.
	if inicializador != nil {
		bucle = &sintaxis.Block{Statements: []sintaxis.Stmt{inicializador, bucle}}
	}

	return bucle
}

// forInitializer reconoce la primera de las tres partes de un for, que
// puede ser una declaración, una expresión o nada.
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

// funDeclaration reconoce "fun nombre(parametros) { cuerpo }".
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

// parameters reconoce la lista de nombres entre paréntesis. Se llama con el
// paréntesis de apertura ya consumido y consume el de cierre.
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

// returnStatement reconoce "return;" y "return expresion;".
func (p *Parser) returnStatement() sintaxis.Stmt {
	palabraClave := p.previous()

	// El valor es opcional: un "return" pelado devuelve nil.
	var valor sintaxis.Expr
	if !p.check(token.SEMICOLON) {
		valor = p.expression()
	}

	p.consume(token.SEMICOLON, "se esperaba ';' después del retorno")

	return &sintaxis.Return{Keyword: palabraClave, Value: valor}
}
