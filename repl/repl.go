package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/MohamTahaB/interpreter-go/lexer"
	"github.com/MohamTahaB/interpreter-go/token"
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
		if strings.TrimSpace(line) == "exit" {
			return
		}
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%v\n", tok)
		}
	}
}
