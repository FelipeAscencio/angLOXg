package interprete

import (
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

func TestErrorRuntime(t *testing.T) {
	tk := token.Token{Lexema: "+", Linea: 10}
	err := NewErrorRuntime(tk, "Prueba de error de ejecución")
	
	msg := err.Error()
	if !strings.Contains(msg, "Prueba de error de ejecución") {
		t.Errorf("El error no contiene el mensaje esperado")
	}
	if !strings.Contains(msg, "[línea 10]") {
		t.Errorf("El error no contiene la línea esperada")
	}
}

func TestValorRetorno(t *testing.T) {
	ret := &ValorRetorno{Valor: "test"}
	if ret.Error() != "retorno" {
		t.Errorf("El mensaje de ValorRetorno debe ser estático para corte de control")
	}
}
