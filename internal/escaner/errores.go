package escaner

import "fmt"

// Representa una falla léxica asociada a una línea del código fuente.
type Error struct {
	Linea   int
	Mensaje string
}

// Implementa la interfaz error de Go.
func (e *Error) Error() string {
	return fmt.Sprintf("[línea %d] Error de escaneo: %s", e.Linea, e.Mensaje)
}
