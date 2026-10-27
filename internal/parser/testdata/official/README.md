# Arquivos oficiais de comparação

Origem: [nand2tetris/projects, projeto 10](https://github.com/nand2tetris/projects/tree/5bd2184edfd4945ff8a28eb0b237c12cdbe9e498/projects/10).
Obtidos em 27/10/2026. A licença original acompanha estes arquivos em `LICENSE`.

Foram incluídos apenas programas `.jack` e seus XML de análise sintática,
sem implementar ou copiar um parser externo. Os XML são as referências oficiais,
não resultados gerados pelo JackCompiler.

- `Square`: Main, Square e SquareGame.
- `ExpressionLessSquare`: as três classes, na versão sem expressões complexas.
- `ArrayTest`: Main, com indexação de arrays.

O pacote oficial chama as referências `Main.xml`, `Square.xml` e
`SquareGame.xml`. Elas correspondem aos arquivos chamados `MainP.xml`,
`SquareP.xml` e `SquareGameP.xml` no enunciado. O teste `TestOfficialXML`
compara a árvore e os terminais de cada referência com a saída do parser,
ignorando a indentação e preservando espaços dentro de strings.
