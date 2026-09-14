package sintaxis

import "github.com/FelipeAscencio/angLOXg/internal/token"

// Expr es un nodo del árbol que produce un valor.
//
// La interfaz es marcadora: su único método es privado, así que ningún tipo de
// afuera de este paquete puede hacerse pasar por un nodo del árbol. Quien
// recorre el árbol distingue los nodos con un "switch" de tipo.
type Expr interface {
	isExpr()
}

// Literal es un valor escrito directamente en el código: un número, una cadena,
// "true", "false" o "nil".
//
// "Value" guarda el dato ya resuelto, con los mismos tipos dinámicos que usa el
// escáner: float64, string, bool o nil.
type Literal struct {
	Value any
}

// Grouping es una expresión encerrada entre paréntesis.
//
// Se guarda como nodo propio, en vez de devolver directamente lo de adentro,
// para no perder la forma original del código.
type Grouping struct {
	Expression Expr
}

// Unary es un operador de un solo operando: "-" o "!".
type Unary struct {
	Operator token.Token
	Right    Expr
}

// Binary es un operador de dos operandos: aritméticos, de comparación o de
// igualdad.
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// Logical es un "and" o un "or".
//
// Va separado de Binary porque no evalúa siempre los dos lados: corta apenas el
// izquierdo alcanza para saber el resultado.
type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// Variable es la lectura de una variable por su nombre.
type Variable struct {
	Name token.Token
}

// Assign le da un valor nuevo a una variable ya declarada.
type Assign struct {
	Name  token.Token
	Value Expr
}

// Call invoca al resultado de "Callee" con una lista de argumentos.
//
// "Callee" es una expresión y no un nombre porque lo que se llama puede ser el
// resultado de otra llamada, como en "f()()".
type Call struct {
	Callee Expr

	// Paréntesis de cierre, que se guarda para poder ubicar los errores de
	// ejecución de la llamada en la línea correcta.
	Paren token.Token

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
