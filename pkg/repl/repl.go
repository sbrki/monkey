package repl

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/sbrki/monkey/pkg/evaluator"
	"github.com/sbrki/monkey/pkg/lexer"
	"github.com/sbrki/monkey/pkg/object"
	"github.com/sbrki/monkey/pkg/parser"
)

const (
	PROMPT = ">>"
)

func Start(in io.ReadCloser, out io.Writer) {
	term, err := readline.NewEx(
		&readline.Config{
			Prompt: PROMPT,
			Stdin:  in,
			Stdout: out,
		},
	)
	if err != nil {
		panic(err)
	}
	defer term.Close()

	env := object.NewEnvironment()
	for {
		line, err := term.Readline()
		if err != nil { // EOF
			break
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
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
