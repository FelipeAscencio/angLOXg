package compilador

import "testing"

func TestNuevoChunkArrancaVacio(t *testing.T) {
	ch := NuevoChunk()

	if ch == nil {
		t.Fatal("NuevoChunk devolvió nil")
	}

	if len(ch.Instrucciones) != 0 || len(ch.Constantes) != 0 || len(ch.Lineas) != 0 {
		t.Errorf("el chunk nuevo no está vacío: %d instrucciones, %d constantes, %d líneas",
			len(ch.Instrucciones), len(ch.Constantes), len(ch.Lineas))
	}
}

func TestEscribirAcumulaInstruccionesYLineas(t *testing.T) {
	ch := NuevoChunk()

	ch.Escribir(uint8(OpConstant), 1)
	ch.Escribir(7, 1)
	ch.Escribir(uint8(OpReturn), 3)

	esperadas := []uint8{uint8(OpConstant), 7, uint8(OpReturn)}
	for i, esperada := range esperadas {
		if ch.Instrucciones[i] != esperada {
			t.Errorf("instrucción %d: se obtuvo %d y se esperaba %d", i, ch.Instrucciones[i], esperada)
		}
	}

	if len(ch.Lineas) != len(ch.Instrucciones) {
		t.Fatalf("hay %d líneas para %d instrucciones", len(ch.Lineas), len(ch.Instrucciones))
	}

	for i, esperada := range []int{1, 1, 3} {
		if ch.Lineas[i] != esperada {
			t.Errorf("línea %d: se obtuvo %d y se esperaba %d", i, ch.Lineas[i], esperada)
		}
	}
}

func TestAgregarConstanteDevuelveElIndice(t *testing.T) {
	ch := NuevoChunk()

	casos := []any{float64(1), "hola", true, nil}
	for i, valor := range casos {
		if indice := ch.AgregarConstante(valor); indice != i {
			t.Errorf("la constante %v quedó en el índice %d y se esperaba %d", valor, indice, i)
		}
	}

	if len(ch.Constantes) != len(casos) {
		t.Fatalf("hay %d constantes y se esperaban %d", len(ch.Constantes), len(casos))
	}

	for i, esperada := range casos {
		if ch.Constantes[i] != esperada {
			t.Errorf("constante %d: se obtuvo %v y se esperaba %v", i, ch.Constantes[i], esperada)
		}
	}
}
