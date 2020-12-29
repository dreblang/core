package ast

import (
	"bytes"

	"github.com/dreblang/core/token"
)

type ClassDefinition struct {
	Token token.Token // the { token
	Name  *Identifier
	Block *BlockStatement
}

func (sd *ClassDefinition) statementNode()       {}
func (sd *ClassDefinition) TokenLiteral() string { return sd.Token.Literal }
func (sd *ClassDefinition) String() string {
	var out bytes.Buffer

	out.WriteString("Class ")
	out.WriteString(sd.Name.String())
	out.WriteString(" {\n")

	for _, s := range sd.Block.Statements {
		out.WriteString(s.String() + "\n")
	}

	out.WriteString("}\n")

	return out.String()
}

func (sd *ClassDefinition) ConvertToFunc() []Statement {
	return []Statement{
		&LetStatement{
			Name: sd.Name,
			Value: &FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ScopeDefinition{
							Name: &Identifier{
								Value: sd.Name.Value + "Class",
							},
							Block: &BlockStatement{
								Statements: append([]Statement{}, sd.Block.Statements...),
							},
						},
						&ReturnStatement{
							ReturnValue: &Identifier{
								Value: sd.Name.Value + "Class",
							},
						},
					},
				},
			},
		},
	}
}
