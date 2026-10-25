package compilador

import (
	"fmt"
	"strings"
)

// operandosPorOpcode dice cuántos bytes de operando lleva cada instrucción.
var operandosPorOpcode = map[OpCode]int{
	OpNil:          0,
	OpTrue:         0,
	OpFalse:        0,
	OpEqual:        0,
	OpGreater:      0,
	OpLess:         0,
	OpNot:          0,
	OpAdd:          0,
	OpSubtract:     0,
	OpMultiply:     0,
	OpDivide:       0,
	OpNegate:       0,
	OpModulo:       0,
	OpPop:          0,
	OpPrint:        0,
	OpCloseUpvalue: 0,
	OpReturn:       0,

	OpConstant:     1,
	OpDefineGlobal: 1,
	OpGetGlobal:    1,
	OpSetGlobal:    1,
	OpGetLocal:     1,
	OpSetLocal:     1,
	OpGetUpvalue:   1,
	OpSetUpvalue:   1,
	OpCall:         1,
	OpClosure:      1,

	OpJump:        2,
	OpJumpIfFalse: 2,
	OpLoop:        2,
}

var nombresDeOpcode = map[OpCode]string{
	OpConstant: "OpConstant", OpNil: "OpNil", OpTrue: "OpTrue", OpFalse: "OpFalse",
	OpEqual: "OpEqual", OpGreater: "OpGreater", OpLess: "OpLess", OpNot: "OpNot",
	OpAdd: "OpAdd", OpSubtract: "OpSubtract", OpMultiply: "OpMultiply",
	OpDivide: "OpDivide", OpNegate: "OpNegate", OpModulo: "OpModulo",
	OpPop: "OpPop", OpPrint: "OpPrint",
	OpDefineGlobal: "OpDefineGlobal", OpGetGlobal: "OpGetGlobal", OpSetGlobal: "OpSetGlobal",
	OpGetLocal: "OpGetLocal", OpSetLocal: "OpSetLocal",
	OpJump: "OpJump", OpJumpIfFalse: "OpJumpIfFalse", OpLoop: "OpLoop",
	OpCall: "OpCall", OpClosure: "OpClosure",
	OpGetUpvalue: "OpGetUpvalue", OpSetUpvalue: "OpSetUpvalue", OpCloseUpvalue: "OpCloseUpvalue",
	OpReturn: "OpReturn",
}

// desensamblar recorre el bytecode de una función y devuelve una instrucción
func desensamblar(fn *Funcion) []string {
	instrucciones := []string{}
	chunk := fn.Chunk

	for i := 0; i < len(chunk.Instrucciones); {
		op := OpCode(chunk.Instrucciones[i])

		nombre, conocido := nombresDeOpcode[op]
		if !conocido {
			instrucciones = append(instrucciones, fmt.Sprintf("DESCONOCIDO(%d)", op))
			break
		}

		operandos := operandosPorOpcode[op]
		renglon := nombre
		for k := 1; k <= operandos; k++ {
			renglon += fmt.Sprintf(" %d", chunk.Instrucciones[i+k])
		}

		if indexaConstante(op) {
			renglon += fmt.Sprintf(" -> %s", valorDeConstante(chunk.Constantes[chunk.Instrucciones[i+1]]))
		}

		instrucciones = append(instrucciones, renglon)
		i += 1 + operandos

		if op == OpClosure {
			anidada := chunk.Constantes[chunk.Instrucciones[i-1]].(*Funcion)
			i += anidada.UpvalueCount * 2
		}
	}

	return instrucciones
}

// indexaConstante indica si el operando de la instrucción es un índice al pool
func indexaConstante(op OpCode) bool {
	switch op {
	case OpConstant, OpDefineGlobal, OpGetGlobal, OpSetGlobal, OpClosure:
		return true
	default:
		return false
	}
}

// valorDeConstante escribe una constante de forma comparable en un test.
func valorDeConstante(valor any) string {
	switch v := valor.(type) {
	case *Funcion:
		return v.String()
	case string:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// funcionConstante busca en el pool la función compilada con ese nombre, para
func funcionConstante(fn *Funcion, nombre string) *Funcion {
	for _, constante := range fn.Chunk.Constantes {
		if anidada, esFuncion := constante.(*Funcion); esFuncion && anidada.Nombre == nombre {
			return anidada
		}
	}

	return nil
}

// soloOpcodes se queda con el nombre de cada instrucción y descarta operandos,
func soloOpcodes(instrucciones []string) []string {
	nombres := make([]string, 0, len(instrucciones))
	for _, instruccion := range instrucciones {
		nombres = append(nombres, strings.Fields(instruccion)[0])
	}

	return nombres
}
