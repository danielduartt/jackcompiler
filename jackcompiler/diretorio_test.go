package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Compilar um diretório gera um .vm por classe, na MESMA pasta, e ignora
// subpastas e arquivos que não são .jack (requisito "entrada por diretório").
func TestCompilaDiretorioGeraUmVMPorClasse(t *testing.T) {
	dir := t.TempDir()
	origem := filepath.Join("testdata", "11", "Pong")
	classes := []string{"Ball", "Bat", "Main", "PongGame"}
	for _, c := range classes {
		src, err := os.ReadFile(filepath.Join(origem, c+".jack"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, c+".jack"), src, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Uma subpasta (outro programa) e um arquivo qualquer não podem ser compilados.
	sub := filepath.Join(dir, "subpasta")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "Outro.jack"), []byte("class Outro {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "LEIAME.txt"), []byte("notas"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Mesmo fluxo da CLI: lista os .jack da pasta e compila cada um.
	arquivos, err := listarJack(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, arq := range arquivos {
		if err := compilarVM(arq); err != nil {
			t.Fatalf("%s: %v", arq, err)
		}
	}

	var gerados []string
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entradas {
		if filepath.Ext(e.Name()) == ".vm" {
			gerados = append(gerados, strings.TrimSuffix(e.Name(), ".vm"))
		}
	}
	sort.Strings(gerados)
	if strings.Join(gerados, ",") != strings.Join(classes, ",") {
		t.Errorf(".vm gerados = %v, esperado um por classe: %v", gerados, classes)
	}

	// Cada .vm contém as funções da sua própria classe.
	for _, c := range classes {
		vm, err := os.ReadFile(filepath.Join(dir, c+".vm"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(string(vm), "function "+c+".") {
			t.Errorf("%s.vm não começa com uma função de %s", c, c)
		}
	}

	if _, err := os.Stat(filepath.Join(sub, "Outro.vm")); err == nil {
		t.Error("a subpasta não deveria ter sido compilada")
	}
}
