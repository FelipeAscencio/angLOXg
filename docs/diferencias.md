---
layout: default
title: Diferencias de implementación y agregados
nav_order: 9
---

# 🧬 Diferencias de implementación y agregados

## Funcionalidades extra (Nuestras innovaciones)

A lo largo del desarrollo del trabajo práctico, se dió un paso más allá de los requisitos básicos y se sumaron tanto herramientas de entorno como extensiones al lenguaje que enriquecen la experiencia de uso y demuestran la flexibilidad de ambas arquitecturas (*Tree-Walk* y *Bytecode VM*).

### Agregados de la entrega inicial (Entorno e inspección)

*   **Intérprete web (Playground vía WebAssembly):** Nuestra primera gran innovación fue llevar **angLOXg** al navegador. Gracias a la capacidad de Go para compilar a WebAssembly (WASM), logramos que el intérprete corra 100% del lado del cliente sin necesitar un servidor backend. Armamos una interfaz web con un editor de scripts completo y una consola interactiva (REPL) en vivo.
*   **Flags de inspección (Modos de ejecución):** Tanto en la consola tradicional (CLI) como en la página web, agregamos la posibilidad de "pausar" o enrutar la ejecución en distintas etapas del motor. Esto es ideal para entender qué pasa por debajo del capó. Implementamos los siguientes modos:
    1.  **Ejecución completa (Tree-Walk):** El comportamiento normal del intérprete evaluando el árbol sintáctico.
    2.  **Modo Máquina Virtual (`--vm`):** Ejecuta el código a través del compilador *Single-Pass* y la Máquina Virtual de Bytecode (Este modo fue agregado para la entrega final).
    3.  **Modo escáner (`--escaneo`):** Muestra únicamente la lista de Tokens generada por el análisis léxico.
    4.  **Modo parser (`--arbol`):** Muestra la estructura del Árbol de Sintaxis Abstracta (AST) formateada como texto.
    5.  **Modo resolver (`--resolucion`):** Ejecuta hasta la fase de resolución estática de ámbitos y referencias, frenando antes de evaluar el código.

### Agregados de la entrega final (Extensiones del lenguaje)

Para el cierre del proyecto, expandimos la gramática y el motor de ejecución incorporando nuevas capacidades de alto nivel, garantizando su compatibilidad total tanto en el intérprete de recorrido de árboles (*Tree-Walk*) como en el compilador de *Bytecode*:

*   **Programación Orientada a Objetos (Clases y herencia):** Implementamos el sistema completo de POO de Lox. Incluye la declaración de clases (`class`), instanciación dinámica, lectura y escritura de propiedades, constructores mediante el método reservado `init()`, ligadura de contexto con `this` y herencia simple (`<`) con resolución de métodos de la superclase a través de la palabra clave `super`.
*   **Soporte nativo de Arrays:** Incorporamos estructuras de listas dinámicas al lenguaje mediante la sintaxis de corchetes (`[a, b, c]`) y el acceso posicional por índice (`arreglo[i]`). Se implementó tanto la evaluación de nodos en el AST como los *opcodes* específicos (`OpArray`, `OpGetIndex`) en la Máquina Virtual, incluyendo validaciones de tipo y límites de rango en tiempo de ejecución.
*   **Operador condicional ternario (`? :`):** Añadimos soporte léxico, sintáctico y semántico para expresiones condicionales compactas de la forma `condicion ? verdadero : falso` con asociatividad a la derecha. En la Máquina Virtual, esto se tradujo en la emisión optimizada de instrucciones de salto condicional en cortocircuito directamente desde el *Pratt Parser*.

## El impacto del lenguaje: Go vs. Python

La implementación de referencia de la cátedra (`plox`) está escrita en Python. Al pasar toda esa lógica a **Go**, nos topamos con diferencias arquitectónicas enormes debido a la naturaleza de cada lenguaje. Estas fueron las principales decisiones de diseño que tuvimos que adaptar:

*   **Tipado estático vs. Tipado dinámico:** 
    *   *En Python:* Lox es dinámico y Python también, por lo que guardar un número, un string o un booleano en una variable de Python es directo.
    *   *En Go:* Al ser de tipado estático estricto, tuvimos que representar los tipos dinámicos de Lox usando el tipo `any` (o `interface{}`). Esto nos obligó a usar *Type Assertions* (`switch v := valor.(type)`) constantemente en el Evaluador y en la VM para verificar si estábamos sumando dos números o concatenando dos strings, lo que hace al código de Go un poco más verboso pero mucho más seguro en tiempo de compilación.
*   **Manejo de errores (Excepciones vs. Valores de retorno):**
    *   *En Python:* Se abusa bastante de los bloques `try/except`. Si hay un error de sintaxis o de ejecución, simplemente se "lanza" (`raise`) una excepción personalizada (como `RuntimeError`). Incluso la sentencia `return` de las funciones de Lox se suele implementar lanzando una excepción `Return` para cortar la pila de llamadas.
    *   *En Go:* Go no tiene excepciones tradicionales. Tuvimos que propagar los errores manualmente devolviendo múltiples valores (por ejemplo, `valor, err := ...`) y chequeando `if err != nil` en cada paso. Para simular el `return` profundo de las funciones de Lox en el *Tree-Walk* sin romper la filosofía de Go, tuvimos que usar un manejo muy cuidadoso de estructuras de error personalizadas.
*   **Programación orientada a objetos (Clases vs. Structs e interfaces):**
    *   *En Python:* El AST se construye usando herencia clásica (una clase base `Expr` y subclases como `Binary`, `Literal`, etc.).
    *   *En Go:* Como no hay clases ni herencia, modelamos los nodos del AST usando `structs` e implementamos el despacho mediante **Interfaces selladas** y *Type Switches*. En Go, si un struct tiene los métodos que pide una interfaz, la implementa automáticamente. Esto nos dio un diseño súper modular y desacoplado, muy diferente a la jerarquía de clases de Python.
