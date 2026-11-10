package vmwriter

import (
	"fmt"

	"github.com/danielduartt/jackcompiler/scanner"
)

// WriteBinaryOp combina os dois operandos já empilhados, deixando o resultado no topo.
func (w *VMWriter) WriteBinaryOp(op scanner.TokenType) error {
	switch op {
	case scanner.PLUS:
		w.WriteArithmetic(Add)
	case scanner.MINUS:
		w.WriteArithmetic(Sub)
	case scanner.AND:
		w.WriteArithmetic(And)
	case scanner.OR:
		w.WriteArithmetic(Or)
	case scanner.LT:
		w.WriteArithmetic(Lt)
	case scanner.GT:
		w.WriteArithmetic(Gt)
	case scanner.EQ:
		w.WriteArithmetic(Eq)
	case scanner.ASTERISK:
		w.WriteCall("Math.multiply", 2)
	case scanner.SLASH:
		w.WriteCall("Math.divide", 2)
	default:
		return fmt.Errorf("token %v não é operador binário Jack", op)
	}
	return nil
}
