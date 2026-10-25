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

// resolverBloque resuelve las sentencias dentro de un ámbito local y devuelve
// el error, si hubo.
func resolverBloque(sentencias ...sintaxis.Stmt) error {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	return res.Resolver([]sintaxis.Stmt{&sintaxis.Block{Statements: sentencias}})
}

// exigirErrorDeInicializadorPropio verifica que el resolvedor haya detectado la
// lectura de una variable dentro de su propio inicializador.
//
// Cada caso mete esa lectura en un lugar distinto de un nodo nuevo: si el
// resolvedor no desciende hasta ahí, el error no aparece y el test falla.
func exigirErrorDeInicializadorPropio(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("se esperaba un error semántico: el resolvedor no descendió al nodo")
	}

	if !strings.Contains(err.Error(), "propio inicializador") {
		t.Errorf("mensaje de error inesperado: %v", err)
	}
}

// Comprueba que declarar una clase en un ámbito local no rompa la resolución.
func TestResolucionDeClaseLocal(t *testing.T) {
	tkClase := token.Token{Lexema: "Punto", Linea: 1}

	err := resolverBloque(
		&sintaxis.Class{Name: tkClase},
		&sintaxis.ExpressionStmt{Expression: &sintaxis.Variable{Name: tkClase}},
	)

	if err != nil {
		t.Errorf("no se esperaba un error al resolver la clase: %v", err)
	}
}

// Comprueba que el resolvedor descienda al objeto de una lectura.
func TestResolucionDelObjetoDeUnaLectura(t *testing.T) {
	tkA := token.Token{Lexema: "a", Linea: 7}

	err := resolverBloque(&sintaxis.Var{
		Name: tkA,
		Initializer: &sintaxis.Get{
			Object: &sintaxis.Variable{Name: tkA},
			Name:   token.Token{Lexema: "x", Linea: 7},
		},
	})

	exigirErrorDeInicializadorPropio(t, err)
}

// Comprueba que el resolvedor descienda al objeto de una escritura.
func TestResolucionDelObjetoDeUnaEscritura(t *testing.T) {
	tkA := token.Token{Lexema: "a", Linea: 7}

	err := resolverBloque(&sintaxis.Var{
		Name: tkA,
		Initializer: &sintaxis.Set{
			Object: &sintaxis.Variable{Name: tkA},
			Name:   token.Token{Lexema: "x", Linea: 7},
			Value:  &sintaxis.Literal{Value: 1.0},
		},
	})

	exigirErrorDeInicializadorPropio(t, err)
}

// Comprueba que el resolvedor descienda al valor de una escritura.
func TestResolucionDelValorDeUnaEscritura(t *testing.T) {
	tkA := token.Token{Lexema: "a", Linea: 7}
	tkP := token.Token{Lexema: "p", Linea: 7}

	err := resolverBloque(
		&sintaxis.Var{Name: tkP, Initializer: &sintaxis.Literal{Value: 1.0}},
		&sintaxis.Var{
			Name: tkA,
			Initializer: &sintaxis.Set{
				Object: &sintaxis.Variable{Name: tkP},
				Name:   token.Token{Lexema: "x", Linea: 7},
				Value:  &sintaxis.Variable{Name: tkA},
			},
		},
	)

	exigirErrorDeInicializadorPropio(t, err)
}

// Usar this fuera de una clase es un error semántico: no hay ninguna instancia
// a la que pueda referirse.
func TestErrorDeThisFueraDeUnaClase(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	sentencias := []sintaxis.Stmt{&sintaxis.Print{
		Value: &sintaxis.This{Keyword: token.Token{Lexema: "this", Linea: 4}},
	}}

	err := res.Resolver(sentencias)
	if err == nil {
		t.Fatal("se esperaba un error semántico al usar this fuera de una clase")
	}

	if !strings.Contains(err.Error(), "fuera de una clase") {
		t.Errorf("mensaje de error inesperado: %v", err)
	}
}

// Usar this fuera de una clase también falla dentro de una función suelta.
func TestErrorDeThisEnUnaFuncionSuelta(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	sentencias := []sintaxis.Stmt{&sintaxis.Function{
		Name: token.Token{Lexema: "suelta", Linea: 1},
		Body: []sintaxis.Stmt{&sintaxis.Print{
			Value: &sintaxis.This{Keyword: token.Token{Lexema: "this", Linea: 2}},
		}},
	}}

	if err := res.Resolver(sentencias); err == nil {
		t.Error("se esperaba un error semántico al usar this en una función suelta")
	}
}

// Dentro de un método, this es válido.
func TestThisEsValidoEnUnMetodo(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	sentencias := []sintaxis.Stmt{&sintaxis.Class{
		Name: token.Token{Lexema: "A", Linea: 1},
		Methods: []*sintaxis.Function{{
			Name: token.Token{Lexema: "m", Linea: 2},
			Body: []sintaxis.Stmt{&sintaxis.Print{
				Value: &sintaxis.This{Keyword: token.Token{Lexema: "this", Linea: 2}},
			}},
		}},
	}}

	if err := res.Resolver(sentencias); err != nil {
		t.Errorf("no se esperaba un error al usar this en un método: %v", err)
	}
}

// Después de cerrar la clase, this vuelve a ser inválido.
func TestThisVuelveASerInvalidoTrasLaClase(t *testing.T) {
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	sentencias := []sintaxis.Stmt{
		&sintaxis.Class{Name: token.Token{Lexema: "A", Linea: 1}},
		&sintaxis.Print{Value: &sintaxis.This{Keyword: token.Token{Lexema: "this", Linea: 2}}},
	}

	if err := res.Resolver(sentencias); err == nil {
		t.Error("se esperaba un error semántico: el ámbito de la clase ya se cerró")
	}
}

// El cuerpo de un método se recorre: un error adentro tiene que salir a la luz.
func TestResolucionDelCuerpoDeUnMetodo(t *testing.T) {
	tkA := token.Token{Lexema: "a", Linea: 3}
	intp := interprete.NuevoInterprete(&bytes.Buffer{})
	res := semantica.Nuevo(intp)

	sentencias := []sintaxis.Stmt{&sintaxis.Class{
		Name: token.Token{Lexema: "A", Linea: 1},
		Methods: []*sintaxis.Function{{
			Name: token.Token{Lexema: "m", Linea: 2},
			Body: []sintaxis.Stmt{&sintaxis.Var{
				Name:        tkA,
				Initializer: &sintaxis.Variable{Name: tkA},
			}},
		}},
	}}

	exigirErrorDeInicializadorPropio(t, res.Resolver(sentencias))
}
