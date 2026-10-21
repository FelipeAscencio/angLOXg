package compilador

// Representa un código de operación de un solo byte para la Máquina Virtual.
type OpCode uint8

const (
	// Carga una constante del pool de constantes al tope de la pila.
	OpConstant OpCode = iota

	// Valores literales directos (ahorran espacio en el pool de constantes)
	OpNil
	OpTrue
	OpFalse

	// Operaciones de comparación y lógica
	OpEqual
	OpGreater
	OpLess
	OpNot

	// Operaciones aritméticas
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNegate

	// Manejo de sentencias
	OpPop   // Extrae y descarta el tope de la pila
	OpPrint // Extrae e imprime el tope de la pila

	// Finaliza la ejecución del chunk actual y retorna el control.
	OpReturn
)
