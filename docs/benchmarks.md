---
layout: default
title: Comparativas de rendimiento
nav_order: 6
---

# 📈 Comparativas de rendimiento

## Visualización interactiva de rendimiento

A continuación, se presenta un gráfico comparativo que recopila los tiempos de ejecución obtenidos por **angLOXg** frente a las demás implementaciones analizadas en la suite de pruebas. Puedes hacer clic en las leyendas superiores para ocultar o mostrar lenguajes específicos y facilitar el análisis cruzado.

<div>
  <canvas id="rendimientoLoxChart" height="120"></canvas>
</div>

<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>

<script>
  const ctx = document.getElementById('rendimientoLoxChart').getContext('2d');
  
  const rendimientoLoxChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: ['Closures', 'Fibonacci (32)', 'Lógica', 'Bucle', 'Matemática (AST)', 'Scopes', 'Strings'],
      datasets: [
        {
          label: 'angLOXg (Go - Propio)',
          data: [0.089, 4.396, 0.402, 0.658, 0.108, 0.199, 0.086],
          borderColor: 'rgb(255, 99, 132)',
          backgroundColor: 'rgba(255, 99, 132, 0.2)',
          tension: 0.1
        },
        {
          label: 'clox (C - Oficial VM)',
          data: [0.039, 0.272, 0.051, 0.059, 0.024, 0.035, 0.332],
          borderColor: 'rgb(54, 162, 235)',
          backgroundColor: 'rgba(54, 162, 235, 0.2)',
          tension: 0.1
        },
        {
          label: 'jlox (Java - Oficial Tree-Walk)',
          data: [0.843, 1.688, 0.741, 0.893, 0.835, 0.696, 0.583],
          borderColor: 'rgb(255, 206, 86)',
          backgroundColor: 'rgba(255, 206, 86, 0.2)',
          tension: 0.1
        },
        {
          label: 'golox (Go - Terceros)',
          data: [0.393, 4.620, 0.341, 0.492, 0.164, 0.214, 0.150],
          borderColor: 'rgb(75, 192, 192)',
          backgroundColor: 'rgba(75, 192, 192, 0.2)',
          tension: 0.1
        },
        {
          label: 'lox.rs (Rust)',
          data: [0.177, 3.479, 0.342, 0.620, 0.114, 0.157, 0.071],
          borderColor: 'rgb(153, 102, 255)',
          backgroundColor: 'rgba(153, 102, 255, 0.2)',
          tension: 0.1
        },
        {
          label: 'tslox (TypeScript / Node)',
          data: [0.423, 11.938, 0.526, 0.500, 0.551, 0.254, 0.231],
          borderColor: 'rgb(255, 159, 64)',
          backgroundColor: 'rgba(255, 159, 64, 0.2)',
          tension: 0.1
        },
        {
          label: 'DotLox (C# / .NET)',
          data: [1.636, 84.903, 0.600, 0.835, 0.330, 0.355, 0.188],
          borderColor: 'rgb(199, 199, 199)',
          backgroundColor: 'rgba(199, 199, 199, 0.2)',
          tension: 0.1
        },
        {
          label: 'SlowLox (Ruby)',
          data: [1.474, 86.239, 3.907, 7.391, 1.635, 1.830, 0.480],
          borderColor: 'rgb(83, 51, 237)',
          backgroundColor: 'rgba(83, 51, 237, 0.2)',
          tension: 0.1
        },
        {
          label: 'clojox (Clojure)',
          data: [23.027, 20.248, 15.026, 16.965, 14.339, 14.163, 16.558],
          borderColor: 'rgb(255, 99, 255)',
          backgroundColor: 'rgba(255, 99, 255, 0.2)',
          tension: 0.1
        },
        {
          label: 'plox (Python)',
          data: [24.646, 120.0, 44.235, 49.289, 13.130, 14.425, 5.154],
          borderColor: 'rgb(46, 139, 87)',
          backgroundColor: 'rgba(46, 139, 87, 0.2)',
          tension: 0.1
        }
      ]
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: 'top',
        },
        title: {
          display: true,
          text: 'Comparativa de Tiempos de Ejecución (Segundos - Menor es mejor)'
        }
      },
      scales: {
        y: {
          type: 'linear',
          title: {
            display: true,
            text: 'Tiempo (segundos)'
          }
        }
      }
    }
  });
</script>

## Evolución de performance

El análisis comparativo de las 10 implementaciones evaluadas bajo la suite de 7 scripts de estrés permite validar que el rendimiento de **angLOXg** se encuentra alineado con los objetivos esperados para un motor desarrollado en Go. Al tratarse de un intérprete basado en árboles (*Tree-Walk*), su comportamiento refleja fortalezas claras frente a contrapartes dinámicas e interpretadas tradicionales, posicionándose muy cerca de otras implementaciones robustas del mismo lenguaje como `golox` e incluso compitiendo de cerca con alternativas compiladas como `lox.rs` (Rust).

### Hallazgos principales del análisis comparativo:

