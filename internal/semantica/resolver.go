package semantica

import (
	"fmt"

	"github.com/FelipeAscencio/angLOXg/internal/interprete"
	"github.com/FelipeAscencio/angLOXg/internal/sintaxis"
	"github.com/FelipeAscencio/angLOXg/internal/token"
)

// Analizador semántico (Resolver) encargado de la resolución estática de variables y ámbitos.
type Resolver struct {
	intp   *interprete.Interprete
	scopes []map[string]bool
}

// Crea una nueva instancia del resolvedor asociada al intérprete.
func Nuevo(intp *interprete.Interprete) *Resolver {
	return &Resolver{
		intp:   intp,
		scopes: make([]map[string]bool, 0),
	}
}

// Resuelve una lista de sentencias de forma secuencial.
func (r *Resolver) Resolver(sentencias []sintaxis.Stmt) error {
	for _, s := range sentencias {
		if err := r.ResolverSentencia(s); err != nil {
			return err
		}
	}

	return nil
}

// Resuelve una sentencia individual según su tipo sintáctico.
func (r *Resolver) ResolverSentencia(stmt sintaxis.Stmt) error {
	switch s := stmt.(type) {
	// Bloque de código con ámbito propio.
	case *sintaxis.Block:
		r.beginScope()
		if err := r.Resolver(s.Statements); err != nil {
			return err
		}

		r.endScope()

	// Declaración de variable (declarar, resolver inicializador y definir).
	case *sintaxis.Var:
		r.declarar(s.Name)
		if s.Initializer != nil {
			if err := r.ResolverExpresion(s.Initializer); err != nil {
				return err
			}
		}

		r.definir(s.Name)

	// Declaración de funciones con su respectivo ámbito de parámetros y cuerpo.
	case *sintaxis.Function:
		r.declarar(s.Name)
		r.definir(s.Name)
		if err := r.resolverFuncion(s); err != nil {
			return err
		}

	// Sentencia de expresión suelta.
	case *sintaxis.ExpressionStmt:
		return r.ResolverExpresion(s.Expression)

	// Estructura condicional if/else.
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

	// Sentencia de impresión.
	case *sintaxis.Print:
		return r.ResolverExpresion(s.Value)

	// Sentencia de retorno.
	case *sintaxis.Return:
		if s.Value != nil {
			return r.ResolverExpresion(s.Value)
		}

	// Bucle while.
	case *sintaxis.While:
		if err := r.ResolverExpresion(s.Condition); err != nil {
			return err
		}

		return r.ResolverSentencia(s.Body)
	}

	return nil
}

// Resuelve una expresión determinando alcances estáticos y validando restricciones semánticas.
func (r *Resolver) ResolverExpresion(expr sintaxis.Expr) error {
	switch e := expr.(type) {
	// Referencia a variable (valida uso en propio inicializador y resuelve distancia local).
	case *sintaxis.Variable:
		if len(r.scopes) > 0 {
			if definido, existe := r.scopes[len(r.scopes)-1][e.Name.Lexema]; existe && !definido {
				return fmt.Errorf("[línea %d] No se puede leer una variable local en su propio inicializador", e.Name.Linea)
			}
		}

		r.resolverLocal(e, e.Name)

	// Asignación de variable.
	case *sintaxis.Assign:
		if err := r.ResolverExpresion(e.Value); err != nil {
			return err
		}

		r.resolverLocal(e, e.Name)

	// Expresión binaria.
	case *sintaxis.Binary:
		if err := r.ResolverExpresion(e.Left); err != nil {
			return err
		}

		return r.ResolverExpresion(e.Right)

	// Invocación o llamada a función.
	case *sintaxis.Call:
		if err := r.ResolverExpresion(e.Callee); err != nil {
			return err
		}

		for _, arg := range e.Arguments {
			if err := r.ResolverExpresion(arg); err != nil {
				return err
			}
		}

	// Expresión agrupada entre paréntesis.
	case *sintaxis.Grouping:
		return r.ResolverExpresion(e.Expression)

	// Expresión lógica (AND/OR).
	case *sintaxis.Logical:
		if err := r.ResolverExpresion(e.Left); err != nil {
			return err
		}

		return r.ResolverExpresion(e.Right)

	// Expresión unaria.
	case *sintaxis.Unary:
		return r.ResolverExpresion(e.Right)

	// Literales (sin resolución estática requerida).
	case *sintaxis.Literal:
	}
	return nil
}

// Calcula la distancia estática de una variable en los ámbitos anidados.
func (r *Resolver) resolverLocal(expr sintaxis.Expr, nombre token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][nombre.Lexema]; ok {
			r.intp.ResolverLocal(expr, len(r.scopes)-1-i)
			return
		}
	}
}

// Resuelve una función declarada creando un nuevo ámbito para sus parámetros.
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

// Inicia un nuevo ámbito local.
func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

// Cierra y descarta el ámbito local actual.
func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

// Declara una variable en el ámbito actual como no definida (útil para detectar inicialización propia).
func (r *Resolver) declarar(nombre token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[nombre.Lexema] = false
}

// Define una variable en el ámbito actual como completamente inicializada.
func (r *Resolver) definir(nombre token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][nombre.Lexema] = true
}
