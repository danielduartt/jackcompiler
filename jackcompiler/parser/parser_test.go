package parser

import (
	"strings"
	"testing"
)

func TestParseMinimalClass(t *testing.T) {
	got, err := New("class Main {}").Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	want := strings.Join([]string{
		"<class>",
		"  <keyword> class </keyword>",
		"  <identifier> Main </identifier>",
		"  <symbol> { </symbol>",
		"  <symbol> } </symbol>",
		"</class>",
		"",
	}, "\n")
	if got != want {
		t.Errorf("unexpected XML:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestParseClassDeclarationsAndStatements(t *testing.T) {
	src := `class Main {
  field int x, y;
  static boolean ready;

  function void test(int a, boolean b) {
    var int i;
    let x = a + 1;
    let y[2] = -x;
    if (b) {
      do Output.printInt(x);
    } else {
      while (i < 10) {
        let i = i + 1;
      }
    }
    return;
  }
}`

	got, err := New(src).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	for _, tag := range []string{
		"<classVarDec>",
		"<subroutineDec>",
		"<parameterList>",
		"<varDec>",
		"<letStatement>",
		"<ifStatement>",
		"<whileStatement>",
		"<doStatement>",
		"<returnStatement>",
		"<expressionList>",
	} {
		if !strings.Contains(got, tag) {
			t.Errorf("XML does not contain %s:\n%s", tag, got)
		}
	}

	assertXMLOrder(t, got,
		"<classVarDec>",
		"<subroutineDec>",
		"<varDec>",
		"<statements>",
		"<letStatement>",
		"<ifStatement>",
		"<whileStatement>",
		"<returnStatement>",
	)
}

func TestParseExpressionsAndSubroutineCalls(t *testing.T) {
	src := `class Main {
  function int calc() {
    return ~(1 + 2) * -3;
  }

  method void run() {
    do Worker.run(true, null, this);
    return;
  }
}`

	got, err := New(src).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	for _, fragment := range []string{
		"<keyword> true </keyword>",
		"<keyword> null </keyword>",
		"<keyword> this </keyword>",
		"<symbol> ~ </symbol>",
		"<symbol> ( </symbol>",
		"<symbol> + </symbol>",
		"<symbol> * </symbol>",
		"<expressionList>",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("XML does not contain %s:\n%s", fragment, got)
		}
	}
}

func TestParseSyntaxErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		line int
	}{
		{
			name: "missing semicolon",
			src:  "class Main {\n  function void main() {\n    return\n  }\n}",
			line: 4,
		},
		{
			name: "invalid term",
			src:  "class Main {\n  function void main() {\n    let x = ;\n  }\n}",
			line: 3,
		},
		{
			name: "trailing token",
			src:  "class Main {} let",
			line: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.src).Parse()
			if err == nil {
				t.Fatal("expected parse error")
			}

			syntaxErr, ok := err.(*SyntaxError)
			if !ok {
				t.Fatalf("expected *SyntaxError, got %T: %v", err, err)
			}
			if syntaxErr.Line != test.line {
				t.Errorf("error line = %d, want %d", syntaxErr.Line, test.line)
			}
			if syntaxErr.Msg == "" {
				t.Error("syntax error message is empty")
			}
		})
	}
}

func assertXMLOrder(t *testing.T, xml string, fragments ...string) {
	t.Helper()

	offset := 0
	for _, fragment := range fragments {
		index := strings.Index(xml[offset:], fragment)
		if index < 0 {
			t.Errorf("XML does not contain %s after offset %d", fragment, offset)
			return
		}
		offset += index + len(fragment)
	}
}
