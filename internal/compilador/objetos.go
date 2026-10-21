package compilador

import "fmt"

// Funcion representa la estructura estática de una función compilada.
type Funcion struct {
	Aridad int
	Chunk  *Chunk
	Nombre string
}

// Crea una nueva función vacía con aridad cero y un chunk vacío.
func NuevaFuncion() *Funcion {
	return &Funcion{
		Aridad: 0,
		Chunk:  NuevoChunk(),
		Nombre: "",
	}
}

// Devuelve una representación legible de la función, mostrando su nombre o "<script>" si es anónima.
func (f *Funcion) String() string {
	if f.Nombre == "" {
		return "<script>"
	}
	return fmt.Sprintf("<fn %s>", f.Nombre)
}

// Closure envuelve una Función en tiempo de ejecución para capturar su entorno.
type Closure struct {
	Funcion *Funcion
	// Me falta implementar la captura de variables locales para cierres (closures).
}

// Crea un nuevo closure a partir de una función dada.
func NuevoClosure(f *Funcion) *Closure {
	return &Closure{
		Funcion: f,
	}
}
