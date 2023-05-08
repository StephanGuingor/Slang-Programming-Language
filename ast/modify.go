package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {

	// TODO: Handle errors
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ForExpression:
		node.Init, _ = Modify(node.Init, modifier).(Statement)
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Post, _ = Modify(node.Post, modifier).(Expression)
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *FunctionLiteral:
		for i, param := range node.Parameters {
			node.Parameters[i], _ = Modify(param, modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, elem := range node.Elements {
			node.Elements[i], _ = Modify(elem, modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, value := range node.Pairs {
			newKey, _ := Modify(key, modifier).(Expression)
			newValue, _ := Modify(value, modifier).(Expression)
			newPairs[newKey] = newValue
		}
		node.Pairs = newPairs
	case *CallExpression:
		node.Function, _ = Modify(node.Function, modifier).(Expression)
		for i, arg := range node.Arguments {
			node.Arguments[i], _ = Modify(arg, modifier).(Expression)
		}
	case *BlockStatement:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	}

	return modifier(node)
}
