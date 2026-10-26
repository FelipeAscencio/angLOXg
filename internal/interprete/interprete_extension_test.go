package interprete

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Verifica la correcta evaluación del operador ternario
func TestEvaluarTernario(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	verdadero := &sintaxis.Literal{Value: "verdad"}
	falso := &sintaxis.Literal{Value: "mentira"}

	// Caso 1: Condición Verdadera
	ternarioTrue := &sintaxis.Ternary{
		Condition: &sintaxis.Literal{Value: true},
		TrueExpr:  verdadero,
		FalseExpr: falso,
	}
	if val, _ := i.Evaluar(ternarioTrue); val != "verdad" {
		t.Errorf("El ternario con condición true debió retornar la rama verdadera, retornó %v", val)
	}

	// Caso 2: Condición Falsa
	ternarioFalse := &sintaxis.Ternary{
		Condition: &sintaxis.Literal{Value: false},
		TrueExpr:  verdadero,
		FalseExpr: falso,
	}
	if val, _ := i.Evaluar(ternarioFalse); val != "mentira" {
		t.Errorf("El ternario con condición false debió retornar la rama falsa, retornó %v", val)
	}

	// Caso 3: Error propagado desde la condición
	exprMala := &sintaxis.Variable{Name: token.Token{Lexema: "invalida"}}
	ternarioErr := &sintaxis.Ternary{
		Condition: exprMala,
		TrueExpr:  verdadero,
		FalseExpr: falso,
	}
	if _, err := i.Evaluar(ternarioErr); err == nil {
		t.Errorf("Se esperaba propagación de error desde la condición del ternario")
	}
}

// Verifica la creación y el indexado de Arrays
func TestEvaluarArrays(t *testing.T) {
	i := NuevoInterprete(&bytes.Buffer{})
	elementos := []sintaxis.Expr{
		&sintaxis.Literal{Value: 10.0},
		&sintaxis.Literal{Value: "texto"},
	}

	// Creación de ArrayLiteral
	arrayExpr := &sintaxis.ArrayLiteral{Elements: elementos}
	val, err := i.Evaluar(arrayExpr)
	if err != nil {
		t.Fatalf("Error inesperado al crear array: %v", err)
	}

	array, ok := val.(*LoxArray)
	if !ok {
		t.Fatalf("La expresión ArrayLiteral no devolvió un *LoxArray")
	}

	if len(array.Elementos) != 2 || array.Elementos[0] != 10.0 || array.Elementos[1] != "texto" {
		t.Errorf("El array no se inicializó con los elementos correctos")
	}

	// Representación de cadena
	if array.String() != "[10, texto]" {
		t.Errorf("La representación de cadena del array no es la esperada: %s", array.String())
	}

	// Error al intentar propagar desde los elementos durante la creación
	exprMala := &sintaxis.Variable{Name: token.Token{Lexema: "invalida"}}
	arrayMalo := &sintaxis.ArrayLiteral{Elements: []sintaxis.Expr{exprMala}}
	if _, err := i.Evaluar(arrayMalo); err == nil {
		t.Errorf("Se esperaba error al crear un array con un elemento inválido")
	}

	// Indexado válido
	indexExpr := &sintaxis.Index{
		Object:  arrayExpr,
		Bracket: token.Token{Tipo: token.LEFT_BRACKET},
		Index:   &sintaxis.Literal{Value: 1.0},
	}
	valIdx, _ := i.Evaluar(indexExpr)
	if valIdx != "texto" {
		t.Errorf("El acceso al array falló, esperado 'texto', obtenido %v", valIdx)
	}

	// Indexado inválido: objeto no es array
	indexInvalido1 := &sintaxis.Index{
		Object:  &sintaxis.Literal{Value: "no_soy_array"},
		Bracket: token.Token{Tipo: token.LEFT_BRACKET},
		Index:   &sintaxis.Literal{Value: 1.0},
	}
	if _, err := i.Evaluar(indexInvalido1); err == nil {
		t.Errorf("Se esperaba error al indexar algo que no es un array")
	}

	// Indexado inválido: índice no es un número entero
	indexInvalido2 := &sintaxis.Index{
		Object:  arrayExpr,
		Bracket: token.Token{Tipo: token.LEFT_BRACKET},
		Index:   &sintaxis.Literal{Value: 1.5},
	}
	if _, err := i.Evaluar(indexInvalido2); err == nil {
		t.Errorf("Se esperaba error al usar un índice fraccionario")
	}

	// Indexado fuera de rango
	indexInvalido3 := &sintaxis.Index{
		Object:  arrayExpr,
		Bracket: token.Token{Tipo: token.LEFT_BRACKET},
		Index:   &sintaxis.Literal{Value: 5.0},
	}
	if _, err := i.Evaluar(indexInvalido3); err == nil {
		t.Errorf("Se esperaba error al usar un índice fuera de rango")
	}
}

