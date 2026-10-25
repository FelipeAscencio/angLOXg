package compilador

import "testing"

// Las globales se guardan por nombre en una tabla, así que su nombre viaja en
var casosVariables = []casoDeCompilacion{
	{
		nombre: "declaración y lectura de una global",
		fuente: "var x = 1; print x;",
		esperado: []string{
			"OpConstant 0 -> 1",
			`OpDefineGlobal 1 -> "x"`,
			`OpGetGlobal 2 -> "x"`,
			"OpPrint", "OpReturn",
		},
	},
	{
		nombre: "global sin valor inicial arranca en nil",
		fuente: "var x;",
		esperado: []string{
			"OpNil",
			`OpDefineGlobal 0 -> "x"`,
			"OpReturn",
		},
	},
	{
		nombre: "asignación a una global",
		fuente: "var x = 1; x = 2;",
		esperado: []string{
			"OpConstant 0 -> 1",
			`OpDefineGlobal 1 -> "x"`,
			"OpConstant 3 -> 2",
			`OpSetGlobal 2 -> "x"`,
			"OpPop", "OpReturn",
		},
	},
	{
		nombre: "una local no toca la tabla de globales",
		fuente: "{ var y = 2; print y; }",
		esperado: []string{
			"OpConstant 0 -> 2",
			"OpGetLocal 1",
			"OpPrint",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre: "asignación a una local",
		fuente: "{ var y = 1; y = 2; }",
		esperado: []string{
			"OpConstant 0 -> 1",
			"OpConstant 1 -> 2",
			"OpSetLocal 1",
			"OpPop",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "dos locales en el mismo bloque usan índices distintos",
		fuente:      "{ var a = 1; var b = 2; print a; print b; }",
		soloNombres: true,
		esperado: []string{
			"OpConstant", "OpConstant",
			"OpGetLocal", "OpPrint",
			"OpGetLocal", "OpPrint",
			"OpPop", "OpPop",
			"OpReturn",
		},
	},
}

func TestBytecodeDeVariables(t *testing.T) {
	correrCasosDeCompilacion(t, casosVariables)
}

// Los índices de las locales tienen que ser distintos entre sí; si no, dos
func TestIndicesDeLocalesDistintos(t *testing.T) {
	instrucciones := desensamblar(compilarFuente(t, "{ var a = 1; var b = 2; print a; print b; }"))

	lecturas := []string{}
	for _, instruccion := range instrucciones {
		if len(instruccion) > 10 && instruccion[:10] == "OpGetLocal" {
			lecturas = append(lecturas, instruccion)
		}
	}

	if len(lecturas) != 2 {
		t.Fatalf("se esperaban 2 lecturas de local y hubo %d: %v", len(lecturas), instrucciones)
	}

	if lecturas[0] == lecturas[1] {
		t.Errorf("las dos locales leen el mismo índice (%s)", lecturas[0])
	}
}
