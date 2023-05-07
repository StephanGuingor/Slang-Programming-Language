package repl

import (
	"bufio"
	"compiler-book/evaluator"
	"compiler-book/lexer"
	"compiler-book/object"
	"compiler-book/parser"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

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

		evaluated := evaluator.Eval(program, env)
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
