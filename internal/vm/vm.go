package vm

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/compilador"
)

// Representa la máquina virtual basada en pila que ejecuta el bytecode.
type VM struct {
	chunk    *compilador.Chunk
	ip       int	// Puntero de instrucción (Instruction Pointer).
	pila     []any
	globales map[string]any
	salida   io.Writer
}

// Inicializa una nueva Máquina Virtual.
func Nueva(salida io.Writer) *VM {
	return &VM{
		pila:     make([]any, 0, 256),
		globales: make(map[string]any),
		salida:   salida,
	}
}

// Arranca el ciclo de la máquina virtual sobre un Chunk compilado.
func (v *VM) Ejecutar(chunk *compilador.Chunk) bool {
	v.chunk = chunk
	v.ip = 0

	return v.correr()
}

// Ciclo principal (loop) que despacha las instrucciones.
func (v *VM) correr() bool {
	for v.ip < len(v.chunk.Instrucciones) {
		instruccion := compilador.OpCode(v.leerByte())

		switch instruccion {
		case compilador.OpConstant:
			constante := v.leerConstante()
			v.push(constante)

		case compilador.OpAdd:
			b := v.pop().(float64)
			a := v.pop().(float64)
			v.push(a + b)

		case compilador.OpSubtract:
			b := v.pop().(float64)
			a := v.pop().(float64)
			v.push(a - b)

		case compilador.OpMultiply:
			b := v.pop().(float64)
			a := v.pop().(float64)
			v.push(a * b)

		case compilador.OpDivide:
			b := v.pop().(float64)
			a := v.pop().(float64)
			v.push(a / b)

		case compilador.OpNegate:
			// Extraemos el tope, lo negamos y lo volvemos a meter
			valor := v.pop().(float64)
			v.push(-valor)

		case compilador.OpNil:
			v.push(nil)

		case compilador.OpTrue:
			v.push(true)

		case compilador.OpFalse:
			v.push(false)

		case compilador.OpEqual:
			b := v.pop()
			a := v.pop()
			v.push(sonIguales(a, b))

		case compilador.OpGreater:
			b := v.pop().(float64)
			a := v.pop().(float64)
			v.push(a > b)

		case compilador.OpLess:
			b := v.pop().(float64)
			a := v.pop().(float64)
			v.push(a < b)

		case compilador.OpNot:
			valor := v.pop()
			v.push(esFalso(valor))

		case compilador.OpPop:
			v.pop()

		case compilador.OpPrint:
			fmt.Fprintln(v.salida, v.pop())

		case compilador.OpDefineGlobal:
			nombre := v.leerConstante().(string)
			v.globales[nombre] = v.pop()

		case compilador.OpGetGlobal:
			nombre := v.leerConstante().(string)
			valor, existe := v.globales[nombre]
			if !existe {
				fmt.Fprintf(v.salida, "Error en ejecución: Variable no definida '%s'.\n", nombre)
				return false
			}

			v.push(valor)

		case compilador.OpSetGlobal:
			nombre := v.leerConstante().(string)
			if _, existe := v.globales[nombre]; !existe {
				fmt.Fprintf(v.salida, "Error en ejecución: Variable no definida '%s'.\n", nombre)
				return false
			}

			// La asignación deja el valor en la pila, por eso usamos peek(0) en vez de pop().
			v.globales[nombre] = v.peek(0)

		case compilador.OpReturn:
			return true
		}
	}

	return false
}

// Consume el byte actual y avanza el IP.
func (v *VM) leerByte() uint8 {
	b := v.chunk.Instrucciones[v.ip]
	v.ip++
	return b
}

// Lee el siguiente byte como índice y obtiene la constante del pool.
func (v *VM) leerConstante() any {
	indice := v.leerByte()
	return v.chunk.Constantes[indice]
}

// Agrega un valor al tope de la pila.
func (v *VM) push(valor any) {
	v.pila = append(v.pila, valor)
}

// Retira y devuelve el valor en el tope de la pila.
func (v *VM) pop() any {
	if len(v.pila) == 0 {
		return nil 
	}

	tope := len(v.pila) - 1
	valor := v.pila[tope]
	v.pila = v.pila[:tope]
	return valor
}

// Lox sigue la convención de Ruby: false y nil son falsos, todo lo demás es verdadero.
func esFalso(valor any) bool {
	if valor == nil {
		return true
	}
	if b, ok := valor.(bool); ok {
		return !b
	}
	return false
}

// Comprueba igualdad tolerando distintos tipos.
func sonIguales(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

// Lee un valor de la pila sin sacarlo. La distancia es desde el tope (0 es el tope actual).
func (v *VM) peek(distancia int) any {
	return v.pila[len(v.pila)-1-distancia]
}
