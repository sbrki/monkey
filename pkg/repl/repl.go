package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sbrki/monkey/pkg/evaluator"
	"github.com/sbrki/monkey/pkg/lexer"
	"github.com/sbrki/monkey/pkg/parser"
)

const (
	PROMPT = ">>"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		if !scanner.Scan() {
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

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Whoops! We ran into some monkey business when parsing the input!\n")
	for idx, msg := range errors {
		io.WriteString(out, fmt.Sprintf("[%d]\t%s\n", idx+1, msg))
	}
}
