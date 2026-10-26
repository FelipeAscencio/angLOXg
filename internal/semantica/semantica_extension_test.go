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

// Verifica la propagación de errores semánticos en las nuevas estructuras (Ternario, Arrays, Índices).
func TestResolucionTernarioYArrays(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})

	// Variable diseñada para fallar al resolverse (usada en su propio inicializador)
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{
		Name:        tkX,
		Initializer: &sintaxis.Variable{Name: tkX},
	}

	exprMala := &sintaxis.Variable{Name: tkX}
	lit := &sintaxis.Literal{Value: 1.0}

	casosExpr := []sintaxis.Expr{
		// Ternario: error en condición
		&sintaxis.Ternary{Condition: exprMala, TrueExpr: lit, FalseExpr: lit},
		// Ternario: error en rama verdadera
		&sintaxis.Ternary{Condition: lit, TrueExpr: exprMala, FalseExpr: lit},
		// Ternario: error en rama falsa
		&sintaxis.Ternary{Condition: lit, TrueExpr: lit, FalseExpr: exprMala},
		// Array: error en uno de los elementos
		&sintaxis.ArrayLiteral{Elements: []sintaxis.Expr{lit, exprMala}},
		// Indexado: error en el objeto referenciado
		&sintaxis.Index{Object: exprMala, Index: lit},
		// Indexado: error en el cálculo del índice
		&sintaxis.Index{Object: lit, Index: exprMala},
	}

	for _, expr := range casosExpr {
		r := semantica.Nuevo(intp)
		stmt := &sintaxis.ExpressionStmt{Expression: expr}
		
		// Encapsulamos en un bloque para que la variable se defina y el error estalle al leerla
		bloque := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, stmt}}
		
		if err := r.ResolverSentencia(bloque); err == nil {
			t.Errorf("Se esperaba propagación de error en la expresión %T", expr)
		}
	}
}

// Verifica las restricciones semánticas estrictas de la POO (Herencia, Init, Super).
func TestRestriccionesDeClases(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})

	// 1. Clase que hereda de sí misma
	res1 := semantica.Nuevo(intp)
	claseHeredaDeSiMisma := &sintaxis.Class{
		Name:       token.Token{Lexema: "Auto", Linea: 1},
		Superclass: &sintaxis.Variable{Name: token.Token{Lexema: "Auto", Linea: 1}},
	}
	
	err1 := res1.ResolverSentencia(claseHeredaDeSiMisma)
	if err1 == nil || !strings.Contains(err1.Error(), "heredar de sí misma") {
		t.Errorf("Falló detección de herencia de sí misma: %v", err1)
	}

	// 2. Intentar devolver un valor dentro de un método init()
	res2 := semantica.Nuevo(intp)
	metodoInit := &sintaxis.Function{
		Name:   token.Token{Lexema: "init", Linea: 2},
		Params: []token.Token{},
		Body: []sintaxis.Stmt{
			&sintaxis.Return{
				Keyword: token.Token{Lexema: "return", Linea: 3},
				Value:   &sintaxis.Literal{Value: "un valor"},
			},
		},
	}
	claseConInitMalo := &sintaxis.Class{
		Name:    token.Token{Lexema: "Moto", Linea: 2},
		Methods: []*sintaxis.Function{metodoInit},
	}
	
	err2 := res2.ResolverSentencia(claseConInitMalo)
	if err2 == nil || !strings.Contains(err2.Error(), "devolver un valor desde un inicializador") {
		t.Errorf("Falló detección de return con valor en init: %v", err2)
	}

	// 3. Uso de 'super' en una clase que no tiene superclase
	res3 := semantica.Nuevo(intp)
	metodoMaloSuper := &sintaxis.Function{
		Name:   token.Token{Lexema: "acelerar", Linea: 4},
		Params: []token.Token{},
		Body: []sintaxis.Stmt{
			&sintaxis.ExpressionStmt{
				Expression: &sintaxis.Super{
					Keyword: token.Token{Lexema: "super", Linea: 5},
					Method:  token.Token{Lexema: "acelerar"},
				},
			},
		},
	}
	claseSinSuper := &sintaxis.Class{
		Name:    token.Token{Lexema: "Bici", Linea: 4},
		Methods: []*sintaxis.Function{metodoMaloSuper},
	}
	
	err3 := res3.ResolverSentencia(claseSinSuper)
	if err3 == nil || !strings.Contains(err3.Error(), "clase sin superclase") {
		t.Errorf("Falló detección de super en clase sin superclase: %v", err3)
	}
}

