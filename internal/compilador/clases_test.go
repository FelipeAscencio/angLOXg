package compilador

import "testing"

var casosClases = []casoDeCompilacion{
	{
		nombre: "clase sin métodos",
		fuente: "class P {}",
		esperado: []string{
			`OpClass 0 -> "P"`,
			`OpDefineGlobal 1 -> "P"`,
			`OpGetGlobal 2 -> "P"`,
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre: "clase con un método",
		fuente: "class P { m() { return 1; } }",
		esperado: []string{
			`OpClass 0 -> "P"`,
			`OpDefineGlobal 1 -> "P"`,
			`OpGetGlobal 2 -> "P"`,
			"OpClosure 4 -> <fn m>",
			`OpMethod 3 -> "m"`,
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "clase con varios métodos",
		fuente:      "class P { uno() {} dos() {} }",
		soloNombres: true,
		esperado: []string{
			"OpClass", "OpDefineGlobal", "OpGetGlobal",
			"OpClosure", "OpMethod",
			"OpClosure", "OpMethod",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre: "una clase local no va a la tabla de globales",
		fuente: "{ class P {} }",
		esperado: []string{
			`OpClass 0 -> "P"`,
			"OpGetLocal 1",
			"OpPop",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "escritura de una propiedad",
		fuente:      "var p = P(); p.x = 3;",
		soloNombres: true,
		esperado: []string{
			"OpGetGlobal", "OpCall", "OpDefineGlobal",
			"OpGetGlobal", "OpConstant", "OpSetProperty", "OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "lectura de una propiedad",
		fuente:      "print p.x;",
		soloNombres: true,
		esperado:    []string{"OpGetGlobal", "OpGetProperty", "OpPrint", "OpReturn"},
	},
	{
		nombre:      "lecturas encadenadas",
		fuente:      "print a.b.c;",
		soloNombres: true,
		esperado:    []string{"OpGetGlobal", "OpGetProperty", "OpGetProperty", "OpPrint", "OpReturn"},
	},
}

var casosHerenciaBytecode = []casoDeCompilacion{
	{
		nombre:      "clase con superclase",
		fuente:      "class B < A {}",
		soloNombres: true,
		esperado: []string{
			"OpClass", "OpDefineGlobal",
			"OpGetGlobal", "OpGetGlobal", "OpInherit",
			"OpGetGlobal", "OpPop",
			"OpPop",
			"OpReturn",
		},
	},
	{
		nombre:      "llamada a super dentro de un método",
		fuente:      "class B < A { m() { super.m(); } }",
		soloNombres: true,
		esperado: []string{
			"OpClass", "OpDefineGlobal",
			"OpGetGlobal", "OpGetGlobal", "OpInherit",
			"OpGetGlobal", "OpClosure", "OpMethod", "OpPop",
			"OpCloseUpvalue",
			"OpReturn",
		},
	},
}

func TestBytecodeDeHerencia(t *testing.T) {
	correrCasosDeCompilacion(t, casosHerenciaBytecode)
}

func TestSuperSeCapturaComoUpvalue(t *testing.T) {
	programa := compilarFuente(t, "class B < A { m() { super.m(); } }")

	metodo := funcionConstante(programa, "m")
	if metodo == nil {
		t.Fatal("el método compilado no quedó en el pool de constantes")
	}

	if metodo.UpvalueCount != 1 {
		t.Errorf("el método captura %d upvalues y se esperaba 1, que es super", metodo.UpvalueCount)
	}

	if !contieneInstruccion(desensamblar(metodo), "OpGetSuper") {
		t.Errorf("el método no emite OpGetSuper: %v", desensamblar(metodo))
	}
}

func TestBytecodeDeClases(t *testing.T) {
	correrCasosDeCompilacion(t, casosClases)
}

func TestCuerpoDeUnMetodo(t *testing.T) {
	programa := compilarFuente(t, "class P { ver() { return this.x; } }")

	ver := funcionConstante(programa, "ver")
	if ver == nil {
		t.Fatal("el método compilado no quedó en el pool de constantes")
	}

	esperado := []string{
		"OpGetLocal 0",
		`OpGetProperty 0 -> "x"`,
		"OpReturn",
		"OpNil", "OpReturn",
	}

	obtenidas := desensamblar(ver)
	if len(obtenidas) != len(esperado) {
		t.Fatalf("se obtuvieron %d instrucciones y se esperaban %d\nobtenidas: %v",
			len(obtenidas), len(esperado), obtenidas)
	}

	for i, instruccion := range esperado {
		if obtenidas[i] != instruccion {
			t.Errorf("instrucción %d: se obtuvo %q y se esperaba %q", i, obtenidas[i], instruccion)
		}
	}
}

func TestAridadDeUnMetodo(t *testing.T) {
	programa := compilarFuente(t, "class P { sumar(a, b) { return a + b; } }")

	sumar := funcionConstante(programa, "sumar")
	if sumar == nil {
		t.Fatal("el método compilado no quedó en el pool de constantes")
	}

	if sumar.Aridad != 2 {
		t.Errorf("la aridad es %d y se esperaba 2", sumar.Aridad)
	}
}

func TestRetornoImplicitoDelInicializador(t *testing.T) {
	programa := compilarFuente(t, "class P { init() {} }")

	init := funcionConstante(programa, NombreDelInicializador)
	if init == nil {
		t.Fatal("el inicializador compilado no quedó en el pool de constantes")
	}

	esperado := []string{"OpGetLocal 0", "OpReturn"}
	obtenidas := desensamblar(init)

	if len(obtenidas) != len(esperado) || obtenidas[0] != esperado[0] || obtenidas[1] != esperado[1] {
		t.Errorf("se obtuvo %v y se esperaba %v", obtenidas, esperado)
	}
}

func TestRetornoPeladoDentroDelInicializador(t *testing.T) {
	programa := compilarFuente(t, "class P { init() { return; } }")

	init := funcionConstante(programa, NombreDelInicializador)
	if init == nil {
		t.Fatal("el inicializador compilado no quedó en el pool de constantes")
	}

	obtenidas := desensamblar(init)
	if len(obtenidas) == 0 || obtenidas[0] != "OpGetLocal 0" {
		t.Errorf("la primera instrucción es %v y se esperaba OpGetLocal 0", obtenidas)
	}
}

func TestUpvalueDeThisEnUnaFuncionAnidada(t *testing.T) {
	programa := compilarFuente(t, "class P { m() { fun interna() { return this.x; } return interna(); } }")

	metodo := funcionConstante(programa, "m")
	if metodo == nil {
		t.Fatal("el método compilado no quedó en el pool de constantes")
	}

	interna := funcionConstante(metodo, "interna")
	if interna == nil {
		t.Fatal("la función anidada no quedó en el pool de constantes")
	}

	if interna.UpvalueCount != 1 {
		t.Errorf("la función anidada captura %d upvalues y se esperaba 1, que es this", interna.UpvalueCount)
	}
}
