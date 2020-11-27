package parser

import (
	"fmt"
	"strconv"

	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/token"
)

const (
	_ int = iota
	Lowest
	Assign        // =
	Equals        // ==
	LessOrGreater // < or >
	Sum           // +
	Product       // *
	Prefix        // -X or !X
	Dot           // obj.member
	Scope         // scope::item
	Call          // myFunction(X)
	Index         // array[index]
)

var precedences = map[token.TokenType]int{
	token.Assign:         Assign,
	token.Equal:          Equals,
	token.NotEqual:       Equals,
	token.LessThan:       LessOrGreater,
	token.LessOrEqual:    LessOrGreater,
	token.GreaterThan:    LessOrGreater,
	token.GreaterOrEqual: LessOrGreater,
	token.Plus:           Sum,
	token.Minus:          Sum,
	token.Slash:          Product,
	token.Asterisk:       Product,
	token.Percent:        Product,
	token.LeftParen:      Call,
	token.LeftBracket:    Index,
	token.Dot:            Dot,
	token.DoubleColon:    Scope,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Identifier, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.Float, p.parseFloatLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.LeftParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Loop, p.parseLoopExpression)
	p.registerPrefix(token.Function, p.parseFunctionLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.LeftBracket, p.parseArrayLiteral)
	p.registerPrefix(token.LeftBrace, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Assign, p.parseInfixExpression)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Percent, p.parseInfixExpression)
	p.registerInfix(token.Equal, p.parseInfixExpression)
	p.registerInfix(token.NotEqual, p.parseInfixExpression)
	p.registerInfix(token.LessThan, p.parseInfixExpression)
	p.registerInfix(token.LessOrEqual, p.parseInfixExpression)
	p.registerInfix(token.GreaterThan, p.parseInfixExpression)
	p.registerInfix(token.GreaterOrEqual, p.parseInfixExpression)
	p.registerInfix(token.Dot, p.parseMemberExpression)
	p.registerInfix(token.DoubleColon, p.parseScopeResolutionExpression)
	p.registerInfix(token.LeftParen, p.parseCallExpression)
	p.registerInfix(token.LeftBracket, p.parseIndexExpression)

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			switch stmtT := stmt.(type) {
			case *ast.IterStatement:
				stmts := stmtT.ConvertToLoop()
				program.Statements = append(program.Statements, stmts...)

			default:
				program.Statements = append(program.Statements, stmtT)
			}
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Statements

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Scope:
		return p.parseScopeDefinition()
	case token.Export:
		return p.parseExportStatement()
	case token.Load:
		return p.parseLoadStatement()
	case token.Iter:
		return p.parseIterStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	// get value
	p.nextToken()
	stmt.Value = p.parseExpression(Lowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	// get value
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(Lowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseScopeDefinition() *ast.ScopeDefinition {
	sd := &ast.ScopeDefinition{Token: p.currentToken}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	sd.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.LeftBrace) {
		return nil
	}
	sd.Block = p.parseBlockStatement()
	return sd
}

func (p *Parser) parseExportStatement() *ast.ExportStatement {
	stmt := &ast.ExportStatement{Token: p.currentToken}

	// get value
	p.nextToken()
	stmt.Identifier = p.parseIdentifier().(*ast.Identifier)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLoadStatement() *ast.LoadStatement {
	stmt := &ast.LoadStatement{Token: p.currentToken}

	// get value
	p.nextToken()
	stmt.Identifier = p.parseIdentifier().(*ast.Identifier)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(Lowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIterStatement() *ast.IterStatement {
	stmt := ast.IterStatement{Token: p.currentToken}
	p.nextToken()

	stmt.Identifier = p.parseIdentifier().(*ast.Identifier)

	if p.peekTokenIs(token.Over) {
		p.nextToken()
		p.nextToken()
	}

	stmt.Expression = p.parseExpression(Lowest)

	if p.peekTokenIs(token.LeftBrace) {
		p.nextToken()
	}

	stmt.Statements = p.parseBlockStatement()
	return &stmt
}

// Expressions

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix expressions

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, _ := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.currentToken}

	value, _ := strconv.ParseFloat(p.currentToken.Literal, 64)

	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	// get expression
	p.nextToken()
	expression.Right = p.parseExpression(Prefix)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.True)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// get expression
	p.nextToken()
	exp := p.parseExpression(Lowest)

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LeftParen) {
		return nil
	}

	// get condition
	p.nextToken()
	expression.Condition = p.parseExpression(Lowest)

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	if !p.expectPeek(token.LeftBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// TODO parse else if
	if p.peekTokenIs(token.Else) {
		p.nextToken()

		if !p.expectPeek(token.LeftBrace) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(token.RightBrace) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LeftParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LeftBrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.RightParen) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}

	array.Elements = p.parseExpressionList(token.RightBracket)

	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(Lowest))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(Lowest))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}

	if !p.peekTokenIs(token.Comma) {
		p.nextToken()
		exp.Index = p.parseExpression(Lowest)
	}

	if p.peekTokenIs(token.Comma) {
		p.nextToken() // Skip to comma

		exp.HasUpper = true
		if !(p.peekTokenIs(token.RightBracket) || p.peekTokenIs(token.Comma)) {
			p.nextToken()
			exp.IndexUpper = p.parseExpression(Lowest)
		}
	}

	if p.peekTokenIs(token.Comma) {
		p.nextToken() // Skip to comma
		exp.HasSkip = true
		if !p.peekTokenIs(token.RightBracket) {
			p.nextToken()
			exp.IndexSkip = p.parseExpression(Lowest)
		}
	}

	if !p.expectPeek(token.RightBracket) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RightBrace) {
		p.nextToken()
		key := p.parseExpression(Lowest)

		if _, ok := key.(*ast.Identifier); ok {
			key = &ast.StringLiteral{
				Value: key.String(),
			}
		}

		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(Lowest)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RightBrace) && !p.expectPeek(token.Comma) {
			return nil
		}
	}

	p.expectPeek(token.RightBrace)
	return hash
}

// Infix expressions

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	p.nextToken()
	right := p.parseStringLiteral()
	expression.Right = right
	return expression
}

func (p *Parser) parseScopeResolutionExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	p.nextToken()
	right := p.parseStringLiteral()
	expression.Right = right
	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RightParen)
	return exp
}

// Token

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Precedence

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return Lowest
}

// Error

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// Operators

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseLoopExpression() ast.Expression {
	expression := &ast.LoopExpression{Token: p.currentToken}

	if !p.expectPeek(token.LeftParen) {
		return nil
	}

	// get condition
	p.nextToken()
	expression.Condition = p.parseExpression(Lowest)

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	if !p.expectPeek(token.LeftBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	return expression
}
