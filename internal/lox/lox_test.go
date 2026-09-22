package lox_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// TestBugDelClosure verifica que el scope estático funcione correctamente y que
// la redefinición de una variable en un scope externo no rompa el closure interno.
func TestBugDelClosure(t *testing.T) {
	fuente := `
	var a = "global";
	{
		fun showA() {
			print a;
		}
		
		showA();
		var a = "block";
		showA();
	}
	`
	var salida bytes.Buffer
	ok := lox.Ejecutar(fuente, &salida)

	if !ok {
		t.Fatalf("La ejecución falló inesperadamente.")
	}

	// Como el scope es estático, ambas llamadas deben imprimir "global"
	esperado := "global\nglobal\n"
	obtenido := salida.String()

	if obtenido != esperado {
		t.Errorf("Bug del Closure detectado.\nEsperado:\n%q\nObtenido:\n%q", esperado, obtenido)
	}
}

// TestErrorVariableEnPropioInicializador verifica que el Resolver atrape
// el uso de una variable local en su propia declaración.
func TestErrorVariableEnPropioInicializador(t *testing.T) {
	fuente := `
	{
		var a = "outer";
		{
			var a = a;
		}
	}
	`
	var salida bytes.Buffer
	ok := lox.Ejecutar(fuente, &salida)

	// El análisis semántico debe fallar y devolver false
	if ok {
		t.Fatalf("Se esperaba que la ejecución fallara por error semántico, pero dio OK.")
	}

	obtenido := salida.String()
	if !strings.Contains(obtenido, "propio inicializador") {
		t.Errorf("Se esperaba un error de inicializador, pero se obtuvo: %s", obtenido)
	}
}

// TestShadowingAnidado verifica la correcta resolución de saltos en múltiples niveles.
func TestShadowingAnidado(t *testing.T) {
	fuente := `
	var a = "global";
	{
		var a = "outer";
		{
			var a = "inner";
			print a;
		}
		print a;
	}
	print a;
	`
	var salida bytes.Buffer
	ok := lox.Ejecutar(fuente, &salida)

	if !ok {
		t.Fatalf("La ejecución falló inesperadamente.")
	}

	esperado := "inner\nouter\nglobal\n"
	obtenido := salida.String()

	if obtenido != esperado {
		t.Errorf("El shadowing anidado falló.\nEsperado:\n%q\nObtenido:\n%q", esperado, obtenido)
	}
}

// TestAsignacionDesdeClosure verifica que un closure pueda mutar una variable de un scope superior.
func TestAsignacionDesdeClosure(t *testing.T) {
	fuente := `
	var a = "antes";
	fun mutar() {
		a = "despues";
	}
	mutar();
	print a;
	`
	var salida bytes.Buffer
	ok := lox.Ejecutar(fuente, &salida)

	if !ok {
		t.Fatalf("La ejecución falló inesperadamente.")
	}

	esperado := "despues\n"
	obtenido := salida.String()

	if obtenido != esperado {
		t.Errorf("La mutación desde el closure falló.\nEsperado:\n%q\nObtenido:\n%q", esperado, obtenido)
	}
}
