package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"jackcompiler/internal/parser"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("uso: jackcompiler entrada.jack [saida.xml]")
	}
	inputPath := args[0]
	if !strings.EqualFold(filepath.Ext(inputPath), ".jack") {
		return fmt.Errorf("o arquivo de entrada deve ter extensão .jack")
	}
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "P.xml"
	if len(args) == 2 {
		outputPath = args[1]
	}
	if !strings.EqualFold(filepath.Ext(outputPath), ".xml") {
		return fmt.Errorf("o arquivo de saída deve ter extensão .xml")
	}
	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("ler %s: %w", inputPath, err)
	}
	outputInfo, err := os.Stat(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("verificar %s: %w", outputPath, err)
	}
	if err == nil && os.SameFile(inputInfo, outputInfo) {
		return fmt.Errorf("entrada e saída devem ser arquivos diferentes")
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("ler %s: %w", inputPath, err)
	}
	output, err := parser.Parse(string(input))
	if err != nil {
		return fmt.Errorf("%s: %w", inputPath, err)
	}
	// Erros sintáticos não criam XML parcial nem sobrescrevem uma saída anterior.
	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("escrever %s: %w", outputPath, err)
	}
	_, err = fmt.Fprintf(stdout, "Análise sintática concluída: %s\n", outputPath)
	return err
}
