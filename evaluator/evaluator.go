package evaluator

import (
	"monk/ast"
	"monk/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(left, node.Operator, right)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return NULL
	}
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftval := left.(*object.Integer).Value
	rightval := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftval + rightval}
	case "-":
		return &object.Integer{Value: leftval - rightval}
	case "*":
		return &object.Integer{Value: leftval * rightval}
	case "/":
		return &object.Integer{Value: leftval / rightval}

	case ">":
		return nativeBoolToBooleanObject(leftval > rightval)
	case "<":
		return nativeBoolToBooleanObject(leftval < rightval)
	case "==":
		return nativeBoolToBooleanObject(leftval == rightval)
	case "!=":
		return nativeBoolToBooleanObject(leftval != rightval)

	default:
		return NULL
	}
}