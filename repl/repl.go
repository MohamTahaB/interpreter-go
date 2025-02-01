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

const ERROR_HEADER = `
                             ud$$$**$$$$$$$bc.
                          u@**"        4$$$$$$$Nu
                        J                ""#$$$$$$r
                       @                       $$$$b
                     .F                        ^*3$$$
                    :% 4                         J$$$N
                    $  :F                       :$$$$$
                   4F  9                       J$$$$$$$
                   4$   k             4$$$$bed$$$$$$$$$
                   $$r  'F            $$$$$$$$$$$$$$$$$r
                   $$$   b.           $$$$$$$$$$$$$$$$$N
                   $$$$$k 3eeed$$b    $$$Euec."$$$$$$$$$
    .@$**N.        $$$$$" $$$$$$F'L $$$$$$$$$$$  $$$$$$$
    :$$L  'L       $$$$$ 4$$$$$$  * $$$$$$$$$$F  $$$$$$F         edNc
   @$$$$N  ^k      $$$$$  3$$$$*%   $F4$$$$$$$   $$$$$"        d"  z$N
   $$$$$$   ^k     '$$$"   #$$$F   .$  $$$$$c.u@$$$          J"  @$$$$r
   $$$$$$$b   *u    ^$L            $$  $$$$$$$$$$$$u@       $$  d$$$$$$
    ^$$$$$$.    "NL   "N. z@*     $$$  $$$$$$$$$$$$$P      $P  d$$$$$$$
       ^"*$$$$b   '*L   9$E      4$$$  d$$$$$$$$$$$"     d*   J$$$$$r
            ^$$$$u  '$.  $$$L     "#" d$$$$$$".@$$    .@$"  z$$$$*"
              ^$$$$. ^$N.3$$$       4u$$$$$$$ 4$$$  u$*" z$$$"
                '*$$$$$$$$ *$b      J$$$$$$$b u$$P $"  d$$P
                   #$$$$$$ 4$ 3*$"$*$ $"$'c@@$$$$ .u@$$$P
                     "$$$$  ""F~$ $uNr$$$^&J$$$$F $$$$#
                       "$$    "$$$bd$.$W$$$$$$$$F $$"
                         ?k         ?$$$$$$$$$$$F'*
                          9$$bL     z$$$$$$$$$$$F
                           $$$$    $$$$$$$$$$$$$
                            '#$$c  '$$$$$$$$$"
                             .@"#$$$$$$$$$$$$b
                           z*      $$$$$$$$$$$$N.
                         e"      z$$"  #$$$k  '*$$.
                     .u*      u@$P"      '#$$c   "$$c
              u@$*"""       d$$"            "$$$u  ^*$$b.
            :$F           J$P"                ^$$$c   '"$$$$$$bL
           d$$  ..      @$#                      #$$b         '#$
           9$$$$$$b   4$$                          ^$$k         '$
            "$$6""$b u$$                             '$    d$$$$$P
              '$F $$$$$"                              ^b  ^$$$$b$
               '$W$$$$"                                'b@$$$$"
                                                        ^$$$*

------------------------------------------------

ITS NOT LOOKING GOOD BRUV !!! Had some issues:
`

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
	io.WriteString(out, fmt.Sprintf("\n%s\n", ERROR_HEADER))
	for _, errorMsg := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", errorMsg))
	}
}
