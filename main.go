// Package main ...
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sqlit/generator"
	"sqlit/parser"
	"sqlit/tokenizer"
	"strings"
)

var scriptPtr *string

var cleanPtr *bool

// DebugPtr toggles global printing of debugging statements
var DebugPtr *bool

func main() {

	scriptPtr = flag.String("script", "", "run a SQL script from file in dir sqlit/")

	cleanPtr = flag.Bool("clean", false, "deletes all previously created databases in sqlit/tmp")

	DebugPtr = flag.Bool("debug", false, "displays debugging info")

	flag.Parse()

	if *cleanPtr {
		removeContents("tmp")
	}

	if *scriptPtr != "" {
		SQLFilename := *scriptPtr
		SQLFileContents, err := ioutil.ReadFile(SQLFilename)
		check(err)

		runSQLScript(SQLFileContents)
	} else {
		launchConsole()
	}

}

// launchConsole loops for standard input, does some basic
// filtering to allow for multi line statements
func launchConsole() {
	reader := bufio.NewReader(os.Stdin)
	partialStatementBuffer := make([]string, 0, 1000)

	fmt.Println("🔥")

	for {
		fmt.Print("sqlit> ")

		line, _ := reader.ReadString('\n')

		line = filterComment(line)
		line = filterNewline(line)
		line = filterReturn(line)

		if strings.EqualFold(".EXIT", line) == true {
			fmt.Println("All done.")
			os.Exit(0)
		}

		if len(line) <= 1 {
			continue
		}

		if includesDelimiter(line) == false {
			partialStatementBuffer = append(partialStatementBuffer, line)
		} else {
			line = filterDelimiter(line)

			statement := line + strings.Join(partialStatementBuffer, " ")

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

	// Generate a set of assertions and a set of executions for our query
	operation := generator.Generate(statement)

	// Make sure our query is valid before we request resources
	err := operation.Assert()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Finally, execute our query
	success, err := operation.Invoke()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(success)
	}

	if *DebugPtr {
		tokenizer.PrintStatement(statement)
	}
}

//
//			Helper functions below
//

func filterComment(line string) string {
	if strings.HasPrefix(line, "--") {
		return ""
	}

	return line
}

func filterReturn(line string) string {
	return strings.Replace(line, "\r", "", -1)
}

func filterNewline(line string) string {
	return strings.Replace(line, "\n", "", -1)
}

func filterDelimiter(line string) string {
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

// runSQLScript loads a SQL script
// TODO: this needs to be refactored and might not work
func runSQLScript(fileContents []byte) {

	// Make string array from fileData by line
	linesCommentsRemoved := make([]string, 0, 1000)
	lines := strings.Split(string(fileContents), "\n")
	for _, line := range lines {

		// Strip comment lines
		if strings.HasPrefix(line, "--") == false {

			// Strip empty lines
			if len(line) > 1 {
				linesCommentsRemoved = append(linesCommentsRemoved, line)
			}
		}
	}

	// remove carridge returns from lines
	linesCarridgeReturnsRemoved := make([]string, 0, 1000)
	for _, line := range linesCommentsRemoved {
		// Strip empty lines
		if line[len(line)-1] == 13 {
			line = strings.Replace(line, "\r", "", -1)
			linesCarridgeReturnsRemoved = append(linesCarridgeReturnsRemoved, line)
		}
	}

	// join lines that aren't delimited
	linesBreaksRemoved := make([]string, 0, 1000)
	i := 0
	for _, line := range linesCarridgeReturnsRemoved {

		fmt.Println(line)

		if strings.HasSuffix(line, ";") == true {
			line = strings.Replace(line, ";", "", -1)
			linesBreaksRemoved = append(linesBreaksRemoved, line)
		} else if strings.HasSuffix(line, ";") == false {
			if strings.EqualFold(line, ".exit") == false {
				a := []string{line, linesCarridgeReturnsRemoved[i+1]}
				linesCarridgeReturnsRemoved[i+1] = strings.Join(a, "")
			}
		}
		i++
	}

	fmt.Println("")
	for _, line := range linesBreaksRemoved {
		processLine(line)
	}

	fmt.Println("All done.")

	launchConsole()
}
