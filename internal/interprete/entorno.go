package interprete

import (
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Administra el estado y las variables de un ámbito (scope) particular.
type Entorno struct {
	valores  map[string]any
	Ancestro *Entorno
}

// Crea una nueva instancia de entorno, opcionalmente vinculada a un ancestro.
func NuevoEntorno(ancestro *Entorno) *Entorno {
	return &Entorno{
		valores:  make(map[string]any),
		Ancestro: ancestro,
	}
}

// Define una variable en el entorno actual.
func (e *Entorno) Definir(nombre string, valor any) {
	e.valores[nombre] = valor
}

// Busca y retorna una variable, recorriendo recursivamente los ancestros si es necesario.
func (e *Entorno) Obtener(nombre token.Token) (any, error) {
	if valor, ok := e.valores[nombre.Lexema]; ok {
		return valor, nil
	}

	if e.Ancestro != nil {
		return e.Ancestro.Obtener(nombre)
	}

	return nil, NewErrorRuntime(nombre, "Variable no definida '"+nombre.Lexema+"'.")
}

// Modifica una variable existente buscando en la jerarquía de entornos.
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

func (e *Entorno) ancestro(distancia int) *Entorno {
	actual := e
	for i := 0; i < distancia; i++ {
		actual = actual.Ancestro
	}
	
	return actual
}

// Obtiene una variable ubicada a una distancia estática de entornos.
func (e *Entorno) ObtenerEn(distancia int, nombre string) any {
	return e.ancestro(distancia).valores[nombre]
}

// Modifica una variable ubicada a una distancia estática de entornos.
func (e *Entorno) AsignarEn(distancia int, nombre string, valor any) {
	e.ancestro(distancia).valores[nombre] = valor
}
