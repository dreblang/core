package ast

import (
	"bytes"

	"github.com/dreblang/core/token"
)

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString(token.LeftParen)
	out.WriteString(ie.Left.String())
	out.WriteString(token.LeftBracket)
	out.WriteString(ie.Index.String())
	out.WriteString(token.RightBracket)
	out.WriteString(token.RightParen)

	return out.String()
}
