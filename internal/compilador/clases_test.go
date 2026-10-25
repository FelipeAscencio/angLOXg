package compilador

import "testing"

var casosClases = []casoDeCompilacion{
	{
		nombre: "clase sin métodos",
		fuente: "class P {}",
		esperado: []string{
			`OpClass 0 -> "P"`,
			`OpDefineGlobal 1 -> "P"`,
			"OpReturn",
		},
	},
	{
		nombre: "clase con un método",
		fuente: "class P { m() { return 1; } }",
		// La clase queda en la pila mientras se le cuelga el método, y recién
		// después se define como global.
		esperado: []string{
			`OpClass 0 -> "P"`,
			"OpClosure 2 -> <fn m>",
			`OpMethod 1 -> "m"`,
			`OpDefineGlobal 3 -> "P"`,
			"OpReturn",
		},
	},
	{
		nombre:      "clase con varios métodos",
		fuente:      "class P { uno() {} dos() {} }",
		soloNombres: true,
		esperado: []string{
			"OpClass",
			"OpClosure", "OpMethod",
			"OpClosure", "OpMethod",
			"OpDefineGlobal",
			"OpReturn",
		},
	},
	{
		nombre: "una clase local no va a la tabla de globales",
		fuente: "{ class P {} }",
		// El OpPop del final es el que cierra el ámbito del bloque.
		esperado: []string{
			`OpClass 0 -> "P"`,
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

func TestBytecodeDeClases(t *testing.T) {
	correrCasosDeCompilacion(t, casosClases)
}

func TestCuerpoDeUnMetodo(t *testing.T) {
	programa := compilarFuente(t, "class P { ver() { return this.x; } }")

	ver := funcionConstante(programa, "ver")
	if ver == nil {
		t.Fatal("el método compilado no quedó en el pool de constantes")
	}

	// "this" es la local del slot cero: la máquina deja ahí el receptor antes
	// de entrar al método, así que se lee como cualquier otra local.
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

	// Un inicializador devuelve la instancia, que vive en el slot cero, en
	// lugar del nil que devuelve cualquier otro cuerpo.
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

	// El return pelado también devuelve la instancia, no nil.
	obtenidas := desensamblar(init)
	if len(obtenidas) == 0 || obtenidas[0] != "OpGetLocal 0" {
		t.Errorf("la primera instrucción es %v y se esperaba OpGetLocal 0", obtenidas)
	}
}

func TestUpvalueDeThisEnUnaFuncionAnidada(t *testing.T) {
	// Una función declarada dentro de un método captura "this" como upvalue.
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
