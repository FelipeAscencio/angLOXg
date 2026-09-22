package interprete

import (
	"bytes"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

func TestReglasSemanticasBasicas(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	if i.esVerdadero(nil) != false { t.Errorf("nil debe ser falso") }
	if i.esVerdadero(false) != false { t.Errorf("false debe ser falso") }
	if i.esVerdadero(true) != true { t.Errorf("true debe ser verdadero") }
	if i.esVerdadero(0) != true { t.Errorf("0 debe ser verdadero en Lox") }
	if !i.esIgual(nil, nil) { t.Errorf("nil == nil") }
	if i.esIgual(nil, 5) { t.Errorf("nil != 5") }
	if !i.esIgual(5.0, 5.0) { t.Errorf("5 == 5") }
	if i.esIgual(5.0, 6.0) { t.Errorf("5 != 6") }
	tk := token.Token{Tipo: token.MINUS}
	if i.checkNumeroOperando(tk, 5.0) != nil { t.Errorf("5.0 es num") }
	if i.checkNumeroOperando(tk, "5") == nil { t.Errorf("str no es num") }
	if i.checkNumerosOperandos(tk, 5.0, 5.0) != nil { t.Errorf("num, num es valido") }
	if i.checkNumerosOperandos(tk, 5.0, "5") == nil { t.Errorf("num, str es invalido") }
}

func TestEvaluarExpresionesAritmeticasYLogicas(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	num := &sintaxis.Literal{Value: 10.0}
	str := &sintaxis.Literal{Value: "hola"}
	tkPlus := token.Token{Tipo: token.PLUS}
	tkMinus := token.Token{Tipo: token.MINUS}
	tkSlash := token.Token{Tipo: token.SLASH}
	tkPercent := token.Token{Tipo: token.PERCENT}
	if val, _ := i.Evaluar(&sintaxis.Unary{Operator: tkMinus, Right: num}); val != -10.0 {
		t.Errorf("Unario falló")
	}

	if val, _ := i.Evaluar(&sintaxis.Binary{Left: num, Operator: tkPlus, Right: num}); val != 20.0 {
		t.Errorf("Suma falló")
	}

	if val, _ := i.Evaluar(&sintaxis.Binary{Left: str, Operator: tkPlus, Right: str}); val != "holahola" {
		t.Errorf("Concatenación falló")
	}

	if _, err := i.Evaluar(&sintaxis.Binary{Left: num, Operator: tkPlus, Right: str}); err == nil {
		t.Errorf("Suma mixta debe fallar")
	}

	if _, err := i.Evaluar(&sintaxis.Binary{Left: num, Operator: tkSlash, Right: &sintaxis.Literal{Value: 0.0}}); err == nil {
		t.Errorf("División por cero debe fallar")
	}

	if _, err := i.Evaluar(&sintaxis.Binary{Left: num, Operator: tkPercent, Right: &sintaxis.Literal{Value: 0.0}}); err == nil {
		t.Errorf("Módulo por cero debe fallar")
	}

	valOr, _ := i.Evaluar(&sintaxis.Logical{Left: &sintaxis.Literal{Value: true}, Operator: token.Token{Tipo: token.OR}, Right: str})
	if valOr != true { t.Errorf("OR cortocircuito falló") }
}

func TestEjecutarSentenciasYControlDeFlujo(t *testing.T) {
	var buf bytes.Buffer
	i := NuevoInterprete(&buf)
	i.Ejecutar(&sintaxis.Print{Value: &sintaxis.Literal{Value: "salida"}})
	if buf.String() != "salida\n" { t.Errorf("Print falló") }
	
	i.Ejecutar(&sintaxis.Var{Name: token.Token{Lexema: "x"}, Initializer: &sintaxis.Literal{Value: 99.0}})
	if val, _ := i.globales.Obtener(token.Token{Lexema: "x"}); val != 99.0 {
		t.Errorf("Declaración de variable falló")
	}

	i.Ejecutar(&sintaxis.While{
		Condition: &sintaxis.Literal{Value: false},
		Body: &sintaxis.Print{Value: &sintaxis.Literal{Value: "no"}},
	})
}

func TestLlamadasYRetornos(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	if _, err := i.Evaluar(&sintaxis.Call{Callee: &sintaxis.Literal{Value: 5.0}}); err == nil {
		t.Errorf("Llamar a un número debe fallar")
	}

	varNode := &sintaxis.Variable{Name: token.Token{Lexema: "param"}}
	i.ResolverLocal(varNode, 0)
	fnNode := &sintaxis.Function{
		Name: token.Token{Lexema: "testFn"},
		Params: []token.Token{{Lexema: "param"}},
		Body: []sintaxis.Stmt{
			&sintaxis.Return{Value: varNode},
		},
	}
	
	loxFn := &LoxFunction{Declaracion: fnNode, Closure: i.entorno}
	if loxFn.Aridad() != 1 { t.Errorf("Aridad incorrecta") }
	res, err := loxFn.Llamar(i, []any{"argumento inyectado"})
	if err != nil { 
		t.Errorf("Llamar devolvió error inesperado: %v", err) 
	}

	if res != "argumento inyectado" { 
		t.Errorf("El retorno de la función no fue capturado correctamente") 
	}
}

func TestResolverLocalEIntreprete(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	expr := &sintaxis.Variable{Name: token.Token{Lexema: "z"}}
	i.ResolverLocal(expr, 1)
	if dist, ok := i.locales[expr]; !ok || dist != 1 {
		t.Errorf("ResolverLocal falló")
	}
}

func TestInterpretarYFlujoDeControl(t *testing.T) {
	var buf bytes.Buffer
	i := NuevoInterprete(&buf)
	exprStmt := &sintaxis.ExpressionStmt{Expression: &sintaxis.Literal{Value: "hola"}}
	ifStmt := &sintaxis.If{
		Condition: &sintaxis.Literal{Value: true},
		Then:      &sintaxis.Print{Value: &sintaxis.Literal{Value: 1.0}},
		Else:      &sintaxis.Print{Value: &sintaxis.Literal{Value: 2.0}},
	}

	ifFalseStmt := &sintaxis.If{
		Condition: &sintaxis.Literal{Value: false},
		Then:      &sintaxis.Print{Value: &sintaxis.Literal{Value: 1.0}},
		Else:      &sintaxis.Print{Value: &sintaxis.Literal{Value: 2.0}},
	}

	err := i.Interpretar([]sintaxis.Stmt{exprStmt, ifStmt, ifFalseStmt})
	if err != nil {
		t.Errorf("Interpretar devolvió error inesperado: %v", err)
	}

	esperado := "1\n2\n"
	if buf.String() != esperado {
		t.Errorf("Flujo If/Else falló. Esperado %q, obtenido %q", esperado, buf.String())
	}
}

func TestImpresionDeNil(t *testing.T) {
	var buf bytes.Buffer
	i := NuevoInterprete(&buf)
	i.Ejecutar(&sintaxis.Print{Value: &sintaxis.Literal{Value: nil}})
	if buf.String() != "nil\n" {
		t.Errorf("Esperado 'nil\\n', obtenido %q", buf.String())
	}
}

func TestEvaluarComparacionesYAgrupaciones(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	num5 := &sintaxis.Literal{Value: 5.0}
	num10 := &sintaxis.Literal{Value: 10.0}
	grp := &sintaxis.Grouping{Expression: num5}
	if val, _ := i.Evaluar(grp); val != 5.0 {
		t.Errorf("Grouping falló")
	}

	notTrue := &sintaxis.Unary{Operator: token.Token{Tipo: token.BANG}, Right: &sintaxis.Literal{Value: true}}
	if val, _ := i.Evaluar(notTrue); val != false {
		t.Errorf("BANG falló")
	}

	casos := []struct {
		op       token.TipoDeToken
		izq      sintaxis.Expr
		der      sintaxis.Expr
		esperado any
	}{
		{token.GREATER, num10, num5, true},
		{token.GREATER_EQUAL, num5, num5, true},
		{token.LESS, num5, num10, true},
		{token.LESS_EQUAL, num10, num10, true},
		{token.EQUAL_EQUAL, num5, num5, true},
		{token.BANG_EQUAL, num5, num10, true},
	}

	for _, c := range casos {
		expr := &sintaxis.Binary{Left: c.izq, Operator: token.Token{Tipo: c.op}, Right: c.der}
		if val, _ := i.Evaluar(expr); val != c.esperado {
			t.Errorf("Fallo en operador %v", c.op)
		}
	}
}

func TestOperadoresLogicosCompletos(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	tTrue := &sintaxis.Literal{Value: true}
	tFalse := &sintaxis.Literal{Value: false}
	t1 := &sintaxis.Literal{Value: 1.0}
	orFalse := &sintaxis.Logical{Left: tFalse, Operator: token.Token{Tipo: token.OR}, Right: t1}
	if val, _ := i.Evaluar(orFalse); val != 1.0 {
		t.Errorf("OR camino largo falló")
	}

	andFalse := &sintaxis.Logical{Left: tFalse, Operator: token.Token{Tipo: token.AND}, Right: t1}
	if val, _ := i.Evaluar(andFalse); val != false {
		t.Errorf("AND cortocircuito falló")
	}

	andTrue := &sintaxis.Logical{Left: tTrue, Operator: token.Token{Tipo: token.AND}, Right: t1}
	if val, _ := i.Evaluar(andTrue); val != 1.0 {
		t.Errorf("AND camino largo falló")
	}
}

func TestAsignacionYVariablesNoDeclaradas(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	asignarInvalido := &sintaxis.Assign{
		Name:  token.Token{Lexema: "noexisto"},
		Value: &sintaxis.Literal{Value: 5.0},
	}

	if _, err := i.Evaluar(asignarInvalido); err == nil {
		t.Errorf("Asignar a una variable no declarada debió lanzar error")
	}

	varInvalida := &sintaxis.Variable{Name: token.Token{Lexema: "noexisto"}}
	if _, err := i.Evaluar(varInvalida); err == nil {
		t.Errorf("Evaluar variable no declarada debió lanzar error")
	}
}

func TestPropagacionDeErroresEnEjecutar(t *testing.T) {
	var buf bytes.Buffer
	i := NuevoInterprete(&buf)
	exprMala := &sintaxis.Variable{Name: token.Token{Lexema: "invalida"}}
	casos := []sintaxis.Stmt{
		&sintaxis.ExpressionStmt{Expression: exprMala},
		&sintaxis.Print{Value: exprMala},
		&sintaxis.Var{Name: token.Token{Lexema: "a"}, Initializer: exprMala},
		&sintaxis.If{Condition: exprMala, Then: &sintaxis.Block{}},
		&sintaxis.While{Condition: exprMala, Body: &sintaxis.Block{}},
		&sintaxis.While{
			// Entra al bucle, pero falla evaluando el cuerpo
			Condition: &sintaxis.Literal{Value: true}, 
			Body: &sintaxis.ExpressionStmt{Expression: exprMala}, 
		},
		&sintaxis.Block{Statements: []sintaxis.Stmt{&sintaxis.ExpressionStmt{Expression: exprMala}}},
		&sintaxis.Return{Value: exprMala},
	}

	for _, stmt := range casos {
		if err := i.Ejecutar(stmt); err == nil {
			t.Errorf("Se esperaba que Ejecutar propagara el error de la expresión mala en %T", stmt)
		}
	}
}

func TestPropagacionDeErroresEnEvaluar(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	exprMala := &sintaxis.Variable{Name: token.Token{Lexema: "invalida"}}
	exprBuena := &sintaxis.Literal{Value: 1.0}
	casos := []sintaxis.Expr{
		&sintaxis.Unary{Operator: token.Token{Tipo: token.MINUS}, Right: exprMala},
		&sintaxis.Binary{Left: exprMala, Operator: token.Token{Tipo: token.PLUS}, Right: exprBuena},
		&sintaxis.Binary{Left: exprBuena, Operator: token.Token{Tipo: token.PLUS}, Right: exprMala},
		&sintaxis.Logical{Left: exprMala, Operator: token.Token{Tipo: token.OR}, Right: exprBuena},
		&sintaxis.Assign{Name: token.Token{Lexema: "a"}, Value: exprMala},
		&sintaxis.Call{Callee: exprMala},
		&sintaxis.Call{Callee: &sintaxis.Literal{Value: 5.0}, Arguments: []sintaxis.Expr{exprMala}},
	}

	for _, expr := range casos {
		if _, err := i.Evaluar(expr); err == nil {
			t.Errorf("Se esperaba que Evaluar propagara el error en %T", expr)
		}
	}
}

func TestErrorEnCuerpoDeFuncion(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	fnNode := &sintaxis.Function{
		Name: token.Token{Lexema: "f"},
		Params: []token.Token{},
		Body: []sintaxis.Stmt{
			&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: token.Token{Lexema: "noexiste"}}},
		},
	}
	
	loxFn := &LoxFunction{Declaracion: fnNode, Closure: i.entorno}
	_, err := loxFn.Llamar(i, []any{})
	if err == nil {
		t.Errorf("La función debería propagar el error fatal, no silenciarlo")
	}
}