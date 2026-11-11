package parser

import "testing"

func TestCompileExpression(t *testing.T) {
	for _, tc := range []struct{ src, want string }{
		{"2 + 3 * 4", "push constant 2\npush constant 3\nadd\npush constant 4\ncall Math.multiply 2\n"},
		{"2 + (3 * -4)", "push constant 2\npush constant 3\npush constant 4\nneg\ncall Math.multiply 2\nadd\n"},
		{"~false = true", "push constant 0\nnot\npush constant 1\nneg\neq\n"},
		{"null", "push constant 0\n"},
	} {
		got, err := CompileExpression(tc.src)
		if err != nil || got != tc.want {
			t.Errorf("%q: VM=%q, esperado=%q, erro=%v", tc.src, got, tc.want, err)
		}
	}
}

func TestCompileExpressionRejectsIncompleteOutput(t *testing.T) {
	for _, src := range []string{"1 +", "1 2", "(3", "1 + x", `"abc"`} {
		if got, err := CompileExpression(src); err == nil || got != "" {
			t.Errorf("%q: esperado erro sem VM parcial, obtido %q, %v", src, got, err)
		}
	}
}
