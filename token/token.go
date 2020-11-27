package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	// Identifiers + Literals
	Identifier = "Identifier" // add, x ,y, ...
	Int        = "Int"        // 123456
	Float      = "Float"
	String     = "String" // "x", "y"

	// Operators
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"
	Percent  = "%"
	Equal    = "=="
	NotEqual = "!="

	LessThan       = "<"
	LessOrEqual    = "<="
	GreaterThan    = ">"
	GreaterOrEqual = ">="

	// Delimiters
	Comma       = ","
	Semicolon   = ";"
	Colon       = ":"
	DoubleColon = "::"
	Dot         = "."

	LeftParen    = "("
	RightParen   = ")"
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"

	// Keywords
	Function = "Function"
	Let      = "Let"
	True     = "True"
	False    = "False"
	If       = "If"
	Else     = "Else"
	Return   = "Return"
	Loop     = "Loop"
	Scope    = "Scope"
	Export   = "Export"
	Load     = "Load"
	Iter     = "Iter"
	Over     = "Over"
)

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
	"loop":   Loop,
	"scope":  Scope,
	"export": Export,
	"load":   Load,
	"iter":   Iter,
	"over":   Over,
}

func LookupIdentifierType(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return Identifier
}
