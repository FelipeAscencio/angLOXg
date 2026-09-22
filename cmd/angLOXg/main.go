package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Códigos de salida del programa.
const (
	salidaOK           = 0
	salidaUso          = 64 // Error en la invocación del programa.
	salidaErrorDeDatos = 65 // El código Lox tenía errores.
	salidaSinArchivo   = 66 // No se pudo leer el script.
)

func main() {
	modoEscaneo := flag.Bool("escaneo", false,
		"mostrar los tokens producidos por el escáner en lugar de ejecutar el código")

	modoArbol := flag.Bool("arbol", false,
		"mostrar el árbol de sintaxis (AST) en lugar de ejecutar el código")
		
	modoResolucion := flag.Bool("resolucion", false,
		"ejecutar hasta la fase de análisis semántico (Resolver) sin interpretar el código")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "uso: angloxg [--escaneo | --arbol | --resolucion] [script.lox]")
		fmt.Fprintln(os.Stderr, "\nSin argumentos, abre la consola interactiva (REPL).")
		fmt.Fprintln(os.Stderr, "\nOpciones:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(salidaUso)
	}

	if flag.NArg() == 1 {
		os.Exit(ejecutarArchivo(flag.Arg(0), *modoEscaneo, *modoArbol, *modoResolucion))
	}

	os.Exit(ejecutarConsola(*modoEscaneo, *modoArbol, *modoResolucion))
}

func ejecutarArchivo(ruta string, modoEscaneo, modoArbol, modoResolucion bool) int {
	fuente, err := os.ReadFile(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo leer %q: %v\n", ruta, err)
		return salidaSinArchivo
	}

	if !despachar(string(fuente), os.Stdout, modoEscaneo, modoArbol, modoResolucion) {
		return salidaErrorDeDatos
	}

	return salidaOK
}

func ejecutarConsola(modoEscaneo, modoArbol, modoResolucion bool) int {
	fmt.Println("=== angLOXg ===")
	fmt.Println("Intérprete de Lox escrito en Go")
	entrada := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !entrada.Scan() {
			break
		}

		despachar(entrada.Text(), os.Stdout, modoEscaneo, modoArbol, modoResolucion)
	}

	if err := entrada.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error leyendo la entrada: %v\n", err)
		return salidaErrorDeDatos
	}

	fmt.Println()
	return salidaOK
}

func despachar(fuente string, salida io.Writer, modoEscaneo, modoArbol, modoResolucion bool) bool {
	switch {
	case modoEscaneo:
		return lox.Escanear(fuente, salida)

	case modoArbol:
		return lox.Parsear(fuente, salida)
		
	case modoResolucion:
		return lox.Resolver(fuente, salida)

	default:
		return lox.Ejecutar(fuente, salida)
	}
}
