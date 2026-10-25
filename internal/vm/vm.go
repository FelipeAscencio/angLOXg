package vm

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/compilador"
)

// CallFrame representa una función en ejecución junto con su estado.
type CallFrame struct {
	closure *compilador.Closure
	ip      int // Puntero de instrucción (Instruction Pointer).
	slots   int
}

// Representa la máquina virtual basada en pila que ejecuta el bytecode.
type VM struct {
	frames           []*CallFrame
	pila             []any
	globales         map[string]any
	upvaluesAbiertos []*compilador.ObjUpvalue
	salida           io.Writer
}

// Inicializa una nueva Máquina Virtual.
func Nueva(salida io.Writer) *VM {
	return &VM{
		frames:           make([]*CallFrame, 0, 64),
		pila:             make([]any, 0, 256),
		globales:         make(map[string]any),
		upvaluesAbiertos: make([]*compilador.ObjUpvalue, 0),
		salida:           salida,
	}
}

// Arranca el ciclo de la máquina virtual sobre un Chunk compilado.
func (v *VM) Ejecutar(funcion *compilador.Funcion) bool {
	closure := compilador.NuevoClosure(funcion)
	v.push(closure)
	v.frames = append(v.frames, &CallFrame{
		closure: closure,
		ip:      0,
		slots:   0,
	})

	return v.correr()
}

// Atrapa una variable local, o devuelve el upvalue si ya estaba atrapada.
func (v *VM) capturarUpvalue(local int) *compilador.ObjUpvalue {
	for _, upvalue := range v.upvaluesAbiertos {
		if upvalue.Locacion == local {
			return upvalue
		}
	}

	nuevoUpvalue := &compilador.ObjUpvalue{Locacion: local}
	v.upvaluesAbiertos = append(v.upvaluesAbiertos, nuevoUpvalue)
	return nuevoUpvalue
}

// Mueve los upvalues de la pila al valor "Cerrado" permanente y limpia la lista de abiertos.
func (v *VM) cerrarUpvalues(ultimo int) {
	nuevosAbiertos := make([]*compilador.ObjUpvalue, 0)
	for _, upvalue := range v.upvaluesAbiertos {
		if upvalue.Locacion >= ultimo {
			upvalue.Cerrado = v.pila[upvalue.Locacion]
			upvalue.Locacion = -1
		} else {
			nuevosAbiertos = append(nuevosAbiertos, upvalue)
		}
	}

	v.upvaluesAbiertos = nuevosAbiertos
}

