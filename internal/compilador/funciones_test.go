package compilador

import (
	"strings"
	"testing"
)

// Una función se compila a su propio chunk, que queda guardado como una
var casosFunciones = []casoDeCompilacion{
	{
		nombre:      "declaración y llamada",
		fuente:      "fun f(a) { return a; } print f(1);",
		soloNombres: true,
		esperado: []string{
			"OpClosure", "OpDefineGlobal",
			"OpGetGlobal", "OpConstant", "OpCall",
			"OpPrint", "OpReturn",
		},
	},
	{
		nombre:      "llamada sin argumentos",
		fuente:      "fun f() { return 1; } f();",
		soloNombres: true,
		esperado: []string{
			"OpClosure", "OpDefineGlobal",
			"OpGetGlobal", "OpCall", "OpPop",
			"OpReturn",
		},
	},
}

func TestBytecodeDeFunciones(t *testing.T) {
	correrCasosDeCompilacion(t, casosFunciones)
}

func TestLaFuncionQuedaComoConstante(t *testing.T) {
	programa := compilarFuente(t, "fun sumar(a, b) { return a + b; }")

	sumar := funcionConstante(programa, "sumar")
	if sumar == nil {
		t.Fatal("la función compilada no quedó en el pool de constantes")
	}

	if sumar.Aridad != 2 {
		t.Errorf("la aridad es %d y se esperaba 2", sumar.Aridad)
	}

	if sumar.String() != "<fn sumar>" {
		t.Errorf("la representación es %q y se esperaba \"<fn sumar>\"", sumar.String())
	}
}

func TestCuerpoDeUnaFuncion(t *testing.T) {
	programa := compilarFuente(t, "fun sumar(a, b) { return a + b; }")
	sumar := funcionConstante(programa, "sumar")
	if sumar == nil {
		t.Fatal("la función compilada no quedó en el pool de constantes")
	}

	esperado := []string{
		"OpGetLocal 1", "OpGetLocal 2", "OpAdd", "OpReturn",
		"OpNil", "OpReturn",
	}

	obtenidas := desensamblar(sumar)
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

func TestRetornoImplicito(t *testing.T) {
	programa := compilarFuente(t, "fun nada() {}")
	nada := funcionConstante(programa, "nada")
	if nada == nil {
		t.Fatal("la función compilada no quedó en el pool de constantes")
	}

	esperado := []string{"OpNil", "OpReturn"}
	obtenidas := desensamblar(nada)

	if len(obtenidas) != len(esperado) || obtenidas[0] != esperado[0] || obtenidas[1] != esperado[1] {
		t.Errorf("se obtuvo %v y se esperaba %v", obtenidas, esperado)
	}
}

func TestUpvalueDeUnClosure(t *testing.T) {
	programa := compilarFuente(t, `
		fun externa() {
			var capturada = 1;
			fun interna() { return capturada; }
			return interna();
		}
	`)

	externa := funcionConstante(programa, "externa")
	if externa == nil {
		t.Fatal("no se compiló la función externa")
	}

	interna := funcionConstante(externa, "interna")
	if interna == nil {
		t.Fatal("no se compiló la función interna")
	}

	if interna.UpvalueCount != 1 {
		t.Errorf("la función interna captura %d upvalues y se esperaba 1", interna.UpvalueCount)
	}

	obtenidas := desensamblar(interna)
	if len(obtenidas) == 0 || obtenidas[0] != "OpGetUpvalue 0" {
		t.Errorf("la primera instrucción es %v y se esperaba OpGetUpvalue 0", obtenidas)
	}
}

func TestSinUpvaluesCuandoNoCapturaNada(t *testing.T) {
	programa := compilarFuente(t, "fun externa() { fun interna() { return 1; } return interna(); }")

	externa := funcionConstante(programa, "externa")
	if externa == nil {
		t.Fatal("no se compiló la función externa")
	}

	interna := funcionConstante(externa, "interna")
	if interna == nil {
		t.Fatal("no se compiló la función interna")
	}

	if interna.UpvalueCount != 0 {
		t.Errorf("captura %d upvalues y no debería capturar ninguno", interna.UpvalueCount)
	}
}

// Una función declarada dentro de otra es local: su nombre no puede terminar en
func TestLaFuncionAnidadaNoSeDefineComoGlobal(t *testing.T) {
	programa := compilarFuente(t, "fun externa() { fun interna() { return 1; } return interna(); }")

	externa := funcionConstante(programa, "externa")
	if externa == nil {
		t.Fatal("no se compiló la función externa")
	}

	instrucciones := desensamblar(externa)
	for _, instruccion := range instrucciones {
		if strings.HasPrefix(instruccion, "OpDefineGlobal") {
			t.Errorf("la función anidada se define como global: %v", instrucciones)
		}

		if strings.HasPrefix(instruccion, "OpGetGlobal") {
			t.Errorf("la función anidada se lee como global: %v", instrucciones)
		}
	}

	if !contieneInstruccion(instrucciones, "OpGetLocal") {
		t.Errorf("la función anidada no se lee como local: %v", instrucciones)
	}
}

// Una función local se tiene que poder llamar a sí misma: su nombre queda
func TestLaFuncionAnidadaPuedeSerRecursiva(t *testing.T) {
	programa := compilarFuente(t, `
		fun externa() {
			fun fact(n) { if (n <= 1) return 1; return n * fact(n - 1); }
			return fact(5);
		}
	`)

	externa := funcionConstante(programa, "externa")
	if externa == nil {
		t.Fatal("no se compiló la función externa")
	}

	fact := funcionConstante(externa, "fact")
	if fact == nil {
		t.Fatal("no se compiló la función fact")
	}

	if fact.UpvalueCount != 1 {
		t.Errorf("fact captura %d upvalues y se esperaba 1, que es ella misma", fact.UpvalueCount)
	}
}

// contieneInstruccion indica si alguna instrucción arranca con ese opcode.
func contieneInstruccion(instrucciones []string, opcode string) bool {
	for _, instruccion := range instrucciones {
		if strings.HasPrefix(instruccion, opcode) {
			return true
		}
	}

	return false
}
