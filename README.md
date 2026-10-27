# JackCompiler

Implementação do analisador léxico (scanner) e do analisador sintático (parser)
da linguagem **Jack**, do projeto
[nand2tetris](https://www.nand2tetris.org/), desenvolvida em **Go** para a
disciplina de Compiladores.

O parser utiliza **recursive descent parsing**, reutiliza o scanner e gera a
árvore sintática em XML no formato exigido pelo nand2tetris.

## Objetivo da unidade

Este repositório reúne duas etapas do compilador Jack:

1. **Scanner:** transforma o código-fonte em tokens;
2. **Parser:** consome os tokens, verifica a gramática Jack e produz a árvore
   sintática XML.

O scanner e o parser são executados no mesmo fluxo pela aplicação principal.
Esta versão ainda não implementa a Compilation Engine, a geração de código VM
ou a execução na plataforma Hack; essas etapas pertencem às unidades
seguintes.

## Equipe

- **Daniel Nunes Duarte** — Matrícula: 20250013550
- **João Victor Oliveira** — Matrícula: 20260001375

## Requisitos

| Requisito | Informação |
|---|---|
| Linguagem | Go |
| Versão mínima | Go 1.18 |
| Dependências externas | Nenhuma |
| Sistema operacional | Linux, macOS ou Windows |
| Módulo Go | `github.com/danielduartt/jackcompiler` |

O projeto usa apenas a biblioteca padrão. Confira a instalação com:

```bash
go version
```

## Estrutura do projeto

Os comandos devem ser executados na pasta `jackcompiler/`, onde está o
`go.mod`:

```text
jackcompiler/
├── go.mod
├── main.go
├── parser/
│   ├── parser.go
│   ├── parser_test.go
│   └── testdata/
├── scanner/
│   ├── scanner.go
│   ├── token.go
│   ├── xml.go
│   ├── scanner_test.go
│   ├── token_test.go
│   └── oficial_test.go
├── testdata/
│   ├── ArrayTest/
│   ├── ExpressionLessSquare/
│   └── Square/
└── out/
```

- `main.go`: CLI e leitura/escrita de arquivos.
- `parser/`: analisador sintático recursive descent e testes.
- `scanner/`: analisador léxico, tokens e serialização XML.
- `testdata/`: fontes e gabaritos de referência.
- `parser/testdata/`: casos adicionais específicos do parser.
- `out/`: saída gerada localmente, ignorada pelo Git.

## Compilar e executar

A partir da pasta em que deseja clonar o repositório:

```bash
git clone https://github.com/danielduartt/jackcompiler
cd jackcompiler/jackcompiler
```

Se já estiver na raiz do repositório, entre apenas em `jackcompiler`. Os comandos abaixo devem ser executados dentro dessa pasta, onde está o `go.mod`.

Para processar um diretório ou apenas um arquivo:

```bash
go run . ./testdata/Square
go run . ./testdata/Square/Main.jack
```

O argumento de entrada pode ser um arquivo `.jack` ou um diretório contendo
arquivos `.jack` diretamente. A leitura de diretórios não é recursiva. Quando
a entrada não existe, é um diretório vazio ou não contém arquivos `.jack`, a
CLI retorna erro.

Para compilar e executar o binário no Windows:

```powershell
go build -o jackcompiler.exe .
./jackcompiler.exe ./testdata/Square
```

O diretório de saída padrão é `out/`. Para escolher outro diretório, passe o segundo argumento:

```bash
go run . ./testdata/Square ./out/Square
```

A CLI gera dois arquivos para cada fonte:

- `<Nome>T.xml`: tokens produzidos pelo scanner;
- `<Nome>P.xml`: árvore sintática produzida pelo parser.

Por exemplo, `Main.jack` gera `MainT.xml` e `MainP.xml`. Use uma pasta de
saída separada de `testdata` para preservar os gabaritos: arquivos existentes
no destino são sobrescritos. A leitura de diretórios não é recursiva.

Para processar um caso de teste específico do parser:

```bash
go run . ./parser/testdata/MainTest ./out
```

O parser aceita classes Jack completas, incluindo declarações de classe e
subrotinas, parâmetros, variáveis locais, comandos `let`, `if`, `while`, `do`
e `return`, expressões, operadores, constantes, chamadas de subrotinas e
acesso a arrays.

### Formato da linha de comando

```text
jackcompiler <arquivo.jack | diretorio> [pasta-de-saida]
```

Exemplos:

```bash
# Diretório inteiro, saída padrão em ./out
go run . ./testdata/Square

# Arquivo individual
go run . ./testdata/Square/Main.jack

# Diretório com saída em uma pasta específica
go run . ./testdata/Square ./out/Square

# Caso adicional do parser
go run . ./parser/testdata/MainTest ./out
```

Em caso de erro léxico ou sintático, a CLI escreve a mensagem no stderr e
termina com código de saída diferente de zero. Arquivos de saída existentes
com o mesmo nome são sobrescritos.

## Gramática implementada

O parser segue a gramática Jack por descida recursiva e implementa:

```text
class           = 'class' className '{' classVarDec* subroutineDec* '}'
classVarDec    = ('static' | 'field') type varName (',' varName)* ';'
subroutineDec  = ('constructor' | 'function' | 'method')
                  ('void' | type) subroutineName
                  '(' parameterList ')' subroutineBody
parameterList  = ((type varName) (',' type varName)*)?
subroutineBody = '{' varDec* statements '}'
varDec         = 'var' type varName (',' varName)* ';'
statements     = statement*
statement      = letStatement | ifStatement | whileStatement
                | doStatement | returnStatement
letStatement   = 'let' varName ('[' expression ']')? '=' expression ';'
ifStatement    = 'if' '(' expression ')' '{' statements '}'
                  ('else' '{' statements '}')?
whileStatement = 'while' '(' expression ')' '{' statements '}'
doStatement    = 'do' subroutineCall ';'
returnStatement = 'return' expression? ';'
expression     = term (op term)*
term           = integerConstant | stringConstant | keywordConstant
                | varName | varName '[' expression ']'
                | subroutineCall | '(' expression ')' | unaryOp term
subroutineCall = subroutineName '(' expressionList ')'
                | (className | varName) '.' subroutineName
                  '(' expressionList ')'
expressionList = (expression (',' expression)*)?
```

São aceitos:

- tipos `int`, `char`, `boolean` e tipos definidos pelo usuário;
- constantes `true`, `false`, `null` e `this`;
- operadores binários `+`, `-`, `*`, `/`, `&`, `|`, `<`, `>` e `=`;
- operadores unários `-` e `~`;
- chamadas qualificadas, como `Output.printInt(x)`;
- acesso indexado, como `array[i]`.

As expressões são representadas como sequências de termos e operadores,
conforme a gramática Jack e o formato do nand2tetris.

## Estratégia de testes

A validação usa XML no padrão nand2tetris, testes unitários do scanner e do
parser, além de verificação estática. Execute no diretório do módulo:

```bash
go test ./... -count=1 -cover
go test ./parser -v -count=1
go test ./scanner -run TestXMLOficial -v -count=1
go vet ./...
```

Para executar todos os testes com saída detalhada:

```bash
go test ./... -v -count=1
```

Para gerar um relatório de cobertura:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

`TestXMLOficial` lê cada fonte `.jack` em `testdata`, executa o scanner e compara o XML produzido com o respectivo gabarito `*T.xml`. A comparação normaliza CRLF para LF e remove espaços e quebras nas extremidades do documento; o conteúdo das linhas de tokens é comparado exatamente. Em caso de divergência, informa a primeira linha diferente.

Os fontes e gabaritos de `Square/Main`, `Square/Square` e `Square/SquareGame` são obrigatórios: a ausência de qualquer um faz o teste falhar. Todo outro fonte `.jack` encontrado também precisa de seu gabarito. Assim, a falta dos arquivos de teste não resulta em uma validação silenciosamente ignorada.

Há sete pares de referência versionados: três em `Square`, três em `ExpressionLessSquare` e um em `ArrayTest`. Seus gabaritos contêm **1.911 tokens** no total. A execução bem-sucedida dos testes confirma a correspondência da saída com esses gabaritos.

Os testes unitários complementam os programas de referência com verificações das 21 palavras reservadas, dos 19 símbolos, das categorias e dos nomes dos tokens; identificadores; strings e comentários; espaços e contagem de linhas; limite de inteiros; caracteres inválidos; strings e comentários não terminados; e EOF. O XML inclui espaços internos nas tags, escapa `&`, `<`, `>` e aspas duplas, fica envolto por `<tokens>` e `</tokens>` e omite EOF. Há um teste específico para a omissão de EOF na escrita direta de XML e casos de rejeição de quebras LF, CR e CRLF em strings.

Os testes do parser verificam classes mínimas, declarações, subrotinas,
parâmetros, variáveis locais, comandos, expressões, chamadas de subrotinas,
acesso a arrays e erros sintáticos com identificação da linha. O parser foi
validado contra os programas oficiais `Main`, `Square` e `SquareGame`; os
arquivos gerados foram comparados com os respectivos gabaritos XML e ficaram
idênticos.

Para comparar manualmente um arquivo do parser com o gabarito:

```bash
go run . ./testdata/Square ./out
diff -u -Z -B ./out/MainP.xml ./testdata/Square/Main.xml
diff -u -Z -B ./out/SquareP.xml ./testdata/Square/Square.xml
diff -u -Z -B ./out/SquareGameP.xml ./testdata/Square/SquareGame.xml
```

Ausência de saída no `diff` significa que os arquivos são idênticos. Os
gabaritos do parser neste repositório usam os nomes `Main.xml`, `Square.xml` e
`SquareGame.xml`; a CLI mantém o sufixo `P` nos arquivos gerados para
diferenciá-los dos arquivos de tokens.

### Casos cobertos pelos testes

Os testes do scanner cobrem palavras reservadas, símbolos, identificadores,
strings, comentários, espaços, contagem de linhas, limites de inteiros,
caracteres inválidos, entradas não terminadas, EOF, serialização XML e escape
de caracteres especiais.

Os testes do parser cobrem classes mínimas, declarações de classe, subrotinas,
parâmetros, variáveis locais, os cinco comandos (`let`, `if`, `while`, `do` e
`return`), expressões, operadores binários e unários, constantes, chamadas de
subrotinas, listas de argumentos, acesso a arrays e erros sintáticos com
identificação da linha.

### Validação contra os arquivos oficiais

A validação reproduzível dos três programas oficiais do diretório `Square` é:

```bash
rm -rf /tmp/jackcompiler-out
go run . ./testdata/Square /tmp/jackcompiler-out

for f in Main Square SquareGame; do
  if diff -u -Z -B \
      "/tmp/jackcompiler-out/${f}P.xml" \
      "./testdata/Square/${f}.xml" > "/tmp/${f}.diff"; then
    echo "$f: IDENTICO"
  else
    echo "$f: DIFERENTE"
    cat "/tmp/${f}.diff"
  fi
done
```

Sem saída no `diff` significa que os arquivos são idênticos após ignorar
diferenças de espaços e linhas em branco. Os gabaritos do parser usam os
nomes `Main.xml`, `Square.xml` e `SquareGame.xml`; a CLI usa `MainP.xml`,
`SquareP.xml` e `SquareGameP.xml` para distingui-los dos arquivos do scanner.

### Diagnóstico de erros

O parser interrompe a análise no primeiro erro e retorna um `SyntaxError`
contendo a linha do token inesperado. São rejeitados, entre outros:

- palavra-chave `class` ausente;
- delimitadores ou ponto e vírgula ausentes;
- expressão sem termo;
- chamada de subrotina incompleta;
- token sobrando depois do fechamento da classe;
- erros léxicos encontrados durante o consumo do parser.

## Organização

- `jackcompiler/main.go`: CLI e leitura/escrita de arquivos.
- `jackcompiler/parser/`: parser recursive descent, erros sintáticos e testes.
- `jackcompiler/scanner/`: scanner, tipos de tokens, serialização XML e testes.
- `jackcompiler/testdata/`: fontes e gabaritos de referência versionados.
- `jackcompiler/parser/testdata/`: casos adicionais específicos do parser.
- `jackcompiler/out/`: saída gerada, ignorada pelo Git.

O scanner expõe `NextToken()` para consumo incremental e `Tokens()` para obter
todos os tokens. O parser recebe o código-fonte, reutiliza o scanner e expõe
`Parse()` para devolver a árvore sintática XML.

## API principal

Scanner:

```go
scanner.New(src).NextToken()
scanner.New(src).Tokens()
```

Parser:

```go
parser.New(src).Parse()
```

`Parse()` retorna `(string, error)`. Em caso de sucesso, a string contém a
árvore XML completa. Em caso de falha, o chamador deve tratar o erro e não
usar a saída parcial.

## Formato da saída XML

Para uma classe simples:

```xml
<class>
  <keyword> class </keyword>
  <identifier> Main </identifier>
  <symbol> { </symbol>
  <symbol> } </symbol>
</class>
```

Cada token folha é emitido com a categoria do formato nand2tetris:

```xml
<keyword> class </keyword>
<identifier> Main </identifier>
<integerConstant> 10 </integerConstant>
<stringConstant> Hello </stringConstant>
<symbol> &lt; </symbol>
```

Os não-terminais são representados por tags como `<class>`,
`<subroutineDec>`, `<statements>`, `<expression>` e `<term>`. O token EOF não
é incluído no XML.

| Código-fonte | XML |
|---|---|
| `&` | `&amp;` |
| `<` | `&lt;` |
| `>` | `&gt;` |
| `"` | `&quot;` |

## Limitações e decisões de implementação

- O parser processa uma classe Jack completa por arquivo.
- A CLI processa somente arquivos `.jack` diretamente dentro do diretório
  informado.
- A CLI gera os XML do scanner e do parser no mesmo processamento.
- O parser informa linha, mas não coluna, nos erros sintáticos atuais.
- O parser para no primeiro erro e não realiza recuperação para continuar
  analisando o restante do arquivo.
- A validação oficial automatizada existente cobre diretamente o scanner; a
  comparação dos XML oficiais do parser é reproduzível pelo comando `diff`
  documentado acima.
- O diretório `out/` é uma pasta de saída local. Seus arquivos não devem ser
  usados como gabaritos versionados.

## Status da implementação

- Scanner: implementado e validado com os gabaritos oficiais.
- Parser: implementado com recursive descent e integrado à CLI.
- XML do parser: validado contra `Main`, `Square` e `SquareGame`.
- Testes: scanner com 95,0% de cobertura e parser com 90,6% de cobertura
  durante a validação atual.
- Verificação estática: `go vet ./...` sem erros.

## Etapas do compilador Jack

O projeto completo de um compilador Jack pode ser desenvolvido nas seguintes
etapas:

### Etapa 1 — Análise léxica

**Status: concluída.**

O scanner lê o código-fonte e produz tokens classificados como:

- palavras reservadas;
- símbolos;
- identificadores;
- constantes inteiras;
- constantes de string;
- EOF.

Também são tratados comentários, espaços, contagem de linhas, caracteres
inválidos, limites de inteiros e erros em strings ou comentários não
terminados. A saída de tokens pode ser escrita em `<Nome>T.xml`.

### Etapa 2 — Análise sintática

**Status: concluída.**

O parser consome os tokens do scanner por recursive descent, valida a
estrutura da classe Jack e gera a árvore sintática XML em `<Nome>P.xml`.
Essa etapa inclui declarações, subrotinas, comandos, expressões, chamadas de
subrotinas e acesso a arrays.

### Etapa 3 — Geração de código VM

**Status: pendente.**

O compilador deverá traduzir a linguagem Jack para a linguagem intermediária
VM do nand2tetris, incluindo:

- constantes e operações aritméticas;
- operações booleanas;
- acesso e atribuição de variáveis;
- acesso a arrays;
- comandos condicionais `if/else`;
- laços `while`;
- chamadas e retornos de subrotinas;
- construção de objetos;
- métodos, funções e construtores;
- manipulação de `this` e `that`.

### Próximos passos recomendados

1. criar testes automatizados para a comparação dos XML do parser;
2. adicionar posição de coluna aos tokens e às mensagens de erro;
3. implementar tabelas de símbolos;

## Checklist de reprodução da entrega

Execute, dentro de `jackcompiler/`:

```bash
gofmt -w .
go test ./... -count=1 -cover
go test ./scanner -run TestXMLOficial -v -count=1
go vet ./...
go build ./...
go run . ./testdata/Square /tmp/jackcompiler-out
```
## Referências

- [Nand2Tetris](https://www.nand2tetris.org/)
- [Software oficial do Nand2Tetris](https://www.nand2tetris.org/software)
- [Gramática Jack — Apêndice A](https://lambdageo-edu.github.io/book-compiladores-com-nand2tetris/apendices/a-gramatica-jack/)
- [Repositório do projeto](https://github.com/danielduartt/jackcompiler)
