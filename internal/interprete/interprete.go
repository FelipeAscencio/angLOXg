package interprete

import (
	"fmt"
	"io"

	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// ====================
// Funciones y closures
// ====================

type LoxCallable interface {
	Aridad() int
	Llamar(intp *Interprete, argumentos []any) (any, error)
}

type LoxFunction struct {
	Declaracion *sintaxis.Function
	Closure     *Entorno
}

func (f *LoxFunction) Aridad() int {
	return len(f.Declaracion.Params)
}

func (f *LoxFunction) Llamar(intp *Interprete, argumentos []any) (any, error) {
	entorno := NuevoEntorno(f.Closure)
	for i, param := range f.Declaracion.Params {
		entorno.Definir(param.Lexema, argumentos[i])
	}

	err := intp.ejecutarBloque(f.Declaracion.Body, entorno)
	if ret, ok := err.(*ValorRetorno); ok {
		return ret.Valor, nil
	}
	
	return nil, err
}

type ValorRetorno struct {
	Valor any
}

func (v *ValorRetorno) Error() string {
	return "retorno"
}

// ==========
// Intérprete
// ==========

type Interprete struct {
	Salida  io.Writer
	entorno *Entorno
	globales *Entorno
	locales map[sintaxis.Expr]int
}

func NuevoInterprete(salida io.Writer) *Interprete {
	entornoGlobal := NuevoEntorno(nil)
	return &Interprete{
		Salida:   salida,
		entorno:  entornoGlobal,
		globales: entornoGlobal,
		locales:  make(map[sintaxis.Expr]int),
	}
}

// ResolverLocal es llamado por el Analizador Semántico para guardar la distancia de una variable.
func (i *Interprete) ResolverLocal(expr sintaxis.Expr, profundidad int) {
	i.locales[expr] = profundidad
}

// buscarVariable usa la distancia estática si existe.
func (i *Interprete) buscarVariable(nombre token.Token, expr sintaxis.Expr) (any, error) {
	if distancia, ok := i.locales[expr]; ok {
		return i.entorno.ObtenerEn(distancia, nombre.Lexema), nil
	}
	
	return i.globales.Obtener(nombre)
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

	case *sintaxis.If:
		condicion, err := i.Evaluar(s.Condition)
		if err != nil {
			return err
		}

		if i.esVerdadero(condicion) {
			return i.Ejecutar(s.Then)
		} else if s.Else != nil {
			return i.Ejecutar(s.Else)
		}

		return nil

	case *sintaxis.While:
		for {
			condicion, err := i.Evaluar(s.Condition)
			if err != nil {
				return err
			}

			if !i.esVerdadero(condicion) {
				break
			}

			if err := i.Ejecutar(s.Body); err != nil {
				return err
			}
		}

		return nil

	case *sintaxis.Function:
		funcion := &LoxFunction{
			Declaracion: s,
			Closure:     i.entorno,
		}

		i.entorno.Definir(s.Name.Lexema, funcion)
		return nil

	case *sintaxis.Return:
		var valor any
		var err error
		if s.Value != nil {
			valor, err = i.Evaluar(s.Value)
			if err != nil {
				return err
			}
		}

		return &ValorRetorno{Valor: valor}
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
		case token.PERCENT:
			if err := i.checkNumerosOperandos(e.Operator, izquierda, derecha); err != nil {
				return nil, err
			}
			der := derecha.(float64)
			if der == 0 {
				return nil, NewErrorRuntime(e.Operator, "Módulo por cero.")
			}
			izq := izquierda.(float64)
			return float64(int64(izq) % int64(der)), nil
		}

	case *sintaxis.Variable:
		return i.buscarVariable(e.Name, e)

	case *sintaxis.Assign:
		valor, err := i.Evaluar(e.Value)
		if err != nil {
			return nil, err
		}
		
		if distancia, ok := i.locales[e]; ok {
			i.entorno.AsignarEn(distancia, e.Name.Lexema, valor)
		} else {
			if err := i.globales.Asignar(e.Name, valor); err != nil {
				return nil, err
			}
		}
		return valor, nil

	case *sintaxis.Logical:
		izquierda, err := i.Evaluar(e.Left)
		if err != nil {
			return nil, err
		}

		if e.Operator.Tipo == token.OR {
			if i.esVerdadero(izquierda) {
				return izquierda, nil
			}
			
		} else if e.Operator.Tipo == token.AND {
			if !i.esVerdadero(izquierda) {
				return izquierda, nil
			}
		}

		return i.Evaluar(e.Right)

	case *sintaxis.Call:
		callee, err := i.Evaluar(e.Callee)
		if err != nil {
			return nil, err
		}

		var argumentos []any
		for _, argExpr := range e.Arguments {
			arg, err := i.Evaluar(argExpr)
			if err != nil {
				return nil, err
			}
			argumentos = append(argumentos, arg)
		}

		funcion, ok := callee.(LoxCallable)
		if !ok {
			return nil, NewErrorRuntime(e.Paren, "Solo se pueden llamar funciones y clases.")
		}

		if len(argumentos) != funcion.Aridad() {
			mensaje := fmt.Sprintf("Se esperaban %d argumentos pero se obtuvieron %d.", funcion.Aridad(), len(argumentos))
			return nil, NewErrorRuntime(e.Paren, mensaje)
		}

		return funcion.Llamar(i, argumentos)
	}

	return nil, nil
}

// ========================
// Reglas semánticas de Lox
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
