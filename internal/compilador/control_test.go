package compilador

import "testing"

// El control de flujo se arma con saltos. Los desplazamientos concretos son un
var casosControlDeFlujo = []casoDeCompilacion{
	{
		nombre:      "condicional sin rama falsa",
		fuente:      "if (true) print 1;",
		soloNombres: true,
		esperado: []string{
			"OpTrue",
			"OpJumpIfFalse",
			"OpPop",
			"OpConstant", "OpPrint",
			"OpJump",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "condicional con rama falsa",
		fuente:      "if (true) print 1; else print 2;",
		soloNombres: true,
		esperado: []string{
			"OpTrue",
			"OpJumpIfFalse",
			"OpPop", "OpConstant", "OpPrint",
			"OpJump",
			"OpPop", "OpConstant", "OpPrint",
			"OpReturn",
		},
	},
	{
		nombre:      "bucle mientras",
		fuente:      "while (false) print 1;",
		soloNombres: true,
		esperado: []string{
			"OpFalse",
			"OpJumpIfFalse",
			"OpPop", "OpConstant", "OpPrint",
			"OpLoop",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "conjunción: corta sin evaluar el lado derecho",
		fuente:      "print true and false;",
		soloNombres: true,
		esperado: []string{
			"OpTrue",
			"OpJumpIfFalse",
			"OpPop",
			"OpFalse",
			"OpPrint", "OpReturn",
		},
	},
	{
		nombre:      "disyunción: corta sin evaluar el lado derecho",
		fuente:      "print true or false;",
		soloNombres: true,
		esperado: []string{
			"OpTrue",
			"OpJumpIfFalse",
			"OpJump",
			"OpPop",
			"OpFalse",
			"OpPrint", "OpReturn",
		},
	},
}

// El for no tiene instrucciones propias: se compila con los mismos saltos que
var casosBucleFor = []casoDeCompilacion{
	{
		nombre:      "for completo",
		fuente:      "for (var i = 0; i < 3; i = i + 1) print i;",
		soloNombres: true,
		esperado: []string{
			"OpConstant", // i = 0

			"OpGetLocal", "OpConstant", "OpLess",
			"OpJumpIfFalse", "OpPop",

			"OpJump",

			"OpGetLocal", "OpConstant", "OpAdd", "OpSetLocal", "OpPop",
			"OpLoop", // vuelve a la condición

			"OpGetLocal", "OpPrint",
			"OpLoop", // vuelve al incremento

			"OpPop", // descarte de la condición
			"OpPop", // cierre del ámbito de i
			"OpReturn",
		},
	},
}

func TestBytecodeDeControlDeFlujo(t *testing.T) {
	correrCasosDeCompilacion(t, casosControlDeFlujo)
}

func TestBytecodeDelBucleFor(t *testing.T) {
	correrCasosDeCompilacion(t, casosBucleFor)
}
