package evaluator

import (
	"compiler-book/ast"
	"compiler-book/object"
	"compiler-book/token"
	"fmt"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpression.Function.TokenLiteral() == "unquote"
}

// could lead to an inconstent AST
func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Float:
		t := token.Token{
			Type:    token.FLOAT,
			Literal: fmt.Sprintf("%f", obj.Value),
		}
		return &ast.FloatLiteral{Token: t, Value: obj.Value}
	case *object.String:
		t := token.Token{
			Type:    token.STRING,
			Literal: obj.Value,
		}
		return &ast.StringLiteral{Token: t, Value: obj.Value}
	case *object.Rune:
		t := token.Token{
			Type:    token.RUNE,
			Literal: fmt.Sprintf("%c", obj.Value),
		}
		return &ast.RuneLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		if obj.Value {
			t := token.Token{
				Type:    token.TRUE,
				Literal: fmt.Sprintf("%t", obj.Value),
			}
			return &ast.Boolean{Token: t, Value: obj.Value}
		} else {
			t := token.Token{
				Type:    token.FALSE,
				Literal: fmt.Sprintf("%t", obj.Value),
			}
			return &ast.Boolean{Token: t, Value: obj.Value}
		}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
