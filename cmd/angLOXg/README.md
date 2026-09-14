# Uso Local de angLOXg

Esta sección detalla cómo compilar y probar el intérprete "Lox" de forma local.

El motor está construido 100% en "Go", sin dependencias externas.

## Prerrequisitos

* Go 1.22 o superior instalado en el sistema.

## Instalación

Cloná el repositorio y navegá hasta la raíz del proyecto:

```bash
git clone https://github.com/FelipeAscencio/angLOXg.git
cd angLOXg
```

## Modos de ejecución

### 1. Consola interactiva (REPL)

Si ejecutás el programa sin argumentos, se abrirá la consola interactiva donde podés escribir sentencias de Lox línea por línea:

```bash
go run ./cmd/angLOXg
```

Para salir, `Ctrl+D` (o `Ctrl+Z` seguido de `Enter` en Windows). Un error no corta la sesión: se reporta y la consola sigue esperando la próxima sentencia.

### 2. Ejecución de scripts

Para ejecutar un archivo con código Lox, pasá la ruta del archivo como argumento:

```bash
go run ./cmd/angLOXg script.lox
```

### 3. Modo de escaneo

Con `--escaneo` el programa se detiene después del escáner y muestra los tokens que produjo, en lugar de ejecutar el código. Sirve para depurar la primera fase del intérprete:

```bash
go run ./cmd/angLOXg --escaneo script.lox
```

Por ejemplo, `var x = 10;` produce:

```
VAR
IDENTIFIER<x>
EQUAL
NUMBER<10>
SEMICOLON
EOF
```

El flag también funciona en el REPL.

## Compilación

Si preferís generar el binario ejecutable para no depender del comando `go run`:

```bash
go build -o angloxg ./cmd/angLOXg
```

```bash
./angloxg [--escaneo] [ruta_al_script.lox]
```

## Códigos de salida

| Código | Significado |
|:---|:---|
| `0` | Todo salió bien. |
| `64` | Error en la invocación del programa (argumentos de más). |
| `65` | El código Lox tenía errores. |
| `66` | No se pudo leer el script indicado. |

## Tests

Los tests del motor se corren desde la raíz del repositorio:

```bash
go test ./...
```

Para ver el detalle de cada caso, o la cobertura alcanzada:

```bash
go test ./... -v
```

```bash
go test ./... -cover
```
