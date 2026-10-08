package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"jackcompiler/scanner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: jackcompiler <arquivo.jack | diretorio> [pasta-de-saida]")
		os.Exit(1)
	}
	entrada := os.Args[1]
	saida := "out"
	if len(os.Args) > 2 {
		saida = os.Args[2]
	}

	arquivos, err := listarJack(entrada)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(saida, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	falhou := false
	for _, arq := range arquivos {
		if err := processar(arq, saida); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", arq, err)
			falhou = true
		}
	}
	if falhou {
		os.Exit(1)
	}
}

// listarJack devolve o próprio arquivo, ou todos os .jack de um diretório.
func listarJack(caminho string) ([]string, error) {
	info, err := os.Stat(caminho)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{caminho}, nil
	}
	arquivos, err := filepath.Glob(filepath.Join(caminho, "*.jack"))
	if err != nil {
		return nil, err
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
	tokens, err := scanner.New(string(src)).Tokens()
	if err != nil {
		return err
	}
	nome := strings.TrimSuffix(filepath.Base(arquivo), ".jack")
	destino := filepath.Join(pastaSaida, nome+"T.xml")

	f, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer f.Close() // executa quando a função terminar, mesmo se houver erro

	return scanner.WriteXML(f, tokens)
}
