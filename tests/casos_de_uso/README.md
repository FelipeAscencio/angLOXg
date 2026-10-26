# Casos de uso reales en angLOXg

Esta suite contiene una serie de programas demostrativos desarrollados íntegramente en **Lox**.

Su propósito es validar la robustez del motor, su capacidad de procesamiento y, fundamentalmente, poner a prueba las extensiones y características incorporadas durante el desarrollo del TP en escenarios de ejecución reales.

## 1. Duelo Lógico: Adiviná el Número con "IA" (`duelo_logico.lox`)

**Prueba Integradora de innovaciones del TP**

Este programa actúa como la demostración principal de las extensiones arquitectónicas implementadas sobre el lenguaje original.

Simula un juego de deducción donde un agente "IA" evalúa un historial precalculado de jugadas para encontrar un número secreto, integrando simultáneamente todas las innovaciones requeridas.

- **Qué demuestra:**
  * **Soporte de arrays:** Almacenamiento dinámico de estado y acceso iterativo por índice para procesar el historial de intentos.
  * **Operador ternario:** Evaluación condicional de una sola línea para definir estados de validación y formatear mensajes de salida compactos.
  * **POO (Clases y herencia):** Encapsulamiento del estado en instancias, uso de constructores `init` con la referencia `this`, y enrutamiento dinámico de la lógica de validación hacia la clase base mediante `super`.

**Ejecución:**

```bash

go run ../../cmd/angLOXg duelo_logico.lox

```

## 2. Motor analítico para procesamiento de series de datos (`sensor_analiticas.lox`)

Simula un pipeline de procesamiento de métricas recolectadas por sensores industriales o dispositivos IoT, demostrando la capacidad del lenguaje para manejar flujos de datos continuos.

- **Qué demuestra:** 
  * Manipulación intensiva de variables numéricas y expresiones lógicas.
  * Modularidad mediante funciones dedicadas a operaciones de agregación matemática (mínimo, máximo y promedio aritmético).
  * Lógica de sanitización y filtrado condicional para detectar anomalías o valores fuera de rangos seguros.

**Ejecución:**

```bash

go run ../../cmd/angLOXg sensor_analiticas.lox

```

## 3. Sistema de gestión de descuentos y transacciones (`motor_ecommerce.lox`)

Calcula el costo final de órdenes de compra aplicando reglas de negocio comerciales complejas, poniendo a prueba el manejo de memoria y los entornos del intérprete.

- **Qué demuestra:**
  * **Cierres Léxicos (Closures):** Uso de funciones generadoras (`crearMotorDeDescuento`) que capturan, retienen y mutan estados de su entorno exterior, como el factor de lealtad de un cliente.
  * Lógica condicional de múltiples niveles para la evaluación de reglas anidadas, combinando descuentos por volumen e impuestos regionales.

**Ejecución:**

```bash

go run ../../cmd/angLOXg motor_ecommerce.lox

```

## 4. Motor de validación de riesgo crediticio (`motor_financiero.lox`)

Evalúa la elegibilidad de solicitudes financieras bajo estrictas normativas de riesgo. Es una prueba de estrés para el anidamiento de scopes y la pila de llamadas.

- **Qué demuestra:**
  * Ámbitos léxicos de alta complejidad y validación cruzada de variables en distintos niveles de profundidad.
  * **Funciones recursivas:** Ejecución de algoritmos iterativos avanzados (`proyectarRiesgo`) para la proyección matemática del estado financiero.
  * Estructuras de control de flujo robustas con múltiples ramificaciones funcionales.

**Ejecución:**

```bash

go run ../../cmd/angLOXg motor_financiero.lox

```