// 1. Uso de 'super' completamente fuera de una clase.
func TestSuperFueraDeClase(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	stmt := &sintaxis.ExpressionStmt{
		Expression: &sintaxis.Super{
			Keyword: token.Token{Lexema: "super", Linea: 1},
			Method:  token.Token{Lexema: "metodo"},
		},
	}
	
	err := res.ResolverSentencia(stmt)
	if err == nil || !strings.Contains(err.Error(), "fuera de una clase") {
		t.Errorf("Falló detección de super fuera de clase: %v", err)
	}
}

// 2. Error al intentar resolver la expresión de la superclase.
func TestErrorResolviendoSuperclase(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{Name: tkX, Initializer: &sintaxis.Variable{Name: tkX}}

	clase := &sintaxis.Class{
		Name:       token.Token{Lexema: "ClaseHija"},
		Superclass: &sintaxis.Variable{Name: tkX}, // Falla al resolverse
	}

	bloque := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, clase}}
	if err := res.ResolverSentencia(bloque); err == nil {
		t.Errorf("Se esperaba propagación del error al resolver la superclase")
	}
}

// 3. Error semántico dentro de un método (activa la limpieza de scopes).
func TestErrorResolviendoMetodoActivaLimpieza(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	tkX := token.Token{Lexema: "x", Linea: 1}
	varMalo := &sintaxis.Var{Name: tkX, Initializer: &sintaxis.Variable{Name: tkX}}

	metodo := &sintaxis.Function{
		Name:   token.Token{Lexema: "m"},
		Params: []token.Token{},
		Body:   []sintaxis.Stmt{&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: tkX}}}, // Detona al leer 'x'
	}

	clase := &sintaxis.Class{
		Name:       token.Token{Lexema: "A"},
		Superclass: &sintaxis.Variable{Name: token.Token{Lexema: "B"}}, // Superclase válida genérica
		Methods:    []*sintaxis.Function{metodo},
	}

	bloque := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, clase}}
	if err := res.ResolverSentencia(bloque); err == nil {
		t.Errorf("Se esperaba propagación del error desde el método de la clase")
	}
}

// 4. Return pelado (sin valor de retorno).
func TestReturnSinValor(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)
	
	ret := &sintaxis.Return{Keyword: token.Token{Lexema: "return"}, Value: nil}
	if err := res.ResolverSentencia(ret); err != nil {
		t.Errorf("No se esperaba error para un return vacío, pero dio: %v", err)
	}
}

// 5. Propagación de errores en asignaciones y llamadas (Callee).
func TestUltimosCasosDeErrorExpresion(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	tkX := token.Token{Lexema: "x", Linea: 1}
	
	// Variable diseñada para fallar
	varMalo := &sintaxis.Var{Name: tkX, Initializer: &sintaxis.Variable{Name: tkX}}
	exprMala := &sintaxis.Variable{Name: tkX}

	casosExpr := []sintaxis.Expr{
		// Asignación donde el valor de la derecha falla al resolverse
		&sintaxis.Assign{Name: token.Token{Lexema: "y"}, Value: exprMala},
		
		// Llamada a función donde el Callee (la función a invocar) falla
		&sintaxis.Call{Callee: exprMala, Arguments: []sintaxis.Expr{}},
	}

	for _, expr := range casosExpr {
		r := semantica.Nuevo(intp)
		stmt := &sintaxis.ExpressionStmt{Expression: expr}
		bloque := &sintaxis.Block{Statements: []sintaxis.Stmt{varMalo, stmt}}
		
		if err := r.ResolverSentencia(bloque); err == nil {
			t.Errorf("Se esperaba propagación del error en la expresión %T", expr)
		}
	}
}

