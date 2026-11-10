package parser

import "testing"

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
