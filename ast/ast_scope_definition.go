package ast

import (
	"bytes"

	"github.com/dreblang/core/token"
)

type ScopeDefinition struct {
	Token token.Token // the { token
	Name  *Identifier
	Block *BlockStatement
}

func (sd *ScopeDefinition) statementNode()       {}
func (sd *ScopeDefinition) TokenLiteral() string { return sd.Token.Literal }
func (sd *ScopeDefinition) String() string {
	var out bytes.Buffer

	out.WriteString("scope ")
	out.WriteString(sd.Name.String())
	out.WriteString(" {\n")

	for _, s := range sd.Block.Statements {
		out.WriteString(s.String() + "\n")
	}

	out.WriteString("}\n")

	return out.String()
}