// Verifica toda la estructura y comportamiento de POO: Clases, Instancias, Propiedades, Métodos, Herencia y Super
func TestPOOCompleto(t *testing.T) {
	var buf bytes.Buffer
	i := NuevoInterprete(&buf)

	// --- 1. Definición de Clase sin superclase ---
	claseStmt := &sintaxis.Class{
		Name:    token.Token{Lexema: "MiClase"},
		Methods: []*sintaxis.Function{},
	}
	err := i.Ejecutar(claseStmt)
	if err != nil {
		t.Fatalf("Fallo al definir la clase: %v", err)
	}

	valClase, _ := i.entorno.Obtener(token.Token{Lexema: "MiClase"})
	loxClase, ok := valClase.(*LoxClass)
	if !ok || loxClase.Nombre != "MiClase" {
		t.Errorf("La clase no se registró correctamente en el entorno")
	}

	if loxClase.String() != "MiClase" {
		t.Errorf("String() de la clase falló")
	}

	if loxClase.Aridad() != 0 {
		t.Errorf("Aridad de clase sin init debe ser 0")
	}

	// --- 2. Instanciación y Propiedades (Get / Set) ---
	instanciaVal, _ := loxClase.Llamar(i, []any{})
	instancia, ok := instanciaVal.(*LoxInstance)
	if !ok {
		t.Fatalf("La instanciación no devolvió un *LoxInstance")
	}

	if !strings.Contains(instancia.String(), "instancia de MiClase") {
		t.Errorf("String() de la instancia falló")
	}

	// Set de una propiedad
	setExpr := &sintaxis.Set{
		Object: &sintaxis.Literal{Value: instancia},
		Name:   token.Token{Lexema: "edad"},
		Value:  &sintaxis.Literal{Value: 25.0},
	}
	valSet, _ := i.Evaluar(setExpr)
	if valSet != 25.0 {
		t.Errorf("Set falló al asignar propiedad")
	}

	// Get de una propiedad
	getExpr := &sintaxis.Get{
		Object: &sintaxis.Literal{Value: instancia},
		Name:   token.Token{Lexema: "edad"},
	}
	valGet, _ := i.Evaluar(getExpr)
	if valGet != 25.0 {
		t.Errorf("Get falló al leer propiedad, obtenido: %v", valGet)
	}

	// Get de una propiedad inexistente
	getInexistente := &sintaxis.Get{
		Object: &sintaxis.Literal{Value: instancia},
		Name:   token.Token{Lexema: "noExiste"},
	}
	if _, err := i.Evaluar(getInexistente); err == nil {
		t.Errorf("Se esperaba error al leer propiedad inexistente")
	}

	// Set/Get sobre algo que no es instancia
	if _, err := i.Evaluar(&sintaxis.Set{Object: &sintaxis.Literal{Value: 5.0}}); err == nil {
		t.Errorf("Set sobre un número debe fallar")
	}
	if _, err := i.Evaluar(&sintaxis.Get{Object: &sintaxis.Literal{Value: 5.0}}); err == nil {
		t.Errorf("Get sobre un número debe fallar")
	}

	// --- 3. Inicializador y This ---
	thisExpr := &sintaxis.This{Keyword: token.Token{Lexema: "this"}}
	
	// Guardamos la referencia exacta en memoria de la variable "x"
	varXExpr := &sintaxis.Variable{Name: token.Token{Lexema: "x"}} 
	
	setThisExpr := &sintaxis.Set{
		Object: thisExpr,
		Name:   token.Token{Lexema: "x"},
		Value:  varXExpr, // Usamos la referencia
	}
	
	metodoInit := &sintaxis.Function{
		Name:   token.Token{Lexema: "init"},
		Params: []token.Token{{Lexema: "x"}},
		Body:   []sintaxis.Stmt{&sintaxis.ExpressionStmt{Expression: setThisExpr}},
	}

	claseConInit := &sintaxis.Class{
		Name:    token.Token{Lexema: "ClaseInit"},
		Methods: []*sintaxis.Function{metodoInit},
	}
	i.Ejecutar(claseConInit)
	
	valClaseInit, _ := i.entorno.Obtener(token.Token{Lexema: "ClaseInit"})
	loxClaseInit := valClaseInit.(*LoxClass)

	// Inyectamos las distancias usando LAS MISMAS claves en memoria
	i.ResolverLocal(thisExpr, 1)
	i.ResolverLocal(varXExpr, 0) 

	instanciaInitVal, err := loxClaseInit.Llamar(i, []any{99.0})
	if err != nil {
		t.Fatalf("Falló llamar a clase con init: %v", err)
	}

	instanciaInit := instanciaInitVal.(*LoxInstance)
	valPropX, _ := instanciaInit.Obtener(token.Token{Lexema: "x"})
	if valPropX != 99.0 {
		t.Errorf("El constructor (init) falló en mutar el estado con this")
	}
}