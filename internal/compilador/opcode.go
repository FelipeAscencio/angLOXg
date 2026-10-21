package compilador

// Representa un código de operación de un solo byte para la Máquina Virtual.
type OpCode uint8

const (
	// Carga una constante del pool de constantes al tope de la pila.
	OpConstant OpCode = iota

	// Operaciones aritméticas y lógicas
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNegate

	// Finaliza la ejecución del chunk actual y retorna el control.
	OpReturn
)
