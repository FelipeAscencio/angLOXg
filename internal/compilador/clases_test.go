package compilador

import "testing"

var casosClases = []casoDeCompilacion{
	{
		nombre: "clase sin métodos",
		fuente: "class P {}",
		// La clase se define como variable antes de los métodos, y después se
		// la vuelve a traer a la pila para colgárselos. El OpPop la saca.
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
		// Sin OpDefineGlobal: la clase se queda en la pila como local. El
		// primer OpPop saca la copia que se trajo para los métodos y el
		// segundo cierra el ámbito del bloque.
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
		// La superclase se lee y queda en la pila como la local "super". El
		// OpInherit copia sus métodos y el OpPop del final cierra ese ámbito.
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
		// Al cerrar el ámbito sale OpCloseUpvalue y no OpPop: el método se
		// llevó capturada la superclase, así que no alcanza con sacarla de la
		// pila, hay que dejarla viva en el heap.
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

	// El método captura la superclase: es lo que hace que cada "super" apunte
	// al nivel donde fue escrito y no a la clase de la instancia.
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
