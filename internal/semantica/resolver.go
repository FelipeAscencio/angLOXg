package semantica

import (
	"fmt"

	"github.com/FelipeAscencio/angLOXg/internal/interprete"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

type Resolver struct {
	intp   *interprete.Interprete
	scopes []map[string]bool
}

func Nuevo(intp *interprete.Interprete) *Resolver {
	return &Resolver{
		intp:   intp,
		scopes: make([]map[string]bool, 0),
	}
}

func (r *Resolver) Resolver(sentencias []sintaxis.Stmt) error {
	for _, s := range sentencias {
		if err := r.ResolverSentencia(s); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) ResolverSentencia(stmt sintaxis.Stmt) error {
	switch s := stmt.(type) {
	case *sintaxis.Block:
		r.beginScope()
		if err := r.Resolver(s.Statements); err != nil {
			return err
		}

		r.endScope()

	case *sintaxis.Var:
		r.declarar(s.Name)
		if s.Initializer != nil {
			if err := r.ResolverExpresion(s.Initializer); err != nil {
				return err
			}
		}

		r.definir(s.Name)

	case *sintaxis.Function:
		r.declarar(s.Name)
		r.definir(s.Name)
		if err := r.resolverFuncion(s); err != nil {
			return err
		}

	case *sintaxis.ExpressionStmt:
		return r.ResolverExpresion(s.Expression)

	case *sintaxis.If:
		if err := r.ResolverExpresion(s.Condition); err != nil {
			return err
		}

		if err := r.ResolverSentencia(s.Then); err != nil {
			return err
		}

		if s.Else != nil {
			if err := r.ResolverSentencia(s.Else); err != nil {
				return err
			}
		}

	case *sintaxis.Print:
		return r.ResolverExpresion(s.Value)

	case *sintaxis.Return:
		if s.Value != nil {
			return r.ResolverExpresion(s.Value)
		}

	case *sintaxis.While:
		if err := r.ResolverExpresion(s.Condition); err != nil {
			return err
		}

		return r.ResolverSentencia(s.Body)
	}

	return nil
}

func (r *Resolver) ResolverExpresion(expr sintaxis.Expr) error {
	switch e := expr.(type) {
	case *sintaxis.Variable:
		if len(r.scopes) > 0 {
			if definido, existe := r.scopes[len(r.scopes)-1][e.Name.Lexema]; existe && !definido {
				return fmt.Errorf("[línea %d] No se puede leer una variable local en su propio inicializador", e.Name.Linea)
			}
		}

		r.resolverLocal(e, e.Name)

	case *sintaxis.Assign:
		if err := r.ResolverExpresion(e.Value); err != nil {
			return err
		}

		r.resolverLocal(e, e.Name)

	case *sintaxis.Binary:
		if err := r.ResolverExpresion(e.Left); err != nil {
			return err
		}

		return r.ResolverExpresion(e.Right)

	case *sintaxis.Call:
		if err := r.ResolverExpresion(e.Callee); err != nil {
			return err
		}

		for _, arg := range e.Arguments {
			if err := r.ResolverExpresion(arg); err != nil {
				return err
			}
		}

	case *sintaxis.Grouping:
		return r.ResolverExpresion(e.Expression)

	case *sintaxis.Logical:
		if err := r.ResolverExpresion(e.Left); err != nil {
			return err
		}

		return r.ResolverExpresion(e.Right)

	case *sintaxis.Unary:
		return r.ResolverExpresion(e.Right)

	case *sintaxis.Literal:
		// No hay nada que resolver.
	}
	return nil
}

func (r *Resolver) resolverLocal(expr sintaxis.Expr, nombre token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][nombre.Lexema]; ok {
			r.intp.ResolverLocal(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolverFuncion(funcion *sintaxis.Function) error {
	r.beginScope()
	for _, param := range funcion.Params {
		r.declarar(param)
		r.definir(param)
	}

	if err := r.Resolver(funcion.Body); err != nil {
		return err
	}

	r.endScope()
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declarar(nombre token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[nombre.Lexema] = false
}

func (r *Resolver) definir(nombre token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][nombre.Lexema] = true
}
