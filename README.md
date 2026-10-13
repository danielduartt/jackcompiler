# JackCompiler

Compilador da linguagem **Jack** (projeto [nand2tetris](https://www.nand2tetris.org/)), desenvolvido em **Go** para a disciplina de Compiladores. Esta é a **Entrega Parcial 1: o analisador léxico (scanner)**.

## Equipe

- **Daniel Nunes Duarte** — Matrícula: 20250013550
- **João Victor Oliveira** — Matrícula: 20260001375

## Requisitos

Go 1.18 ou superior, conforme a versão mínima declarada em `jackcompiler/go.mod`. O projeto usa apenas a biblioteca padrão. A compatibilidade deve ser confirmada executando os testes na versão de Go utilizada pela dupla.

## Compilar e executar

A partir da pasta em que deseja clonar o repositório:

```powershell
git clone https://github.com/danielduartt/jackcompiler
cd jackcompiler/jackcompiler
```

Se já estiver na raiz do repositório, entre apenas em `jackcompiler`. Os comandos abaixo devem ser executados dentro dessa pasta, onde está o `go.mod`.

Para processar um diretório ou apenas um arquivo:

```powershell
go run . ./testdata/Square
go run . ./testdata/Square/Main.jack
```

Para compilar e executar o binário no Windows:

```powershell
go build -o jackcompiler.exe .
./jackcompiler.exe ./testdata/Square
```

O diretório de saída padrão é `out/`. Para escolher outro diretório, passe o segundo argumento:

```powershell
go run . ./testdata/Square ./out/Square
```

A CLI gera `<Nome>T.xml` para cada fonte. Use uma pasta de saída separada de `testdata` para preservar os gabaritos: arquivos existentes no destino são sobrescritos. A leitura de diretórios não é recursiva.

## Estratégia de testes

A validação usa XML no padrão nand2tetris e testes unitários. Execute no diretório do módulo:

```powershell
go test ./... -count=1 -cover
go test ./scanner -run TestXMLOficial -v -count=1
go vet ./...
```

`TestXMLOficial` lê cada fonte `.jack` em `testdata`, executa o scanner e compara o XML produzido com o respectivo gabarito `*T.xml`. A comparação normaliza CRLF para LF e remove espaços e quebras nas extremidades do documento; o conteúdo das linhas de tokens é comparado exatamente. Em caso de divergência, informa a primeira linha diferente.

Os fontes e gabaritos de `Square/Main`, `Square/Square` e `Square/SquareGame` são obrigatórios: a ausência de qualquer um faz o teste falhar. Todo outro fonte `.jack` encontrado também precisa de seu gabarito. Assim, a falta dos arquivos de teste não resulta em uma validação silenciosamente ignorada.

Há sete pares de referência versionados: três em `Square`, três em `ExpressionLessSquare` e um em `ArrayTest`. Seus gabaritos contêm **1.911 tokens** no total. A execução bem-sucedida dos testes confirma a correspondência da saída com esses gabaritos.

Os testes unitários complementam os programas de referência com verificações das 21 palavras reservadas, dos 19 símbolos, das categorias e dos nomes dos tokens; identificadores; strings e comentários; espaços e contagem de linhas; limite de inteiros; caracteres inválidos; strings e comentários não terminados; e EOF. O XML inclui espaços internos nas tags, escapa `&`, `<`, `>` e aspas duplas, fica envolto por `<tokens>` e `</tokens>` e omite EOF. Há um teste específico para a omissão de EOF na escrita direta de XML e casos de rejeição de quebras LF, CR e CRLF em strings.

## Organização

- `jackcompiler/main.go`: CLI e leitura/escrita de arquivos.
- `jackcompiler/scanner/`: scanner, tipos de tokens, serialização XML e testes.
- `jackcompiler/testdata/`: fontes e gabaritos de referência versionados.
- `jackcompiler/out/`: saída gerada, ignorada pelo Git.

O scanner expõe `NextToken()` para consumo incremental e `Tokens()` para obter todos os tokens, permitindo reutilização pelo futuro analisador sintático.
