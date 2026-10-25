package compilador

import (
	"fmt"
	"strings"
)

// Funcion representa la estructura estática de una función compilada.
type Funcion struct {
	Aridad       int
	UpvalueCount int
	Chunk        *Chunk
	Nombre       string
}

// Crea una nueva función vacía con aridad cero y un chunk vacío.
func NuevaFuncion() *Funcion {
	return &Funcion{
		Aridad:       0,
		UpvalueCount: 0,
		Chunk:        NuevoChunk(),
		Nombre:       "",
	}
}

// Devuelve una representación legible de la función, mostrando su nombre o "<script>" si es anónima.
func (f *Funcion) String() string {
	if f.Nombre == "" {
		return "<script>"
	}

	return fmt.Sprintf("<fn %s>", f.Nombre)
}

// Nombre del método que Lox usa como constructor.
const NombreDelInicializador = "init"

// Clase representa una clase de Lox en tiempo de ejecución.
type Clase struct {
	Nombre  string
	Metodos map[string]*Closure
}

// Crea una clase sin métodos.
func NuevaClase(nombre string) *Clase {
	return &Clase{Nombre: nombre, Metodos: make(map[string]*Closure)}
}

// Devuelve la representación de la clase tal como la muestra Lox.
func (c *Clase) String() string {
	return c.Nombre
}

// Instancia representa una instancia de una clase con sus propiedades.
type Instancia struct {
	Clase       *Clase
	Propiedades map[string]any
}

// Crea una instancia vacía de la clase indicada.
func NuevaInstancia(clase *Clase) *Instancia {
	return &Instancia{Clase: clase, Propiedades: make(map[string]any)}
}

// Devuelve la representación de la instancia tal como la muestra Lox.
func (i *Instancia) String() string {
	return "instancia de " + i.Clase.Nombre
}

// MetodoLigado es un método junto con la instancia sobre la que se va a correr.
type MetodoLigado struct {
	Receptor *Instancia
	Metodo   *Closure
}

// Devuelve la representación del método ligado tal como la muestra Lox.
func (m *MetodoLigado) String() string {
	return m.Metodo.String()
}

// ObjUpvalue representa una variable capturada en tiempo de ejecución.
type ObjUpvalue struct {
	Locacion int
	Cerrado  any
}

// Closure envuelve una Función en tiempo de ejecución para capturar su entorno.
type Closure struct {
	Funcion  *Funcion
	Upvalues []*ObjUpvalue
}

// Crea un nuevo closure a partir de una función dada.
func NuevoClosure(f *Funcion) *Closure {
	return &Closure{
		Funcion:  f,
		Upvalues: make([]*ObjUpvalue, f.UpvalueCount),
	}
}

// Devuelve la representación de la función que envuelve el closure.
func (c *Closure) String() string {
	return c.Funcion.String()
}

// Array representa un array de Lox en tiempo de ejecución.
type Array struct {
	Elementos []any
}

// Crea un array con los elementos indicados.
func NuevoArray(elementos []any) *Array {
	return &Array{Elementos: elementos}
}

// Devuelve la representación del array tal como la muestra Lox.
func (a *Array) String() string {
	partes := make([]string, 0, len(a.Elementos))
	for _, elemento := range a.Elementos {
		if elemento == nil {
			partes = append(partes, "nil")
			continue
		}

		partes = append(partes, fmt.Sprintf("%v", elemento))
	}

	return "[" + strings.Join(partes, ", ") + "]"
}
