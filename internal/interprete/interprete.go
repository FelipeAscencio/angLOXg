package interprete

import (
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

type Interprete struct{}

func NuevoInterprete() *Interprete {
	return &Interprete{}
}

// Punto de entrada para procesar una expresión.
func (i *Interprete) Evaluar(expr sintaxis.Expresion) (any, error) {
	return expr.Aceptar(i)
}

// ===================================
// Implementación del VisitorExpresion
// ===================================

func (i *Interprete) VisitarLiteral(expr *sintaxis.Literal) (any, error) {
	return expr.Valor, nil
}

func (i *Interprete) VisitarAgrupacion(expr *sintaxis.Agrupacion) (any, error) {
	return i.Evaluar(expr.Expresion)
}

func (i *Interprete) VisitarUnaria(expr *sintaxis.Unaria) (any, error) {
	derecha, err := i.Evaluar(expr.Derecha)
	if err != nil {
		return nil, err
	}

	switch expr.Operador.Tipo {
	case token.BANG:
		return !i.esVerdadero(derecha), nil
	case token.MENOS:
		if err := i.checkNumeroOperando(expr.Operador, derecha); err != nil {
			return nil, err
		}

		return -derecha.(float64), nil
	}

	return nil, nil
}

// Evaluación post-order: Evaluamos hijos de izquierda a derecha.
func (i *Interprete) VisitarBinaria(expr *sintaxis.Binaria) (any, error) {
	izquierda, err := i.Evaluar(expr.Izquierda)
	if err != nil {
		return nil, err
	}

	derecha, err := i.Evaluar(expr.Derecha)
	if err != nil {
		return nil, err
	}

	switch expr.Operador.Tipo {
	case token.MAYOR:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		return izquierda.(float64) > derecha.(float64), nil
	case token.MAYOR_IGUAL:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		return izquierda.(float64) >= derecha.(float64), nil
	case token.MENOR:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		return izquierda.(float64) < derecha.(float64), nil
	case token.MENOR_IGUAL:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		return izquierda.(float64) <= derecha.(float64), nil
	case token.MENOS:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		return izquierda.(float64) - derecha.(float64), nil
	case token.SLASH:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		if derecha.(float64) == 0 {
			return nil, NewErrorRuntime(expr.Operador, "División por cero.")
		}

		return izquierda.(float64) / derecha.(float64), nil
	case token.ESTRELLA:
		if err := i.checkNumerosOperandos(expr.Operador, izquierda, derecha); err != nil {
			return nil, err
		}

		return izquierda.(float64) * derecha.(float64), nil
	case token.SUMA:
		if izq, ok := izquierda.(float64); ok {
			if der, ok := derecha.(float64); ok {
				return izq + der, nil
			}
		}

		if izq, ok := izquierda.(string); ok {
			if der, ok := derecha.(string); ok {
				return izq + der, nil
			}
		}

		return nil, NewErrorRuntime(expr.Operador, "Los operandos deben ser dos números o dos cadenas.")
	case token.BANG_IGUAL:
		return !i.esIgual(izquierda, derecha), nil
	case token.IGUAL_IGUAL:
		return i.esIgual(izquierda, derecha), nil
	}

	return nil, nil
}

// ========================
// Reglas Semánticas de Lox
// ========================

// Característica de Lox: "false" y "nil" son falsos, todo lo demás es verdadero.
func (i *Interprete) esVerdadero(objeto any) bool {
	if objeto == nil {
		return false
	}

	if b, ok := objeto.(bool); ok {
		return b
	}

	return true
}

func (i *Interprete) esIgual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

func (i *Interprete) checkNumeroOperando(operador token.Token, operando any) error {
	if _, ok := operando.(float64); ok {
		return nil
	}

	return NewErrorRuntime(operador, "El operando debe ser un número.")
}

func (i *Interprete) checkNumerosOperandos(operador token.Token, izq, der any) error {
	_, ok1 := izq.(float64)
	_, ok2 := der.(float64)
	if ok1 && ok2 {
		return nil
	}
	
	return NewErrorRuntime(operador, "Los operandos deben ser números.")
}
