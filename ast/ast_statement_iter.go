package ast

import (
	"math/rand"

	"github.com/dreblang/core/token"
)

type IterStatement struct {
	Token      token.Token // the token.Let token
	Identifier *Identifier
	Expression Expression
	Statements *BlockStatement
}

func (is *IterStatement) statementNode()       {}
func (ls *IterStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *IterStatement) String() string {
	return ""
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (is *IterStatement) ConvertToLoop() []Statement {
	stmts := []Statement{}

	identName := RandStringBytes(16)
	ident := &Identifier{Value: identName}

	stmt1 := &LetStatement{
		Token: token.Token{Type: token.Let},
		Name:  ident,
		Value: &IntegerLiteral{
			Token: token.Token{Type: token.Int, Literal: "0"},
			Value: 0,
		},
	}

	consequence := []Statement{}
	consequence = append(
		consequence,
		&LetStatement{
			Name: is.Identifier,
			Value: &IndexExpression{
				Left:     is.Expression,
				Index:    ident,
				HasSkip:  false,
				HasUpper: false,
			},
		},
	)
	consequence = append(consequence, is.Statements.Statements...)
	consequence = append(
		consequence,
		&LetStatement{
			Name: ident,
			Value: &InfixExpression{
				Left:     ident,
				Operator: token.Plus,
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.Int, Literal: "1"},
					Value: 1,
				},
			},
		},
	)

	stmt2 := &ExpressionStatement{
		Expression: &LoopExpression{
			Token: token.Token{Type: token.Loop},
			Condition: &InfixExpression{
				Left:     &Identifier{Value: identName},
				Operator: token.LessThan,
				Right: &CallExpression{
					Function:  &Identifier{Value: "len"},
					Arguments: []Expression{is.Expression},
				},
			},
			Consequence: &BlockStatement{
				Statements: consequence,
			},
		},
	}

	stmts = append(stmts, stmt1)
	stmts = append(stmts, stmt2)

	return stmts
}
