package parser

import (
	"strings"
	"testing"
)

func TestCompileSeven(t *testing.T) {
	src := "class Main { function void main() { do Output.printInt(1 + (2 * 3)); return; } }"
	want := "function Main.main 0\npush constant 1\npush constant 2\npush constant 3\n" +
		"call Math.multiply 2\nadd\ncall Output.printInt 1\npop temp 0\npush constant 0\nreturn\n"
	got, err := CompileClass(src)
	if err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

func TestDoArgumentCounts(t *testing.T) {
	src := "class App { function void run() { do Screen.clearScreen(); do Output.moveCursor(3, 4); return; } }"
	want := "function App.run 0\ncall Screen.clearScreen 0\npop temp 0\npush constant 3\npush constant 4\n" +
		"call Output.moveCursor 2\npop temp 0\npush constant 0\nreturn\n"
	if got, err := CompileClass(src); err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

func TestDoRejectsObjectReceiver(t *testing.T) {
	got, err := CompileClass("class App { function void run(Game g) { do g.play(); return; } }")
	if err == nil || got != "" {
		t.Fatalf("chamada de método aceita como função de classe: VM = %q, erro = %v", got, err)
	}
}

func TestCompileVoidFunction(t *testing.T) {
	got, err := CompileClass("class Main { function void main() { return; } }")
	want := "function Main.main 0\npush constant 0\nreturn\n"
	if err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

func TestCompileClassRejectsPartialVM(t *testing.T) {
	got, err := CompileClass("class App { function void start() { return; } method void move() { return; } }")
	if err == nil || got != "" {
		t.Fatalf("esperado erro sem VM parcial, obtido %q, %v", got, err)
	}
}

// O "function" leva o número de variáveis locais, contado depois de todos os var.
func TestFunctionLocalsCount(t *testing.T) {
	src := "class Main { function int f(int a, int b) { var int x, y; var boolean ok; return; } }"
	want := "function Main.f 3\npush constant 0\nreturn\n"
	if got, err := CompileClass(src); err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

// Variáveis viram push/pop no segmento do seu Kind: arg -> argument,
// var -> local, static -> static.
func TestLetAndReturnWithVariables(t *testing.T) {
	src := `class Main {
  static int count;
  function int f(int a) {
    var int x;
    let x = a + 1;
    let count = count + x;
    return Math.abs(x);
  }
}`
	want := "function Main.f 1\n" +
		"push argument 0\npush constant 1\nadd\npop local 0\n" +
		"push static 0\npush local 0\nadd\npop static 0\n" +
		"push local 0\ncall Math.abs 1\nreturn\n"
	if got, err := CompileClass(src); err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

func TestUndeclaredVariable(t *testing.T) {
	for _, src := range []string{
		"class Main { function void f() { let y = 1; return; } }",
		"class Main { function int f() { return y; } }",
	} {
		got, err := CompileClass(src)
		if err == nil || got != "" {
			t.Errorf("%q: esperado erro de variável não declarada, obtido %q, %v", src, got, err)
			continue
		}
		if !strings.Contains(err.Error(), "'y' não declarada") {
			t.Errorf("%q: mensagem inesperada: %v", src, err)
		}
	}
}

// Arrays são da Etapa 2: por enquanto a geração recusa com erro, em vez de
// produzir código errado.
func TestArraysNotYetSupported(t *testing.T) {
	for _, src := range []string{
		"class Main { function void f(Array a) { let a[0] = 1; return; } }",
		"class Main { function int f(Array a) { return a[0]; } }",
	} {
		if got, err := CompileClass(src); err == nil || got != "" {
			t.Errorf("%q: esperado erro, obtido %q, %v", src, got, err)
		}
	}
}

func TestIfElse(t *testing.T) {
	src := "class Main { function int max(int a, int b) { if (a > b) { return a; } else { return b; } } }"
	want := "function Main.max 0\n" +
		"push argument 0\npush argument 1\ngt\nnot\nif-goto IF_ELSE0\n" +
		"push argument 0\nreturn\ngoto IF_END0\n" +
		"label IF_ELSE0\npush argument 1\nreturn\n" +
		"label IF_END0\n"
	if got, err := CompileClass(src); err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

// if sem else: o rótulo IF_ELSE cai direto no IF_END.
func TestIfWithoutElse(t *testing.T) {
	src := "class Main { function void f(int a) { if (a) { let a = 0; } return; } }"
	want := "function Main.f 0\n" +
		"push argument 0\nnot\nif-goto IF_ELSE0\n" +
		"push constant 0\npop argument 0\ngoto IF_END0\n" +
		"label IF_ELSE0\nlabel IF_END0\n" +
		"push constant 0\nreturn\n"
	if got, err := CompileClass(src); err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

func TestWhile(t *testing.T) {
	src := "class Main { function void f() { var int i; while (i < 10) { let i = i + 1; } return; } }"
	want := "function Main.f 1\n" +
		"label WHILE_EXP0\npush local 0\npush constant 10\nlt\nnot\nif-goto WHILE_END0\n" +
		"push local 0\npush constant 1\nadd\npop local 0\n" +
		"goto WHILE_EXP0\nlabel WHILE_END0\n" +
		"push constant 0\nreturn\n"
	if got, err := CompileClass(src); err != nil || got != want {
		t.Fatalf("VM = %q, erro = %v, esperado %q", got, err, want)
	}
}

// Comandos aninhados e em sequência recebem rótulos diferentes; outra
// subrotina recomeça a contagem do zero.
func TestNestedLabelsAreUnique(t *testing.T) {
	src := `class Main {
  function void f(int a) {
    while (a) {
      if (a) { while (a) { let a = 0; } } else { let a = 1; }
      if (a) { let a = 2; }
    }
    return;
  }
  function void g(int a) {
    while (a) { let a = 0; }
    return;
  }
}`
	got, err := CompileClass(src)
	if err != nil {
		t.Fatal(err)
	}
	var funcs []map[string]bool
	for _, line := range strings.Split(got, "\n") {
		if strings.HasPrefix(line, "function ") {
			funcs = append(funcs, map[string]bool{})
		}
		if label, ok := strings.CutPrefix(line, "label "); ok {
			seen := funcs[len(funcs)-1]
			if seen[label] {
				t.Errorf("rótulo %s repetido na mesma função:\n%s", label, got)
			}
			seen[label] = true
		}
	}
	if len(funcs) != 2 || len(funcs[0]) != 8 || len(funcs[1]) != 2 {
		t.Fatalf("esperado 8 rótulos em f e 2 em g:\n%s", got)
	}
	if !strings.Contains(got, "function Main.g 0\nlabel WHILE_EXP0\n") {
		t.Errorf("a contagem de rótulos deveria recomeçar em g:\n%s", got)
	}
}
