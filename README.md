# JackCompiler

Compilador da linguagem **Jack** (projeto [nand2tetris](https://www.nand2tetris.org/)), desenvolvido em **Go** para a disciplina de Compiladores. Esta é a **Entrega Parcial 1: o analisador léxico (scanner)**.

## 👥 Equipe
* **Daniel Nunes Duarte** - Matrícula: [preencher matrícula]
* **João Victor Oliveira** - Matrícula: [preencher matrícula]

## 💻 Linguagem de Programação
* **Golang (Go)** 

## ✅ Evidência de validação
A validação oficial do scanner foi feita com os programas de referência do projeto nand2tetris. O resultado obtido foi:

* **1.911 tokens idênticos em 7 programas oficiais**

Também foi verificado com o comando:

```bash
go test ./...
```

saindo com sucesso para o pacote do scanner.

## 🚀 Como Compilar e Executar

1. Certifique-se de ter o Go instalado (versão 1.18 ou superior).
2. Clone este repositório:
   ```bash
   git clone https://github.com/danielduartt/jackcompiler
   ```
3. Acesse o diretório do módulo:
   ```bash
   cd jackcompiler
   ```
4. Execute o compilador com um arquivo `.jack` ou um diretório:
   ```bash
   go run .\main.go ..\testdata\Square
   ```
   ou
   ```bash
   go run .\main.go ..\testdata\Square\Main.jack
   ```

## 📌 Observações
* O diretório de saída padrão é `out/`.
* Os arquivos gerados seguem o formato XML do projeto nand2tetris e são armazenados no diretório de saída indicado pela CLI.
