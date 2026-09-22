---
layout: default
title: Justificación de stack
nav_order: 2
---

# 🛠️ Justificación del stack utilizado

Elegimos armar el intérprete angLOXg (y su versión web con WebAssembly) usando **Go (Golang)** por tres motivos principales. Buscábamos algo que no solo nos sirviera para el trabajo práctico, sino que nos dejara experiencia real y tuviera buenas ventajas técnicas:

*   **Ganas de aprender algo súper actual:** Go se usa muchísimo hoy en día para armar sistemas, infraestructura en la nube y herramientas de consola. Hacer el proyecto en este lenguaje fue la excusa perfecta para meter las manos en la masa, entender cómo funciona su ecosistema, probar sus tests nativos y agarrarle la mano a una tecnología muy solicitada actualmente.
*   **Rendimiento y velocidad:** Como Go es un lenguaje tipado que compila directo a código máquina (además de tener un *garbage collector* muy bien optimizado), intuiamos que tendría un muy buen rendimiento. Nos daba mucha intriga ver qué tan eficiente y rápido iba a terminar siendo nuestro motor a la hora de procesar los tokens, armar el árbol de sintaxis (AST) y ejecutar todo en tiempo real.
*   **Salir de la zona de confort:** Como el código de referencia que vimos en clase estaba en Python, elegir Go nos obligó a pensar distinto. Son lenguajes con filosofías tan diferentes que era imposible traducir el código línea por línea. Esto nos vino genial para el aprendizaje: nos forzó a entender de verdad la teoría que hay de fondo para después armarlo desde cero "a la manera de Go", usando sus interfaces, estructuras y su forma de controlar los errores.
