package ast

import (
	"github.com/dreblang/core/token"
)

type Identifier struct {
	Token token.Token // the token.Identifier token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
