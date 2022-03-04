package parser

import (
	"testing"

	"github.com/sbrki/monkey/pkg/ast"
	"github.com/sbrki/monkey/pkg/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("len(program.Statements) = %d, expected = 3",
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral = %q, expected = 'let'", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("Could not downcast ast.Statement to ast.LetStatement. got = %q", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value = %s, expected = %s",
			letStmt.Name.Value, name)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() = %s, expected = %s",
			letStmt.Name.TokenLiteral(), name)
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
	return 10;
	return 42;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("len(program.Statements) = %d, expected = 3",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Could not downcast ast.Statement to ast.ReturnStatement. got = %q", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() = '%q', expected = 'return'",
				returnStmt.TokenLiteral())
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) = %d, expected = 1",
			len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Could not downcast ast.Statement to ast.ExpressionStatement. got = %q", program.Statements[0])
	}

	ident, ok := exprStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.Identifier. got = %q", exprStmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value = '%q', expected = 'foobar'", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() = '%q', exptected = 'foobar'", ident.TokenLiteral())
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser encountered %d errors:", len(errors))
	for idx, e := range errors {
		t.Errorf("\t%d: %s", idx, e)
	}
	t.FailNow()
}
