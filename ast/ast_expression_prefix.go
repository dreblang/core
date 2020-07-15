package ast

import (
	"bytes"

	"github.com/dreblang/core/token"
)

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(token.LeftParen)
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(token.RightParen)

	return out.String()
}
