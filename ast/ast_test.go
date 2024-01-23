package ast

import (
	"monkeyinterpreter/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name:  &Indentifier{Token: token.Token{Type: token.IDENT, Literal: "myVar"}, Value: "myVar"},
				Value: &Indentifier{Token: token.Token{Type: token.IDENT, Literal: "anotherVar"}, Value: "anotherVar"},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
