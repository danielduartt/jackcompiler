package parser

import "testing"

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
