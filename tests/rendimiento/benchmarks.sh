#!/bin/bash

# Asegura que el script se ejecute parado en tests/rendimiento-tree_walk, sin importar desde dónde se llame.
cd "$(dirname "$0")"

MODO=$1
NOMBRE_MODO="Tree-Walk (AST)"
FLAG=""

if [ "$MODO" == "--vm" ]; then
    NOMBRE_MODO="Bytecode VM"
    FLAG="--vm"
fi

# Archivo de salida para las métricas.
OUTPUT="resultados_benchmarks.txt"

echo "========================================" > $OUTPUT
echo "  RESULTADOS DE RENDIMIENTO - angLOXg   " >> $OUTPUT
echo "  Modo: $NOMBRE_MODO                    " >> $OUTPUT
echo "  Fecha: $(date)                        " >> $OUTPUT
echo "========================================" >> $OUTPUT
echo "" >> $OUTPUT

# 1. Compilar el proyecto.
echo "Compilando el motor Go para las pruebas ($NOMBRE_MODO)..."
go build -o angloxg_bin ../../cmd/angLOXg

# Verificar si compiló bien
if [ ! -f ./angloxg_bin ]; then
    echo "Error: Falló la compilación. Verificá que las rutas sean correctas."
    exit 1
fi

echo "Motor compilado. Ejecutando benchmarks..."

# 2. Iterar sobre los archivos de prueba en el directorio actual y medir tiempos.
for script in *.lox; do
    echo "Corriendo: $script"
    
    echo "----------------------------------------" >> $OUTPUT
    echo "Script: $script" >> $OUTPUT
    echo "----------------------------------------" >> $OUTPUT
    
    # Redirigimos la salida normal (1) y la salida de error/tiempo (2) al archivo.
    if [ -n "$FLAG" ]; then
        { time ./angloxg_bin "$FLAG" "$script" ; } 2>> $OUTPUT 1>> $OUTPUT
    else
        { time ./angloxg_bin "$script" ; } 2>> $OUTPUT 1>> $OUTPUT
    fi
    
    echo "" >> $OUTPUT
done

# 3. Limpiar el binario temporal.
rm angloxg_bin

echo "¡Pruebas finalizadas! Podés revisar el archivo $OUTPUT en esta misma carpeta."
