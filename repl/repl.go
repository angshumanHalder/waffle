package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// env := object.NewEnvironment()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GLOBALSSIZE)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprintf(out, "%s", PROMPT)
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

		/* evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		} */

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing Bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := machine.LastPopppedStackElem()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

const WAFFLE = `            
                                  ad88    ad88 88             
                                d8"     d8"   88             
                                88      88    88             
8b      db      d8 ,adPPYYba, MM88MMM MM88MMM 88  ,adPPYba,
 8b    d88b    d8' ""      Y8   88      88    88 a8P_____88  
  8b  d8' 8b  d8'  ,adPPPPP88   88      88    88 8PP""""""" 
   8bd8'   8bd8'   88,    ,88   88      88    88 "8b,   ,aa
    YP      YP      "8bbdP"Y8   88      88    88   "Ybbd8"'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, WAFFLE)
	io.WriteString(out, "Woops! We ran into some sticky business here!\n\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
