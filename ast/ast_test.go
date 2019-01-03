package ast

import (
	"monk/token"
	"testing"
)

func TestString001(t *testing.T) {
	program := &Program {
		Statements: []Statement {
			&LetStatement {
				Token: token.Token {Type: token.LET, Literal: "let"},
				Name: &Identifier {
					Token: token.Token {Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier {
					Token: token.Token {Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program String() wrong, got=%q", program.String())
	}
}

func TestString002(t *testing.T) {
	program := &Program {
		Statements: []Statement {
			&ReturnStatement {
				Token: token.Token {Type: token.RETURN, Literal: "return"},
				ReturnValue: &Identifier {
					Token: token.Token {Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
			},
		},
	}

	if program.String() != "return foo;" {
		t.Errorf("program String() wrong, got=%q", program.String())
	}
}
