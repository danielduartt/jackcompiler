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
