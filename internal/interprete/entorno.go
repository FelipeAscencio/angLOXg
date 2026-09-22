package interprete

import (
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Entorno guarda el estado (variables) de un scope particular.
type Entorno struct {
	valores  map[string]any
	Ancestro *Entorno // Puntero al entorno padre.
}

// NuevoEntorno crea un entorno. Si su ancestro no es nil, será un sub-scope.
func NuevoEntorno(ancestro *Entorno) *Entorno {
	return &Entorno{
		valores:  make(map[string]any),
		Ancestro: ancestro,
	}
}

// Definir crea una variable en el entorno actual.
func (e *Entorno) Definir(nombre string, valor any) {
	e.valores[nombre] = valor
}

// Obtener busca una variable por su token. Si no la encuentra, sube recursivamente por el árbol de ancestros.
func (e *Entorno) Obtener(nombre token.Token) (any, error) {
	if valor, ok := e.valores[nombre.Lexema]; ok {
		return valor, nil
	}

	if e.Ancestro != nil {
		return e.Ancestro.Obtener(nombre)
	}

	return nil, NewErrorRuntime(nombre, "Variable no definida '"+nombre.Lexema+"'.")
}

// Asignar modifica una variable existente. A diferencia de Definir, lanza un error si la variable no existía previamente.
func (e *Entorno) Asignar(nombre token.Token, valor any) error {
	if _, ok := e.valores[nombre.Lexema]; ok {
		e.valores[nombre.Lexema] = valor
		return nil
	}

	if e.Ancestro != nil {
		return e.Ancestro.Asignar(nombre, valor)
	}

	return NewErrorRuntime(nombre, "Variable no definida '"+nombre.Lexema+"'.")
}
