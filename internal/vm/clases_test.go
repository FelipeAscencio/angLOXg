package vm

import "testing"

var casosClases = []casoDeEjecucion{
	{
		nombre:   "una clase se imprime con su nombre",
		fuente:   "class Punto {} print Punto;",
		esperado: []string{"Punto"},
		exito:    true,
	},
	{
		nombre:   "llamar a la clase devuelve una instancia",
		fuente:   "class Punto {} print Punto();",
		esperado: []string{"instancia de Punto"},
		exito:    true,
	},
	{
		nombre:   "la instancia se guarda en una variable",
		fuente:   "class P {} var p = P(); print p;",
		esperado: []string{"instancia de P"},
		exito:    true,
	},
	{
		nombre:   "una clase declarada dentro de un bloque",
		fuente:   "{ class P {} print P(); }",
		esperado: []string{"instancia de P"},
		exito:    true,
	},
	{
		nombre:   "escribir y leer una propiedad",
		fuente:   "class P {} var p = P(); p.x = 3; print p.x;",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "varias propiedades",
		fuente:   "class P {} var p = P(); p.x = 3; p.y = 4; print p.x + p.y;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "dos instancias no comparten propiedades",
		fuente:   "class P {} var a = P(); var b = P(); a.x = 1; b.x = 2; print a.x; print b.x;",
		esperado: []string{"1", "2"},
		exito:    true,
	},
	{
		nombre:   "la escritura devuelve el valor asignado",
		fuente:   "class P {} var p = P(); print p.x = 5;",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre:   "una propiedad puede guardar otra instancia",
		fuente:   "class P {} var a = P(); var b = P(); a.hijo = b; a.hijo.v = 42; print a.hijo.v;",
		esperado: []string{"42"},
		exito:    true,
	},
	{
		nombre:   "leer una propiedad que no existe",
		fuente:   "class P {} print P().z;",
		esperado: []string{"Error en ejecución: Propiedad no definida 'z'."},
		exito:    false,
	},
	{
		nombre:   "leer una propiedad de un número",
		fuente:   "print (1).x;",
		esperado: []string{"Error en ejecución: Solo las instancias tienen propiedades."},
		exito:    false,
	},
	{
		nombre:   "escribir una propiedad en un número",
		fuente:   "var a = 1; a.x = 1;",
		esperado: []string{"Error en ejecución: Solo las instancias tienen propiedades."},
		exito:    false,
	},
}

var casosMetodosVM = []casoDeEjecucion{
	{
		nombre:   "llamada a un método",
		fuente:   `class A { saludar() { print "hola"; } } A().saludar();`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "método con parámetros y retorno",
		fuente:   "class A { sumar(x, y) { return x + y; } } print A().sumar(2, 3);",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre:   "un método se imprime como función",
		fuente:   "class A { m() {} } print A().m;",
		esperado: []string{"<fn m>"},
		exito:    true,
	},
	{
		nombre:   "varios métodos",
		fuente:   "class A { uno() { return 1; } dos() { return 2; } } var a = A(); print a.uno(); print a.dos();",
		esperado: []string{"1", "2"},
		exito:    true,
	},
	{
		nombre:   "una propiedad tapa a un método del mismo nombre",
		fuente:   `class A { m() { return "metodo"; } } var a = A(); a.m = "propiedad"; print a.m;`,
		esperado: []string{"propiedad"},
		exito:    true,
	},
	{
		nombre:   "llamar a un método que no existe",
		fuente:   "class A {} A().noExiste();",
		esperado: []string{"Error en ejecución: Propiedad no definida 'noExiste'."},
		exito:    false,
	},
}

var casosThisVM = []casoDeEjecucion{
	{
		nombre:   "un método escribe una propiedad con this",
		fuente:   "class A { poner() { this.x = 7; } } var a = A(); a.poner(); print a.x;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "un método lee una propiedad con this",
		fuente:   "class A { ver() { return this.x; } } var a = A(); a.x = 3; print a.ver();",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "un método llama a otro con this",
		fuente:   "class A { doble() { return this.valor() * 2; } valor() { return 21; } } print A().doble();",
		esperado: []string{"42"},
		exito:    true,
	},
	{
		nombre: "this queda ligado a la instancia correcta",
		fuente: `class A { nombre() { return this.n; } }
			var a = A(); var b = A();
			a.n = "primera"; b.n = "segunda";
			print a.nombre(); print b.nombre();`,
		esperado: []string{"primera", "segunda"},
		exito:    true,
	},
	{
		nombre: "un método extraído se lleva su instancia",
		fuente: `class A { nombre() { return this.n; } }
			var a = A(); a.n = "primera";
			var m = a.nombre;
			print m();`,
		esperado: []string{"primera"},
		exito:    true,
	},
	{
		nombre:   "un método acumula estado",
		fuente:   "class C { subir() { this.n = this.n + 1; return this.n; } } var c = C(); c.n = 0; print c.subir(); print c.subir();",
		esperado: []string{"1", "2"},
		exito:    true,
	},
}

