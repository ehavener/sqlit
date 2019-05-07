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
	"sqlit/diskio"
	"sqlit/generator"
	"sqlit/parser"
	"sqlit/tokenizer"
	"strings"
)

var cleanPtr *bool

// DebugPtr toggles global printing of debugging statements
var DebugPtr *bool

var inTransactionMode bool

var transactionStack []generator.Operation

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func main() {
	inTransactionMode = false

	createTmpDirectory()

	cleanPtr = flag.Bool("clean", false, "deletes all previously created databases in sqlit/tmp")

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

	lineNumber := 0

	for {
		lineNumber++
		if lineNumber > 31 {
			fmt.Println("All done.")
			os.Exit(0)
		}

		if lastLineWasEmpty == false {
			fmt.Printf(WarningColor, "sqlit> ")
		}

		line, _ := reader.ReadString('\n')

		line = removeComment(line)
		line = removeNewline(line)
		line = removeReturn(line)

		// fmt.Println("line:" + " " + line)

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

	if *DebugPtr {
		tokenizer.PrintStatement(statement)
	}

	// Give them some syntactical meaning
	statement = parser.ParseStatement(statement)

	// interpert transaction mode entry, break if entering
	if statement.Type == "BEGIN" {
		inTransactionMode = true
		fmt.Printf(DebugColor, "Transaction starts.")
		fmt.Println()
		return
	}

	// interpert transaction commit
	if statement.Type == "COMMIT" {

		// 1. assert that all operations in the transaction stack are valid
		//	(including table lock checks)
		//	return on error
		for _, operation := range transactionStack {
			err := operation.Assert()

			if err != nil {
				fmt.Printf(ErrorColor, err)
				fmt.Println()
				fmt.Printf(ErrorColor, "Transaction abort.")
				fmt.Println()

				// empty out our transaction stack
				transactionStack = []generator.Operation{}

				// unlock all associated resources
				diskio.UnlockTable("flights")

				return
			}
		}

		// 2. invoke all operations in transaction
		for _, operation := range transactionStack {
			success, err := operation.Invoke()
			if err != nil {
				fmt.Printf(ErrorColor, err)
				fmt.Println()
			} else {
				fmt.Printf(DebugColor, success)
				fmt.Println()
			}
		}

		// empty out our transaction stack
		transactionStack = []generator.Operation{}

		// unlock all associated resources
		diskio.UnlockTable("flights")

		return
	}

	if *DebugPtr {
		tokenizer.PrintStatement(statement)
	}

	// Generate a function of assertions and a function of operations for our query
	operation := generator.Generate(statement, inTransactionMode)

	// Make sure our query is valid before we request resources
	err := operation.Assert()
	if err != nil {
		fmt.Printf(ErrorColor, err)
		fmt.Println()

		if inTransactionMode {
			fmt.Printf(ErrorColor, "Transaction abort.")
			fmt.Println()
		}

		return
	}

	// if we're in transaction mode, and assertions pass, we store the operation on the transaction stack rather then executing it immediately
	if inTransactionMode {
		transactionStack = append(transactionStack, operation)
		return
	}

	// Finally, execute our query
	success, err := operation.Invoke()
	if err != nil {
		fmt.Printf(ErrorColor, err)
		fmt.Println()
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

	if strings.Contains(line, "--") {
		splitByComment := strings.Split(line, "--")
		//	fmt.Println("hmm")
		//	fmt.Println(splitByComment)
		// fmt.Println(splitByComment[0])
		return splitByComment[0]
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
