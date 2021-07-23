package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// Have our interpreter had an error?
var hadError bool

func main() {
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

	for _, t := range tokens {
		fmt.Println(t)
	}
}

func loxerr(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "%s on line %d: %s", where, line, message)
	hadError = true
}
