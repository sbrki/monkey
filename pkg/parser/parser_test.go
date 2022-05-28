package parser

import (
	"fmt"
	"testing"

	"github.com/sbrki/monkey/pkg/ast"
	"github.com/sbrki/monkey/pkg/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) = %d, expected = 1",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
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

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) = %d, expected = 1",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		val := stmt.(*ast.ReturnStatement).ReturnValue
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"foo bar"`

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

	strLit, ok := exprStmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.StringLiteral. got = %q", exprStmt.Expression)
	}

	if strLit.Value != "foo bar" {
		t.Errorf("strLit.Value = %s, expected = \"foo bar\"", strLit.Value)
	}

	if strLit.TokenLiteral() != "foo bar" {
		t.Errorf("strLit.TokenLiteral() = '%s', exptected = \"foo bar\"", strLit.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		if !testLiteralExpression(t, prefixExpr.Right, tt.value) {
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
		{
			"1 + (2 + 3) + 4",
			"((1+(2+3))+4)",
		},
		{
			"(5 + 5) * 2",
			"((5+5)*2)",
		},
		{
			"2 / (5 + 5)",
			"(2/(5+5))",
		},
		{
			"-(5 + 5)",
			"(-(5+5))",
		},
		{
			"!(true == true)",
			"(!(true==true))",
		},
		{
			"a + add(b * c) + d",
			"((a+add((b*c)))+d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a,b,1,(2*3),(4+5),add(6,(7*8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a+b)+((c*d)/f))+g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a*([1,2,3,4][(b*c)]))*d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1,2][1])",
			"add((a*(b[2])),(b[1]),(2*([1,2][1])))",
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

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

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.IfExpression. got = %q", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("len(ifExpr.Consequence.Statements) = %d, expected = 1",
			len(ifExpr.Consequence.Statements))
	}

	consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if ifExpr.Alternative != nil {
		t.Errorf("ifExpr.Alternative is not nil, got = %+v", ifExpr.Alternative)
	}
}

func TestIfExpressionWithAlternative(t *testing.T) {
	input := `if (x < y) { x } else { y }`

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

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.IfExpression. got = %q", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	// consequence

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("len(ifExpr.Consequence.Statements) = %d, expected = 1",
			len(ifExpr.Consequence.Statements))
	}

	consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// alternative

	if ifExpr.Alternative == nil {
		t.Error("ifExpr.Alternative is nil")
	}

	if len(ifExpr.Alternative.Statements) != 1 {
		t.Fatalf("len(ifExpr.Alternative.Statements) = %d, expected = 1",
			len(ifExpr.Alternative.Statements))
	}

	alternative, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", ifExpr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x,y) { x + y; }`

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

	functionLiteral, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.FunctionLiteral. got = %q", exprStmt.Expression)
	}

	if len(functionLiteral.Parameters) != 2 {
		t.Fatalf("len(functionLiteral.Parameters) = %d, want = 2", len(functionLiteral.Parameters))
	}

	testLiteralExpression(t, functionLiteral.Parameters[0], "x")
	testLiteralExpression(t, functionLiteral.Parameters[1], "y")

	if len(functionLiteral.Body.Statements) != 1 {
		t.Fatalf("len(functionLiteral.Body.Statements = %d, want = 1", len(functionLiteral.Body.Statements))
	}

	bodyStmt, ok := functionLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Could not downcast ast.Statement to ast.ExpressionStatement, got = %q", functionLiteral.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5, foo)`

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

	callExpr, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.CallExpression. got = %q", exprStmt.Expression)
	}

	if !testIdentifier(t, callExpr.Function, "add") {
		return
	}

	if len(callExpr.Arguments) != 4 {
		t.Fatalf("len(callExpr.Arguments) = %d, want = 4", len(callExpr.Arguments))
	}

	testLiteralExpression(t, callExpr.Arguments[0], 1)
	testInfixExpression(t, callExpr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpr.Arguments[2], 4, "+", 5)
	testLiteralExpression(t, callExpr.Arguments[3], "foo")
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

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

	arrayLit, ok := exprStmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.ArrayLiteral. got = %q", exprStmt.Expression)
	}

	if len(arrayLit.Elements) != 3 {
		t.Fatalf("len(arrayLit.Elements) != 3, got = %d", len(arrayLit.Elements))
	}

	testIntegerLiteral(t, arrayLit.Elements[0], 1)
	testInfixExpression(t, arrayLit.Elements[1], 2, "*", 2)
	testInfixExpression(t, arrayLit.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

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

	indexExpr, ok := exprStmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.IndexExpression. got = %q", exprStmt.Expression)
	}

	if !testIdentifier(t, indexExpr.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExpr.Index, 1, "+", 1) {
		return
	}

}

func TestParsingHashLiteralStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

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

	hashLiteral, ok := exprStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.HashLiteral. got = %q", exprStmt.Expression)
	}

	if len(hashLiteral.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got=%d, want=3", len(hashLiteral.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got=%T", key)
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

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

	hashLiteral, ok := exprStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.HashLiteral. got = %q", exprStmt.Expression)
	}

	if len(hashLiteral.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length, got=%d, want=0", len(hashLiteral.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15/5 }`

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

	hashLiteral, ok := exprStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Could not downcast ast.Expression to ast.HashLiteral. got = %q", exprStmt.Expression)
	}

	if len(hashLiteral.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got=%d, want=3", len(hashLiteral.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
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
