<div align="center">
  <img src="data/images/banner.png" alt="angLOXg Banner">
</div>

Un intérprete del lenguaje Lox escrito puramente en Go, desarrollado como "Trabajo práctico" para la materia "Lenguajes y compiladores I (TB027)" de la Universidad de Buenos Aires (FIUBA).

Esta implementación sigue los conceptos y la arquitectura del libro *Crafting Interpreters* de Robert Nystrom.

**Enlaces rápidos**

* 📖 [Documentación oficial](https://FelipeAscencio.github.io/angLOXg/)
* ▶️ [Playground web interactivo](https://FelipeAscencio.github.io/angLOXg/web/)

## Prerrequisitos e instalación

Para compilar, ejecutar y testear el proyecto localmente, asegurate de contar con:
* **Go 1.22 o superior** instalado en el sistema.
* **gotestsum** (herramienta recomendada para la ejecución formateada de tests y métricas de cobertura).

Cloná el repositorio, navegá a la raíz del proyecto e instalá las dependencias de testeo global:

```bash

go install gotest.tools/gotestsum@latest

```

## Modos de Ejecución

*Nota: Las flags de depuración (`--escaneo`, `--arbol`, `--resolucion`) se pueden utilizar tanto pasando un archivo de script como en el formato de ejecución directo de consola (REPL). Al no pasar una flag, se inicia automáticamente en modo interprete completo.*

### 1. Consola interactiva (REPL)

Si ejecutás el programa sin argumentos, se abrirá la consola interactiva donde podés escribir sentencias de Lox línea por línea:

```bash

go run ./cmd/angLOXg [FLAG_DESEADA]

```

*(Para salir, usá `Ctrl+D` o `Ctrl+Z`).*

### 2. Ejecución de scripts

Para ejecutar un archivo de código Lox (`.lox`), pasá la ruta como argumento:

```bash

go run ./cmd/angLOXg [ruta_al_script.lox]

```

### 3. Modo de escaneo (`--escaneo`)

Detiene la ejecución tras la fase léxica y muestra los tokens generados en lugar de evaluar el código:

```bash

go run ./cmd/angLOXg --escaneo [ruta_al_script.lox]

```

### 4. Modo de árbol de sintaxis (`--arbol`)

Detiene la ejecución tras el parser y muestra el árbol sintáctico en formato de S-expresiones:

```bash

go run ./cmd/angLOXg --arbol [ruta_al_script.lox]

```

### 5. Modo de resolución (`--resolucion`)

Detiene la ejecución tras el análisis semántico (etapa del "Resolver") para verificar que las variables estén correctamente ligadas a sus scopes estáticos:

```bash

go run ./cmd/angLOXg --resolucion [ruta_al_script.lox]

```

### Compilación del binario

Si preferís generar el ejecutable para no depender de `go run`:

```bash

go build -o angloxg ./cmd/angLOXg

./angloxg [--escaneo | --arbol | --resolucion] [ruta_al_script.lox]

```

### Códigos de salida

| Código | Significado |
|:---|:---|
| `0` | Todo salió bien. |
| `64` | Error en la invocación del programa (argumentos de más). |
| `65` | El código Lox tenía errores léxicos o de sintaxis. |
| `66` | No se pudo leer el script indicado. |

## Pruebas y cobertura

El proyecto mantiene un estándar de **cobertura mínima del 80%** para todos los módulos internos.

Para correr la suite completa de pruebas utilizando la interfaz limpia de `gotestsum`:

```bash

gotestsum -- -v ./internal/... ./tests/...

```

Para calcular la **cobertura real unificada** de todo el motor (excluyendo el punto de entrada principal):

```bash

gotestsum -- -cover ./internal/... ./tests/...

```

## Distribución de Archivos y Estructura del Repositorio

X.