var casosInitVM = []casoDeEjecucion{
	{
		nombre:   "el init corre al instanciar",
		fuente:   `class A { init() { print "creada"; } } A();`,
		esperado: []string{"creada"},
		exito:    true,
	},
	{
		nombre:   "el init recibe argumentos",
		fuente:   "class P { init(x, y) { this.x = x; this.y = y; } } var p = P(3, 4); print p.x + p.y;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "instanciar devuelve la instancia",
		fuente:   "class P { init(x) { this.x = x; } } print P(1);",
		esperado: []string{"instancia de P"},
		exito:    true,
	},
	{
		nombre:   "la aridad de la clase es la del init",
		fuente:   "class P { init(x) {} } P();",
		esperado: []string{"Error en ejecución: Se esperaban 1 argumentos pero se obtuvieron 0."},
		exito:    false,
	},
	{
		nombre:   "sin init la clase no acepta argumentos",
		fuente:   "class P {} P(1);",
		esperado: []string{"Error en ejecución: Se esperaban 0 argumentos pero se obtuvieron 1."},
		exito:    false,
	},
	{
		nombre:   "un return pelado dentro del init devuelve la instancia",
		fuente:   "class P { init() { this.x = 1; return; this.x = 2; } } print P().x;",
		esperado: []string{"1"},
		exito:    true,
	},
	{
		nombre:   "el init puede llamar a otros métodos",
		fuente:   "class P { init() { this.preparar(); } preparar() { this.listo = true; } } print P().listo;",
		esperado: []string{"true"},
		exito:    true,
	},
	{
		nombre: "cada instancia corre su propio init",
		fuente: `class C { init(n) { this.n = n; } }
			var a = C(1); var b = C(2);
			print a.n; print b.n;`,
		esperado: []string{"1", "2"},
		exito:    true,
	},
}

func TestClases(t *testing.T) {
	correrCasos(t, casosClases)
}

func TestMetodos(t *testing.T) {
	correrCasos(t, casosMetodosVM)
}

func TestThis(t *testing.T) {
	correrCasos(t, casosThisVM)
}

func TestInit(t *testing.T) {
	correrCasos(t, casosInitVM)
}

var casosHerenciaVM = []casoDeEjecucion{
	{
		nombre:   "un método se hereda de la superclase",
		fuente:   `class A { saludar() { print "hola"; } } class B < A {} B().saludar();`,
		esperado: []string{"hola"},
		exito:    true,
	},
	{
		nombre:   "la subclase puede sobrescribir un método",
		fuente:   `class A { m() { print "base"; } } class B < A { m() { print "derivada"; } } B().m();`,
		esperado: []string{"derivada"},
		exito:    true,
	},
	{
		nombre:   "super llama al método de la superclase",
		fuente:   `class A { m() { print "base"; } } class B < A { m() { super.m(); print "derivada"; } } B().m();`,
		esperado: []string{"base", "derivada"},
		exito:    true,
	},
	{
		nombre:   "super ve this de la instancia real",
		fuente:   "class A { ver() { return this.x; } } class B < A {} var b = B(); b.x = 5; print b.ver();",
		esperado: []string{"5"},
		exito:    true,
	},
	{
		nombre: "herencia de tres niveles",
		fuente: `class A { m() { return "A"; } }
			class B < A { m() { return super.m() + "B"; } }
			class C < B { m() { return super.m() + "C"; } }
			print C().m();`,
		// Cada super apunta al nivel donde se escribió, no a la clase de la
		// instancia: si no, esto sería una recursión infinita.
		esperado: []string{"ABC"},
		exito:    true,
	},
	{
		nombre:   "el init se hereda",
		fuente:   "class A { init(x) { this.x = x; } } class B < A {} print B(7).x;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "la subclase puede llamar al init de la superclase",
		fuente:   "class A { init(x) { this.x = x; } } class B < A { init(x) { super.init(x * 2); } } print B(3).x;",
		esperado: []string{"6"},
		exito:    true,
	},
	{
		nombre:   "las propiedades son de la instancia",
		fuente:   "class A {} class B < A {} var a = A(); var b = B(); a.x = 1; b.x = 2; print a.x; print b.x;",
		esperado: []string{"1", "2"},
		exito:    true,
	},
	{
		nombre:   "la superclase tiene que ser una clase",
		fuente:   "var noEsClase = 1; class A < noEsClase {}",
		esperado: []string{"Error en ejecución: La superclase debe ser una clase."},
		exito:    false,
	},
	{
		nombre:   "un método heredado que no existe",
		fuente:   "class A {} class B < A {} B().noExiste();",
		esperado: []string{"Error en ejecución: Propiedad no definida 'noExiste'."},
		exito:    false,
	},
	{
		nombre:   "super sobre un método que la superclase no tiene",
		fuente:   "class A {} class B < A { m() { super.noExiste(); } } B().m();",
		esperado: []string{"Error en ejecución: Propiedad no definida 'noExiste'."},
		exito:    false,
	},
}

func TestHerenciaVM(t *testing.T) {
	correrCasos(t, casosHerenciaVM)
}
