# JackCompiler

Compilador da linguagem **Jack** ([nand2tetris](https://www.nand2tetris.org/))
escrito em **Go** para a disciplina de Compiladores. Ele traduz cada classe Jack
para código da máquina virtual (`.vm`) executável no **VM Emulator** (Project 11).

O compilador tem quatro partes: o **scanner** (`scanner/`) transforma o fonte em
tokens; o **parser** (`parser/`), por descida recursiva, verifica a gramática e
gera o código enquanto analisa; a **tabela de símbolos** (`symtab/`) resolve cada
variável para seu segmento e índice; e o **`VMWriter`** (`vmwriter/`) escreve os
comandos VM, um por linha.

## Equipe

- **Daniel Nunes Duarte** — Matrícula: 20250013550
- **João Victor Oliveira** — Matrícula: 20260001375

## Vídeo de apresentação

<!-- Substituir pelo link do vídeo -->
🎥 **Link do vídeo:** _em breve_

## Compilar e executar

Requisito: Go 1.18 ou superior, sem dependências externas. Os comandos são
executados na pasta `jackcompiler/`, onde está o `go.mod`:

```bash
git clone https://github.com/danielduartt/jackcompiler
cd jackcompiler/jackcompiler
```

**Diretório** (um programa Jack): compila todos os `.jack` da pasta, sem descer
em subpastas, e grava um `.vm` por classe **na mesma pasta**:

```bash
go run . caminho/Square          # gera Main.vm, Square.vm e SquareGame.vm em caminho/Square
```

**Arquivo** individual:

```bash
go run . caminho/Seven/Main.jack # gera caminho/Seven/Main.vm
```

Também é possível gerar um executável:

```bash
go build -o jackcompiler.exe .   # no Linux/macOS: go build -o jackcompiler .
./jackcompiler.exe caminho/Square
```

Depois, abra a pasta do programa no VM Emulator. Os `.vm` do sistema operacional
(`Math`, `Output`, `Memory`...) não precisam ser copiados, pois o emulador já os
traz embutidos. Erros léxicos, sintáticos ou de variável não declarada são
informados com linha e coluna, e a CLI termina com código diferente de zero.

A análise das unidades anteriores continua disponível com `--xml`. Ela gera
`<Nome>T.xml` (tokens) e `<Nome>P.xml` (árvore sintática) em `out/`:

```bash
go run . --xml ./testdata/Square
```

Testes automatizados (unitários do `VMWriter`, da tabela de símbolos e da
geração de código, mais a comparação com os gabaritos XML do Project 10):

```bash
go test ./... -count=1
go vet ./...
```

## Status da validação

Validação do Project 11: executar cada programa no VM Emulator e conferir o
comportamento esperado.

| Etapa | Programa | O que exercita | Como conferir | Resultado |
|---|---|---|---|---|
| 1 | **Seven** | expressão aritmética, `do`, chamada ao SO | imprime `7` no canto da tela | ⏳ compila; falta confirmar no VM Emulator |
| 1 | **ConvertToBin** | `while`, `if` aninhados, `&` e `~` | com um valor em `RAM[8000]`, `RAM[8001..8016]` recebem os bits | ⏳ compila; falta confirmar no VM Emulator |
| 2 | **Square** | construtores, métodos, campos | quadrado se move com as setas; `z`/`x` mudam o tamanho | ⏳ compila; falta confirmar no VM Emulator |
| 2 | **Average** | arrays, strings, teclado | pede N números e imprime a média | ⏳ compila; falta confirmar no VM Emulator |
| 2 | **Pong** | programa orientado a objetos completo | bola, raquete e placar funcionam | ⏳ compila; falta confirmar no VM Emulator |
| 2 | **ComplexArrays** | índices aninhados, arrays de arrays | valor esperado igual ao obtido em cada teste | ⏳ compila; falta confirmar no VM Emulator |

<!-- TODO (dupla): trocar cada ⏳ por ✅ (ou ❌) depois de executar no VM Emulator. -->

O compilador implementa toda a linguagem Jack: expressões, constantes,
variáveis (`static`, `field`, `arg`, `var`), `let`, `do`, `return`,
`if`/`else` e `while` com rótulos únicos, construtores, métodos, `this`,
chamadas de função e de método, arrays e strings constantes. O teste
`TestCompilaProgramasProject11` garante que os seis programas compilam sem erro.

Unidades anteriores: o XML do scanner e do parser é idêntico aos gabaritos
oficiais nos 7 arquivos de `ExpressionLessSquare`, `Square` e `ArrayTest`
(testes `TestXMLOficial` e `TestParserXMLOficial`).

## Desafios e aprendizados

<!-- TODO (dupla): revisar e completar com a experiência de vocês. -->

- **Passada única:** o código é gerado dentro do próprio parser, sem árvore
  intermediária. Por isso a ordem importa. O `function Classe.f nLocals`, por
  exemplo, só pode ser escrito depois de lidas todas as declarações `var`,
  porque antes disso o número de locais ainda não é conhecido.
- **Sem precedência:** Jack avalia `2 + 3 * 4` da esquerda para a direita, como
  `(2 + 3) * 4`. A gramática plana `term (op term)*` gera exatamente essa
  ordem, em notação pós-fixa.
- **Rótulos únicos:** cada `if`/`while` reserva um número antes de analisar seu
  corpo, de modo que comandos aninhados nunca repetem rótulos dentro da função.
- **Escopos e segmentos:** um `static` e um `field` podem ter o mesmo índice 0
  porque ficam em segmentos diferentes. Em métodos, `argument 0` é o próprio
  objeto e os parâmetros começam em `argument 1`.
- **Trabalho em dupla com Git:** integrar o trabalho dos dois no mesmo arquivo
  (`parser.go`) exigiu commits pequenos e atualizar a branch antes de começar.

## Uso de Inteligência Artificial

| Item | Descrição |
|---|---|
| Ferramenta | Claude Code (extensão do VS Code) |
| Modelo | Claude Opus 5.5 (`claude-opus-5-5`) |

**Finalidade:**

- revisar o repositório em busca de pendências das entregas anteriores e
  escrever o teste `TestParserXMLOficial`;
- explicar a teoria da máquina virtual (pilha, segmentos, desvios e chamadas) e
  planejar a Unidade 2 em etapas menores;
- implementar, com testes unitários, o `VMWriter`, a tabela de símbolos e a
  geração de código de variáveis, `let`, `return`, `if`/`else` e `while`;
- revisar este README.

<!-- TODO (dupla): acrescentar outros usos de IA (inclusive nas entregas do
scanner e do parser), ou declarar que não houve. -->

**Como o resultado foi validado:** o código sugerido foi revisado pela dupla e
verificado com testes automatizados (`go test ./...`), que comparam o código VM
gerado com o esperado em cada construção da linguagem; com `go vet ./...`; e
executando os programas do Project 11 no VM Emulator (tabela acima). A dupla é
responsável por todo o código entregue e sabe explicá-lo.

## Vídeo de demonstração

<!-- TODO (dupla): substituir pelo link do vídeo (10 a 20 min). -->
🎥 *Link do vídeo: a ser adicionado.*
