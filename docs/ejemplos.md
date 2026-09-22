---
layout: default
title: Ejemplos de uso
nav_order: 5
---

# 🎯 Ejemplos de uso real

Esta sección presenta una serie de programas demostrativos desarrollados íntegramente en el lenguaje **Lox**. Su propósito es validar tanto la capacidad del intérprete como su aplicabilidad en escenarios de procesamiento reales, yendo más allá de las pruebas tradicionales de sintaxis.

## 1. Motor analítico para procesamiento y estadísticas de series de datos (`sensor_analiticas.lox`)

* **Objetivo:** Simular un pipeline de procesamiento de métricas recolectadas por sensores industriales o dispositivos IoT.
* **Qué demuestra:** 
  * Manipulación de variables numéricas y lógicas.
  * Funciones modulares dedicadas a operaciones matemáticas de agregación (mínimo, máximo y promedio aritmético).
  * Lógica de sanitización y filtrado de datos anómalos (detección de valores fuera de rangos seguros y aplicación de acotamientos).

## 2. Sistema de gestión de descuentos escalados y transacciones (`motor_ecommerce.lox`)

* **Objetivo:** Calcular el costo final de órdenes de compra aplicando reglas de negocio comerciales complejas.
* **Qué demuestra:**
  * **Closures:** Uso de funciones generadoras (`crearMotorDeDescuento`) que capturan y recuerdan estados de su entorno exterior (como el factor de lealtad del cliente).
  * Lógica condicional de múltiples niveles y evaluación de reglas de negocio anidadas (descuentos por volumen e impuestos regionales geográficos).

## 3. Motor de validación de reglas de negocio y riesgo crediticio (`motor_financiero.lox`)

* **Objetivo:** Evaluar la elegibilidad de solicitudes financieras bajo estrictas normativas de cumplimiento y riesgo.
* **Qué demuestra:**
  * Ámbitos léxicos complejos y validación cruzada de variables institucionales.
  * **Funciones recursivas:** Algoritmos iterativos avanzados (`proyectarRiesgo`) para la proyección matemática de variables financieras a futuro.
  * Estructuras de control de flujo robustas con ramificaciones anidadas (`if/else`).

> **Nota:** Los 3 archivos fuente de estos ejemplos pueden visualizarse y ejecutarse directamente desde el directorio `tests/casos_de_uso` del repositorio.
