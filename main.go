package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
)

//go:generate go run acceptgen/gen.go structs visiters
//go:generate gofmt -w visiters.go

// Have our interpreter had an error?
var (
	interpreter     = NewInterpreter()
	hadError        bool
	hadRuntimeError bool
)

func main() {
	f, err := os.Create("cpu.profile")
	if err != nil {
		log.Fatal(err)
	}
	if pprof.StartCPUProfile(f) != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// apmain()
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, `Usage: jlox <script>`)
		os.Exit(64)
	} else if len(args) == 1 {
		runfile(os.Args[0])
	} else {
		runprompt()
	}
}

func runfile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	run(bs)
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runprompt() {
	rr := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := rr.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		run(line)
		hadError = false
	}
}

func run(source []byte) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)

	stmts, err := parser.Parse()
	if err != nil {
		loxerr2(err)
	}
	interpreter.Interpret(stmts)
}

func loxerr(line int, message string) {
	report(line, "", message)
}

func loxparseerr(tok Token, message string) {
	if tok.Type == tokenEOF {
		report(tok.Line, " at end", message)
	} else {
		report(tok.Line, " at `"+tok.String()+"`", message)
	}
}

func loxerr2(e *Error) {
	fmt.Printf("at line %d: %s\n", e.Token.Line, e.Message)
	hadRuntimeError = true
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "%s on line %d: %s\n", where, line, message)
	hadError = true
}
