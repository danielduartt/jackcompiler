package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesXML(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		name := "nome padrão"
		if explicit {
			name = "nome explícito"
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "Main.jack")
			output := filepath.Join(dir, "MainP.xml")
			if err := os.WriteFile(input, []byte("class Main {}"), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{input}
			if explicit {
				output = filepath.Join(dir, "result.xml")
				args = append(args, output)
			}
			var stdout bytes.Buffer
			if err := run(args, &stdout); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(output)
			if err != nil || !strings.Contains(string(got), "<class>") {
				t.Fatalf("XML não gerado: %q, %v", got, err)
			}
			if !strings.Contains(stdout.String(), output) {
				t.Fatalf("saída não informada: %s", stdout.String())
			}
		})
	}
}

func TestRunPreservesOutputOnSyntaxError(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "Main.jack")
	output := filepath.Join(dir, "MainP.xml")
	if err := os.WriteFile(input, []byte("class Main {"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, []byte("saída anterior"), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := run([]string{input}, &stdout); err == nil || !strings.Contains(err.Error(), "linha") {
		t.Fatalf("esperado diagnóstico, obteve %v", err)
	}
	got, err := os.ReadFile(output)
	if err != nil || string(got) != "saída anterior" || stdout.Len() != 0 {
		t.Fatalf("saída alterada após erro: %q, %v", got, err)
	}
	if err := os.Remove(output); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{input}, &stdout); err == nil {
		t.Fatal("entrada inválida foi aceita")
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("arquivo parcial criado: %v", err)
	}
}

func TestRunArgumentAndFileErrors(t *testing.T) {
	for _, args := range [][]string{nil, {"a", "b", "c"}, {"input.txt"}, {"input.jack", "input.jack"}, {"ausente.jack"}} {
		var stdout bytes.Buffer
		if err := run(args, &stdout); err == nil {
			t.Errorf("argumentos inválidos aceitos: %v", args)
		}
	}
}
