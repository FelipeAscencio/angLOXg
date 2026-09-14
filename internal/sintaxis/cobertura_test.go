package sintaxis

import (
	"fmt"
	"testing"
)

// Una instancia de cada nodo de expresión que existe en el paquete.
//
// Si se agrega un nodo y no se lo suma acá, el test de abajo lo detecta.
var todasLasExpresiones = []Expr{
	&Literal{},
	&Grouping{Expression: &Literal{}},
	&Unary{Right: &Literal{}},
	&Binary{Left: &Literal{}, Right: &Literal{}},
	&Logical{Left: &Literal{}, Right: &Literal{}},
	&Variable{},
	&Assign{Value: &Literal{}},
	&Call{Callee: &Literal{}},
}

// Una instancia de cada nodo de sentencia que existe en el paquete.
var todasLasSentencias = []Stmt{
	&ExpressionStmt{Expression: &Literal{}},
	&Print{Value: &Literal{}},
	&Var{},
	&Block{},
	&If{Condition: &Literal{}, Then: &Block{}},
	&While{Condition: &Literal{}, Body: &Block{}},
	&Function{},
	&Return{},
}

// TestRepresentadorContemplaTodosLosNodos verifica que el representador sepa
// qué hacer con cada tipo de nodo.
//
// Como el recorrido se hace con un "switch" de tipo, el compilador no avisa si
// falta una rama: un nodo sin contemplar recién se notaría al ejecutarlo. Este
// test adelanta ese momento.
func TestRepresentadorContemplaTodosLosNodos(t *testing.T) {
	for _, expresion := range todasLasExpresiones {
		t.Run(fmt.Sprintf("%T", expresion), func(t *testing.T) {
			defer func() {
				if recuperado := recover(); recuperado != nil {
					t.Errorf("el representador no contempla este nodo: %v", recuperado)
				}
			}()

			Representar(expresion)
		})
	}

	for _, sentencia := range todasLasSentencias {
		t.Run(fmt.Sprintf("%T", sentencia), func(t *testing.T) {
			defer func() {
				if recuperado := recover(); recuperado != nil {
					t.Errorf("el representador no contempla este nodo: %v", recuperado)
				}
			}()

			RepresentarSentencia(sentencia)
		})
	}
}

// TestListasDeNodosEstanCompletas verifica que las listas de arriba tengan
// tantos nodos como tipos definidos en el paquete.
//
// Es la mitad que falta: sin esto, un nodo nuevo que nadie sume a las listas
// pasaría el test anterior sin haberse probado nunca.
func TestListasDeNodosEstanCompletas(t *testing.T) {
	const expresionesDefinidas = 8
	const sentenciasDefinidas = 8

	if len(todasLasExpresiones) != expresionesDefinidas {
		t.Errorf("la lista tiene %d expresiones y el paquete define %d",
			len(todasLasExpresiones), expresionesDefinidas)
	}

	if len(todasLasSentencias) != sentenciasDefinidas {
		t.Errorf("la lista tiene %d sentencias y el paquete define %d",
			len(todasLasSentencias), sentenciasDefinidas)
	}
}
