package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/MohamTahaB/interpreter-go/lexer"
	"github.com/MohamTahaB/interpreter-go/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// Exit REPL
		if strings.TrimSpace(line) == "exit" {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, errorMsg := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", errorMsg))
	}
}
