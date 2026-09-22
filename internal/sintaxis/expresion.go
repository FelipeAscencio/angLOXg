package sintaxis

import "github.com/FelipeAscencio/angLOXg/internal/token"

// Interfaz marcadora para los nodos de expresión del árbol de sintaxis.
type Expr interface {
	isExpr()
}

// Representa un valor literal directo (número, cadena, booleano o nil).
type Literal struct {
	Value any
}

// Representa una expresión encerrada entre paréntesis.
type Grouping struct {
	Expression Expr
}

// Representa un operador unario ("-" o "!").
type Unary struct {
	Operator token.Token
	Right    Expr
}

// Representa un operador binario de dos operandos.
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// Representa una expresión lógica ("and" o "or") con evaluación en cortocircuito.
type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// Representa la lectura de una variable por su nombre.
type Variable struct {
	Name token.Token
}

// Representa una asignación de valor a una variable existente.
type Assign struct {
	Name  token.Token
	Value Expr
}

// Representa una llamada o invocación de una función con argumentos.
type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (*Literal) isExpr()  {}
func (*Grouping) isExpr() {}
func (*Unary) isExpr()    {}
func (*Binary) isExpr()   {}
func (*Logical) isExpr()  {}
func (*Variable) isExpr() {}
func (*Assign) isExpr()   {}
func (*Call) isExpr()     {}
