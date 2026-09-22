package sintaxis

import "github.com/FelipeAscencio/angLOXg/internal/token"

// Interfaz marcadora para los nodos de sentencia del árbol de sintaxis.
type Stmt interface {
	isStmt()
}

// Representa una expresión evaluada como sentencia por sus efectos secundarios.
type ExpressionStmt struct {
	Expression Expr
}

// Representa una sentencia de impresión ("print").
type Print struct {
	Value Expr
}

// Representa la declaración de una nueva variable ("var").
type Var struct {
	Name        token.Token
	Initializer Expr
}

// Representa un bloque de sentencias delimitado por llaves que define su propio ámbito.
type Block struct {
	Statements []Stmt
}

// Representa una estructura de control condicional ("if/else").
type If struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

// Representa un bucle iterativo condicional ("while").
type While struct {
	Condition Expr
	Body      Stmt
}

// Representa la declaración de una función nombrada.
type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

// Representa una sentencia de retorno ("return").
type Return struct {
	Keyword token.Token
	Value   Expr
}

func (*ExpressionStmt) isStmt() {}
func (*Print) isStmt()          {}
func (*Var) isStmt()            {}
func (*Block) isStmt()          {}
func (*If) isStmt()             {}
func (*While) isStmt()          {}
func (*Function) isStmt()       {}
func (*Return) isStmt()         {}
