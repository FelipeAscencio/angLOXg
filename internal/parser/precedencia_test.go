package parser

import "testing"

// Operadores unarios: se aplican a un solo operando y se anidan hacia la
// derecha.
var casosUnarios = []casoDeParseo{
	{
		nombre:   "negación aritmética",
		fuente:   "-1",
		esperado: "(- 1)",
	},
	{
		nombre:   "negación lógica",
		fuente:   "!true",
		esperado: "(! true)",
	},
	{
		nombre:   "unarios encadenados",
		fuente:   "--1",
		esperado: "(- (- 1))",
	},
	{
		nombre:   "negaciones lógicas encadenadas",
		fuente:   "!!false",
		esperado: "(! (! false))",
	},
	{
		nombre:  "unario sin operando",
		fuente:  "-",
		errores: []string{"se esperaba una expresión"},
	},
}

// Operadores binarios: la precedencia decide quién se agrupa primero.
var casosPrecedencia = []casoDeParseo{
	{
		nombre:   "producto",
		fuente:   "2 * 3",
		esperado: "(* 2 3)",
	},
	{
		nombre:   "división",
		fuente:   "6 / 3",
		esperado: "(/ 6 3)",
	},
	{
		nombre:   "módulo",
		fuente:   "10 % 3",
		esperado: "(% 10 3)",
	},
	{
		nombre: "el producto agrupa antes que la suma",
		fuente: "1 + 2 * 3",
		// La multiplicación queda adentro: es la prueba de que hay precedencia
		// y no una lectura de izquierda a derecha.
		esperado: "(+ 1 (* 2 3))",
	},
	{
		nombre:   "el producto agrupa antes que la suma, del otro lado",
		fuente:   "1 * 2 + 3",
		esperado: "(+ (* 1 2) 3)",
	},
	{
		nombre:   "el módulo tiene la misma precedencia que el producto",
		fuente:   "1 + 10 % 3",
		esperado: "(+ 1 (% 10 3))",
	},
	{
		nombre:   "el unario agrupa antes que el producto",
		fuente:   "-1 * 2",
		esperado: "(* (- 1) 2)",
	},
	{
		nombre:   "los paréntesis le ganan a la precedencia",
		fuente:   "(1 + 2) * 3",
		esperado: "(* (group (+ 1 2)) 3)",
	},
	{
		nombre:   "la suma agrupa antes que la comparación",
		fuente:   "1 + 2 < 3",
		esperado: "(< (+ 1 2) 3)",
	},
	{
		nombre:   "la comparación agrupa antes que la igualdad",
		fuente:   "1 < 2 == true",
		esperado: "(== (< 1 2) true)",
	},
	{
		nombre:   "todos los operadores de comparación",
		fuente:   "(1 > 2) == (3 >= 4) == (5 < 6) == (7 <= 8)",
		esperado: "(== (== (== (group (> 1 2)) (group (>= 3 4))) (group (< 5 6))) (group (<= 7 8)))",
	},
	{
		nombre:   "desigualdad",
		fuente:   "1 != 2",
		esperado: "(!= 1 2)",
	},
	{
		nombre:   "cadena larga con toda la precedencia",
		fuente:   "1 == 2 < 3 + 4 * -5",
		esperado: "(== 1 (< 2 (+ 3 (* 4 (- 5)))))",
	},
}

// Asociatividad: los binarios de la misma precedencia se agrupan de izquierda a
// derecha.
var casosAsociatividad = []casoDeParseo{
	{
		nombre: "la resta asocia a la izquierda",
		fuente: "1 - 2 - 3",
		// Tiene que ser (1-2)-3 y no 1-(2-3), que daría otro resultado.
		esperado: "(- (- 1 2) 3)",
	},
	{
		nombre:   "la división asocia a la izquierda",
		fuente:   "8 / 4 / 2",
		esperado: "(/ (/ 8 4) 2)",
	},
	{
		nombre:   "la suma asocia a la izquierda",
		fuente:   "1 + 2 + 3",
		esperado: "(+ (+ 1 2) 3)",
	},
	{
		nombre:   "la igualdad asocia a la izquierda",
		fuente:   "1 == 2 == 3",
		esperado: "(== (== 1 2) 3)",
	},
	{
		nombre:  "binario sin operando derecho",
		fuente:  "1 +",
		errores: []string{"se esperaba una expresión"},
	},
	{
		nombre:  "binario sin operando izquierdo",
		fuente:  "* 1",
		errores: []string{"se esperaba una expresión"},
	},
}

func TestOperadoresUnarios(t *testing.T) {
	correrCasosDeExpresion(t, casosUnarios)
}

func TestPrecedenciaDeOperadores(t *testing.T) {
	correrCasosDeExpresion(t, casosPrecedencia)
}

func TestAsociatividad(t *testing.T) {
	correrCasosDeExpresion(t, casosAsociatividad)
}
