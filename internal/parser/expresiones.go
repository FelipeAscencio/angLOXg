package parser

import (
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Las reglas de este archivo forman la cadena de precedencia. Cada una llama a
// la siguiente, de la precedencia más floja a la más fuerte:
//
//	assignment -> or -> and -> equality -> comparison
//	            -> term -> factor -> unary -> call -> primary
//
// El anidamiento es lo que hace que "1 + 2 * 3" agrupe la multiplicación
// primero: cuando "term" va a buscar sus operandos, "factor" ya se llevó el
// producto entero.

// Reconoce expresiones de asignación ("x = valor").
func (p *Parser) assignment() sintaxis.Expr {
	destino := p.or()

	if !p.match(token.EQUAL) {
		return destino
	}

	igual := p.previous()

	// La llamada recursiva es lo que hace que "a = b = 1" se lea como
	// "a = (b = 1)", o sea que asocie a la derecha.
	valor := p.assignment()

	if variable, esVariable := destino.(*sintaxis.Variable); esVariable {
		return &sintaxis.Assign{Name: variable.Name, Value: valor}
	}

	p.fail(igual, "destino de asignación inválido")

	return nil
}

// Reconoce el operador lógico "or".
func (p *Parser) or() sintaxis.Expr {
	return p.logicalLeft(p.and, token.OR)
}

// Reconoce el operador lógico "and".
func (p *Parser) and() sintaxis.Expr {
	return p.logicalLeft(p.equality, token.AND)
}

// Arma una cadena de operadores lógicos, agrupando de izquierda a derecha.
func (p *Parser) logicalLeft(siguiente func() sintaxis.Expr, operadores ...token.TipoDeToken) sintaxis.Expr {
	expression := siguiente()

	for p.match(operadores...) {
		operador := p.previous()
		derecha := siguiente()
		expression = &sintaxis.Logical{
			Left:     expression,
			Operator: operador,
			Right:    derecha,
		}
	}

	return expression
}

// Reconoce operadores de igualdad ("==" y "!=").
func (p *Parser) equality() sintaxis.Expr {
	return p.binaryLeft(p.comparison, token.BANG_EQUAL, token.EQUAL_EQUAL)
}

// Reconoce operadores de comparación (">", ">=", "<" y "<=").
func (p *Parser) comparison() sintaxis.Expr {
	return p.binaryLeft(p.term,
		token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
}

// Reconoce operaciones de suma y resta.
func (p *Parser) term() sintaxis.Expr {
	return p.binaryLeft(p.factor, token.MINUS, token.PLUS)
}

// Reconoce operaciones de multiplicación, división y módulo.
func (p *Parser) factor() sintaxis.Expr {
	return p.binaryLeft(p.unary, token.SLASH, token.STAR, token.PERCENT)
}

// Reconoce operadores unarios ("!" y "-") delante de una expresión.
func (p *Parser) unary() sintaxis.Expr {
	if p.match(token.BANG, token.MINUS) {
		operador := p.previous()

		return &sintaxis.Unary{Operator: operador, Right: p.unary()}
	}

	return p.call()
}

// Reconoce llamadas a funciones y encadenamiento de invocaciones.
func (p *Parser) call() sintaxis.Expr {
	expression := p.primary()

	for p.match(token.LEFT_PAREN) {
		expression = p.finishCall(expression)
	}

	return expression
}

// Consume los argumentos y el paréntesis de cierre de una llamada a función.
func (p *Parser) finishCall(destino sintaxis.Expr) sintaxis.Expr {
	argumentos := []sintaxis.Expr{}

	if !p.check(token.RIGHT_PAREN) {
		for {
			argumentos = append(argumentos, p.expression())

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	parentesis := p.consume(token.RIGHT_PAREN, "se esperaba ')' después de los argumentos")

	return &sintaxis.Call{
		Callee:    destino,
		Paren:     parentesis,
		Arguments: argumentos,
	}
}

// Arma una cadena de operadores binarios con asociatividad a la izquierda.
func (p *Parser) binaryLeft(siguiente func() sintaxis.Expr, operadores ...token.TipoDeToken) sintaxis.Expr {
	expression := siguiente()

	for p.match(operadores...) {
		operador := p.previous()
		derecha := siguiente()
		expression = &sintaxis.Binary{
			Left:     expression,
			Operator: operador,
			Right:    derecha,
		}
	}

	return expression
}
