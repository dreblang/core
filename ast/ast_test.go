package ast

import (
	"testing"

	"github.com/dreblang/core/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "mayVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestSuperProgram(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "mayVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
			&ExpressionStatement{
				Expression: &IfExpression{
					Condition: &InfixExpression{
						Left: &Identifier{
							Token: token.Token{Type: token.Identifier, Literal: "a"},
							Value: "a",
						},
						Right: &Identifier{
							Token: token.Token{Type: token.Identifier, Literal: "b"},
							Value: "b",
						},
						Operator: "<",
					},
					Consequence: &BlockStatement{},
					Alternative: &BlockStatement{},
				},
			},
		},
	}

	if program.String() == "" {
		t.Error("Expected program string")
	}
}

// FIXME: There is no real point in testing this out thoroughly, because this is just a representation of a program
