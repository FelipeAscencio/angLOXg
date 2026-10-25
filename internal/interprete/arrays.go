package interprete

import (
	"fmt"
	"strings"

	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Representa un array de Lox en tiempo de ejecución.
type LoxArray struct {
	Elementos []any
}

// Crea un array con los elementos indicados.
func NuevoArray(elementos []any) *LoxArray {
	return &LoxArray{Elementos: elementos}
}

// Devuelve el elemento que está en una posición del array.
func (a *LoxArray) Obtener(corchete token.Token, indice any) (any, error) {
	posicion, err := posicionValida(corchete, indice)
	if err != nil {
		return nil, err
	}

	if posicion < 0 || posicion >= len(a.Elementos) {
		return nil, NewErrorRuntime(corchete, "Índice fuera de rango.")
	}

	return a.Elementos[posicion], nil
}

// Devuelve la representación del array tal como la muestra Lox.
func (a *LoxArray) String() string {
	partes := make([]string, 0, len(a.Elementos))
	for _, elemento := range a.Elementos {
		partes = append(partes, formatearValor(elemento))
	}

	return "[" + strings.Join(partes, ", ") + "]"
}

// Convierte un índice de Lox en una posición del array.
//
// Los números de Lox son todos de punto flotante, así que hay que descartar
// los que tengan parte decimal: "a[1.5]" no apunta a ninguna posición.
func posicionValida(corchete token.Token, indice any) (int, error) {
	numero, esNumero := indice.(float64)
	if !esNumero || numero != float64(int(numero)) {
		return 0, NewErrorRuntime(corchete, "El índice debe ser un número entero.")
	}

	return int(numero), nil
}

// Escribe un valor de Lox como lo espera el lenguaje, con el nulo como "nil".
func formatearValor(valor any) string {
	if valor == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", valor)
}
