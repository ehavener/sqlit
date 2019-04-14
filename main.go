/* UNR CS 457 | SPRING 2019 | emerson@nevada.unr.edu */

// Package main is the program's core, providing a REPL
// for query input (launchConsole()) as well the a main
// loop which transforms queries down to disk operations (processLine()).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sqlit/generator"
	"sqlit/parser"
	"sqlit/tokenizer"
	"strings"
)

var cleanPtr *bool

// DebugPtr toggles global printing of debugging statements
var DebugPtr *bool

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func main() {

	createTmpDirectory()

	cleanPtr = flag.Bool("clean", true, "deletes all previously created databases in sqlit/tmp")

	DebugPtr = flag.Bool("debug", false, "displays debugging info")

	flag.Parse()

	if *cleanPtr {
		removeContents("tmp")
	}

	launchConsole()
}

// launchConsole loops for standard input, does some basic
// filtering to allow for multi line statements
func launchConsole() {
	reader := bufio.NewReader(os.Stdin)
	partialStatementBuffer := make([]string, 0, 1000)

	lastLineWasEmpty := false

	fmt.Println("🔥")

	for {
		if lastLineWasEmpty == false {
			fmt.Printf(WarningColor, "sqlit> ")
		}

		line, _ := reader.ReadString('\n')

		line = removeComment(line)
		line = removeNewline(line)
		line = removeReturn(line)

		if strings.EqualFold(".EXIT", line) == true {
			fmt.Println("All done.")
			os.Exit(0)
		}

		if len(line) <= 1 {
			lastLineWasEmpty = true
			continue
		} else {
			lastLineWasEmpty = false
		}

		if includesDelimiter(line) == false {
			partialStatementBuffer = append(partialStatementBuffer, line+" ")
		} else {
			line = removeDelimiter(line)

			statement := strings.Join(partialStatementBuffer, " ") + line

			if len(statement) > 1 {
				processLine(statement)
				partialStatementBuffer = partialStatementBuffer[:0]
			}
		}
	}
}

// processLine goes through all the main functionality by transforming input into operations
func processLine(line string) {

	// Break our line of input up into tokens
	statement := tokenizer.TokenizeStatement(line)

	// Give them some syntactical meaning
	statement = parser.ParseStatement(statement)

	if *DebugPtr {
		tokenizer.PrintStatement(statement)
	}

	// Generate a function of assertions and a function of operations for our query
	operation := generator.Generate(statement)

	// Make sure our query is valid before we request resources
	err := operation.Assert()
	if err != nil {
		fmt.Printf(ErrorColor, err)
		return
	}

	// Finally, execute our query
	success, err := operation.Invoke()
	if err != nil {
		fmt.Printf(ErrorColor, err)
	} else {
		fmt.Printf(DebugColor, success)
		fmt.Println()
	}
}

//
//			Helper functions
//

func createTmpDirectory() {
	_, err := os.Stat("tmp/")
	if os.IsNotExist(err) {
		err := os.Mkdir("tmp/", os.ModePerm)
		check(err)
	}
}

func removeComment(line string) string {
	if strings.HasPrefix(line, "--") {
		return ""
	}

	return line
}

func removeReturn(line string) string {
	return strings.Replace(line, "\r", "", -1)
}

func removeNewline(line string) string {
	return strings.Replace(line, "\n", "", -1)
}

func removeDelimiter(line string) string {
	return strings.Replace(line, ";", "", -1)
}

func includesDelimiter(line string) bool {
	return strings.HasSuffix(line, ";")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// removeContents deletes all files in a directory
// source: https://stackoverflow.com/a/33451503/7977208
func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