// Ciclo principal (loop) que despacha las instrucciones.
func (v *VM) correr() bool {
	for {
		frame := v.frameActual()
		if frame.ip >= len(frame.closure.Funcion.Chunk.Instrucciones) {
			break
		}

		instruccion := compilador.OpCode(v.leerByte())
		switch instruccion {
		case compilador.OpConstant:
			constante := v.leerConstante()
			v.push(constante)

		case compilador.OpAdd:
			b := v.pop()
			a := v.pop()
			if numA, okA := a.(float64); okA {
				if numB, okB := b.(float64); okB {
					v.push(numA + numB)
				} else {
					fmt.Fprintln(v.salida, "Error en ejecución: Los operandos deben ser dos números o dos cadenas.")
					return false
				}
			} else if strA, okA := a.(string); okA {
				strB, okB := b.(string)
				if !okB {
					fmt.Fprintln(v.salida, "Error en ejecución: Los operandos deben ser dos números o dos cadenas.")
					return false
				}

				v.push(strA + strB)
			} else {
				fmt.Fprintln(v.salida, "Error en ejecución: Los operandos deben ser dos números o dos cadenas.")
				return false
			}

		case compilador.OpSubtract:
			a, b, ok := v.operandosNumericos()
			if !ok {
				return false
			}

			v.push(a - b)

		case compilador.OpMultiply:
			a, b, ok := v.operandosNumericos()
			if !ok {
				return false
			}

			v.push(a * b)

		case compilador.OpDivide:
			a, b, ok := v.operandosNumericos()
			if !ok {
				return false
			}

			if b == 0 {
				fmt.Fprintln(v.salida, "Error en ejecución: División por cero.")
				return false
			}

			v.push(a / b)

		case compilador.OpModulo:
			a, b, ok := v.operandosNumericos()
			if !ok {
				return false
			}

			if b == 0 {
				fmt.Fprintln(v.salida, "Error en ejecución: Módulo por cero.")
				return false
			}

			v.push(float64(int(a) % int(b)))

		case compilador.OpNegate:
			valor, esNumero := v.pop().(float64)
			if !esNumero {
				fmt.Fprintln(v.salida, "Error en ejecución: El operando debe ser un número.")
				return false
			}

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
			a, b, ok := v.operandosNumericos()
			if !ok {
				return false
			}

			v.push(a > b)

		case compilador.OpLess:
			a, b, ok := v.operandosNumericos()
			if !ok {
				return false
			}

			v.push(a < b)

		case compilador.OpNot:
			valor := v.pop()
			v.push(esFalso(valor))

		case compilador.OpPop:
			v.pop()

		case compilador.OpPrint:
			fmt.Fprintln(v.salida, formatearValor(v.pop()))

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

			v.globales[nombre] = v.peek(0)

		case compilador.OpGetLocal:
			slot := v.leerByte()
			v.push(v.pila[v.frameActual().slots+int(slot)])

		case compilador.OpSetLocal:
			slot := v.leerByte()
			v.pila[v.frameActual().slots+int(slot)] = v.peek(0)

		case compilador.OpGetUpvalue:
			slot := v.leerByte()
			upvalue := v.frameActual().closure.Upvalues[slot]
			if upvalue.Locacion != -1 {
				v.push(v.pila[upvalue.Locacion])
			} else {
				v.push(upvalue.Cerrado)
			}

		case compilador.OpSetUpvalue:
			slot := v.leerByte()
			upvalue := v.frameActual().closure.Upvalues[slot]
			if upvalue.Locacion != -1 {
				v.pila[upvalue.Locacion] = v.peek(0)
			} else {
				upvalue.Cerrado = v.peek(0)
			}

		case compilador.OpCloseUpvalue:
			v.cerrarUpvalues(len(v.pila) - 1)
			v.pop()

		case compilador.OpJump:
			offset := int(v.leerShort())
			v.frameActual().ip += offset

		case compilador.OpJumpIfFalse:
			offset := int(v.leerShort())
			if esFalso(v.peek(0)) {
				v.frameActual().ip += offset
			}

		case compilador.OpLoop:
			offset := int(v.leerShort())
			v.frameActual().ip -= offset

		case compilador.OpClosure:
			funcion := v.leerConstante().(*compilador.Funcion)
			closure := compilador.NuevoClosure(funcion)
			v.push(closure)
			for i := 0; i < funcion.UpvalueCount; i++ {
				esLocal := v.leerByte()
				indice := v.leerByte()
				if esLocal == 1 {
					closure.Upvalues[i] = v.capturarUpvalue(v.frameActual().slots + int(indice))
				} else {
					closure.Upvalues[i] = v.frameActual().closure.Upvalues[indice]
				}
			}

		case compilador.OpCall:
			cantidadArgs := int(v.leerByte())
			callee := v.peek(cantidadArgs)
			closure, ok := callee.(*compilador.Closure)
			if !ok {
				fmt.Fprintln(v.salida, "Error en ejecución: Solo se pueden llamar funciones.")
				return false
			}

			if cantidadArgs != closure.Funcion.Aridad {
				fmt.Fprintf(v.salida, "Error en ejecución: Se esperaban %d argumentos pero se obtuvieron %d.\n", closure.Funcion.Aridad, cantidadArgs)
				return false
			}

			v.frames = append(v.frames, &CallFrame{
				closure: closure,
				ip:      0,
				slots:   len(v.pila) - cantidadArgs - 1,
			})

		case compilador.OpReturn:
			resultado := v.pop()
			frameActual := v.frameActual()
			v.cerrarUpvalues(frameActual.slots)
			v.frames = v.frames[:len(v.frames)-1]
			if len(v.frames) == 0 {
				v.pop()
				return true
			}

			v.pila = v.pila[:frameActual.slots]
			v.push(resultado)
		}
	}

	return false
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

// ==================
// Helpers de Lectura
// ==================

func (v *VM) frameActual() *CallFrame {
	return v.frames[len(v.frames)-1]
}

func (v *VM) leerByte() uint8 {
	frame := v.frameActual()
	b := frame.closure.Funcion.Chunk.Instrucciones[frame.ip]
	frame.ip++
	return b
}

func (v *VM) leerShort() uint16 {
	frame := v.frameActual()
	frame.ip += 2
	instrucciones := frame.closure.Funcion.Chunk.Instrucciones
	return uint16(instrucciones[frame.ip-2])<<8 | uint16(instrucciones[frame.ip-1])
}

func (v *VM) leerConstante() any {
	return v.frameActual().closure.Funcion.Chunk.Constantes[v.leerByte()]
}

// Saca los dos operandos del tope de la pila y verifica que ambos sean números.
func (v *VM) operandosNumericos() (izquierda, derecha float64, ok bool) {
	valorDerecha := v.pop()
	valorIzquierda := v.pop()

	derecha, okDerecha := valorDerecha.(float64)
	izquierda, okIzquierda := valorIzquierda.(float64)
	if !okIzquierda || !okDerecha {
		fmt.Fprintln(v.salida, "Error en ejecución: Los operandos deben ser números.")
		return 0, 0, false
	}

	return izquierda, derecha, true
}

// Escribe un valor de Lox como lo espera el lenguaje, con el nulo como "nil".
func formatearValor(valor any) string {
	if valor == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", valor)
}
