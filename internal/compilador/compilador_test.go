package compilador

import (
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/escaner"
)

// casoDeCompilacion es una entrada de una tabla de pruebas: un fuente y el
type casoDeCompilacion struct {
	nombre string
	fuente string

	esperado    []string
	soloNombres bool
}

// compilarFuente escanea y compila un fuente.
func compilarFuente(t *testing.T, fuente string) (funcion *Funcion) {
	t.Helper()

	tokens, erroresLexicos := escaner.Nuevo(fuente).EscanearTokens()
	if len(erroresLexicos) > 0 {
		t.Fatalf("el fuente de la prueba no escanea: %v", erroresLexicos)
	}

	defer func() {
		if recuperado := recover(); recuperado != nil {
			t.Fatalf("el compilador entró en pánico con %q: %v", fuente, recuperado)
		}
	}()

	funcion = Compilar(tokens)
	if funcion == nil {
		t.Fatal("Compilar devolvió nil")
	}

	return funcion
}

// correrCasosDeCompilacion ejecuta una tabla de pruebas de bytecode.
func correrCasosDeCompilacion(t *testing.T, casos []casoDeCompilacion) {
	t.Helper()

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			obtenidas := desensamblar(compilarFuente(t, caso.fuente))
			if caso.soloNombres {
				obtenidas = soloOpcodes(obtenidas)
			}

			if len(obtenidas) != len(caso.esperado) {
				t.Fatalf("se obtuvieron %d instrucciones y se esperaban %d\nobtenidas: %v",
					len(obtenidas), len(caso.esperado), obtenidas)
			}

			for i, esperada := range caso.esperado {
				if obtenidas[i] != esperada {
					t.Errorf("instrucción %d: se obtuvo %q y se esperaba %q\ncompleto: %v",
						i, obtenidas[i], esperada, obtenidas)
				}
			}
		})
	}
}

// Todo programa compilado termina en un OpReturn, que devuelve el control.
var casosBase = []casoDeCompilacion{
	{
		nombre:   "programa vacío",
		fuente:   "",
		esperado: []string{"OpReturn"},
	},
	{
		nombre:   "una expresión suelta se descarta con OpPop",
		fuente:   "1;",
		esperado: []string{"OpConstant 0 -> 1", "OpPop", "OpReturn"},
	},
}

// Los literales: los números y las cadenas van al pool de constantes, pero
var casosLiterales = []casoDeCompilacion{
	{
		nombre:   "número entero",
		fuente:   "print 1;",
		esperado: []string{"OpConstant 0 -> 1", "OpPrint", "OpReturn"},
	},
	{
		nombre:   "número decimal",
		fuente:   "print 1.5;",
		esperado: []string{"OpConstant 0 -> 1.5", "OpPrint", "OpReturn"},
	},
	{
		nombre:   "cadena",
		fuente:   `print "hola";`,
		esperado: []string{`OpConstant 0 -> "hola"`, "OpPrint", "OpReturn"},
	},
	{
		nombre:   "verdadero",
		fuente:   "print true;",
		esperado: []string{"OpTrue", "OpPrint", "OpReturn"},
	},
	{
		nombre:   "falso",
		fuente:   "print false;",
		esperado: []string{"OpFalse", "OpPrint", "OpReturn"},
	},
	{
		nombre:   "nulo",
		fuente:   "print nil;",
		esperado: []string{"OpNil", "OpPrint", "OpReturn"},
	},
}

// La aritmética queda en notación polaca inversa: primero los operandos y
var casosAritmetica = []casoDeCompilacion{
	{
		nombre:   "suma",
		fuente:   "print 1 + 2;",
		esperado: []string{"OpConstant 0 -> 1", "OpConstant 1 -> 2", "OpAdd", "OpPrint", "OpReturn"},
	},
	{
		nombre:   "módulo",
		fuente:   "print 10 % 3;",
		esperado: []string{"OpConstant 0 -> 10", "OpConstant 1 -> 3", "OpModulo", "OpPrint", "OpReturn"},
	},
	{
		nombre:   "negación",
		fuente:   "print -1;",
		esperado: []string{"OpConstant 0 -> 1", "OpNegate", "OpPrint", "OpReturn"},
	},
	{
		nombre: "el producto se emite antes que la suma",
		fuente: "print 1 + 2 * 3;",
		esperado: []string{
			"OpConstant 0 -> 1", "OpConstant 1 -> 2", "OpConstant 2 -> 3",
			"OpMultiply", "OpAdd", "OpPrint", "OpReturn",
		},
	},
	{
		nombre: "los paréntesis cambian el orden de emisión",
		fuente: "print (1 + 2) * 3;",
		esperado: []string{
			"OpConstant 0 -> 1", "OpConstant 1 -> 2", "OpAdd",
			"OpConstant 2 -> 3", "OpMultiply", "OpPrint", "OpReturn",
		},
	},
	{
		nombre: "la resta asocia a la izquierda",
		fuente: "print 1 - 2 - 3;",
		esperado: []string{
			"OpConstant 0 -> 1", "OpConstant 1 -> 2", "OpSubtract",
			"OpConstant 2 -> 3", "OpSubtract", "OpPrint", "OpReturn",
		},
	},
}

// Tres comparaciones no tienen opcode propio: se arman negando a su opuesta.
var casosComparaciones = []casoDeCompilacion{
	{
		nombre:      "igualdad",
		fuente:      "print 1 == 2;",
		soloNombres: true,
		esperado:    []string{"OpConstant", "OpConstant", "OpEqual", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "desigualdad: es la igualdad negada",
		fuente:      "print 1 != 2;",
		soloNombres: true,
		esperado:    []string{"OpConstant", "OpConstant", "OpEqual", "OpNot", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "mayor",
		fuente:      "print 1 > 2;",
		soloNombres: true,
		esperado:    []string{"OpConstant", "OpConstant", "OpGreater", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "mayor o igual: es el menor negado",
		fuente:      "print 1 >= 2;",
		soloNombres: true,
		esperado:    []string{"OpConstant", "OpConstant", "OpLess", "OpNot", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "menor",
		fuente:      "print 1 < 2;",
		soloNombres: true,
		esperado:    []string{"OpConstant", "OpConstant", "OpLess", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "menor o igual: es el mayor negado",
		fuente:      "print 1 <= 2;",
		soloNombres: true,
		esperado:    []string{"OpConstant", "OpConstant", "OpGreater", "OpNot", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "negación lógica",
		fuente:      "print !true;",
		soloNombres: true,
		esperado:    []string{"OpTrue", "OpNot", "OpPrint", "OpReturn"},
	},
}

func TestBytecodeBase(t *testing.T) {
	correrCasosDeCompilacion(t, casosBase)
}

func TestBytecodeDeLiterales(t *testing.T) {
	correrCasosDeCompilacion(t, casosLiterales)
}

func TestBytecodeDeAritmetica(t *testing.T) {
	correrCasosDeCompilacion(t, casosAritmetica)
}

func TestBytecodeDeComparaciones(t *testing.T) {
	correrCasosDeCompilacion(t, casosComparaciones)
}
