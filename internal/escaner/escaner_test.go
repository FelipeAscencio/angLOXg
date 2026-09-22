package escaner

import (
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Describe un token esperado en el resultado del escaneo.
type tokenEsperado struct {
	tipo    token.TipoDeToken
	lexema  string
	literal any
	linea   int
}

// Representa un caso de prueba individual para las tablas del escáner.
type casoDeEscaneo struct {
	nombre   string
	fuente   string
	esperado []tokenEsperado
	errores  []string
}

// Ejecuta una tabla de pruebas completa.
func correrCasos(t *testing.T, casos []casoDeEscaneo) {
	t.Helper()

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			tokens, errores := Nuevo(caso.fuente).EscanearTokens()
			verificarErrores(t, errores, caso.errores)
			verificarTokens(t, tokens, caso.esperado)
		})
	}
}

// Compara los tokens obtenidos contra los esperados.
func verificarTokens(t *testing.T, obtenidos []token.Token, esperados []tokenEsperado) {
	t.Helper()

	if len(obtenidos) == 0 {
		t.Fatalf("el escáner no produjo ni siquiera el token EOF")
	}

	ultimo := obtenidos[len(obtenidos)-1]
	if ultimo.Tipo != token.EOF {
		t.Errorf("el último token es %s, se esperaba EOF", ultimo.Tipo)
	}

	if ultimo.Lexema != "" {
		t.Errorf("el token EOF tiene lexema %q, se esperaba vacío", ultimo.Lexema)
	}

	reales := obtenidos[:len(obtenidos)-1]
	if len(reales) != len(esperados) {
		t.Fatalf("se obtuvieron %d tokens y se esperaban %d\nobtenidos: %v",
			len(reales), len(esperados), reales)
	}

	for i, esperado := range esperados {
		obtenido := reales[i]

		if obtenido.Tipo != esperado.tipo {
			t.Errorf("token %d: tipo %s, se esperaba %s", i, obtenido.Tipo, esperado.tipo)
		}

		if obtenido.Lexema != esperado.lexema {
			t.Errorf("token %d (%s): lexema %q, se esperaba %q",
				i, esperado.tipo, obtenido.Lexema, esperado.lexema)
		}

		if obtenido.Literal != esperado.literal {
			t.Errorf("token %d (%s): literal %#v, se esperaba %#v",
				i, esperado.tipo, obtenido.Literal, esperado.literal)
		}

		if obtenido.Linea != esperado.linea {
			t.Errorf("token %d (%s): línea %d, se esperaba %d",
				i, esperado.tipo, obtenido.Linea, esperado.linea)
		}
	}
}

// Verifica que los errores léxicos coincidan con los esperados.
func verificarErrores(t *testing.T, obtenidos []error, esperados []string) {
	t.Helper()

	if len(obtenidos) != len(esperados) {
		t.Fatalf("se obtuvieron %d errores y se esperaban %d\nobtenidos: %v",
			len(obtenidos), len(esperados), obtenidos)
	}

	for i, esperado := range esperados {
		if !strings.Contains(obtenidos[i].Error(), esperado) {
			t.Errorf("error %d: %q, se esperaba que contuviera %q",
				i, obtenidos[i].Error(), esperado)
		}
	}
}
