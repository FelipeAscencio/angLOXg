package lox_test

import "testing"

var casosArrays = []casoDeClases{
	{
		nombre:   "array vacío",
		fuente:   "print [];",
		esperado: []string{"[]"},
		exito:    true,
	},
	{
		nombre:   "array con elementos",
		fuente:   "print [1, 2, 3];",
		esperado: []string{"[1, 2, 3]"},
		exito:    true,
	},
	{
		nombre:   "elementos de tipos distintos",
		fuente:   `print [1, "dos", true, nil];`,
		esperado: []string{`[1, dos, true, nil]`},
		exito:    true,
	},
	{
		nombre:   "lectura por posición",
		fuente:   "var a = [10, 20, 30]; print a[0];",
		esperado: []string{"10"},
		exito:    true,
	},
	{
		nombre:   "lectura de la última posición",
		fuente:   "var a = [10, 20, 30]; print a[2];",
		esperado: []string{"30"},
		exito:    true,
	},
	{
		nombre:   "el índice puede ser una expresión",
		fuente:   "var a = [10, 20, 30]; var i = 1; print a[i + 1];",
		esperado: []string{"30"},
		exito:    true,
	},
	{
		nombre:   "lectura sobre un literal",
		fuente:   "print [7, 8][1];",
		esperado: []string{"8"},
		exito:    true,
	},
	{
		nombre:   "los elementos se evalúan al crear el array",
		fuente:   "var x = 5; var a = [x, x * 2]; x = 99; print a[0]; print a[1];",
		esperado: []string{"5", "10"},
		exito:    true,
	},
	{
		nombre:   "arrays anidados",
		fuente:   "var m = [[1, 2], [3, 4]]; print m[1][0];",
		esperado: []string{"3"},
		exito:    true,
	},
	{
		nombre:   "recorrer un array con un for",
		fuente:   "var a = [10, 20, 30]; for (var i = 0; i < 3; i = i + 1) print a[i];",
		esperado: []string{"10", "20", "30"},
		exito:    true,
	},
	{
		nombre:   "un array puede guardar funciones",
		fuente:   "fun f() { return 42; } var a = [f]; print a[0]();",
		esperado: []string{"42"},
		exito:    true,
	},
	{
		nombre:   "un array puede guardar instancias",
		fuente:   "class P {} var a = [P()]; a[0].x = 7; print a[0].x;",
		esperado: []string{"7"},
		exito:    true,
	},
	{
		nombre:   "un array se puede pasar a una función",
		fuente:   "fun primero(lista) { return lista[0]; } print primero([9, 8]);",
		esperado: []string{"9"},
		exito:    true,
	},
	{
		nombre:   "un array se puede guardar en una propiedad",
		fuente:   "class P {} var p = P(); p.lista = [1, 2]; print p.lista[1];",
		esperado: []string{"2"},
		exito:    true,
	},
	{
		nombre:   "dos arrays con los mismos valores no son el mismo",
		fuente:   "print [1] == [1];",
		esperado: []string{"false"},
		exito:    true,
	},
	{
		nombre:   "un array es igual a sí mismo",
		fuente:   "var a = [1]; var b = a; print a == b;",
		esperado: []string{"true"},
		exito:    true,
	},
}

var casosErroresDeArrays = []casoDeClases{
	{
		nombre:   "índice fuera de rango por arriba",
		fuente:   "var a = [1, 2]; print a[2];",
		esperado: []string{"Índice fuera de rango.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "índice negativo",
		fuente:   "var a = [1, 2]; print a[-1];",
		esperado: []string{"Índice fuera de rango.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "índice sobre un array vacío",
		fuente:   "print [][0];",
		esperado: []string{"Índice fuera de rango.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "indexar algo que no es un array",
		fuente:   "var a = 1; print a[0];",
		esperado: []string{"Solo los arrays se pueden indexar.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "indexar una cadena",
		fuente:   `print "hola"[0];`,
		esperado: []string{"Solo los arrays se pueden indexar.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "índice que no es un número",
		fuente:   `var a = [1]; print a["x"];`,
		esperado: []string{"El índice debe ser un número entero.", "[línea 1]"},
		exito:    false,
	},
	{
		nombre:   "índice con parte decimal",
		fuente:   "var a = [1, 2]; print a[1.5];",
		esperado: []string{"El índice debe ser un número entero.", "[línea 1]"},
		exito:    false,
	},
}

func TestArrays(t *testing.T) {
	correrCasosDeClases(t, casosArrays)
}

func TestErroresDeArrays(t *testing.T) {
	correrCasosDeClases(t, casosErroresDeArrays)
}
