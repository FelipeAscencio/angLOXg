//go:build js && wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// ejecutar corre código Lox capturando toda la salida en un buffer, para poder
// devolvérsela a JavaScript como un string.
func ejecutar(args []js.Value) string {
	if len(args) == 0 {
		return ""
	}

	var buf bytes.Buffer
	lox.Ejecutar(args[0].String(), &buf)

	return buf.String()
}

// ejecutarDesdeEditor alimenta el panel del editor del playground.
func ejecutarDesdeEditor(this js.Value, args []js.Value) any {
	return ejecutar(args)
}

// ejecutarDesdeRepl alimenta el panel del REPL del playground.
func ejecutarDesdeRepl(this js.Value, args []js.Value) any {
	return ejecutar(args)
}

func main() {
	// Los nombres globales quedan en inglés porque son los que busca la página
	// web en "web/app.js".
	js.Global().Set("runLox", js.FuncOf(ejecutarDesdeEditor))
	js.Global().Set("evalLox", js.FuncOf(ejecutarDesdeRepl))

	// Bloqueamos para siempre: si main termina, el runtime de Go se apaga y las
	// funciones que acabamos de exponer dejan de existir.
	select {}
}
