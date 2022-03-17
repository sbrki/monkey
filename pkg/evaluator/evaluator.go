package evaluator

import (
	"github.com/sbrki/monkey/pkg/ast"
	"github.com/sbrki/monkey/pkg/object"
)

var (
	NULL = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

// evalStatements evaluates all statements in the passed list.
// It returns the evaluated result of the last statement of the list.
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(in bool) *object.Boolean {
	if in {
		return TRUE
	}
	return FALSE
}
