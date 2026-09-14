package escaner

import "fmt"

// Error representa una falla léxica, atada a la línea donde se detectó.
type Error struct {
	Linea   int
	Mensaje string
}

// Error arma el mensaje de la falla.
//
// Se llama así, y no "Mensaje", porque es el método que Go exige para que el
// tipo sirva como error.
func (e *Error) Error() string {
	return fmt.Sprintf("[línea %d] Error de escaneo: %s", e.Linea, e.Mensaje)
}
