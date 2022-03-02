package parser

import (
	"fmt"

	"github.com/sbrki/monkey/pkg/ast"
	"github.com/sbrki/monkey/pkg/lexer"
	"github.com/sbrki/monkey/pkg/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
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
		return nil
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
			"expected token =  '%s' , got = '%s'",
			expectedType,
			p.peekToken.Type,
		),
	)
}
