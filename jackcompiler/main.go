package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielduartt/jackcompiler/parser"
	"github.com/danielduartt/jackcompiler/scanner"
)

func main() {
	args := os.Args[1:]
	xml := len(args) > 0 && args[0] == "--xml"
	if xml {
		args = args[1:]
	}
	if len(args) < 1 || len(args) > 2 {
		fmt.Fprintln(os.Stderr, "uso: jackcompiler [--xml] <arquivo.jack | diretorio> [pasta-de-saida-xml]")
		os.Exit(1)
	}
	entrada := args[0]
	saida := "out"
	if len(args) == 2 {
		xml = true // A chamada antiga com pasta de saída seleciona a análise XML.
		saida = args[1]
	}

	arquivos, err := listarJack(entrada)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	if xml {
		if err := os.MkdirAll(saida, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
	}

	falhou := false
	for _, arq := range arquivos {
		if xml {
			err = processar(arq, saida)
		} else {
			err = compilarVM(arq)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", arq, err)
			falhou = true
		}
	}
	if falhou {
		os.Exit(1)
	}
}

func compilarVM(arquivo string) error {
	src, err := os.ReadFile(arquivo)
	if err != nil {
		return err
	}
	vm, err := parser.CompileClass(string(src))
	if err != nil {
		return err
	}
	destino := strings.TrimSuffix(arquivo, ".jack") + ".vm"
	return os.WriteFile(destino, []byte(vm), 0o644)
}

// listarJack devolve o próprio arquivo, ou todos os .jack de um diretório.
func listarJack(caminho string) ([]string, error) {
	info, err := os.Stat(caminho)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if !info.Mode().IsRegular() || filepath.Ext(caminho) != ".jack" {
			return nil, fmt.Errorf("entrada deve ser um arquivo .jack ou diretório: %s", caminho)
		}
		return []string{caminho}, nil
	}
	entradas, err := os.ReadDir(caminho)
	if err != nil {
		return nil, err
	}
	var arquivos []string
	for _, entrada := range entradas {
		if entrada.Type().IsRegular() && filepath.Ext(entrada.Name()) == ".jack" {
			arquivos = append(arquivos, filepath.Join(caminho, entrada.Name()))
		}
	}
	if len(arquivos) == 0 {
		return nil, fmt.Errorf("nenhum .jack em %s", caminho)
	}
	return arquivos, nil
}

// processar lê um .jack e grava <pasta>/<Nome>T.xml (nunca sobrescreve o oficial).
func processar(arquivo, pastaSaida string) error {
	src, err := os.ReadFile(arquivo)
	if err != nil {
		return err
	}
	nome := strings.TrimSuffix(filepath.Base(arquivo), ".jack")

	// 1) scanner -> tokens
	tokens, err := scanner.New(string(src)).Tokens()
	if err != nil {
		return err
	}
	if err := gravarTokens(filepath.Join(pastaSaida, nome+"T.xml"), tokens); err != nil {
		return err
	}

	// 2) scanner + parser -> árvore sintática
	arvore, err := parser.New(string(src)).Parse()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(pastaSaida, nome+"P.xml"), []byte(arvore), 0o644)
}

func gravarTokens(destino string, tokens []scanner.Token) error {
	f, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer f.Close()

	return scanner.WriteXML(f, tokens)
}
