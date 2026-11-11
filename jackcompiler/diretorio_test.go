package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Compilar o diretório de cada programa do Project 11 gera um .vm por classe,
// na MESMA pasta, e ignora subpastas e arquivos que não são .jack
// (requisito "entrada por diretório").
func TestCompilaDiretorioGeraUmVMPorClasse(t *testing.T) {
	for _, programa := range programasProject11 {
		t.Run(programa, func(t *testing.T) {
			classes := copiarPrograma(t, filepath.Join("testdata", "11", programa))
			dir := filepath.Dir(classes[0].caminho)

			// Uma subpasta (outro programa) e um arquivo qualquer não podem ser compilados.
			sub := filepath.Join(dir, "subpasta")
			if err := os.Mkdir(sub, 0o755); err != nil {
				t.Fatal(err)
			}
			escrever(t, filepath.Join(sub, "Outro.jack"), "class Outro {}")
			escrever(t, filepath.Join(dir, "LEIAME.txt"), "notas")

			// Mesmo fluxo da CLI: lista os .jack da pasta e compila cada um.
			arquivos, err := listarJack(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, arq := range arquivos {
				if err := compilarVM(arq); err != nil {
					t.Fatalf("%s: %v", filepath.Base(arq), err)
				}
			}

			var esperados []string
			for _, c := range classes {
				esperados = append(esperados, c.nome)
			}
			if gerados := vmsDaPasta(t, dir); strings.Join(gerados, ",") != strings.Join(esperados, ",") {
				t.Errorf(".vm gerados = %v, esperado um por classe: %v", gerados, esperados)
			}

			// Cada .vm contém as funções da sua própria classe.
			for _, c := range esperados {
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
		})
	}
}

type classeJack struct{ nome, caminho string }

// copiarPrograma copia os .jack de origem para uma pasta temporária (para não
// gerar .vm dentro de testdata) e devolve as classes em ordem alfabética.
func copiarPrograma(t *testing.T, origem string) []classeJack {
	t.Helper()
	dir := t.TempDir()
	fontes, err := filepath.Glob(filepath.Join(origem, "*.jack"))
	if err != nil || len(fontes) == 0 {
		t.Fatalf("nenhum .jack em %s (%v)", origem, err)
	}
	var classes []classeJack
	for _, f := range fontes {
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		destino := filepath.Join(dir, filepath.Base(f))
		escrever(t, destino, string(src))
		classes = append(classes, classeJack{strings.TrimSuffix(filepath.Base(f), ".jack"), destino})
	}
	sort.Slice(classes, func(i, j int) bool { return classes[i].nome < classes[j].nome })
	return classes
}

// vmsDaPasta devolve, em ordem alfabética, os nomes dos .vm da pasta (sem extensão).
func vmsDaPasta(t *testing.T, dir string) []string {
	t.Helper()
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var nomes []string
	for _, e := range entradas {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".vm" {
			nomes = append(nomes, strings.TrimSuffix(e.Name(), ".vm"))
		}
	}
	sort.Strings(nomes)
	return nomes
}

func escrever(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}
