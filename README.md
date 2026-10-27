# JackCompiler

Analisadores léxico e sintático da linguagem Jack, em Go, para a disciplina de
Compiladores (nand2tetris, projeto 10). O parser reutiliza o scanner da primeira
atividade e gera a árvore sintática em XML por descida recursiva. Esta etapa
não gera código VM nem realiza análise semântica de nomes e tipos.

## Equipe

- **Daniel Nunes Duarte** — Matrícula: [Matrícula aqui]
- **João Victor Oliveira** — Matrícula: [Matrícula aqui]

As matrículas ainda precisam ser preenchidas pelos integrantes.

## Compilar e executar

Requisito: Go **1.27.1 ou superior**, conforme o `go.mod` existente.
O projeto usa somente a biblioteca padrão, sem dependências externas.

Execute os comandos na raiz do repositório:

```sh
go run ./cmd/jackcompiler examples/Main.jack
```

Esse comando gera **`examples/MainP.xml`**. Por padrão, a saída fica ao lado
da entrada e recebe o sufixo **`P.xml`**: `Square.jack` gera `SquareP.xml` e
`SquareGame.jack` gera `SquareGameP.xml`.

Para escolher o arquivo de saída:

```sh
go run ./cmd/jackcompiler examples/Main.jack output.xml
```

Para compilar um executável no Windows:

```powershell
go build -o jackcompiler.exe ./cmd/jackcompiler
.\jackcompiler.exe examples/Main.jack output.xml
```

O programa recebe um arquivo `.jack` e, opcionalmente, um caminho `.xml`.
Para analisar várias classes, execute uma vez por arquivo. O diretório de
destino deve existir. A saída XML existente no caminho escolhido é substituída
somente após uma análise bem-sucedida. Entrada inválida produz diagnóstico em
stderr, código de saída 1 e não cria XML parcial. Sem argumentos, o programa
mostra a sintaxe de uso e termina com código 1.

Exemplo de diagnóstico para um `;` ausente:

```text
Main.jack: linha 7, coluna 3: esperado ";"; encontrado KEYWORD "return"
```

As posições começam em 1. Colunas contam bytes; tabulação conta como um byte.
O parser informa o primeiro erro encontrado.

## Organização

```text
cmd/jackcompiler/       interface de linha de comando e testes de integração
internal/lexer/         scanner existente, tokens e testes léxicos
internal/parser/        parser recursivo, geração XML e testes sintáticos
  testdata/official/   programas e XML esperados do nand2tetris
examples/Main.jack     exemplo pequeno de execução
```

O scanner mantém `New(input)` e `NextToken() Token`. O método adicional
`TokenColumn()` informa a coluna do último token sem mudar a estrutura `Token`.
O parser usa `advance`, `peek`, `match` e `consume` e possui funções para as
regras de classes, declarações, parâmetros, corpos, comandos, expressões,
termos e chamadas.

Uso do parser em Go, dentro deste módulo:

```go
xml, err := parser.Parse(source)
// Alternativa: parser.New(lexer.New(source)).Parse()
```

## Gramática e formato XML

São reconhecidos campos `static`/`field`, construtores, funções, métodos,
parâmetros, variáveis locais, os cinco comandos Jack, arrays, chamadas
qualificadas ou locais, constantes e todos os operadores da linguagem.
Inteiros devem estar entre 0 e 32767. Strings seguem a validação do scanner:
sem escapes e sem atravessar linhas.

A regra de expressão Jack é `term (op term)*`: os operadores binários aparecem
no mesmo nível da árvore, em ordem da esquerda para a direita. Parênteses
criam expressões aninhadas e operadores unários se ligam ao termo seguinte.
Não se introduz precedência aritmética de outra linguagem, pois isso mudaria
a árvore oficial. Nesta etapa as expressões não são avaliadas.

O XML segue o [formato do projeto 10](https://www.nand2tetris.org/project10):
tags `class`, `classVarDec`, `subroutineDec`, `parameterList`, `subroutineBody`,
`varDec`, `statements`, `letStatement`, `ifStatement`, `whileStatement`,
`doStatement`, `returnStatement`, `expression`, `term` e `expressionList`.
Listas vazias também possuem tags de abertura e fechamento.

Os terminais recebem `keyword`, `symbol`, `identifier`, `integerConstant` e
`stringConstant`. Alternativas como `statement`, `op`, `unaryOp`, `type`,
`keywordConstant` e `subroutineCall` não recebem tags extras: o formato oficial
representa seus componentes diretamente. Isso preserva a hierarquia exigida
pelo comparador. Os caracteres `&`, `<`, `>` e `"` são escapados como
`&amp;`, `&lt;`, `&gt;` e `&quot;`. O conteúdo das strings mantém seus espaços.

## Testes e validação

```sh
go fmt ./...
go test ./...
go vet ./...
go build ./...
go run ./cmd/jackcompiler examples/Main.jack
```

Os testes usam `testing` e `encoding/xml` da biblioteca padrão. Incluem a
gramática completa, estrutura de expressões, conteúdo e escapes XML,
rejeição de entradas inválidas, localização de erros e execução com arquivos.

`TestOfficialXML` compara estrutura, ordem e terminais com sete referências
oficiais: três de `Square`, três de `ExpressionLessSquare` e uma de `ArrayTest`.
As comparações ignoram indentação e preservam espaços dentro de strings.
Não há download durante os testes. A origem e a licença dos arquivos estão
documentadas em `internal/parser/testdata/official/README.md`.

No pacote oficial, as referências chamam-se `Main.xml`, `Square.xml` e
`SquareGame.xml`; elas correspondem aos nomes `MainP.xml`, `SquareP.xml` e
`SquareGameP.xml` usados no enunciado. O JackCompiler gera o sufixo `P.xml`
para distinguir sua saída dos arquivos de referência.

Validação executada em **27/10/2026**, com **Go 1.27.3 windows/amd64**:

- `go fmt ./...`: concluído.
- `go test ./...`: aprovado nos três pacotes, incluindo as sete comparações oficiais.
- `go vet ./...`: aprovado.
- `go build ./...`: aprovado.
- `go run ./cmd/jackcompiler examples/Main.jack`: gerou `examples/MainP.xml`.

O Go não estava no PATH deste ambiente. Para essa validação foi utilizada uma
distribuição oficial portátil em pasta temporária, com SHA-256 conferido,
sem alterar o `go.mod` nem instalar Go no sistema. Para repetir os comandos em
outro terminal, é necessário disponibilizar Go no PATH.

## Desafios da implementação

Um desafio foi distinguir variável, acesso a array e chamada de subrotina,
que começam por um identificador. Um token de antecipação resolve a escolha
sem duplicar a análise léxica. Outro foi representar a gramática no XML exato
do nand2tetris, incluindo listas vazias e expressões sem tags adicionais para
operadores. Os testes oficiais verificam essas decisões de estrutura.

As mensagens de erro precisaram guardar a posição do token atual antes de
ler a antecipação. O scanner oferece a coluna separadamente para preservar
a API de tokens já utilizada na primeira atividade.

Os testes de localização também revelaram um avanço extra no fechamento de
comentários de bloco do scanner. A remoção desse avanço preserva o primeiro
caractere após `*/`, incluindo quebras de linha. Há testes de regressão para
comentário seguido diretamente de token e de nova linha.

## Entrega acadêmica

Preencher matrículas e conferir o link do repositório antes da entrega.
O enunciado exige pelo menos 15 commits coerentes; nenhum histórico foi
fabricado nem foram criados commits automaticamente. O envio do ZIP e do
link no SIGAA é realizado pelo grupo.
