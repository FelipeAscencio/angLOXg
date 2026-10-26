package compilador

import "testing"

var casosArraysBytecode = []casoDeCompilacion{
	{
		nombre: "array vacío",
		fuente: "[];",
		esperado: []string{"OpArray 0", "OpPop", "OpReturn"},
	},
	{
		nombre: "array con elementos",
		fuente: "[1, 2, 3];",
		esperado: []string{
			"OpConstant 0 -> 1",
			"OpConstant 1 -> 2",
			"OpConstant 2 -> 3",
			"OpArray 3",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "lectura por posición",
		fuente:      "a[0];",
		soloNombres: true,
		esperado:    []string{"OpGetGlobal", "OpConstant", "OpGetIndex", "OpPop", "OpReturn"},
	},
	{
		nombre:      "lecturas encadenadas",
		fuente:      "m[0][1];",
		soloNombres: true,
		esperado: []string{
			"OpGetGlobal", "OpConstant", "OpGetIndex",
			"OpConstant", "OpGetIndex",
			"OpPop", "OpReturn",
		},
	},
	{
		nombre:      "arrays anidados",
		fuente:      "[[1], [2]];",
		soloNombres: true,
		esperado: []string{
			"OpConstant", "OpArray",
			"OpConstant", "OpArray",
			"OpArray",
			"OpPop", "OpReturn",
		},
	},
	{
		nombre:      "el indexado liga más fuerte que la aritmética",
		fuente:      "print a[0] + 1;",
		soloNombres: true,
		esperado: []string{
			"OpGetGlobal", "OpConstant", "OpGetIndex",
			"OpConstant", "OpAdd",
			"OpPrint", "OpReturn",
		},
	},
}

func TestBytecodeDeArrays(t *testing.T) {
	correrCasosDeCompilacion(t, casosArraysBytecode)
}
