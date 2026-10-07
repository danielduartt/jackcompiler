package main

import (
	"fmt"
	"jackcompiler/internal/lexer"
)

func main() {
	// Um trecho de código Jack real para testarmos nosso Lexer
	input := `
	class Main {
		function void main() {
			var int x;
			let x = 10;
			if (x < 20) {
				do Output.printString("Hello, World!");
			}
			return;
		}
	}
	`

	fmt.Println("--- INICIANDO ANALISE LEXICA ---")
	l := lexer.New(input)

	// Loop que consome tokens até chegar ao EOF (End Of File)
	for tok := l.NextToken(); tok.Type != lexer.EOF; tok = l.NextToken() {
		fmt.Printf("Tipo: %-15s | Literal: %s\n", tok.Type, tok.Literal)
	}
	fmt.Println("--- ANALISE CONCLUIDA COM SUCESSO ---")
}