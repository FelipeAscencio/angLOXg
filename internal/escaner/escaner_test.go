package escaner

import (
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// tokenEsperado describe un token que el escáner tiene que producir.
type tokenEsperado struct {
	tipo    token.TipoDeToken
	lexema  string
	literal any
	linea   int
}

// casoDeEscaneo es una entrada de una tabla de pruebas: un fuente y lo que se
// espera obtener al escanearlo.
type casoDeEscaneo struct {
	nombre string
	fuente string

	// Tokens esperados, sin contar el EOF final: ese se verifica siempre.
	esperado []tokenEsperado

	// Subcadenas que tienen que aparecer en los errores, una por error y en
	// orden. Si es nil, el escaneo tiene que terminar sin fallas.
	errores []string
}

// correrCasos ejecuta una tabla de pruebas completa.
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

// verificarTokens compara la lista producida contra la esperada, agregándole el
// EOF con el que el escáner siempre tiene que cerrar.
func verificarTokens(t *testing.T, obtenidos []token.Token, esperados []tokenEsperado) {
	t.Helper()

	if len(obtenidos) == 0 {
		t.Fatalf("el escáner no produjo ni siquiera el token EOF")
	}

	// El último token siempre tiene que ser un EOF con lexema vacío.
	ultimo := obtenidos[len(obtenidos)-1]
	if ultimo.Tipo != token.EOF {
		t.Errorf("el último token es %s, se esperaba EOF", ultimo.Tipo)
	}

	if ultimo.Lexema != "" {
		t.Errorf("el token EOF tiene lexema %q, se esperaba vacío", ultimo.Lexema)
	}

	// El resto se compara campo por campo contra la tabla.
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

// verificarErrores chequea que las fallas léxicas sean las esperadas.
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
