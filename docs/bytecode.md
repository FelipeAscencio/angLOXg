---
layout: default
title: 'Fase 2: Bytecode Compiler'
nav_order: 8
---

# ⚡ Fase 2: Bytecode Compiler

En esta segunda fase, evolucionamos el motor de **angLOXg** hacia una arquitectura orientada al alto rendimiento: un compilador de **Bytecode** acoplado a una **Máquina Virtual de Pila** (*Stack-based Virtual Machine*).

A diferencia del intérprete *Tree-Walk*, que construye y recorre estructuras de árboles dispersas en memoria, este motor traduce el código fuente en una secuencia lineal y densa de instrucciones binarias (*opcodes*), logrando grandes mejoras en velocidad de ejecución.

## Arquitectura y componentes

Para alcanzar este salto de eficiencia, rediseñamos el *backend* del intérprete estructurando las responsabilidades en cuatro pilares fundamentales:

-   **1. Compilador Single-Pass (Pratt Parser):** En lugar de generar un Árbol de Sintaxis Abstracta (AST) intermedio y luego recorrerlo, implementamos un compilador de una sola pasada basado en un **Pratt Parser** (Análisis sintáctico por precedencia de operadores). A medida que consume los tokens, resuelve la gramática, calcula los desplazamientos de variables locales o *upvalues* y emite el bytecode de forma directa.
-   **2. Chunks y pool de constantes (`internal/compilador`):** Cada función o script compilado empaqueta su código en una estructura llamada `Chunk`. Un *Chunk* contiene un arreglo plano y contiguo de bytes (`[]uint8`) con los códigos de operación (*opcodes* como `OpAdd`, `OpJumpIfFalse`, `OpCall`, `OpClosure`, `OpClass`, `OpArray`), acompañado de una tabla de constantes donde se almacenan los literales (números, cadenas y objetos función) indexados para un acceso inmediato.
-   **3. Máquina virtual de pila (`internal/vm`):** Implementa la lógica de ejecución en tiempo real. Funciona bajo un ciclo continuo de *Fetch-Decode-Execute* (buscar instrucción, decodificar el *opcode* y ejecutar). Administra una **pila de valores** para resolver expresiones intermedias y una pila de **Marcos de llamada (`CallFrames`)** que mantiene el puntero de instrucción (*IP*) y la ventana de variables locales de cada función activa sin necesidad de alojar mapas dinámicos por cada invocación.
-   **4. Modelo de objetos y memoria en tiempo de ejecución:** Representa las entidades complejas del lenguaje que viven durante la ejecución en la VM: *Closures* con captura de estado mediante *Upvalues* abiertos y cerrados (`OpCloseUpvalue`), estructuras de **Clases (`OpClass`)**, **Instancias**, **Métodos ligados (`OpMethod`)**, **Herencia (`OpInherit` / `OpGetSuper`)** y **Arrays dinámicos (`OpArray` / `OpGetIndex`)**.

## Reutilización e integración con la Fase 1 (Tree-Walk)

Uno de los objetivos centrales en el diseño de **angLOXg** fue mantener una arquitectura modular que evitara la duplicación de código entre ambas fases.

Gracias al bajo acoplamiento planteado desde la primera entrega, logramos reutilizar componentes críticos de forma transparente:

-   **Frontend Léxico compartido (`internal/token` e `internal/escaner`):** El compilador de Bytecode consume exactamente los mismos paquetes de definición de tokens y el mismo **Escáner** desarrollados para el *Tree-Walk*. Toda la lógica de reconocimiento de lexemas, literales, palabras reservadas y manejo de errores léxicos se comparte al 100% como etapa inicial de ambos motores.
-   **Orquestador e Interfaz Unificada (`internal/lox` y `cmd/`):** El paquete central `lox` actúa como fachada común. Reutilizamos la infraestructura de lectura de archivos, captura de buffers de salida (`io.Writer`), interfaz de línea de comandos (CLI) y puente con **WebAssembly (WASM)**, inyectando el nuevo flujo mediante `lox.EjecutarVM()` sin alterar el contrato externo de la aplicación.
-   **Suite de Validación y Casos de Uso (`tests/`):** Toda la batería de pruebas de integración de la cátedra, los scripts de casos de uso reales y la suite de *benchmarks* fueron reutilizados directamente para validar la equivalencia semántica entre el intérprete *Tree-Walk* y la *Bytecode VM*, garantizando que las innovaciones del lenguaje (POO, Arrays y Operador Ternario) se comporten de forma idéntica en ambos motores.

## Ejecución en modo Máquina Virtual

Para poner en marcha el compilador de Bytecode y ejecutar los programas sobre la Máquina Virtual, incorporamos el flag `--vm` en la interfaz del proyecto:

1.  **Desde la consola (CLI):** Agregando el flag `--vm` al ejecutar cualquier script `.lox`, el orquestador desvía el flujo después del escaneo léxico directamente hacia `compilador.Compilar(tokens)` y `vm.Ejecutar(funcion)`.
2.  **Paridad total de funcionalidades:** Este modo soporta el 100% de la especificación de Lox junto con todas nuestras extensiones (clases, constructores `init`, `this`, herencia con `super`, arreglos indexados y el operador ternario `? :`), ejecutándolos a una fracción del tiempo requerido por el recorrido de árboles tradicional.
