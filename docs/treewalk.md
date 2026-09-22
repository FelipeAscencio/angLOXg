---
layout: default
title: 'Fase 1: Tree-Walk Interpreter'
nav_order: 7
---

# 🌳 Fase 1: Tree-Walk Interpreter

El corazón de angLOXg en esta primera fase es un *Tree-Walk Interpreter* clásico. Esto significa que el motor no compila el código a lenguaje máquina de forma directa, sino que procesa el código fuente a través de una tubería (*pipeline*) de transformaciones hasta armar una estructura de datos en memoria, la cual luego recorre nodo por nodo para ejecutar las instrucciones.

## Arquitectura y componentes

Para lograr esto, estructuramos la arquitectura del intérprete dividiendo las responsabilidades en cuatro módulos principales y bien desacoplados:

*   **1. Escáner (Análisis léxico):** Es el punto de entrada. Su trabajo es leer el texto plano del código fuente carácter por carácter y agruparlos en piezas con significado lógico llamadas **Tokens** (Palabras clave, literales, operadores, etc.). Filtra los espacios en blanco y los comentarios, dejando el código listo para la siguiente etapa.
*   **2. Parser (Análisis sintáctico):** Recibe la lista plana de tokens y construye el **Árbol de Sintaxis Abstracta (AST)** utilizando la técnica de *Recursive Descent* (Descenso recursivo). Este módulo entiende la gramática de Lox, respeta la precedencia de los operadores matemáticos y lógicos, y organiza el código en una jerarquía de sentencias y expresiones.
*   **3. Resolver (Análisis semántico):** Antes de que el código empiece a correr de verdad, este componente hace una "pasada en seco" por todo el AST. Su única responsabilidad es resolver los ámbitos léxicos (*scopes*) y rastrear dónde se declara y se usa cada variable. Esto es fundamental para que características avanzadas como los *closures* funcionen de forma rápida y sin errores.
*   **4. Evaluador (Intérprete):** Es el motor final que le da vida al programa. Recorre el AST validado, evalúa las expresiones matemáticas, administra los entornos en memoria (guardando y actualizando las variables) y ejecuta los efectos colaterales, como imprimir datos en pantalla o invocar funciones.

## 🎛️ Modos de ejecución disponibles

Para poder inspeccionar a fondo qué hace el motor por dentro, diagnosticar problemas y demostrar visualmente cómo funciona cada etapa del *pipeline*, implementamos **4 modos de ejecución distintos**. 

Estos modos se pueden utilizar tanto enviando *flags* específicos desde la consola (CLI), como seleccionándolos desde el menú desplegable en nuestro Playground Web:

1.  **Modo escáner (Tokens):** Detiene la ejecución justo después del análisis léxico e imprime en pantalla todos los tokens detectados secuencialmente, detallando su tipo y lexema.
2.  **Modo parser (Árbol AST):** Frena el proceso luego de construir el árbol y muestra una representación estructurada en texto para entender cómo el intérprete agrupó lógicamente cada pieza de código.
3.  **Modo resolución semántica:** Ejecuta el escáner, el parser y el resolver, pero se detiene antes de evaluar el código. Es ideal para validar estáticamente que no haya variables sueltas, conflictos de ámbito o errores de referencias antes de correr el programa.
4.  **Modo ejecución completa:** Es el flujo de trabajo total. Atraviesa todas las fases anteriores de forma transparente y le entrega el control al Evaluador para que ejecute el script y devuelva el resultado final.
