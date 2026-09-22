---
layout: default
title: Pruebas
nav_order: 4
---

# 🧪 Pruebas y validación

## Pruebas de la cátedra

Utilizamos la suite de scripts provista por la cátedra como estándar principal para validar el correcto funcionamiento general de nuestro intérprete. Superar esta batería de pruebas nos permitió asegurar que la implementación base cumple con todos los requisitos, la sintaxis y los casos límite esperados para la aprobación del trabajo práctico.

## Pruebas propias

Para complementar la validación oficial, desarrollamos nuestra propia batería de pruebas unitarias y de integración enfocadas en el comportamiento interno de cada paquete de Go (escáner, parser, semántica, intérprete, etc.). Mediante este enfoque exhaustivo, logramos garantizar mas de un **80% de cobertura de código (*code coverage*)** a lo largo de todo el proyecto.

Adicionalmente, diseñamos escenarios de validación extra basados en **casos de uso reales** (como los simuladores de e-commerce y riesgo crediticio) y scripts específicos orientados a extraer **métricas de rendimiento**. Estos programas resultaron fundamentales no solo para medir la performance, sino para comprobar empíricamente la estabilidad, la robustez y el comportamiento correcto del motor frente a cargas de trabajo más complejas y extensas.
