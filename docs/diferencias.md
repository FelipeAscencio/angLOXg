---
layout: default
title: Diferencias de implementación y agregados
nav_order: 9
---

# 🧬 Diferencias de implementación y agregados

## Funcionalidades extra (Nuestras innovaciones)

Para esta primera entrega, decidimos ir un paso más allá de los requisitos básicos del trabajo práctico y agregar herramientas que nos faciliten tanto el desarrollo como la demostración del intérprete. Nuestros dos grandes agregados fueron:

*   **Intérprete web (Playground vía WebAssembly):** Nuestra mayor innovación fue llevar angLOXg al navegador. Gracias a la capacidad de Go para compilar a WebAssembly (WASM), logramos que el intérprete corra 100% del lado del cliente sin necesitar un servidor backend. Armamos una interfaz web con un editor de scripts completo y una consola interactiva (REPL) en vivo.
*   **Flags de inspección (Modos de ejecución):** Tanto en la consola tradicional (CLI) como en la página web, agregamos la posibilidad de "pausar" la ejecución en distintas etapas del intérprete. Esto es ideal para entender qué pasa por debajo del capó. Implementamos 4 modos visuales:
    1.  **Ejecución completa:** El comportamiento normal de Lox.
    2.  **Modo escáner (`--escaneo`):** Muestra únicamente la lista de Tokens generada.
    3.  **Modo parser (`--arbol`):** Muestra la estructura del Árbol de Sintaxis Abstracta (AST) formateada como texto.
    4.  **Modo resolver (`--resolucion`):** Ejecuta hasta la fase de resolución de ámbitos y referencias, frenando antes de evaluar el código.

## El impacto del lenguaje: Go vs. Python

La implementación de referencia de la cátedra (`plox`) está escrita en Python. Al pasar toda esa lógica a **Go**, nos topamos con diferencias arquitectónicas enormes debido a la naturaleza de cada lenguaje. Estas fueron las principales decisiones de diseño que tuvimos que adaptar:

*   **Tipado estático vs. Tipado dinámico:** 
    *   *En Python:* Lox es dinámico y Python también, por lo que guardar un número, un string o un booleano en una variable de Python es directo.
    *   *En Go:* Al ser de tipado estático estricto, tuvimos que representar los tipos dinámicos de Lox usando el tipo `any` (o `interface{}`). Esto nos obligó a usar *Type Assertions* (`switch v := valor.(type)`) constantemente en el Evaluador para verificar si estábamos sumando dos números o concatenando dos strings, lo que hace al código de Go un poco más verboso pero mucho más seguro en tiempo de compilación.
*   **Manejo de errores (Excepciones vs. Valores de retorno):**
    *   *En Python:* Se abusa bastante de los bloques `try/except`. Si hay un error de sintaxis o de ejecución, simplemente se "lanza" (`raise`) una excepción personalizada (como `RuntimeError`). Incluso la sentencia `return` de las funciones de Lox se suele implementar lanzando una excepción `Return` para cortar la pila de llamadas.
    *   *En Go:* Go no tiene excepciones tradicionales. Tuvimos que propagar los errores manualmente devolviendo múltiples valores (por ejemplo, `valor, err := ...`) y chequeando `if err != nil` en cada paso. Para simular el `return` profundo de las funciones de Lox sin romper la filosofía de Go, tuvimos que usar un manejo muy cuidadoso de estructuras de error personalizadas.
*   **Programación orientada a objetos (Clases vs. Structs e interfaces):**
    *   *En Python:* El AST se construye usando herencia clásica (una clase base `Expr` y subclases como `Binary`, `Literal`, etc.).
    *   *En Go:* Como no hay clases ni herencia, modelamos los nodos del AST usando `structs` e implementamos el patrón *Visitor* a través de **Interfaces**. En Go, si un struct tiene los métodos que pide una interfaz, la implementa automáticamente. Esto nos dio un diseño súper modular y desacoplado, muy diferente a la jerarquía de clases de Python.
