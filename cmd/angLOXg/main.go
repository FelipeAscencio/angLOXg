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

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "uso: angloxg [--escaneo] [script.lox]")
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
		os.Exit(ejecutarArchivo(flag.Arg(0), *modoEscaneo))
	}

	os.Exit(ejecutarConsola(*modoEscaneo))
}

// ejecutarArchivo corre un script completo y devuelve el código de salida.
func ejecutarArchivo(ruta string, modoEscaneo bool) int {
	fuente, err := os.ReadFile(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo leer %q: %v\n", ruta, err)
		return salidaSinArchivo
	}

	if !despachar(string(fuente), os.Stdout, modoEscaneo) {
		return salidaErrorDeDatos
	}

	return salidaOK
}

// ejecutarConsola abre la consola interactiva, que lee y corre una línea por
// vez.
//
// Un error no corta la sesión: se reporta y se sigue esperando la próxima
// sentencia.
func ejecutarConsola(modoEscaneo bool) int {
	fmt.Println("=== angLOXg ===")
	fmt.Println("Intérprete de Lox escrito en Go")

	entrada := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		// Scan devuelve false al llegar al fin de la entrada.
		if !entrada.Scan() {
			break
		}

		despachar(entrada.Text(), os.Stdout, modoEscaneo)
	}

	if err := entrada.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error leyendo la entrada: %v\n", err)
		return salidaErrorDeDatos
	}

	// El salto de línea deja prolija la consola del sistema después del Ctrl+D.
	fmt.Println()

	return salidaOK
}

// despachar manda el fuente al modo correspondiente e informa si salió todo
// bien.
func despachar(fuente string, salida io.Writer, modoEscaneo bool) bool {
	if modoEscaneo {
		return lox.Escanear(fuente, salida)
	}

	return lox.Ejecutar(fuente, salida)
}
