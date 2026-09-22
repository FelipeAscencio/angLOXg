// Referencias a los elementos del DOM de la interfaz web.
const outputContainer = document.getElementById("output");
const codeInput = document.getElementById("code-input");
const runBtn = document.getElementById("run-btn");
const clearBtn = document.getElementById("clear-btn");
const btnExample1 = document.getElementById("btn-example-1");
const btnExample2 = document.getElementById("btn-example-2");
const replInput = document.getElementById("repl-input");
const replHistory = document.getElementById("repl-history");
const replClearBtn = document.getElementById("repl-clear-btn");

// Inicialización del entorno WebAssembly y carga del motor de Go.
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    outputContainer.textContent = "Sistema Wasm inicializado correctamente. angLOXg listo.";
}).catch((err) => {
    outputContainer.textContent = "Error al cargar el motor Wasm: " + err;
    outputContainer.style.color = "#ef4444";
});

// Ejecuta el código fuente completo ingresado en el editor principal.
runBtn.addEventListener("click", () => {
    const code = codeInput.value;
    if (!code.trim()) {
        outputContainer.textContent = "Error: El código fuente está vacío.";
        return;
    }

    try {
        if (typeof runLox === "function") {
            const result = runLox(code);
            outputContainer.textContent = result || "[Ejecución finalizada sin salida]";
        } else {
            outputContainer.textContent = "Error: El motor Wasm todavía no está listo.";
        }
    } catch (e) {
        outputContainer.textContent = "Excepción de ejecución: " + e;
    }
});

// Limpia el contenido del editor de código y la consola de salida.
clearBtn.addEventListener("click", () => {
    codeInput.value = "";
    outputContainer.textContent = "Consola limpia.";
});

// Carga el ejemplo básico de expresiones y operaciones aritméticas en el editor.
btnExample1.addEventListener("click", () => {
    codeInput.value = 'print "Hola desde la Web!";\nprint 5 + 3 * 2;';
});

// Carga el ejemplo avanzado de función recursiva (Fibonacci) en el editor.
btnExample2.addEventListener("click", () => {
    codeInput.value = 'fun fibonacci(n) {\n    if (n <= 1) return n;\n    return fibonacci(n - 1) + fibonacci(n - 2);\n}\n\nprint fibonacci(10);';
});

// Gestiona la entrada interactiva de comandos línea por línea en la terminal REPL.
replInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        const code = replInput.value.trim();
        if (!code) return;

        appendReplLine(`> ${code}`, "input");
        replInput.value = "";

        try {
            if (typeof evalLox === "function") {
                const result = evalLox(code);
                if (result) {
                    appendReplLine(result, "output");
                }
            } else {
                appendReplLine("Error: El motor Wasm no está inicializado.", "error");
            }
        } catch (e) {
            appendReplLine("Excepción: " + e, "error");
        }

        replHistory.scrollTop = replHistory.scrollHeight;
    }
});

// Agrega una línea formateada al historial visual de la terminal REPL.
function appendReplLine(text, type) {
    const line = document.createElement("div");
    line.className = `repl-line ${type}`;
    line.textContent = text;
    replHistory.appendChild(line);
}

// Restablece el historial de la terminal interactiva a su estado inicial.
replClearBtn.addEventListener("click", () => {
    replHistory.innerHTML = '<div class="repl-line text-muted">Terminal limpia. Escribí una nueva sentencia...</div>';
});
