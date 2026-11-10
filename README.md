# JackCompiler

Compilador da linguagem **Jack**, do projeto
[nand2tetris](https://www.nand2tetris.org/), desenvolvido em **Go** para a
disciplina de Compiladores.

O parser utiliza **recursive descent parsing**, reutiliza o scanner e gera a
árvore sintática em XML no formato exigido pelo nand2tetris. A geração de
código VM (Project 11) está **em andamento**; veja a
[Etapa 3](#etapa-3--geração-de-código-vm).

## Objetivo do projeto

Este repositório reúne as etapas do compilador Jack:

1. **Scanner:** transforma o código-fonte em tokens (concluído);
2. **Parser:** consome os tokens, verifica a gramática Jack e produz a árvore
   sintática XML (concluído);
3. **Tabela de símbolos + gerador de código:** traduz cada classe Jack para
   código da máquina virtual (`.vm`) do nand2tetris (em andamento).

O scanner e o parser são executados no mesmo fluxo pela aplicação principal.
A tabela de símbolos (`symtab`) e o escritor de código VM (`vmwriter`) já
estão implementados e testados, mas **ainda não estão integrados ao parser**:
por enquanto a CLI continua gerando apenas os XML.

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
│   ├── oficial_test.go
│   └── testdata/
├── scanner/
│   ├── scanner.go
│   ├── token.go
│   ├── xml.go
│   ├── scanner_test.go
│   ├── token_test.go
│   └── oficial_test.go
├── symtab/
│   ├── symtab.go
│   └── symtab_test.go
├── vmwriter/
│   ├── vmwriter.go
│   └── vmwriter_test.go
├── testdata/
│   ├── ArrayTest/
│   ├── ExpressionLessSquare/
│   └── Square/
└── out/
```

- `main.go`: CLI e leitura/escrita de arquivos.
- `parser/`: analisador sintático recursive descent e testes.
- `scanner/`: analisador léxico, tokens e serialização XML.
- `symtab/`: tabela de símbolos com escopos de classe e de subrotina.
- `vmwriter/`: escrita dos comandos da máquina virtual, um por linha.
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

A validação usa XML no padrão nand2tetris, testes unitários do scanner, do
parser, da tabela de símbolos e do `VMWriter`, além de verificação estática.
Execute no diretório do módulo:

```bash
go test ./... -count=1 -cover
go test ./parser -v -count=1
go test ./scanner -run TestXMLOficial -v -count=1
go test ./symtab ./vmwriter -v -count=1
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
acesso a arrays e erros sintáticos com identificação da linha.

Os testes do `VMWriter` (`vmwriter_test.go`) conferem isoladamente cada
comando gerado: `push`/`pop` em todos os segmentos (por exemplo,
`push constant 5`), os nove comandos aritméticos e lógicos (`add`, `neg`,
`not`...), `label`, `goto`, `if-goto`, `call Math.abs 1`,
`function Main.main 2` e `return`.

Os testes da tabela de símbolos (`symtab_test.go`) conferem os índices por
Kind (um `static` declarado depois de dois `field` recebe índice 0), o `this`
como `argument 0` em métodos, a limpeza do escopo da subrotina sem apagar o da
classe, o sombreamento de um `field` por uma variável local e nomes não
declarados. Os dois pacotes têm 100% de cobertura.

### Teste automatizado do parser contra os gabaritos oficiais

`TestParserXMLOficial` (em `parser/oficial_test.go`) lê cada `X.jack` de
`ExpressionLessSquare`, `Square` e `ArrayTest`, executa o parser e compara a
árvore gerada com o gabarito oficial `X.xml` da mesma pasta. Antes da
comparação, os dois textos são normalizados: os espaços das pontas de cada linha
e as linhas vazias são removidos, o que também descarta o `\r` do Windows.
Diferenças de indentação são ignoradas, mas a sequência de tags e tokens precisa
ser idêntica. A ausência de um fonte ou de um gabarito faz o teste falhar, e uma
divergência informa a primeira linha diferente (nosso × oficial).

```bash
go test ./parser -run TestParserXMLOficial -v -count=1
```

### Status da validação

| Programa | Arquivo | Scanner (`XT.xml`) | Parser (`X.xml`) |
|---|---|---|---|
| ExpressionLessSquare | `Main.jack` | ✅ idêntico | ✅ idêntico |
| ExpressionLessSquare | `Square.jack` | ✅ idêntico | ✅ idêntico |
| ExpressionLessSquare | `SquareGame.jack` | ✅ idêntico | ✅ idêntico |
| Square | `Main.jack` | ✅ idêntico | ✅ idêntico |
| Square | `Square.jack` | ✅ idêntico | ✅ idêntico |
| Square | `SquareGame.jack` | ✅ idêntico | ✅ idêntico |
| ArrayTest | `Main.jack` | ✅ idêntico | ✅ idêntico |

O scanner é verificado por `TestXMLOficial` e o parser por
`TestParserXMLOficial`; ambos rodam com `go test ./... -count=1`.

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
- `jackcompiler/symtab/`: tabela de símbolos e testes.
- `jackcompiler/vmwriter/`: escrita de código VM e testes.
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

Tabela de símbolos:

```go
st := symtab.New()                          // uma tabela por classe
st.Define("x", "int", symtab.Field)         // devolve Symbol{x int field 0}
st.StartSubroutine()                        // limpa arg/var, mantém static/field
sym, ok := st.Resolve("x")                  // subrotina primeiro, depois classe
n := st.VarCount(symtab.Var)                // nº de locais da subrotina
```

VMWriter:

```go
w := vmwriter.New()
w.WriteFunction("Main.main", 0)             // function Main.main 0
w.WritePush(vmwriter.Constant, 7)           // push constant 7
w.WriteCall("Output.printInt", 1)           // call Output.printInt 1
w.WriteReturn()                             // return
codigo := w.String()
```

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
- A comparação com os gabaritos oficiais é automatizada para o scanner
  (`TestXMLOficial`) e para o parser (`TestParserXMLOficial`); o `diff`
  documentado acima serve apenas para conferência manual.
- O diretório `out/` é uma pasta de saída local. Seus arquivos não devem ser
  usados como gabaritos versionados.

## Status da implementação

- Scanner: implementado e validado com os gabaritos oficiais.
- Parser: implementado com recursive descent e integrado à CLI.
- XML do parser: validado automaticamente contra `ExpressionLessSquare`,
  `Square` e `ArrayTest` (7 arquivos, todos idênticos).
- Tabela de símbolos (`symtab`): implementada e testada (100% de cobertura).
- `VMWriter`: implementado e testado (100% de cobertura).
- Geração de código VM no parser: **pendente** (ver checklist da Etapa 3).
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

**Status: em andamento.**

O compilador traduz cada classe Jack para código da máquina virtual do
nand2tetris, gerando um `.vm` por classe na mesma pasta do `.jack`. A geração
é feita em **passada única**: o próprio parser, enquanto analisa, consulta a
tabela de símbolos e emite os comandos pelo `VMWriter`, sem construir uma
árvore intermediária.

| Variável Jack | Kind | Segmento VM |
|---|---|---|
| `static` | `Static` | `static` |
| `field` | `Field` | `this` |
| parâmetro | `Arg` | `argument` |
| `var` | `Var` | `local` |

Em um método, `argument 0` é o próprio objeto (`this`) e os parâmetros
declarados começam em `argument 1`.

#### Já feito

- [x] `VMWriter`: segmentos, comandos aritméticos, `push`/`pop`, `label`,
  `goto`, `if-goto`, `call`, `function` e `return`, com testes unitários;
- [x] tabela de símbolos: escopos de classe e de subrotina, índices por Kind,
  `Define`, `VarCount`, `StartSubroutine` e `Resolve`, com testes unitários.

#### Falta — Etapa 1: expressões e comandos (tag `etapa1`)

- [ ] parser preenche a tabela de símbolos (`static`, `field`, parâmetros e
  `var`);
- [ ] `function Classe.nome nLocals`, emitido depois das declarações `var`;
- [ ] expressões: constantes (`true` = `push constant 1` + `neg`; `false` e
  `null` = `0`), variáveis, operadores binários (`*` e `/` viram
  `Math.multiply` e `Math.divide`) e unários (`neg`, `not`);
- [ ] comandos `let` (variável simples), `do` (com `pop temp 0`) e `return`
  (`void` retorna `push constant 0`);
- [ ] `if`/`else` e `while` com rótulos únicos, inclusive aninhados;
- [ ] chamadas `Classe.funcao(...)` com a contagem de argumentos;
- [ ] CLI grava `Nome.vm` na mesma pasta do `.jack`, para arquivo ou diretório;
- [ ] validar `Seven` e `ConvertToBin` no VM Emulator e criar a tag `etapa1`.

#### Falta — Etapa 2: objetos e arrays

- [ ] construtores: `push constant <nº de fields>`, `call Memory.alloc 1`,
  `pop pointer 0`;
- [ ] métodos: `push argument 0`, `pop pointer 0` na entrada;
- [ ] chamadas de método (`obj.m()` empilha `obj`; `m()` empilha
  `pointer 0`) e `this`;
- [ ] arrays (`a[i]` em expressões e em `let a[i] = ...`, via `pointer 1` e
  `that 0`);
- [ ] strings constantes (`String.new` + `String.appendChar`);
- [ ] validar `Square`, `Average`, `Pong` e `ComplexArrays` no VM Emulator.

#### Falta — entrega final

- [ ] tabela de validação dos seis programas do Project 11 neste README;
- [ ] desafios e aprendizados da Unidade 2;
- [ ] vídeo de demonstração (10 a 20 minutos) e link neste README.

### Melhorias futuras

- adicionar posição de coluna aos tokens e às mensagens de erro.

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

## Uso de Inteligência Artificial

| Item | Descrição |
|---|---|
| Ferramenta | Claude Code (extensão do VS Code) |
| Modelo | Claude Opus 5.5 (`claude-opus-5-5`) |

**Para que a IA foi usada:**

- revisar o repositório para encontrar pendências das entregas anteriores
  (teste automatizado do parser, seção de IA e tabela de validação no README);
- escrever o teste `TestParserXMLOficial` (`parser/oficial_test.go`), que compara
  a saída do parser com os gabaritos oficiais do Project 10;
- atualizar este README;
- planejar a divisão da Unidade 2 (tabela de símbolos, `VMWriter` e geração de
  código VM) em etapas menores;
- explicar a teoria da máquina virtual (pilha, segmentos, desvios e chamadas)
  e implementar, com testes unitários, o pacote `vmwriter` e o pacote
  `symtab`.

<!-- TODO (dupla): descrever aqui se e como a IA foi usada nas entregas do
scanner e do parser, ou declarar que não foi usada. -->

**Papel da dupla:** o scanner e o parser foram implementados pela dupla. As
sugestões da IA foram revisadas pela dupla e validadas executando `go test ./...`
e `go vet ./...` e comparando a saída com os gabaritos oficiais do nand2tetris.
A dupla é responsável por todo o código entregue e sabe explicá-lo.

## Referências

- [Nand2Tetris](https://www.nand2tetris.org/)
- [Software oficial do Nand2Tetris](https://www.nand2tetris.org/software)
- [Gramática Jack — Apêndice A](https://lambdageo-edu.github.io/book-compiladores-com-nand2tetris/apendices/a-gramatica-jack/)
- [Repositório do projeto](https://github.com/danielduartt/jackcompiler)
