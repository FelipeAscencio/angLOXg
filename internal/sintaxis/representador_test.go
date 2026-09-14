package sintaxis

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// tok arma un token mínimo, con la información que hace falta para representar
// el árbol y nada más.
func tok(tipo token.TipoDeToken, lexema string) token.Token {
	return token.Nuevo(tipo, lexema, nil, 1)
}

// num arma un literal numérico, que es el nodo más repetido en las tablas.
func num(valor float64) Expr {
	return &Literal{Value: valor}
}

// variable arma una referencia a una variable por su nombre.
func variable(nombre string) Expr {
	return &Variable{Name: tok(token.IDENTIFIER, nombre)}
}

var casosDeExpresion = []struct {
	nombre   string
	arbol    Expr
	esperado string
}{
	{
		nombre:   "literal entero",
		arbol:    num(1),
		esperado: "1",
	},
	{
		nombre:   "literal decimal",
		arbol:    num(1.5),
		esperado: "1.5",
	},
	{
		nombre: "literal de cadena",
		arbol:  &Literal{Value: "hola"},
		// La cadena se representa entre comillas, para distinguirla de un
		// identificador.
		esperado: `"hola"`,
	},
	{
		nombre:   "literal verdadero",
		arbol:    &Literal{Value: true},
		esperado: "true",
	},
	{
		nombre:   "literal nulo",
		arbol:    &Literal{Value: nil},
		esperado: "nil",
	},
	{
		nombre:   "variable",
		arbol:    variable("contador"),
		esperado: "contador",
	},
	{
		nombre: "negación aritmética",
		arbol: &Unary{
			Operator: tok(token.MINUS, "-"),
			Right:    num(1),
		},
		esperado: "(- 1)",
	},
	{
		nombre: "negación lógica",
		arbol: &Unary{
			Operator: tok(token.BANG, "!"),
			Right:    &Literal{Value: true},
		},
		esperado: "(! true)",
	},
	{
		nombre: "suma",
		arbol: &Binary{
			Left:     num(1),
			Operator: tok(token.PLUS, "+"),
			Right:    num(2),
		},
		esperado: "(+ 1 2)",
	},
	{
		nombre: "módulo",
		arbol: &Binary{
			Left:     num(10),
			Operator: tok(token.PERCENT, "%"),
			Right:    num(3),
		},
		esperado: "(% 10 3)",
	},
	{
		nombre: "agrupación",
		arbol: &Grouping{
			Expression: &Binary{
				Left:     num(1),
				Operator: tok(token.PLUS, "+"),
				Right:    num(2),
			},
		},
		esperado: "(group (+ 1 2))",
	},
	{
		nombre: "árbol anidado",
		arbol: &Binary{
			Left: &Unary{
				Operator: tok(token.MINUS, "-"),
				Right:    num(123),
			},
			Operator: tok(token.STAR, "*"),
			Right:    &Grouping{Expression: num(45.67)},
		},
		esperado: "(* (- 123) (group 45.67))",
	},
	{
		nombre: "disyunción lógica",
		arbol: &Logical{
			Left:     variable("a"),
			Operator: tok(token.OR, "or"),
			Right:    variable("b"),
		},
		esperado: "(or a b)",
	},
	{
		nombre: "asignación",
		arbol: &Assign{
			Name:  tok(token.IDENTIFIER, "x"),
			Value: num(1),
		},
		esperado: "(= x 1)",
	},
	{
		nombre: "asignación anidada a derecha",
		arbol: &Assign{
			Name: tok(token.IDENTIFIER, "a"),
			Value: &Assign{
				Name:  tok(token.IDENTIFIER, "b"),
				Value: num(1),
			},
		},
		esperado: "(= a (= b 1))",
	},
	{
		nombre: "llamada sin argumentos",
		arbol: &Call{
			Callee: variable("f"),
		},
		esperado: "(call f)",
	},
	{
		nombre: "llamada con argumentos",
		arbol: &Call{
			Callee:    variable("f"),
			Arguments: []Expr{num(1), num(2)},
		},
		esperado: "(call f 1 2)",
	},
	{
		nombre: "llamadas encadenadas",
		arbol: &Call{
			Callee:    &Call{Callee: variable("f")},
			Arguments: []Expr{num(1)},
		},
		esperado: "(call (call f) 1)",
	},
}

func TestRepresentarExpresiones(t *testing.T) {
	for _, caso := range casosDeExpresion {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenido := Representar(caso.arbol)
			if obtenido != caso.esperado {
				t.Errorf("se obtuvo %q y se esperaba %q", obtenido, caso.esperado)
			}
		})
	}
}
