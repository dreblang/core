package ast

import (
	"bytes"
	"strings"

	"github.com/dreblang/core/token"
)

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString(token.LeftBracket)
	out.WriteString(strings.Join(elements, token.Comma+" "))
	out.WriteString(token.RightBracket)

	return out.String()
}
