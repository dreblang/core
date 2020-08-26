package ast

import (
	"bytes"

	"github.com/dreblang/core/token"
)

type LoopExpression struct {
	Token       token.Token // The 'loop' token
	Condition   Expression
	Consequence *BlockStatement
}

func (le *LoopExpression) expressionNode()      {}
func (le *LoopExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LoopExpression) String() string {
	var out bytes.Buffer

	out.WriteString("loop")
	out.WriteString(le.Condition.String())
	out.WriteString(" ")
	out.WriteString(le.Consequence.String())

	return out.String()
}
