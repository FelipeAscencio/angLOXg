package lox_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Verifica el funcionamiento del scope estático frente a redefiniciones en closures.
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

	esperado := "global\nglobal\n"
	obtenido := salida.String()

	if obtenido != esperado {
		t.Errorf("Bug del Closure detectado.\nEsperado:\n%q\nObtenido:\n%q", esperado, obtenido)
	}
}

// Verifica que el resolver detecte el uso de una variable en su propio inicializador.
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

	if ok {
		t.Fatalf("Se esperaba que la ejecución fallara por error semántico, pero dio OK.")
	}

	obtenido := salida.String()
	if !strings.Contains(obtenido, "propio inicializador") {
		t.Errorf("Se esperaba un error de inicializador, pero se obtuvo: %s", obtenido)
	}
}

// Verifica la correcta resolución de shadowing en múltiples niveles de anidamiento.
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

// Verifica que un closure pueda modificar una variable perteneciente a un ámbito superior.
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
