//go:build js && wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Ejecuta código Lox y captura la salida para retornarla a JavaScript.
func ejecutar(args []js.Value) string {
	if len(args) == 0 {
		return ""
	}

	var buf bytes.Buffer
	lox.Ejecutar(args[0].String(), &buf)

	return buf.String()
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
