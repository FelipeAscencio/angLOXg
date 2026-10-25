package compilador

import "fmt"

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

// ObjUpvalue representa una variable capturada en tiempo de ejecución[cite: 55].
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
