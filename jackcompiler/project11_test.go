package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielduartt/jackcompiler/parser"
)

// Os seis programas do Project 11 devem compilar sem erro.
// A validação final é executar o .vm no VM Emulator (ver README).
var programasProject11 = []string{"Seven", "ConvertToBin", "Square", "Average", "Pong", "ComplexArrays"}

func TestCompilaProgramasProject11(t *testing.T) {
	for _, programa := range programasProject11 {
		t.Run(programa, func(t *testing.T) {
			arquivos, err := listarJack(filepath.Join("testdata", "11", programa))
			if err != nil {
				t.Fatal(err)
			}
			for _, arq := range arquivos {
				src, err := os.ReadFile(arq)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := parser.CompileClass(string(src)); err != nil {
					t.Errorf("%s: %v", filepath.Base(arq), err)
				}
			}
		})
	}
}
