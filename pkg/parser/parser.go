package parser

import (
	"fmt"
	"strconv"

	"github.com/sbrki/monkey/pkg/ast"
	"github.com/sbrki/monkey/pkg/lexer"
	"github.com/sbrki/monkey/pkg/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // <=
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // foo(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func() ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// read two tokens, so that currToken and peekToken are set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{
		Token: p.currToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	letStmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO(sbrki): parse expressions
	// -- skipping until semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{
		Token: p.currToken,
	}

	p.nextToken()

	// TODO(sbrki): parse expressions
	// -- skipping until semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	exprStmt := &ast.ExpressionStatement{
		Token: p.currToken,
	}

	exprStmt.Expression = p.parseExpression(LOWEST)

	// expression statements have optional semicolons
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return exprStmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

//////////////////////////////
// expression prefix/infix parse fns

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intLit := &ast.IntegerLiteral{
		Token: p.currToken,
	}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse integer literal '%s'", p.currToken.Literal)
		p.errors = append(p.errors, msg)
	}

	intLit.Value = value
	return intLit
}

//////////////////////////////
// parser utilities

func (p *Parser) curTokenIs(targetType token.TokenType) bool {
	return p.currToken.Type == targetType
}

func (p *Parser) peekTokenIs(targetType token.TokenType) bool {
	return p.peekToken.Type == targetType
}

func (p *Parser) expectPeek(targetType token.TokenType) bool {
	if p.peekTokenIs(targetType) {
		p.nextToken()
		return true
	}
	p.peekError(targetType)
	return false
}

func (p *Parser) peekError(expectedType token.TokenType) {
	p.errors = append(
		p.errors,
		fmt.Sprintf(
			"expected token = '%s' , got = '%s'",
			expectedType,
			p.peekToken.Type,
		),
	)
}
