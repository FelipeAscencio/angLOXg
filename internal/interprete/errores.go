package interprete

import (
	"fmt"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Representa un error de ejecución en el intérprete de Lox.
type ErrorRuntime struct {
	Token   token.Token
	Mensaje string
}

func (e *ErrorRuntime) Error() string {
	return fmt.Sprintf("%s \n[línea %d]", e.Mensaje, e.Token.Linea)
}

func NewErrorRuntime(t token.Token, mensaje string) *ErrorRuntime {
	return &ErrorRuntime{Token: t, Mensaje: mensaje}
}
