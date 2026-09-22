package catedra_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FelipeAscencio/angLOXg/internal/lox"
)

func TestCatedra(t *testing.T) {
	archivos, err := filepath.Glob("*.lox")
	if err != nil {
		t.Fatalf("No se pudieron buscar los archivos .lox: %v", err)
	}

	if len(archivos) == 0 {
		t.Fatalf("No se encontraron archivos .lox para probar.")
	}

	for _, archivo := range archivos {
		t.Run(archivo, func(t *testing.T) {
			contenido, err := os.ReadFile(archivo)
			if err != nil {
				t.Fatalf("Error al leer el archivo %s: %v", archivo, err)
			}

			var salida bytes.Buffer
			ok := lox.Ejecutar(string(contenido), &salida)
			obtenido := salida.String()
			if strings.Contains(strings.ToUpper(obtenido), "ERROR") {
				t.Errorf("El script %s reportó un error.\nSalida:\n%s", archivo, obtenido)
			}

			if !ok {
				t.Errorf("La ejecución interna de %s falló.\nSalida:\n%s", archivo, obtenido)
			}
		})
	}
}
