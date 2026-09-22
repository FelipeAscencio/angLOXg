package lox_test

import (
	"bytes"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Verifica el funcionamiento correcto de los diferentes modos de ejecución del CLI.
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

// Verifica que el intérprete detecte y reporte fallos ante errores léxicos y sintácticos.
func TestErroresEjecucion(t *testing.T) {
	var salida bytes.Buffer
	
	if lox.Ejecutar("?", &salida) {
		t.Errorf("Se esperaba fallo por error léxico")
	}

	if lox.Ejecutar("var 123;", &salida) {
		t.Errorf("Se esperaba fallo por error sintáctico")
	}
}
