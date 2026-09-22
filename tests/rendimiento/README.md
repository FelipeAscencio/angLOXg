# Pruebas de rendimiento

Esta carpeta contiene una batería de 7 scripts de Lox diseñados específicamente para estresar y medir el rendimiento del motor **angLOXg** (Y de otras implementaciones para poder generar comparativas).

Estos programas evalúan *qué tan rápido y eficiente* es el intérprete al ejecutar distintas cargas de trabajo.

Cada script ataca un cuello de botella distinto, de los posibles en la arquitectura de los intérpretes *Tree-Walk*:

## 1. `bench_loop.lox` (CPU y variables locales)

* **Objetivo:** Medir la velocidad bruta del Evaluador.
* **Qué hace:** Ejecuta un bucle `for` de un millón de iteraciones realizando una suma aritmética simple. Mide la sobrecarga (*overhead*) de iterar y actualizar variables en el entorno léxico local.

## 2. `bench_fib.lox` (Pila de llamadas y recursividad)

* **Objetivo:** Evaluar el manejo de *Call Frames* (Marcos de llamada).
* **Qué hace:** Calcula el número 32 de la secuencia de Fibonacci de forma puramente recursiva. Obliga al motor a crear, apilar y destruir millones de contextos de ejecución de funciones por segundo.

## 3. `bench_closures.lox` (Asignación en memoria)
* **Objetivo:** Estresar la captura de variables en funciones anónimas (*Closures*).
* **Qué hace:** Instancia decenas de miles de funciones dentro de un bucle que capturan el estado de su entorno superior. Verifica la eficiencia del manejo de memoria y la precisión del recolector de basura.

## 4. `bench_strings.lox` (Manipulación de inmutables)

* **Objetivo:** Estresar la gestión dinámica de strings.
* **Qué hace:** Concatenación masiva de texto. Como los strings en Lox son inmutables, esto fuerza al intérprete a alojar constantemente nuevos bloques de memoria y a descartar los anteriores, saturando el *Garbage Collector*.

## 5. `bench_logic.lox` (Ramificación)

* **Objetivo:** Medir la latencia en las decisiones lógicas y saltos condicionales.
* **Qué hace:** Ejecuta operaciones `and` / `or` y sentencias `if/else` a gran escala. Mide qué tan rápido el intérprete deforma el flujo de ejecución evaluando ramas del AST.

## 6. `bench_scopes.lox` (Resolución de entornos)

* **Objetivo:** Evaluar la eficiencia del `Resolver` y el acceso a variables encadenadas.
* **Qué hace:** Anida múltiples bloques de código (`{ ... }`) creando variables en distintos niveles de profundidad. Obliga al motor a recorrer la cadena de *scopes* para encontrar el valor correcto de las referencias.

## 7. `bench_math.lox` (Sobrecarga del árbol AST)

* **Objetivo:** Medir el peso del patrón en el árbol de sintaxis abstracta.
* **Qué hace:** Resuelve expresiones matemáticas densas (con múltiples paréntesis y operadores unarios/binarios). Evalúa cuánto le cuesta al intérprete saltar entre los distintos nodos del árbol AST de forma recursiva.

## ¿Cómo ejecutar esta suite?

Para correr todos los benchmarks automáticamente y recolectar las métricas de tiempo, ejecutar el siguiente comando:

```bash

./benchmarks.sh

```
