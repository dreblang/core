package ast

import (
	"bytes"
	"strings"

	"github.com/dreblang/core/token"
)

type HashLiteral struct {
	Token token.Token // The '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+token.Colon+value.String())
	}

	out.WriteString(token.LeftBrace)
	out.WriteString(strings.Join(pairs, token.Comma+" "))
	out.WriteString(token.RightBrace)

	return out.String()
}
