package compilador

// Representa un bloque contiguo de bytecode compilado y sus datos asociados.
type Chunk struct {
	Instrucciones []uint8
	Constantes    []any
	Lineas        []int
}

// Inicializa y devuelve un chunk vacío.
func NuevoChunk() *Chunk {
	return &Chunk{
		Instrucciones: make([]uint8, 0),
		Constantes:    make([]any, 0),
		Lineas:        make([]int, 0),
	}
}

// Añade un byte (código de operación u operando) al final del chunk, registrando la línea de código fuente de la que proviene.
func (c *Chunk) Escribir(byteCode uint8, linea int) {
	c.Instrucciones = append(c.Instrucciones, byteCode)
	c.Lineas = append(c.Lineas, linea)
}

// Inserta un valor en el pool de constantes y devuelve su índice.
func (c *Chunk) AgregarConstante(valor any) int {
	c.Constantes = append(c.Constantes, valor)
	return len(c.Constantes) - 1
}
