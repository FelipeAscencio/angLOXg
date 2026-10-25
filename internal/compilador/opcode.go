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
	OpModulo

	// Manejo de sentencias
	OpPop
	OpPrint

	// Manejo de variables globales
	OpDefineGlobal
	OpGetGlobal
	OpSetGlobal

	// Manejo de variables locales
	OpGetLocal
	OpSetLocal

	// Control de flujo
	OpJump
	OpJumpIfFalse
	OpLoop

	// Llamadas a funciones y Closures
	OpCall
	OpClosure

	// Closures y Upvalues
	OpGetUpvalue
	OpSetUpvalue
	OpCloseUpvalue

	// Clases, métodos y propiedades
	OpClass
	OpMethod
	OpGetProperty
	OpSetProperty
	OpInherit
	OpGetSuper

	// Arrays
	OpArray
	OpGetIndex

	// Finaliza la ejecución del chunk actual y retorna el control.
	OpReturn
)
