package interprete

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// ========================
// Intérprete Principal
// ========================

type Interprete struct {
	Salida  io.Writer
	entorno *Entorno
}

func NuevoInterprete(salida io.Writer) *Interprete {
	return &Interprete{
		Salida:  salida,
		entorno: NuevoEntorno(nil),
	}
}

// Interpretar recorre las sentencias y devuelve el primer error de runtime que encuentre.
func (i *Interprete) Interpretar(sentencias []sintaxis.Stmt) error {
	for _, sentencia := range sentencias {
		if err := i.Ejecutar(sentencia); err != nil {
			return err
		}
	}
	return nil
}

// Ejecutar procesa una sentencia usando un switch de tipo.
func (i *Interprete) Ejecutar(stmt sintaxis.Stmt) error {
	switch s := stmt.(type) {
	case *sintaxis.ExpressionStmt:
		_, err := i.Evaluar(s.Expression)
		return err

	case *sintaxis.Print:
		valor, err := i.Evaluar(s.Value)
		if err != nil {
			return err
		}

		if valor == nil {
			fmt.Fprintln(i.Salida, "nil")
		} else {
			fmt.Fprintln(i.Salida, valor)
		}

		return nil
	case *sintaxis.Var:
		var valor any
		var err error
		if s.Initializer != nil {
			valor, err = i.Evaluar(s.Initializer)
			if err != nil {
				return err
			}
		}

		i.entorno.Definir(s.Name.Lexema, valor)
		return nil
	case *sintaxis.Block:
		return i.ejecutarBloque(s.Statements, NuevoEntorno(i.entorno))
	}
	
	return nil
}

// ejecutarBloque ejecuta una lista de sentencias bajo un entorno específico (scope local).
func (i *Interprete) ejecutarBloque(sentencias []sintaxis.Stmt, entornoLocal *Entorno) error {
	entornoAnterior := i.entorno
	defer func() {
		i.entorno = entornoAnterior
	}()

	i.entorno = entornoLocal
	for _, sentencia := range sentencias {
		if err := i.Ejecutar(sentencia); err != nil {
			return err
		}
	}

	return nil
}

// Evaluar procesa una expresión usando un switch de tipo y devuelve su valor resultante.
func (i *Interprete) Evaluar(expr sintaxis.Expr) (any, error) {
	switch e := expr.(type) {
	case *sintaxis.Literal:
		return e.Value, nil

	case *sintaxis.Grouping:
		return i.Evaluar(e.Expression)

	case *sintaxis.Unary:
		derecha, err := i.Evaluar(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Tipo {
		case token.BANG:
			return !i.esVerdadero(derecha), nil
		case token.MINUS:
			if err := i.checkNumeroOperando(e.Operator, derecha); err != nil {
				return nil, err
			}
			return -derecha.(float64), nil
		}

	case *sintaxis.Binary:
		izquierda, err := i.Evaluar(e.Left)
		if err != nil {
			return nil, err
		}

		derecha, err := i.Evaluar(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.Tipo {
		case token.GREATER:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			return izquierda.(float64) > derecha.(float64), nil
		case token.GREATER_EQUAL:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			return izquierda.(float64) >= derecha.(float64), nil
		case token.LESS:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			return izquierda.(float64) < derecha.(float64), nil
		case token.LESS_EQUAL:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			return izquierda.(float64) <= derecha.(float64), nil
		case token.MINUS:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			return izquierda.(float64) - derecha.(float64), nil
		case token.SLASH:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			if derecha.(float64) == 0 {
				return nil, NewErrorRuntime(e.Operator, "División por cero.")
			}
			return izquierda.(float64) / derecha.(float64), nil
		case token.STAR:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			return izquierda.(float64) * derecha.(float64), nil
		case token.PLUS:
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

			return nil, NewErrorRuntime(e.Operator, "Los operandos deben ser dos números o dos cadenas.")
		case token.BANG_EQUAL:
			return !i.esIgual(izquierda, derecha), nil
		case token.EQUAL_EQUAL:
			return i.esIgual(izquierda, derecha), nil
		}
		
	case *sintaxis.Variable:
		return i.entorno.Obtener(e.Name)

	case *sintaxis.Assign:
		valor, err := i.Evaluar(e.Value)
		if err != nil {
			return nil, err
		}
		
		if err := i.entorno.Asignar(e.Name, valor); err != nil {
			return nil, err
		}
		
		return valor, nil
	}

	return nil, nil
}

// ========================
// Reglas Semánticas de Lox
// ========================

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
