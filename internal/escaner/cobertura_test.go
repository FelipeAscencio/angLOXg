package escaner

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Tablas de prueba del escáner para verificación de cobertura.
var todasLasTablas = [][]casoDeEscaneo{
	casosUnCaracter,
	casosDosCaracteres,
	casosComentarios,
	casosCadenas,
	casosNumeros,
	casosIdentificadores,
	casosPalabrasClave,
	casosPrograma,
}

// Verifica que todas las categorías léxicas tengan al menos un caso de prueba asociado.
func TestCoberturaDeTiposDeToken(t *testing.T) {
	cubiertos := make(map[token.TipoDeToken]bool)
	cubiertos[token.EOF] = true

	for _, suite := range todasLasTablas {
		for _, caso := range suite {
			for _, esperado := range caso.esperado {
				cubiertos[esperado.tipo] = true
			}
		}
	}

	var faltantes []token.TipoDeToken
	for tipo := token.LEFT_PAREN; tipo <= token.EOF; tipo++ {
		if !cubiertos[tipo] {
			faltantes = append(faltantes, tipo)
		}
	}

	if len(faltantes) > 0 {
		t.Errorf("hay %d categorías léxicas sin ningún caso de prueba: %v",
			len(faltantes), faltantes)
	}
}
