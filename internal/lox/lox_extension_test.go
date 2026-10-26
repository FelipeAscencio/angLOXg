package lox_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Verifica la correcta ejecución y propagación de errores en la Máquina Virtual.
func TestEjecutarVMCobertura(t *testing.T) {
	var salida bytes.Buffer

	// 1. Ejecución exitosa en la VM
	if !lox.EjecutarVM("var a = 1; print a;", &salida) {
		t.Errorf("EjecutarVM falló con código válido")
	}
	if !strings.Contains(salida.String(), "1") {
		t.Errorf("EjecutarVM no generó la salida esperada")
	}

	// 2. Error léxico propagado en la VM
	salida.Reset()
	if lox.EjecutarVM("var a = @;", &salida) {
		t.Errorf("EjecutarVM debió fallar por error léxico y retornar false")
	}

	// 3. Error en tiempo de ejecución (Runtime Error) en la VM
	salida.Reset()
	if lox.EjecutarVM("print -\"hola\";", &salida) {
		t.Errorf("EjecutarVM debió fallar por error en tiempo de ejecución y retornar false")
	}
}

// Verifica la propagación de false cuando ocurren errores tempranos en los modos del CLI.
func TestFalloEnModosCLI(t *testing.T) {
	var salida bytes.Buffer
	fuenteInvalida := "var a = @;" // El '@' fuerza un error en el escáner

	if lox.Resolver(fuenteInvalida, &salida) {
		t.Errorf("Resolver debió retornar false ante un error léxico")
	}

	if lox.Parsear(fuenteInvalida, &salida) {
		t.Errorf("Parsear debió retornar false ante un error léxico")
	}

	if lox.Escanear(fuenteInvalida, &salida) {
		t.Errorf("Escanear debió retornar false ante un error léxico")
	}
}

// Verifica la falla de Resolver cuando el código es válido léxicamente pero falla en la semántica.
func TestFalloSemanticoEnResolver(t *testing.T) {
	var salida bytes.Buffer
	// Código válido sintácticamente pero con un error semántico (this fuera de clase)
	fuenteSemanticaMala := "print this;"

	if lox.Resolver(fuenteSemanticaMala, &salida) {
		t.Errorf("Resolver debió retornar false ante un error semántico")
	}
}
