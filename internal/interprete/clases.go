package interprete

import (
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Representa una clase de Lox en tiempo de ejecución.
type LoxClass struct {
	Nombre  string
	Metodos map[string]*LoxFunction
}

// Busca un método de la clase por su nombre.
func (c *LoxClass) BuscarMetodo(nombre string) *LoxFunction {
	if metodo, existe := c.Metodos[nombre]; existe {
		return metodo
	}

	return nil
}

// Devuelve la cantidad de argumentos que acepta el constructor de la clase.
func (c *LoxClass) Aridad() int {
	return 0
}

// Crea una instancia nueva de la clase.
func (c *LoxClass) Llamar(intp *Interprete, argumentos []any) (any, error) {
	return NuevaInstancia(c), nil
}

// Devuelve la representación de la clase tal como la muestra Lox.
func (c *LoxClass) String() string {
	return c.Nombre
}

// Representa una instancia de una clase de Lox con sus propiedades.
type LoxInstance struct {
	Clase       *LoxClass
	propiedades map[string]any
}

// Crea una instancia vacía de la clase indicada.
func NuevaInstancia(clase *LoxClass) *LoxInstance {
	return &LoxInstance{
		Clase:       clase,
		propiedades: make(map[string]any),
	}
}

// Devuelve el valor de una propiedad de la instancia, o un método de su clase
// ligado a esta instancia.
func (i *LoxInstance) Obtener(nombre token.Token) (any, error) {
	if valor, existe := i.propiedades[nombre.Lexema]; existe {
		return valor, nil
	}

	if metodo := i.Clase.BuscarMetodo(nombre.Lexema); metodo != nil {
		return metodo.Ligar(i), nil
	}

	return nil, NewErrorRuntime(nombre, "Propiedad no definida '"+nombre.Lexema+"'.")
}

// Define o reemplaza el valor de una propiedad de la instancia.
func (i *LoxInstance) Asignar(nombre token.Token, valor any) {
	i.propiedades[nombre.Lexema] = valor
}

// Devuelve la representación de la instancia tal como la muestra Lox.
func (i *LoxInstance) String() string {
	return "instancia de " + i.Clase.Nombre
}
