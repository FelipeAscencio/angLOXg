package sintaxis

import "github.com/FelipeAscencio/angLOXg/internal/token"

// Stmt es un nodo del árbol que se ejecuta pero no produce un valor.
//
// Igual que Expr, la interfaz es marcadora: el método privado impide que un
// tipo de afuera del paquete se cuele como nodo del árbol.
type Stmt interface {
	isStmt()
}

// ExpressionStmt es una expresión usada como sentencia, por su efecto y no por
// su valor. Es la forma de una llamada suelta, como "f();".
type ExpressionStmt struct {
	Expression Expr
}

// Print muestra el valor de una expresión.
type Print struct {
	Value Expr
}

// Var crea una variable nueva.
//
// "Initializer" queda en nil cuando se declara sin valor, como en "var x;".
type Var struct {
	Name        token.Token
	Initializer Expr
}

// Block es una serie de sentencias entre llaves, que además abre un ámbito
// propio para las variables que se declaren adentro.
type Block struct {
	Statements []Stmt
}

// If ejecuta una rama u otra según una condición.
//
// "Else" queda en nil cuando no hay rama falsa.
type If struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

// While repite un cuerpo mientras la condición se cumpla.
//
// No hay un nodo para "for": el parser lo traduce a un While envuelto en un
// Block, porque hacen exactamente lo mismo.
type While struct {
	Condition Expr
	Body      Stmt
}

// Function crea una función con nombre.
type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

// Return corta la ejecución de una función y devuelve un valor.
//
// "Value" queda en nil cuando se escribe "return;" pelado.
type Return struct {
	// Se guarda el token "return" para poder ubicar en su línea un retorno mal
	// puesto, por ejemplo fuera de una función.
	Keyword token.Token

	Value Expr
}

func (*ExpressionStmt) isStmt() {}
func (*Print) isStmt()          {}
func (*Var) isStmt()            {}
func (*Block) isStmt()          {}
func (*If) isStmt()             {}
func (*While) isStmt()          {}
func (*Function) isStmt()       {}
func (*Return) isStmt()         {}
