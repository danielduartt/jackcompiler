package parser

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestEmptyClassXML(t *testing.T) {
	got, err := Parse("class Empty {}")
	if err != nil {
		t.Fatal(err)
	}
	want := "<class>\n  <keyword> class </keyword>\n  <identifier> Empty </identifier>\n  <symbol> { </symbol>\n  <symbol> } </symbol>\n</class>\n"
	if got != want {
		t.Fatalf("XML esperado:\n%s\nobtido:\n%s", want, got)
	}
}

func TestCompleteGrammar(t *testing.T) {
	input := `class Demo {
        static int count, total;
        field Array data;
        field boolean active;
        field char letter;
        constructor Demo new(int size, boolean enabled, char c, Array source) {
            var int i, value;
            var Demo other;
            let data = source;
            let active = enabled;
            let letter = c;
            let i = 0;
            let other = null;
            while (i < size) {
                let data[i] = -(i + 1) * 2 / 1 - 3;
                let value = data[i] + count;
                if (~(value > 0) | (value = 0) & true) {
                    do set(i, value);
                } else {
                    do other.set(i, Demo.identity(value));
                }
                let i = i + 1;
            }
            return this;
        }
        method void set(int index, int value) {
            let data[index] = value;
            if (false) {} else {}
            while (false) {}
            do Output.printString("// /* */ <>&");
            do Output.printString("");
            do flush();
            return;
        }
        function int identity(int n) { return n; }
        function void flush() {}
    }`
	got, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	xmlEvents(t, got) // Também confirma que o XML é bem formado.
	for _, tag := range []string{"classVarDec", "subroutineDec", "parameterList", "subroutineBody",
		"varDec", "statements", "letStatement", "ifStatement", "whileStatement", "doStatement",
		"returnStatement", "expression", "term", "expressionList"} {
		if !strings.Contains(got, "<"+tag+">") {
			t.Errorf("faltou a tag %s", tag)
		}
	}
	for _, forbidden := range []string{"<statement>", "<op>", "<unaryOp>", "<subroutineCall>", "<type>"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("tag fora do formato oficial: %s", forbidden)
		}
	}
}

func TestExpressionStructure(t *testing.T) {
	got, err := Parse("class Main { function int f() { return 2 + 3 * -4; } }")
	if err != nil {
		t.Fatal(err)
	}
	wantExpression := `<expression>
  <term><integerConstant> 2 </integerConstant></term>
  <symbol> + </symbol>
  <term><integerConstant> 3 </integerConstant></term>
  <symbol> * </symbol>
  <term><symbol> - </symbol><term><integerConstant> 4 </integerConstant></term></term>
</expression>`
	start := strings.Index(got, "<expression>")
	end := strings.Index(got, "</expression>") + len("</expression>")
	if start < 0 || end < len("</expression>") {
		t.Fatal("expressão ausente")
	}
	if !reflect.DeepEqual(xmlEvents(t, got[start:end]), xmlEvents(t, wantExpression)) {
		t.Fatalf("estrutura de operadores incorreta:\n%s", got[start:end])
	}
}

func TestStringContentAndXMLSymbols(t *testing.T) {
	got, err := Parse(`class Main { function void f() {
        do Output.printString("  a & <b> // /* */ \n\t  ");
        if (1 < 2 & 3 > 0) { return; }
    } }`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `<stringConstant>   a &amp; &lt;b&gt; // /* */ \n\t   </stringConstant>`) {
		t.Fatalf("conteúdo da string foi alterado:\n%s", got)
	}
	for _, escaped := range []string{"&lt;", "&gt;", "&amp;"} {
		if !strings.Contains(got, "<symbol> "+escaped+" </symbol>") {
			t.Errorf("símbolo não escapado: %s", escaped)
		}
	}
	if got := xmlEscaper.Replace(`"&<>`); got != "&quot;&amp;&lt;&gt;" {
		t.Errorf("escape XML incorreto: %s", got)
	}
	xmlEvents(t, got)
}

func TestSyntaxErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"entrada vazia", ""},
		{"classe ausente", "function void f() {}"},
		{"nome reservado", "class class {}"},
		{"classe truncada", "class Main {"},
		{"tipo void em campo", "class Main { field void x; }"},
		{"campo após função", "class Main { function void f() {} field int x; }"},
		{"parâmetro sem nome", "class Main { function void f(int) {} }"},
		{"vírgula final em parâmetros", "class Main { function void f(int x,) {} }"},
		{"variável após comando", "class Main { function void f() { return; var int x; } }"},
		{"expressão ausente", "class Main { function void f() { let x = ; } }"},
		{"operador sem termo", "class Main { function int f() { return 1 +; } }"},
		{"colchete ausente", "class Main { function void f() { let a[1 = 2; } }"},
		{"chamada sem parênteses", "class Main { function void f() { do run; } }"},
		{"argumento final vazio", "class Main { function void f() { do run(1,); } }"},
		{"else sem bloco", "class Main { function void f() { if (true) {} else return; } }"},
		{"constante inválida", "class Main { function int f() { return class; } }"},
		{"inteiro fora do limite", "class Main { function int f() { return 32768; } }"},
		{"inteiro enorme", "class Main { function int f() { return 99999999999999999999999; } }"},
		{"token ilegal", "class Main { function int f() { return @; } }"},
		{"string sem fechamento", `class Main { function void f() { do run("texto`},
		{"string multilinha", "class Main { function void f() { do run(\"a\nb\"); } }"},
		{"duas classes", "class A {} class B {}"},
		{"lixo após classe", "class A {} @"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			var syntaxError *Error
			if !errors.As(err, &syntaxError) {
				t.Fatalf("esperado erro localizado, obteve XML=%q, erro=%v", got, err)
			}
			if got != "" || syntaxError.Line < 1 || syntaxError.Column < 1 {
				t.Fatalf("XML parcial ou posição inválida: XML=%q, erro=%v", got, err)
			}
		})
	}
}

func TestErrorLocationAfterComments(t *testing.T) {
	input := "// cabeçalho\nclass Main {\n/* comentário\n em duas linhas */\n function void f() {\n  let x = 1\n  return;\n }\n}"
	_, err := Parse(input)
	var syntaxError *Error
	if !errors.As(err, &syntaxError) {
		t.Fatalf("esperado erro sintático, obteve %v", err)
	}
	if syntaxError.Line != 7 || syntaxError.Column != 3 || syntaxError.Expected != ";" {
		t.Fatalf("posição incorreta: %+v", syntaxError)
	}
}

func TestOfficialXML(t *testing.T) {
	// Fixtures oficiais, não saídas produzidas por este parser.
	for _, suite := range []string{"Square", "ExpressionLessSquare", "ArrayTest"} {
		paths, err := filepath.Glob(filepath.Join("testdata", "official", suite, "*.jack"))
		if err != nil || len(paths) == 0 {
			t.Fatalf("fixtures de %s ausentes: %v", suite, err)
		}
		for _, path := range paths {
			t.Run(suite+"/"+filepath.Base(path), func(t *testing.T) {
				input, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				expectedPath := strings.TrimSuffix(path, ".jack") + ".xml"
				expected, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatal(err)
				}
				got, err := Parse(string(input))
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(xmlEvents(t, got), xmlEvents(t, string(expected))) {
					t.Fatalf("XML difere estruturalmente de %s", expectedPath)
				}
			})
		}
	}
}

// Compara hierarquia, ordem e conteúdo dos terminais sem depender da indentação.
func xmlEvents(t *testing.T, input string) []string {
	t.Helper()
	decoder := xml.NewDecoder(strings.NewReader(input))
	var events []string
	var parent string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return events
		}
		if err != nil {
			t.Fatalf("XML inválido: %v", err)
		}
		switch token := token.(type) {
		case xml.StartElement:
			parent = token.Name.Local
			events = append(events, "<"+parent+">")
		case xml.EndElement:
			parent = ""
			events = append(events, "</"+token.Name.Local+">")
		case xml.CharData:
			text := string(token)
			if parent == "stringConstant" {
				// Retira apenas o espaço de apresentação, preservando o literal.
				events = append(events, strings.TrimSuffix(strings.TrimPrefix(text, " "), " "))
			} else if text = strings.TrimSpace(text); text != "" {
				events = append(events, text)
			}
		}
	}
}
