package interprete

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

func TestEntornoObtenerYAsignar(t *testing.T) {
	global := NuevoEntorno(nil)
	local := NuevoEntorno(global)
	tkA := token.Token{Lexema: "a"}
	tkFantasma := token.Token{Lexema: "fantasma"}
	if _, err := global.Obtener(tkFantasma); err == nil {
		t.Errorf("Se esperaba error al obtener variable no definida")
	}

	if err := global.Asignar(tkFantasma, 10); err == nil {
		t.Errorf("Se esperaba error al asignar variable no definida")
	}

	global.Definir("a", 10)
	if val, _ := local.Obtener(tkA); val != 10 {
		t.Errorf("No se pudo obtener la variable del ancestro")
	}

	local.Asignar(tkA, 20)
	if val, _ := global.Obtener(tkA); val != 20 {
		t.Errorf("La asignación no mutó al ancestro dinámico")
	}
}

func TestEntornoEstatico(t *testing.T) {
	global := NuevoEntorno(nil)
	global.Definir("x", 100)
	intermedio := NuevoEntorno(global)
	local := NuevoEntorno(intermedio)
	if val := local.ObtenerEn(2, "x"); val != 100 {
		t.Errorf("ObtenerEn falló con saltos estáticos")
	}

	local.AsignarEn(2, "x", 999)
	if val := global.valores["x"]; val != 999 {
		t.Errorf("AsignarEn falló al mutar con saltos estáticos")
	}
}