// 6. Test para cubrir los 'return nil' por defecto de los Type Switches.
func TestFallbackNodosDesconocidos(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	r := semantica.Nuevo(intp)
	
	// En Go, al pasar 'nil' a un type switch, no matchea ningún 'case' 
	// y cae directo al final (default), retornando nil.
	var stmtNulo sintaxis.Stmt = nil
	if err := r.ResolverSentencia(stmtNulo); err != nil {
		t.Errorf("Un Stmt nil debería ser ignorado silenciosamente, no retornar error: %v", err)
	}
	
	var exprNula sintaxis.Expr = nil
	if err := r.ResolverExpresion(exprNula); err != nil {
		t.Errorf("Un Expr nil debería ser ignorado silenciosamente, no retornar error: %v", err)
	}
}

// 7. Test de integración pesado para barrer casos borde (shadowing, closures anidados y limpiezas sin superclase).
func TestCoberturaTotalBordesSemanticos(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	// Caso A: Variable local sin inicializador (ej: var a;)
	varSinInit := &sintaxis.Var{
		Name:        token.Token{Lexema: "sinInit", Linea: 1},
		Initializer: nil,
	}

	// Caso B: Función con múltiples parámetros (fuerza la iteración completa del slice de params)
	tkParam1 := token.Token{Lexema: "p1", Linea: 2}
	tkParam2 := token.Token{Lexema: "p2", Linea: 2}
	fnMultiParam := &sintaxis.Function{
		Name:   token.Token{Lexema: "multi", Linea: 2},
		Params: []token.Token{tkParam1, tkParam2},
		Body: []sintaxis.Stmt{
			&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: tkParam1}},
			&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: tkParam2}},
		},
	}

	// Caso C: Clase SIN superclase pero con un método que falla.
	// Esto obliga al resolver a pasar por la rama falsa de cerrarAmbitosDeClase.
	tkX := token.Token{Lexema: "invalida", Linea: 3}
	varMalo := &sintaxis.Var{Name: tkX, Initializer: &sintaxis.Variable{Name: tkX}}
	metodoMalo := &sintaxis.Function{
		Name:   token.Token{Lexema: "metodoMalo", Linea: 3},
		Params: []token.Token{},
		Body:   []sintaxis.Stmt{varMalo},
	}
	claseSinSuperMala := &sintaxis.Class{
		Name:       token.Token{Lexema: "FallaSinSuper", Linea: 3},
		Superclass: nil,
		Methods:    []*sintaxis.Function{metodoMalo},
	}

	// Caso D: Bloque profundo válido para forzar a resolverLocal a iterar varios scopes hacia atrás.
	bloqueProfundo := &sintaxis.Block{
		Statements: []sintaxis.Stmt{
			&sintaxis.Var{Name: token.Token{Lexema: "lejana", Linea: 4}, Initializer: &sintaxis.Literal{Value: 1.0}},
			&sintaxis.Block{
				Statements: []sintaxis.Stmt{
					&sintaxis.Block{
						Statements: []sintaxis.Stmt{
							&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: token.Token{Lexema: "lejana", Linea: 4}}},
						},
					},
				},
			},
		},
	}

	// Validaciones
	if err := res.ResolverSentencia(varSinInit); err != nil {
		t.Errorf("Variable sin inicializador falló: %v", err)
	}
	if err := res.ResolverSentencia(fnMultiParam); err != nil {
		t.Errorf("Función con múltiples params falló: %v", err)
	}
	if err := res.ResolverSentencia(claseSinSuperMala); err == nil {
		t.Errorf("Clase sin superclase debió fallar por su método malo y limpiar el scope")
	}
	if err := res.ResolverSentencia(bloqueProfundo); err != nil {
		t.Errorf("Resolución de variables en bloques profundos falló: %v", err)
	}
}
