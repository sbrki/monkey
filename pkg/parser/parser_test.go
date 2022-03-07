package parser

import (
	"fmt"
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

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

	intLit, ok := exprStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.IntegerLiteral. got = %q", exprStmt.Expression)
	}

	if intLit.Value != 5 {
		t.Errorf("intLit.Value = %d, expected = 5", intLit.Value)
	}

	if intLit.TokenLiteral() != "5" {
		t.Errorf("intLit.TokenLiteral() = '%s', exptected = '5'", intLit.TokenLiteral())
	}

}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
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

		prefixExpr, ok := exprStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Could not downcast ast.Expression to ast.PrefixExpression. got = %q", exprStmt.Expression)
		}

		if prefixExpr.Operator != tt.operator {
			t.Fatalf("prefixExpr.Operator is not '%s'. got = '%s'", tt.operator, prefixExpr.Operator)
		}

		if !testIntegerLiteral(t, prefixExpr.Right, tt.integerValue) {
			return
		}

	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Could not downcast ast.Expression to ast.IntegerLiteral. got = %q", il)
		return false
	}

	if intLit.Value != value {
		t.Errorf("intLit.Value = %d, expected = %d", intLit.Value, value)
		return false
	}

	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intLit.TokenLiteral() = '%s', exptected = '%s'", intLit.TokenLiteral(), fmt.Sprintf("%d", value))
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("Could not downcast ast.Expression to ast.Identifier. got = %q", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value = '%q', expected = '%q'", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() = '%s', expected = '%s'", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	boolean, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("Could not downcast ast.Expression to ast.Boolean. got = %q", expr)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value = '%t', expected = '%t'", boolean.Value, value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() = '%s', expected = '%s'", boolean.TokenLiteral(), fmt.Sprintf("%t", value))
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("type of expr not handled. got = %T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Could not downcast ast.Expression to ast.InfixExpression. got = %q", expr)
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("opExpr.Operator = '%s'. expected = '%s'", opExpr.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}
	return true
}

func TestInfixExpression(t *testing.T) {
	// for now, tailored for int literals as left and right expressions.
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
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

		testInfixExpression(t, exprStmt.Expression, tt.leftValue, tt.operator, tt.rightValue)

		//infixExpr, ok := exprStmt.Expression.(*ast.InfixExpression)
		//if !ok {
		//t.Fatalf("Could not downcast ast.Expression to ast.InfixExpression. got = %q", exprStmt.Expression)
		//}

		//if !testIntegerLiteral(t, infixExpr.Left, tt.leftValue) {
		//return
		//}

		//if infixExpr.Operator != tt.operator {
		//t.Fatalf("infixExpr.Operator is not '%s'. got = '%s'", tt.operator, infixExpr.Operator)
		//}

		//if !testIntegerLiteral(t, infixExpr.Right, tt.rightValue) {
		//return
		//}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a)*b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a+b)+c)",
		},
		{
			"a + b + c",
			"((a+b)+c)",
		},
		{
			"a + b - c",
			"((a+b)-c)",
		},
		{
			"a * b * c",
			"((a*b)*c)",
		},
		{
			"a * b / c",
			"((a*b)/c)",
		},
		{
			"a + b / c",
			"(a+(b/c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a+(b*c))+(d/e))-f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3+4)((-5)*5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5>4)==(3<4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5<4)!=(3>4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3+(4*5))==((3*1)+(4*5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3>5)==false)",
		},
		{
			"3 < 5 == true",
			"((3<5)==true)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected='%q', got='%q'", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
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

		boolean, ok := exprStmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("Could not downcast ast.Expression to ast.Boolean. got = %q", exprStmt.Expression)
		}

		if boolean.Value != tt.expected {
			t.Fatalf("boolean.Value is not '%t'. got = '%t'", tt.expected, boolean.Value)
		}

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
