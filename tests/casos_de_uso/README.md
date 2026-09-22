# Casos de uso reales en angLOXg

Esta carpeta contiene una serie de programas demostrativos desarrollados íntegramente en el lenguaje **Lox** para validar tanto la capacidad del mismo, como su aplicabilidad en escenarios de procesamiento reales más allá de las pruebas tradicionales de sintaxis.

## 1. Motor analítico para procesamiento y estadísticas de series de datos (`sensor_analiticas.lox`)

* **Objetivo:** Simular un pipeline de procesamiento de métricas recolectadas por sensores industriales o dispositivos IoT.
* **Qué demuestra:** 
  * Manipulación de variables numéricas y lógicas.
  * Funciones modulares dedicadas a operaciones matemáticas de agregación (Mínimo, máximo y promedio aritmético).
  * Lógica de sanitización y filtrado de datos anómalos (Detección de valores fuera de rangos seguros y aplicación de acotamientos).

### Ejecución

```bash

go run ../../cmd/angLOXg sensor_analiticas.lox

```

## 2. Sistema de gestión de descuentos escalados y transacciones (`motor_ecommerce.lox`)

* **Objetivo:** Calcular el costo final de órdenes de compra aplicando reglas de negocio comerciales complejas.
* **Qué demuestra:**
  * **Closures:** Uso de funciones generadoras (`crearMotorDeDescuento`) que capturan y recuerdan estados de su entorno exterior (Como el factor de lealtad del cliente).
  * Lógica condicional de múltiples niveles y evaluación de reglas de negocio anidadas (Descuentos por volumen e impuestos regionales geográficos).

### Ejecución

```bash

go run ../../cmd/angLOXg motor_ecommerce.lox

```

## 3. Motor de validación de reglas de negocio y riesgo crediticio (`motor_financiero.lox`)

* **Objetivo:** Evaluar la elegibilidad de solicitudes financieras bajo estrictas normativas de cumplimiento y riesgo.
* **Qué demuestra:**
  * Ámbitos léxicos complejos y validación cruzada de variables institucionales.
  * **Funciones recursivas:** Algoritmos iterativos avanzados (`proyectarRiesgo`) para la proyección matemática de variables financieras a futuro.
  * Estructuras de control de flujo robustas con ramificaciones anidadas (`if/else`).

### Ejecución

```bash

go run ../../cmd/angLOXg motor_financiero.lox

```
