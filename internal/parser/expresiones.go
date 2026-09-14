package parser

import (
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Las reglas de este archivo forman la cadena de precedencia. Cada una llama a
// la siguiente, de la precedencia más floja a la más fuerte:
//
//	assignment -> or -> and -> equality -> comparison
//	           -> term -> factor -> unary -> call -> primary
//
// El anidamiento es lo que hace que "1 + 2 * 3" agrupe la multiplicación
// primero: cuando "term" va a buscar sus operandos, "factor" ya se llevó el
// producto entero.

// assignment reconoce "x = valor".
//
// El problema es que el lado izquierdo no se sabe que era un destino hasta
// haber leído el "=", y para entonces ya se consumió. La salida es parsear
// primero una expresión común y, si aparece el "=", revisar que lo que salió
// haya sido una variable.
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

// or reconoce el operador "or".
func (p *Parser) or() sintaxis.Expr {
	return p.logicalLeft(p.and, token.OR)
}

// and reconoce el operador "and".
func (p *Parser) and() sintaxis.Expr {
	return p.logicalLeft(p.equality, token.AND)
}

// logicalLeft arma una cadena de operadores lógicos, agrupando de izquierda a
// derecha. Es igual que binaryLeft, pero produce nodos Logical.
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

// equality reconoce "==" y "!=".
func (p *Parser) equality() sintaxis.Expr {
	return p.binaryLeft(p.comparison, token.BANG_EQUAL, token.EQUAL_EQUAL)
}

// comparison reconoce ">", ">=", "<" y "<=".
func (p *Parser) comparison() sintaxis.Expr {
	return p.binaryLeft(p.term,
		token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
}

// term reconoce la suma y la resta.
func (p *Parser) term() sintaxis.Expr {
	return p.binaryLeft(p.factor, token.MINUS, token.PLUS)
}

// factor reconoce el producto, la división y el módulo.
func (p *Parser) factor() sintaxis.Expr {
	return p.binaryLeft(p.unary, token.SLASH, token.STAR, token.PERCENT)
}

// unary reconoce "!" y "-" delante de una expresión.
//
// Se llama a sí misma, y no al siguiente nivel, para que "--1" se lea como la
// negación de una negación.
func (p *Parser) unary() sintaxis.Expr {
	if p.match(token.BANG, token.MINUS) {
		operador := p.previous()

		return &sintaxis.Unary{Operator: operador, Right: p.unary()}
	}

	return p.call()
}

// call reconoce una expresión seguida de cero o más listas de argumentos
// entre paréntesis.
//
// El bucle es lo que permite encadenar: en "f()()", el destino de la segunda
// llamada es el nodo que armó la primera.
func (p *Parser) call() sintaxis.Expr {
	expression := p.primary()

	for p.match(token.LEFT_PAREN) {
		expression = p.finishCall(expression)
	}

	return expression
}

// finishCall consume los argumentos y el paréntesis de cierre. Se llama
// con el paréntesis de apertura ya consumido.
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

// binaryLeft arma una cadena de operadores binarios del mismo nivel de
// precedencia, agrupando de izquierda a derecha.
//
// El bucle es lo que da la asociatividad: cada vuelta mete el árbol acumulado
// como operando izquierdo del operador nuevo, así "1 - 2 - 3" queda como
// "(1 - 2) - 3" y no como "1 - (2 - 3)".
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
