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

// launchConsole loops for standard input, does some basic filtering
func launchConsole() {
	reader := bufio.NewReader(os.Stdin)
	partialStatementBuffer := make([]string, 0, 1000)

	for {
		fmt.Print("🔥 ")
		line, _ := reader.ReadString('\n')

		if strings.EqualFold(".EXIT", line) == true {
			fmt.Println("All done.")
			os.Exit(0)
		}

		line = filterComment(line)
		line = filterNewline(line)
		line = filterReturn(line)

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

// processLine preforms all main operations
func processLine(line string) {
	statements := make([]tokenizer.Statement, 1000)

	statement := tokenizer.TokenizeStatement(line)

	statements = append(statements, statement)

	statement = parser.ParseStatement(statement)

	generator.Generate(statement)

	if *DebugPtr {
		tokenizer.PrintStatement(statement)
	}
}

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
