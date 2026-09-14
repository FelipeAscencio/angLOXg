package sintaxis

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

var casosDeSentencia = []struct {
	nombre   string
	arbol    Stmt
	esperado string
}{
	{
		nombre:   "sentencia de expresión",
		arbol:    &ExpressionStmt{Expression: num(1)},
		esperado: "(expr 1)",
	},
	{
		nombre:   "impresión",
		arbol:    &Print{Value: &Literal{Value: "hola"}},
		esperado: `(print "hola")`,
	},
	{
		nombre: "declaración con valor inicial",
		arbol: &Var{
			Name:        tok(token.IDENTIFIER, "x"),
			Initializer: num(5),
		},
		esperado: "(var x 5)",
	},
	{
		nombre:   "declaración sin valor inicial",
		arbol:    &Var{Name: tok(token.IDENTIFIER, "x")},
		esperado: "(var x)",
	},
	{
		nombre:   "bloque vacío",
		arbol:    &Block{},
		esperado: "(block)",
	},
	{
		nombre: "bloque con sentencias",
		arbol: &Block{Statements: []Stmt{
			&Print{Value: num(1)},
			&Print{Value: num(2)},
		}},
		esperado: "(block (print 1) (print 2))",
	},
	{
		nombre: "condicional sin rama falsa",
		arbol: &If{
			Condition: variable("a"),
			Then:      &Print{Value: num(1)},
		},
		esperado: "(if a (print 1))",
	},
	{
		nombre: "condicional con rama falsa",
		arbol: &If{
			Condition: variable("a"),
			Then:      &Print{Value: num(1)},
			Else:      &Print{Value: num(2)},
		},
		esperado: "(if a (print 1) (print 2))",
	},
	{
		nombre: "bucle",
		arbol: &While{
			Condition: variable("a"),
			Body:      &Print{Value: num(1)},
		},
		esperado: "(while a (print 1))",
	},
	{
		nombre: "función sin parámetros",
		arbol: &Function{
			Name: tok(token.IDENTIFIER, "saludar"),
			Body: []Stmt{&Print{Value: &Literal{Value: "hola"}}},
		},
		esperado: `(fun saludar () (print "hola"))`,
	},
	{
		nombre: "función con parámetros",
		arbol: &Function{
			Name: tok(token.IDENTIFIER, "sumar"),
			Params: []token.Token{
				tok(token.IDENTIFIER, "a"),
				tok(token.IDENTIFIER, "b"),
			},
			Body: []Stmt{&Return{Value: variable("a")}},
		},
		esperado: "(fun sumar (a b) (return a))",
	},
	{
		nombre:   "retorno con valor",
		arbol:    &Return{Value: num(1)},
		esperado: "(return 1)",
	},
	{
		nombre:   "retorno sin valor",
		arbol:    &Return{},
		esperado: "(return)",
	},
}

func TestRepresentarSentencias(t *testing.T) {
	for _, caso := range casosDeSentencia {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenido := RepresentarSentencia(caso.arbol)
			if obtenido != caso.esperado {
				t.Errorf("se obtuvo %q y se esperaba %q", obtenido, caso.esperado)
			}
		})
	}
}
