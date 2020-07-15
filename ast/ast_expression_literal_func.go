package ast

import (
	"bytes"
	"strings"

	"github.com/dreblang/core/token"
)

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString(token.LeftParen)
	out.WriteString(strings.Join(params, token.Comma+" "))
	out.WriteString(token.RightParen)
	out.WriteString(fl.Body.String())

	return out.String()
}
