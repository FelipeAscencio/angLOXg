#!/bin/bash

# Asegura que el script se ejecute parado en tests/rendimiento-tree_walk, sin importar desde dónde se llame.
cd "$(dirname "$0")"

OUTPUT="resultados_comparativa.txt"

echo "==========================================================================" > $OUTPUT
echo "        COMPARATIVA DE RENDIMIENTO: TREE-WALK VS BYTECODE VM              " >> $OUTPUT
echo "        Fecha: $(date)                                                    " >> $OUTPUT
echo "==========================================================================" >> $OUTPUT
echo "" >> $OUTPUT

# 1. Compilar el proyecto.
echo "Compilando el motor Go para la comparativa..."
go build -o angloxg_bin ../../cmd/angLOXg

if [ ! -f ./angloxg_bin ]; then
    echo "Error: Falló la compilación. Verificá que las rutas sean correctas."
    exit 1
fi

echo "Motor compilado. Ejecutando ambos motores sobre cada benchmark..."
echo ""

HEADER=$(printf "%-20s | %-12s | %-12s | %-10s | %-12s" "Script" "Tree-Walk" "Bytecode VM" "Speedup" "Optimización")
LINEA="----------------------------------------------------------------------------"

echo "$LINEA"
echo "$HEADER"
echo "$LINEA"

echo "$LINEA" >> $OUTPUT
echo "$HEADER" >> $OUTPUT
echo "$LINEA" >> $OUTPUT

# 2. Iterar sobre cada benchmark y medir ambos modos.
for script in *.lox; do
    # Medir tiempo en modo Tree-Walk (AST)
    start_ast=$(date +%s.%N)
    ./angloxg_bin "$script" > /dev/null 2>&1
    end_ast=$(date +%s.%N)
    time_ast=$(awk -v s="$start_ast" -v e="$end_ast" 'BEGIN {printf "%.4f", e - s}')

    # Medir tiempo en modo Bytecode VM (--vm)
    start_vm=$(date +%s.%N)
    ./angloxg_bin --vm "$script" > /dev/null 2>&1
    end_vm=$(date +%s.%N)
    time_vm=$(awk -v s="$start_vm" -v e="$end_vm" 'BEGIN {printf "%.4f", e - s}')

    # Calcular factor de mejora (Speedup) y porcentaje de reducción de tiempo
    speedup=$(awk -v ast="$time_ast" -v vm="$time_vm" 'BEGIN { if (vm > 0) printf "%.2fx", ast / vm; else print "N/A" }')
    ahorro=$(awk -v ast="$time_ast" -v vm="$time_vm" 'BEGIN { if (ast > 0) printf "-%.1f%%", ((ast - vm) / ast) * 100; else print "N/A" }')

    FILA=$(printf "%-20s | %-10s s | %-10s s | %-10s | %-12s" "$script" "$time_ast" "$time_vm" "$speedup" "$ahorro")
    
    echo "$FILA"
    echo "$FILA" >> $OUTPUT
done

echo "$LINEA"
echo "$LINEA" >> $OUTPUT

# 3. Limpiar el binario temporal.
rm angloxg_bin

echo ""
echo "¡Comparativa finalizada! Resultados guardados en $OUTPUT."
