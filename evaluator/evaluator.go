package evaluator

import (
	"compiler-book/ast"
	"compiler-book/object"
	"fmt"
)

var (
	// TRUE represents the true value.
	TRUE = &object.Boolean{Value: true}
	// FALSE represents the false value.
	FALSE = &object.Boolean{Value: false}
	// NULL represents the null value.
	NULL = &object.Null{}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.RuneLiteral:
		return &object.Rune{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.PostfixExpression:
		ident := node.TokenLiteral()

		return evalPostfixExpression(node.Operator, env, ident)
	case *ast.AssignExpression:
		return evalAssignExpression(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		// store the value in the environment
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		// quote is a special form, so we handle it here
		if node.Function.TokenLiteral() == "quote" {
			return quote(node.Arguments[0], env)
		}

		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	}
	return NULL
}

func evalAssignExpression(exp *ast.AssignExpression, env *object.Environment) object.Object {
	val := Eval(exp.Value, env)
	if isError(val) {
		return val
	}

	switch left := exp.Left.(type) {
	case *ast.Identifier:
		_, ok := env.SetOnFound(left.Value, val) // this sets in the first found scope else returns error
		if !ok {
			return newError("identifier not found: " + left.Value)
		}
		return NULL
	case *ast.IndexExpression:
		structure := Eval(left.Left, env)
		if isError(structure) {
			return structure
		}

		index := Eval(left.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexAssignExpression(structure, index, val)
	}

	return val
}

func evalIndexAssignExpression(structure object.Object, index, val object.Object) object.Object {
	switch structure := structure.(type) {
	case *object.Array:
		idx, ok := index.(*object.Integer)
		if !ok {
			return newError("index must be an integer")
		}

		if idx.Value < 0 || idx.Value > int64(len(structure.Elements)-1) {
			return newError("index out of bounds")
		}

		structure.Elements[idx.Value] = val
		return NULL
	case *object.Hash:
		key, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", val.Type())
		}

		hashed := key.HashKey()
		structure.Pairs[hashed] = object.HashPair{Key: index, Value: val}
		return NULL
	default:
		return newError("index operator not supported: %s", structure.Type())
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	}

	return newError("not a function: %s", fn.Type())
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)

	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	// if index is -n, return the nth element from the end
	if idx < 0 {
		idx = max + idx + 1
	}

	if idx < 0 || idx > max {
		return newError("index out of range: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalIndexExpression(left, index object.Object) object.Object {
	structure := left.Type()
	switch {
	case structure == object.ARRAY:
		return evalArrayIndexExpression(left, index)
	case structure == object.HASH:
		return evalHashIndexExpression(left, index)
	}

	return newError("index operator not supported: %s", structure)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalPostfixExpression(operator string, env *object.Environment, ident string) object.Object {
	val, ok := env.Get(ident)
	if !ok {
		return newError("identifier not found: " + ident)
	}

	switch operator {
	case "++":
		switch val.Type() {
		case object.INTEGER:
			val := val.(*object.Integer)
			val.Value++
		case object.FLOAT:
			val := val.(*object.Float)
			val.Value++
		default:
			return newError("unknown operator: %s%s", operator, val.Type())
		}
	case "--":
		switch val.Type() {
		case object.INTEGER:
			val := val.(*object.Integer)
			val.Value--
		case object.FLOAT:
			val := val.(*object.Float)
			val.Value--
		default:
			return newError("unknown operator: %s%s", operator, val.Type())
		}
	default:
		return newError("unknown operator: %s%s", operator, val.Type())
	}

	env.Set(ident, val)
	return val
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	var bodyResult object.Object

	enclosedEnv := object.NewEnclosedEnvironment(env)

	startResult := Eval(fe.Init, enclosedEnv)
	if isError(startResult) {
		return startResult
	}

	bodyResult = NULL

	for isTruthy(Eval(fe.Condition, enclosedEnv)) {
		bodyResult = Eval(fe.Body, enclosedEnv)
		if isError(bodyResult) {
			return bodyResult
		}

		if bodyResult.Type() == object.RETURN_VALUE {
			return bodyResult
		}

		postResult := Eval(fe.Post, enclosedEnv)
		if isError(postResult) {
			return postResult
		}
	}

	return bodyResult
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	enclosedEnv := object.NewEnclosedEnvironment(env)

	condition := Eval(ie.Condition, enclosedEnv)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, enclosedEnv)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, enclosedEnv)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT && right.Type() == object.FLOAT:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.RUNE && right.Type() == object.RUNE:
		return evalRuneInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	}
	return NULL
}

func evalRuneInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Rune).Value
	rightVal := right.(*object.Rune).Value

	switch operator {
	case "+":
		return &object.Rune{Value: leftVal + rightVal}
	case "-":
		return &object.Rune{Value: leftVal - rightVal}
	}
	return NULL
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
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
	}
	return FALSE
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if result != nil && result.Type() == object.RETURN_VALUE || result.Type() == object.ERROR {
			return result
		}
	}

	return result
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
