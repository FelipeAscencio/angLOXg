package escaner

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Todas las tablas de pruebas del paquete. Si se agrega una nueva, hay que
// sumarla acá para que entre en el chequeo de cobertura.
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

// TestCoberturaDeTiposDeToken verifica que no quede ninguna categoría léxica sin
// al menos un caso de prueba que la produzca.
func TestCoberturaDeTiposDeToken(t *testing.T) {
	cubiertos := make(map[token.TipoDeToken]bool)

	// El EOF no aparece en ninguna tabla porque "verificarTokens" lo chequea en
	// todos los casos, sin necesidad de declararlo.
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
