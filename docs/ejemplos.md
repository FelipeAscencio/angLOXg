---
layout: default
title: Ejemplos de uso
nav_order: 5
---

# 🎯 Ejemplos de uso real

Esta sección presenta una serie de programas demostrativos desarrollados íntegramente en el lenguaje **Lox**.

Su propósito es validar tanto la capacidad del intérprete como su aplicabilidad en escenarios de procesamiento reales, poniendo a prueba las extensiones y características incorporadas durante el desarrollo del TP más allá de las pruebas tradicionales de sintaxis.

## 1. Duelo Lógico: Adiviná el Número con "IA" (`duelo_logico.lox`)

**Prueba integradora de innovaciones del TP**

- **Objetivo:** Simular un juego de deducción donde un agente "IA" evalúa un historial precalculado de jugadas para encontrar un número secreto, actuando como la demostración principal que integra simultáneamente todas las extensiones arquitectónicas implementadas sobre el lenguaje original.
- **Qué demuestra:**
  * **Soporte de arrays:** Almacenamiento dinámico de estado y acceso iterativo por índice para procesar el historial de intentos.
  * **Operador ternario:** Evaluación condicional de una sola línea para definir estados de validación y formatear mensajes de salida compactos.
  * **POO (Clases y herencia):** Encapsulamiento del estado en instancias, uso de constructores `init` con la referencia `this`, y enrutamiento dinámico de la lógica de validación hacia la clase base mediante `super`.

## 2. Motor analítico para procesamiento y estadísticas de series de datos (`sensor_analiticas.lox`)

- **Objetivo:** Simular un pipeline de procesamiento de métricas recolectadas por sensores industriales o dispositivos IoT, demostrando la capacidad del lenguaje para manejar flujos de datos continuos.
- **Qué demuestra:** 
  * Manipulación intensiva de variables numéricas y expresiones lógicas.
  * Modularidad mediante funciones dedicadas a operaciones de agregación matemática (mínimo, máximo y promedio aritmético).
  * Lógica de sanitización y filtrado condicional de datos anómalos (detección de valores fuera de rangos seguros y aplicación de acotamientos).

## 3. Sistema de gestión de descuentos escalados y transacciones (`motor_ecommerce.lox`)

- **Objetivo:** Calcular el costo final de órdenes de compra aplicando reglas de negocio comerciales complejas, poniendo a prueba el manejo de memoria y los entornos del intérprete.
- **Qué demuestra:**
  * **Closures (Cierres léxicos):** Uso de funciones generadoras (`crearMotorDeDescuento`) que capturan, retienen y mutan estados de su entorno exterior (como el factor de lealtad del cliente).
  * Lógica condicional de múltiples niveles y evaluación de reglas de negocio anidadas (descuentos por volumen e impuestos regionales geográficos).

## 4. Motor de validación de reglas de negocio y riesgo crediticio (`motor_financiero.lox`)

- **Objetivo:** Evaluar la elegibilidad de solicitudes financieras bajo estrictas normativas de cumplimiento y riesgo, funcionando como una prueba de estrés para el anidamiento de scopes y la pila de llamadas.
- **Qué demuestra:**
  * Ámbitos léxicos complejos y validación cruzada de variables institucionales en distintos niveles de profundidad.
  * **Funciones recursivas:** Algoritmos iterativos avanzados (`proyectarRiesgo`) para la proyección matemática de variables financieras a futuro.
  * Estructuras de control de flujo robustas con ramificaciones anidadas (`if/else`).

> **Nota:** Los 4 archivos fuente de estos ejemplos pueden visualizarse y ejecutarse directamente desde el directorio `tests/casos_de_uso` del repositorio.
