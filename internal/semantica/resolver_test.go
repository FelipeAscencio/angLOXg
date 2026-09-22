package semantica_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/interprete"
	"github.com/FelipeAscencio/angLOXg/internal/semantica"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Verifica que el analizador semántico detecte el uso de una variable en su propio inicializador.
func TestErrorInicializadorPropio(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	tkA := token.Token{Lexema: "a", Linea: 5}
	stmtVar := &sintaxis.Var{
		Name:        tkA,
		Initializer: &sintaxis.Variable{Name: tkA},
	}

	stmtBlock := &sintaxis.Block{Statements: []sintaxis.Stmt{stmtVar}}
	err := res.Resolver([]sintaxis.Stmt{stmtBlock})
	if err == nil {
		t.Errorf("Se esperaba un error semántico al leer la variable en su propio inicializador")
	} else if !strings.Contains(err.Error(), "propio inicializador") {
		t.Errorf("Mensaje de error inesperado: %v", err)
	}
}

// Comprueba la correcta resolución de variables globales y ámbitos vacíos sin errores.
func TestResolucionGlobalYAmbitoVacio(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	tkG := token.Token{Lexema: "globalVar"}
	stmts := []sintaxis.Stmt{
		&sintaxis.Var{Name: tkG, Initializer: &sintaxis.Literal{Value: 1.0}},
		&sintaxis.Function{
			Name:   token.Token{Lexema: "fnGlobal"},
			Params: []token.Token{},
			Body:   []sintaxis.Stmt{},
		},
		&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: tkG}},
	}

	if err := res.Resolver(stmts); err != nil {
		t.Errorf("La resolución global falló inesperadamente: %v", err)
	}
}

// Verifica la correcta propagación de errores semánticos a través de diferentes sentencias contenedorizadas.
func TestPropagacionDeErroresEnResolver(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{
		Name:        tkX,
		Initializer: &sintaxis.Variable{Name: tkX},
	}

	casos := []sintaxis.Stmt{
		&sintaxis.Block{Statements: []sintaxis.Stmt{varMalo}},
		&sintaxis.If{Condition: &sintaxis.Literal{Value: true}, Then: &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo}}},
		&sintaxis.If{Condition: &sintaxis.Literal{Value: true}, Then: &sintaxis.Block{}, Else: &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo}}},
		&sintaxis.While{Condition: &sintaxis.Literal{Value: true}, Body: &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo}}},
		&sintaxis.Function{
			Name:   token.Token{Lexema: "fnMala", Linea: 1},
			Params: []token.Token{},
			Body:   []sintaxis.Stmt{varMalo},
		},
	}

	for _, stmt := range casos {
		r := semantica.Nuevo(intp)
		bloqueContenedor := &sintaxis.Block{Statements: []sintaxis.Stmt{stmt}}
		if err := r.ResolverSentencia(bloqueContenedor); err == nil {
			t.Errorf("Se esperaba error de resolución en el tipo %T", stmt)
		}
	}
}

// Comprueba que la resolución de expresiones variadas válidas no genere errores.
func TestResolucionExpresionesVariadas(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	tk := token.Token{Lexema: "a"}
	lit := &sintaxis.Literal{Value: 1.0}
	exprs := []sintaxis.Expr{
		&sintaxis.Binary{Left: lit, Operator: token.Token{Tipo: token.PLUS}, Right: lit},
		&sintaxis.Logical{Left: lit, Operator: token.Token{Tipo: token.OR}, Right: lit},
		&sintaxis.Unary{Operator: token.Token{Tipo: token.MINUS}, Right: lit},
		&sintaxis.Grouping{Expression: lit},
		&sintaxis.Call{Callee: lit, Arguments: []sintaxis.Expr{lit}},
		&sintaxis.Assign{Name: tk, Value: lit},
	}

	for _, expr := range exprs {
		if err := res.ResolverExpresion(expr); err != nil {
			t.Errorf("ResolverExpresion falló para %T: %v", expr, err)
		}
	}
}

// Verifica la detección de errores semánticos dentro de condiciones de estructuras de control (if, while).
func TestErrorEnCondicionesDeControl(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{
		Name:        tkX,
		Initializer: &sintaxis.Variable{Name: tkX},
	}

	condicionMala := &sintaxis.Variable{Name: tkX}
	casos := []sintaxis.Stmt{
		&sintaxis.If{
			Condition: condicionMala,
			Then:      &sintaxis.Block{},
		},
		&sintaxis.While{
			Condition: condicionMala,
			Body:      &sintaxis.Block{},
		},
	}

	for _, stmt := range casos {
		r := semantica.Nuevo(intp)
		bloque := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, stmt}}
		if err := r.ResolverSentencia(bloque); err == nil {
			t.Errorf("Se esperaba un error semántico al evaluar la condición en %T", stmt)
		}
	}
}

// Comprueba errores semánticos en cuerpos de funciones e invocaciones de llamadas.
func TestErroresEnFuncionesYLlamadas(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{
		Name:        tkX,
		Initializer: &sintaxis.Variable{Name: tkX},
	}

	fnMala := &sintaxis.Function{
		Name:   token.Token{Lexema: "miFunc"},
		Params: []token.Token{},
		Body:   []sintaxis.Stmt{varMalo},
	}

	r1 := semantica.Nuevo(intp)
	if err := r1.ResolverSentencia(fnMala); err == nil {
		t.Errorf("Se esperaba un error semántico dentro del cuerpo de la función")
	}

	callMalo := &sintaxis.Call{
		Callee:    &sintaxis.Literal{Value: "funcion"},
		Arguments: []sintaxis.Expr{&sintaxis.Variable{Name: tkX}},
	}
	
	r2 := semantica.Nuevo(intp)
	bloqueCall := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, &sintaxis.ExpressionStmt{Expression: callMalo}}}
	if err := r2.ResolverSentencia(bloqueCall); err == nil {
		t.Errorf("Se esperaba un error semántico en los argumentos de la llamada")
	}
}

// Verifica la propagación de errores semánticos en expresiones complejas anidadas.
func TestPropagacionDeErroresEnExpresionesComplejas(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{
		Name:        tkX,
		Initializer: &sintaxis.Variable{Name: tkX},
	}
	
	exprMala := &sintaxis.Variable{Name: tkX}
	lit := &sintaxis.Literal{Value: 1.0}
	casosExpr := []sintaxis.Expr{
		&sintaxis.Binary{Left: lit, Operator: token.Token{Tipo: token.PLUS}, Right: exprMala},
		&sintaxis.Logical{Left: lit, Operator: token.Token{Tipo: token.AND}, Right: exprMala},
		&sintaxis.Unary{Operator: token.Token{Tipo: token.MINUS}, Right: exprMala},
		&sintaxis.Grouping{Expression: exprMala},
	}

	for _, expr := range casosExpr {
		r := semantica.Nuevo(intp)
		stmt := &sintaxis.ExpressionStmt{Expression: expr}
		bloque := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, stmt}}
		if err := r.ResolverSentencia(bloque); err == nil {
			t.Errorf("Se esperaba propagación de error en la expresión %T", expr)
		}
	}
}
