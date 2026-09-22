package lox_test

import (
	"bytes"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

func TestModosCLI(t *testing.T) {
	var salida bytes.Buffer
	fuente := "var a = 1;"
	lox.Escanear(fuente, &salida)
	lox.Parsear(fuente, &salida)
	lox.Resolver(fuente, &salida)
	if salida.Len() == 0 {
		t.Errorf("Los modos CLI no generaron salida")
	}
}

func TestErroresEjecucion(t *testing.T) {
	var salida bytes.Buffer
	
	if lox.Ejecutar("?", &salida) {
		t.Errorf("Se esperaba fallo por error léxico")
	}

	if lox.Ejecutar("var 123;", &salida) {
		t.Errorf("Se esperaba fallo por error sintáctico")
	}
}
