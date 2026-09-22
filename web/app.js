// Referencias a los elementos del DOM de la interfaz web.
const outputContainer = document.getElementById("output");
const codeInput = document.getElementById("code-input");
const runBtn = document.getElementById("run-btn");
const clearBtn = document.getElementById("clear-btn");
const btnExample1 = document.getElementById("btn-example-1");
const btnExample2 = document.getElementById("btn-example-2");
const btnCase1 = document.getElementById("btn-case-1");
const btnCase2 = document.getElementById("btn-case-2");
const btnCase3 = document.getElementById("btn-case-3");
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

// Referencia al selector del modo de ejecución/flag.
const executionModeSelect = document.getElementById("execution-mode");

// Ejecuta el código fuente completo ingresado en el editor principal según el modo elegido.
runBtn.addEventListener("click", () => {
    const code = codeInput.value;
    const mode = executionModeSelect ? executionModeSelect.value : "run";

    if (!code.trim()) {
        outputContainer.textContent = "Error: El código fuente está vacío.";
        return;
    }

    try {
        if (typeof runLox === "function") {
            const result = runLox(code, mode);
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

// Carga el Caso de Uso 1: Motor analítico para procesamiento de series de sensores.
btnCase1.addEventListener("click", () => {
    codeInput.value = 
`// Calcula el promedio (media aritmética) de una lista simulada de 5 valores numéricos.
fun calcularPromedio(v1, v2, v3, v4, v5) {
    var suma = v1 + v2 + v3 + v4 + v5;
    return suma / 5;
}

// Encuentra el valor máximo entre 5 mediciones.
fun obtenerMaximo(v1, v2, v3, v4, v5) {
    var max = v1;
    if (v2 > max) max = v2;
    if (v3 > max) max = v3;
    if (v4 > max) max = v4;
    if (v5 > max) max = v5;
    return max;
}

// Encuentra el valor mínimo entre 5 mediciones.
fun obtenerMinimo(v1, v2, v3, v4, v5) {
    var min = v1;
    if (v2 < min) min = v2;
    if (v3 < min) min = v3;
    if (v4 < min) min = v4;
    if (v5 < min) min = v5;
    return min;
}

// Filtra y recorta valores fuera de un rango operativo seguro [minSeguro, maxSeguro].
// Si un sensor falla y emite un valor anómalo, lo reemplaza por el límite tolerable.
fun sanitizarLectura(valor, minSeguro, maxSeguro) {
    if (valor < minSeguro) {
        return minSeguro;
    }
    if (valor > maxSeguro) {
        return maxSeguro;
    }
    return valor;
}

print "=== INFORME DE ANÁLISIS DE SENSORES INDUSTRIALES ===";

// Datos crudos de temperatura capturados por un dispositivo IoT
var t1 = 22.5;
var t2 = 999.0; // Inyectamos una Anomalía / Ruido en el sensor
var t3 = 24.1;
var t4 = 21.8;
var t5 = 23.0;

print "Aplicando filtros de integridad de datos...";
var s1 = sanitizarLectura(t1, 0.0, 100.0);
var s2 = sanitizarLectura(t2, 0.0, 100.0);
var s3 = sanitizarLectura(t3, 0.0, 100.0);
var s4 = sanitizarLectura(t4, 0.0, 100.0);
var s5 = sanitizarLectura(t5, 0.0, 100.0);

var promedio = calcularPromedio(s1, s2, s3, s4, s5);
var maximo = obtenerMaximo(s1, s2, s3, s4, s5);
var minimo = obtenerMinimo(s1, s2, s3, s4, s5);

print "Métricas calculadas exitosamente:";
print promedio;
print maximo;
print minimo;`;
});

// Carga el Caso de Uso 2: Sistema de gestión de descuentos escalados (E-Commerce).
btnCase2.addEventListener("click", () => {
    codeInput.value = 
`// Generador de políticas de descuento basado en el volumen de unidades
fun crearMotorDeDescuento(factorLealtad) {
    // Esta función interna captura la variable externa "factorLealtad"
    fun calcular(subtotal, cantidadItems) {
        var descuentoBase = 0;
        
        if (cantidadItems >= 10) {
            descuentoBase = 0.20; // 20% de descuento por volumen alto
        } else {
            if (cantidadItems >= 5) {
                descuentoBase = 0.10; // 10% por volumen moderado
            } else {
                descuentoBase = 0.0;
            }
        }

        // Se suma el factor de lealtad del cliente al descuento
        var descuentoTotal = descuentoBase + factorLealtad;
        if (descuentoTotal > 0.35) {
            descuentoTotal = 0.35; // Tope máximo de descuento permitido por negocio
        }

        return subtotal * (1 - descuentoTotal);
    }
    return calcular;
}

print "=== PROCESAMIENTO DE TRANSACCIÓN COMERCIAL ===";

var precioUnitario = 150.0;
var cantidadComprada = 8;
var subtotalBruto = precioUnitario * cantidadComprada;

print "Subtotal bruto de la orden:";
print subtotalBruto;

// Instanciamos el motor de precios para un cliente VIP (otorga un 5% adicional de lealtad)
var calcularPrecioVIP = crearMotorDeDescuento(0.05);

var totalConDescuento = calcularPrecioVIP(subtotalBruto, cantidadComprada);
print "Total tras aplicar reglas de volumen y lealtad VIP:";
print totalConDescuento;

// Impuesto regional según destino (simulación de lógica condicional compleja)
var esEnvioInternacional = true;
var tasaImpuesto = 0.0;

if (esEnvioInternacional) {
    tasaImpuesto = 0.18; // 18% arancel aduanero
} else {
    tasaImpuesto = 0.05; // 5% impuesto local
}

var totalFinal = totalConDescuento * (1 + tasaImpuesto);
print "Monto final a cobrar (incluyendo impuestos aplicables):";
print totalFinal;`;
});

// Carga el Caso de Uso 3: Motor de validación de riesgo crediticio.
btnCase3.addEventListener("click", () => {
    codeInput.value = 
`// Valida si la antigüedad laboral cumple con los requisitos mínimos según el sector
fun validarAntiguedad(mesesLaborales, esSectorPublico) {
    if (esSectorPublico) {
        return mesesLaborales >= 6; // Sector público exige menor historial
    }
    return mesesLaborales >= 12; // Sector privado exige al menos 1 año
}

// Algoritmo recursivo para calcular la proyección de riesgo acumulado en cuotas
fun proyectarRiesgo(cuotaActual, riesgoBase, acumulado) {
    if (cuotaActual <= 0) {
        return acumulado;
    }
    // Incremento exponencial simulado del riesgo financiero por cada período a futuro
    var factorRiesgo = riesgoBase * 1.05;
    return proyectarRiesgo(cuotaActual - 1, riesgoBase, acumulado + factorRiesgo);
}

print "=== AUDITORÍA AUTOMATIZADA DE CRÉDITO ===";

var clienteNombre = "F संtino"; // Nombre de referencia
var ingresosMensuales = 3500.0;
var deudaActual = 800.0;
var mesesEnEmpleo = 18;
var esPublico = false;

print "Verificando estabilidad laboral...";
var cumpleEstabilidad = validarAntiguedad(mesesEnEmpleo, esPublico);

// Cálculo del ratio de endeudamiento (Deuda / Ingresos)
var ratioEndeudamiento = deudaActual / ingresosMensuales;
print "Ratio de endeudamiento calculado:";
print ratioEndeudamiento;

// Reglas de negocio estrictas para aprobación
var deudaSaludable = ratioEndeudamiento < 0.40;
var creditoAprobado = false;

if (cumpleEstabilidad) {
    if (deudaSaludable) {
        if (ingresosMensuales >= 2000.0) {
            creditoAprobado = true;
        } else {
            print "Rechazado: Ingresos por debajo del umbral mínimo de la categoría.";
        }
    } else {
        print "Rechazado: Excede el límite de endeudamiento permitido.";
    }
} else {
    print "Rechazado: Antigüedad laboral insuficiente.";
}

if (creditoAprobado) {
    print "ESTADO DE SOLICITUD: APROBADA EXITOSAMENTE.";
    // Proyectamos el riesgo financiero para un plan de 5 cuotas
    var riesgoTotalProyectado = proyectarRiesgo(5, 12.5, 0.0);
    print "Riesgo financiero total proyectado a 5 plazos:";
    print riesgoTotalProyectado;
} else {
    print "ESTADO DE SOLICITUD: RECHAZADA POR POLÍTICAS DE RIESGO.";
}`;
});

// Gestiona la entrada interactiva de comandos línea por línea en la terminal REPL usando el modo actual.
replInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        const code = replInput.value.trim();
        const mode = executionModeSelect ? executionModeSelect.value : "run";
        if (!code) return;

        appendReplLine(`> [${mode}] ${code}`, "input");
        replInput.value = "";

        try {
            if (typeof evalLox === "function") {
                const result = evalLox(code, mode);
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
