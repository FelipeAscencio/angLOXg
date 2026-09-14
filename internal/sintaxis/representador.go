package sintaxis

import (
	"fmt"
	"strconv"
	"strings"
)

// Representar arma una representación textual del árbol de una expresión, con
// la forma de una S-expresión: cada nodo compuesto queda entre paréntesis, con
// el operador adelante y los operandos atrás.
//
//	1 + 2 * 3   ->   (+ 1 (* 2 3))
//
// Sirve para ver de un vistazo cómo quedó agrupado el árbol, que es
// exactamente lo que no se puede deducir mirando la lista de tokens.
func Representar(expresion Expr) string {
	switch n := expresion.(type) {
	case *Literal:
		return representarValor(n.Value)

	case *Grouping:
		return encerrar("group", Representar(n.Expression))

	case *Unary:
		return encerrar(n.Operator.Lexema, Representar(n.Right))

	case *Binary:
		return encerrar(n.Operator.Lexema, Representar(n.Left), Representar(n.Right))

	case *Logical:
		return encerrar(n.Operator.Lexema, Representar(n.Left), Representar(n.Right))

	case *Variable:
		return n.Name.Lexema

	case *Assign:
		return encerrar("=", n.Name.Lexema, Representar(n.Value))

	case *Call:
		partes := []string{Representar(n.Callee)}
		for _, argumento := range n.Arguments {
			partes = append(partes, Representar(argumento))
		}

		return encerrar("call", partes...)

	default:
		// Si aparece un nodo nuevo y nadie lo contempló acá, conviene que se
		// note enseguida y no que se represente como vacío.
		panic(fmt.Sprintf("expresión no contemplada por el representador: %T", expresion))
	}
}

// RepresentarSentencia arma la representación textual del árbol de una
// sentencia, con el mismo formato de S-expresiones que usa Representar.
func RepresentarSentencia(sentencia Stmt) string {
	switch n := sentencia.(type) {
	case *ExpressionStmt:
		return encerrar("expr", Representar(n.Expression))

	case *Print:
		return encerrar("print", Representar(n.Value))

	case *Var:
		if n.Initializer == nil {
			return encerrar("var", n.Name.Lexema)
		}

		return encerrar("var", n.Name.Lexema, Representar(n.Initializer))

	case *Block:
		return encerrar("block", representarSentencias(n.Statements)...)

	case *If:
		partes := []string{Representar(n.Condition), RepresentarSentencia(n.Then)}
		if n.Else != nil {
			partes = append(partes, RepresentarSentencia(n.Else))
		}

		return encerrar("if", partes...)

	case *While:
		return encerrar("while", Representar(n.Condition), RepresentarSentencia(n.Body))

	case *Function:
		nombres := make([]string, 0, len(n.Params))
		for _, parametro := range n.Params {
			nombres = append(nombres, parametro.Lexema)
		}

		partes := []string{n.Name.Lexema, "(" + strings.Join(nombres, " ") + ")"}
		partes = append(partes, representarSentencias(n.Body)...)

		return encerrar("fun", partes...)

	case *Return:
		if n.Value == nil {
			return encerrar("return")
		}

		return encerrar("return", Representar(n.Value))

	default:
		panic(fmt.Sprintf("sentencia no contemplada por el representador: %T", sentencia))
	}
}

// representarValor escribe el valor de un literal.
func representarValor(valor any) string {
	switch v := valor.(type) {
	case nil:
		return "nil"
	case string:
		// Las comillas distinguen una cadena de un identificador.
		return `"` + v + `"`
	case float64:
		// El formato 'g' con precisión -1 escribe la menor cantidad de dígitos
		// que permita recuperar el mismo número: 1 queda "1" y no "1.000000".
		return strconv.FormatFloat(v, 'g', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// representarSentencias representa una lista de sentencias.
func representarSentencias(sentencias []Stmt) []string {
	partes := make([]string, 0, len(sentencias))
	for _, sentencia := range sentencias {
		partes = append(partes, RepresentarSentencia(sentencia))
	}

	return partes
}

// encerrar arma "(nombre parte parte ...)".
func encerrar(nombre string, partes ...string) string {
	if len(partes) == 0 {
		return "(" + nombre + ")"
	}

	return "(" + nombre + " " + strings.Join(partes, " ") + ")"
}