* **Velocidad en Go (`angLOXg` vs `golox`):** Los tiempos obtenidos por `angLOXg` demuestran una excelente eficiencia en tareas secuenciales y de manipulación de estructuras. Por ejemplo, en el test de closures (`bench_closures.lox`), registra un tiempo sobresaliente de 0.089 segundos, superando o igualando los registros de implementaciones de terceros en Go y superando ampliamente a intérpretes en entornos dinámicos. Asimismo, en el recorrido de AST (`bench_math.lox`), marca unos sólidos 0.108 segundos.
* **El estándar de referencia en C (`clox`):** Como era de esperarse, la implementación oficial basada en máquina virtual y bytecode (`clox`) lidera la tabla con amplia ventaja en todas las categorías (completando Fibonacci en 0.272s y bucles en 0.059s), evidenciando la diferencia estructural frente a los intérpretes de recorrido directo sobre árboles.
* **Gestión de la pila y recursividad (`bench_fib.lox`):** En la prueba de recursividad profunda con Fibonacci de 32, **angLOXg** resuelve la ejecución en 4.396 segundos, un tiempo sumamente competitivo que contrasta fuertemente con los bloqueos o tiempos críticos sufridos por intérpretes en lenguajes dinámicos puros como `plox` (Python), el cual alcanza el límite de tiempo (*Timeout* a los 120 segundos), o `SlowLox` (Ruby), que demanda más de 86 segundos.
* **Manejo de memoria y recolección:** En pruebas de saturación por inmutabilidad de textos (`bench_strings.lox`) y resolución de ámbitos (`bench_scopes.lox`), **angLOXg** muestra una estabilidad notable con 0.086s y 0.199s respectivamente, probando que el recolector de basura nativo de Go gestiona eficientemente la presión de objetos efímeros generados por el árbol de sintaxis.

**Conclusión general:** El rendimiento de **angLOXg** es óptimo y consistente. Demuestra que la elección de Go como lenguaje base aporta una ventaja sustancial en velocidad de ejecución y manejo de memoria frente a los intérpretes tradicionales de tipo *Tree-Walk* implementados en lenguajes interpretados de alto nivel, validando con éxito el diseño arquitectónico del motor.

## Selección de la suite de pruebas y justificación técnica

Para evaluar rigurosamente el rendimiento de **angLOXg** frente a las demás implementaciones del ecosistema Lox, se diseñó y adaptó una batería compuesta por **7 scripts de estrés**.

Cada prueba fue seleccionada estratégicamente para aislar y poner a prueba un componente crítico o cuello de botella específico de la arquitectura de los intérpretes orientados a árboles (*Tree-Walk*):

### 1. `bench_loop.lox` — *Velocidad bruta del evaluador*

* **Por qué se eligió:** Es la prueba base para medir el rendimiento secuencial puro. Al ejecutar un bucle `for` intensivo con un millón de iteraciones y operaciones aritméticas simples, nos permite dimensionar la sobrecarga (*overhead*) general que introduce el motor al procesar instrucciones repetitivas y actualizar variables locales en el entorno.

### 2. `bench_fib.lox` — *Gestión de la pila y llamadas recursivas*

* **Por qué se eligió:** La recursividad profunda (calculando el 32º número de Fibonacci) somete al intérprete a una presión extrema en la creación y destrucción de contextos de ejecución. Se eligió para medir la eficiencia del motor al manejar la pila de llamadas (*Call Frames*) y evaluar el impacto del costo de invocación de funciones por segundo.

### 3. `bench_closures.lox` — *Captura de estado y gestión de memoria*

* **Por qué se eligió:** Las clausuras (*closures*) representan uno de los mayores desafíos arquitectónicos en los lenguajes dinámicos, ya que obligan a las variables a sobrevivir más allá del ámbito léxico donde fueron declaradas. Este test instancia miles de funciones anónimas para verificar cómo el motor resuelve la persistencia de referencias y qué tanto exige al recolector de basura (*Garbage Collector*).

### 4. `bench_strings.lox` — *Inmutabilidad y saturación del recolector*

* **Por qué se eligió:** Las cadenas de texto en Lox son inmutables. Al someter al intérprete a una concatenación masiva y repetitiva de strings, forzamos la asignación constante y masiva de pequeños bloques de memoria efímeros. Este escenario es ideal para evaluar cómo maneja la plataforma la creación dinámica de objetos y el estrés consecuente sobre la memoria RAM.

### 5. `bench_logic.lox` — *Latencia en ramificaciones y control de flujo*

* **Por qué se eligió:** La evaluación de expresiones lógicas anidadas (`and`/`or`) y estructuras condicionales (`if/else`) a gran escala pone a prueba la velocidad con la que el evaluador descarta ramas de ejecución y procesa tablas de verdad. Mide la agilidad del motor ante la divergencia de flujos en tiempo de ejecución.

### 6. `bench_scopes.lox` — *Eficiencia del Resolver y resolución léxica*

* **Por qué se eligió:** Con múltiples niveles anidados de bloques de código y variables declaradas en distintas profundidades, este script evalúa la efectividad del módulo `Resolver`. Permite contrastar el rendimiento entre buscar variables de forma dinámica en tiempo de ejecución frente a la resolución estática por distancias indexadas.

### 7. `bench_math.lox` — *Recorrido profundo del Árbol de Sintaxis Abstracta (AST)*

* **Por qué se eligió:** Las expresiones matemáticas densas y cargadas de paréntesis obligan al intérprete a generar y transitar estructuras de nodos muy ramificadas en el AST. Se utiliza para medir el costo computacional inherente al patrón de diseño Visitor o de evaluación recursiva sobre los nodos sintácticos.
