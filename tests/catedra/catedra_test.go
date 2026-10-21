package catedra_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

// Ejecuta de forma automatizada todos los archivos de prueba (.lox) de la cátedra
// probando tanto el modo Tree-Walk como el motor de Bytecode Compiler & VM.
func TestCatedra(t *testing.T) {
	archivos, err := filepath.Glob("*.lox")
	if err != nil {
		t.Fatalf("No se pudieron buscar los archivos .lox: %v", err)
	}

	if len(archivos) == 0 {
		t.Fatalf("No se encontraron archivos .lox para probar.")
	}

	for _, archivo := range archivos {
		// 1. Test con el intérprete original (Tree-Walk)
		t.Run(archivo+" [Tree-Walk]", func(t *testing.T) {
			contenido, err := os.ReadFile(archivo)
			if err != nil {
				t.Fatalf("Error al leer el archivo %s: %v", archivo, err)
			}

			fmt.Printf("\n[TEST CÁTEDRA] 🌳 Ejecutando '%s' con el motor TREE-WALK...\n", archivo)
			var salida bytes.Buffer
			ok := lox.Ejecutar(string(contenido), &salida)
			obtenido := salida.String()
			if strings.Contains(strings.ToUpper(obtenido), "ERROR") {
				t.Errorf("El script %s (Tree-Walk) reportó un error.\nSalida:\n%s", archivo, obtenido)
			}

			if !ok {
				t.Errorf("La ejecución interna de %s (Tree-Walk) falló.\nSalida:\n%s", archivo, obtenido)
			}
		})

		// 2. Test con la Máquina Virtual y Bytecode Compiler (VM)
		t.Run(archivo+" [Bytecode VM]", func(t *testing.T) {
			contenido, err := os.ReadFile(archivo)
			if err != nil {
				t.Fatalf("Error al leer el archivo %s: %v", archivo, err)
			}

			fmt.Printf("\n[TEST CÁTEDRA] ⚡ Ejecutando '%s' con el motor BYTECODE & VM...\n", archivo)
			var salida bytes.Buffer
			ok := lox.EjecutarVM(string(contenido), &salida)
			obtenido := salida.String()
			if strings.Contains(strings.ToUpper(obtenido), "ERROR") {
				t.Errorf("El script %s (VM) reportó un error.\nSalida:\n%s", archivo, obtenido)
			}

			if !ok {
				t.Errorf("La ejecución interna de %s (VM) falló.\nSalida:\n%s", archivo, obtenido)
			}
		})
	}
}
