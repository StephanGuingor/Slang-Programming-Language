package repl

import (
	"bufio"
	"compiler-book/ast"
	"compiler-book/evaluator"
	"compiler-book/lexer"
	"compiler-book/object"
	"compiler-book/parser"
	"fmt"
	"io"
	"os"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		if expanded == nil {
			continue
		}

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			color := yellow
			if evaluated.Type() == object.NULL {
				color = gray
			}

			io.WriteString(out, formatColor(color, evaluated.Inspect()))
			io.WriteString(out, "\n")
		}
	}
}

func StartFile(filename string) {
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()
	program := parseFile(filename)

	if program == nil {
		return
	}

	evaluator.DefineMacros(program, macroEnv)
	expanded := evaluator.ExpandMacros(program, macroEnv)

	if expanded == nil {
		return
	}

	err := evaluator.Eval(expanded, env)
	if isError(err) {
		fmt.Println(err.Inspect())
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func parseFile(filename string) *ast.Program {
	// read file to string

	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	l := lexer.New(string(file))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return nil
	}

	return program
}

func printParserErrors(out io.Writer, errors []*parser.ParseError) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err)
	}
}

type color string

const (
	gray   color = "30"
	red    color = "31"
	green  color = "32"
	yellow color = "33"
	blue   color = "34"
)

func formatColor(color color, str string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", color, str)
}
