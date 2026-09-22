package sintaxis

import (
	"fmt"
	"testing"
)

// Instancia de cada nodo de expresión existente en el paquete.
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

// Instancia de cada nodo de sentencia existente en el paquete.
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

// Verifica que el representador soporte y contemple todos los tipos de nodos definidos.
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

// Verifica que las listas de nodos incluyan exactamente la cantidad de tipos definidos en el paquete.
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
