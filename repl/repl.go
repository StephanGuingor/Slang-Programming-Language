package repl

import (
	"bufio"
	"compiler-book/lexer"
	"compiler-book/parser"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		fmt.Fprintln(out, program.String())
	}
}

func printParserErrors(out io.Writer, errors []*parser.ParseError) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err)
	}
}
