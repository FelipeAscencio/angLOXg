//go:build js && wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Ejecuta código Lox según el modo seleccionado (run, tokens, ast, semantica) y captura la salida.
// Ejecuta código Lox según el modo seleccionado (run, tokens, ast, semantica) y captura la salida.
func ejecutar(args []js.Value) string {
	if len(args) == 0 {
		return ""
	}

	codigo := args[0].String()
	modo := "run"
	if len(args) > 1 {
		modo = args[1].String()
	}

	var salida bytes.Buffer

	switch modo {
	case "tokens":
		lox.Escanear(codigo, &salida)

	case "ast":
		lox.Parsear(codigo, &salida)

	case "semantica":
		lox.Resolver(codigo, &salida)

	case "run":
		fallthrough
	default:
		lox.Ejecutar(codigo, &salida)
	}

	return salida.String()
}

// Expone la ejecución para el panel del editor.
func ejecutarDesdeEditor(this js.Value, args []js.Value) any {
	return ejecutar(args)
}

// Expone la ejecución para el panel del REPL.
func ejecutarDesdeRepl(this js.Value, args []js.Value) any {
	return ejecutar(args)
}

func main() {
	// Se exponen las funciones globales requeridas por la interfaz web.
	js.Global().Set("runLox", js.FuncOf(ejecutarDesdeEditor))
	js.Global().Set("evalLox", js.FuncOf(ejecutarDesdeRepl))

	// Bloquea el hilo principal para mantener activo el runtime de WebAssembly.
	select {}
}
