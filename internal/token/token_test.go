package token_test

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

func TestResolverIdentificador(t *testing.T) {
	casos := []struct {
		lexema   string
		esperado token.TipoDeToken
	}{
		{"and", token.AND},
		{"while", token.WHILE},
		{"return", token.RETURN},
		{"if", token.IF},
		{"superVariable", token.IDENTIFIER},
		{"_privada", token.IDENTIFIER},
		{"minsky_even", token.IDENTIFIER},
	}

	for _, c := range casos {
		obtenido := token.ResolverIdentificador(c.lexema)
		if obtenido != c.esperado {
			t.Errorf("Para el lexema %q esperaba %v, pero obtuve %v", c.lexema, c.esperado, obtenido)
		}
	}
}

func TestTipoDeTokenStringLimites(t *testing.T) {
	if token.PLUS.String() != "PLUS" {
		t.Errorf("El string del token PLUS falló")
	}

	invalidoNegativo := token.TipoDeToken(-1)
	if invalidoNegativo.String() != "UNKNOWN" {
		t.Errorf("Un token negativo debería devolver UNKNOWN de manera segura")
	}

	invalidoPositivo := token.TipoDeToken(9999)
	if invalidoPositivo.String() != "UNKNOWN" {
		t.Errorf("Un token fuera de rango debería devolver UNKNOWN de manera segura")
	}
}

func TestTokenStringYConstructor(t *testing.T) {
	tkPlus := token.Nuevo(token.PLUS, "+", nil, 1)
	if tkPlus.String() != "PLUS" {
		t.Errorf("Esperado 'PLUS', obtenido %q", tkPlus.String())
	}

	tkId := token.Nuevo(token.IDENTIFIER, "contador", nil, 2)
	if tkId.String() != "IDENTIFIER<contador>" {
		t.Errorf("El formato de string para IDENTIFIER falló. Obtenido: %q", tkId.String())
	}

	tkNum := token.Nuevo(token.NUMBER, "3.14", 3.14, 3)
	if tkNum.String() != "NUMBER<3.14>" {
		t.Errorf("El formato de string para NUMBER falló. Obtenido: %q", tkNum.String())
	}

	tkStr := token.Nuevo(token.STRING, "\"hola\"", "hola", 4)
	if tkStr.String() != "STRING<hola>" {
		t.Errorf("El formato de string para STRING falló. Obtenido: %q", tkStr.String())
	}
}
