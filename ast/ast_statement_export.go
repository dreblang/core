package ast

import (
	"bytes"

	"github.com/dreblang/core/token"
)

type ExportStatement struct {
	Token      token.Token // the 'return' token
	Identifier *Identifier
}

func (rs *ExportStatement) statementNode()       {}
func (rs *ExportStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ExportStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Identifier != nil {
		out.WriteString(rs.Identifier.String())
	}

	out.WriteString(token.Semicolon)
	return out.String()
}
