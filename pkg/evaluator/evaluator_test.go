package evaluator

import (
	"testing"

	"github.com/sbrki/monkey/pkg/lexer"
	"github.com/sbrki/monkey/pkg/object"
	"github.com/sbrki/monkey/pkg/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Could not downcast object.Object to object.Integer, got = %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Integer.Value = %d, expected = %d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Could not downcast object.Object to object.Boolean, got = %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Boolean.Value = %t, expected = %t", result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}
