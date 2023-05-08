package evaluator

import (
	"bytes"
	"compiler-book/object"
	"fmt"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: btLen,
	},
	"print": {
		Fn: btPrint,
	},
	"printf": {
		Fn: btPrintf,
	},
	"push": {
		Fn: btPush,
	},
	"pop": {
		Fn: btPop,
	},
	"first": {
		Fn: btFirst,
	},
	"rest": {
		Fn: btRest,
	},
}

func btLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}
}

func btPrint(args ...object.Object) object.Object {
	var out bytes.Buffer

	for _, arg := range args {
		out.WriteString(arg.Inspect())
		out.WriteString(" ")
	}

	fmt.Println(out.String())

	return NULL
}

func btPrintf(args ...object.Object) object.Object {
	var opts []any

	if len(args) < 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	format := args[0]

	if format.Type() != object.STRING {
		return newError("argument to `printf` not supported, got %s",
			format.Type())
	}

	formatValue := format.(*object.String).Value

	for _, arg := range args[1:] {
		switch arg.Type() {
		case object.INTEGER:
			opts = append(opts, arg.(*object.Integer).Value)
		case object.FLOAT:
			opts = append(opts, arg.(*object.Float).Value)
		case object.STRING:
			opts = append(opts, arg.(*object.String).Value)
		case object.RUNE:
			opts = append(opts, arg.(*object.Rune).Value)
		case object.BOOLEAN:
			opts = append(opts, arg.(*object.Boolean).Value)
		default:
			return newError("argument to `printf` not supported, got %s",
				arg.Type())
		}
	}

	unescape := strings.Replace(formatValue, "\\n", "\n", -1) // FIXME: improve this
	fmt.Printf(unescape+"\n", opts...)

	return NULL
}

func btPush(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if args[0].Type() != object.ARRAY {
		return newError("argument to `push` not supported, got %s",
			args[0].Type())
	}

	if args[1] == args[0] {
		return newError("argument to `push` cannot be the same array")
	}

	array := args[0].(*object.Array)
	array.Elements = append(array.Elements, args[1])

	return NULL
}

func btPop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY {
		return newError("argument to `pop` not supported, got %s",
			args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	if length == 0 {
		return NULL
	}

	last := array.Elements[length-1]
	array.Elements = array.Elements[:length-1]

	return last
}

func btFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY {
		return newError("argument to `first` not supported, got %s",
			args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	if length == 0 {
		return NULL
	}

	return array.Elements[0]
}

func btRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY {
		return newError("argument to `rest` not supported, got %s",
			args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	if length == 0 {
		return NULL
	}

	newElements := make([]object.Object, length-1)
	copy(newElements, array.Elements[1:length])

	return &object.Array{Elements: newElements}
}
