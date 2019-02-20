package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	// "log"
	"os"
	"strings"
	// "io"
	// // "go/scanner"
	// "os/exec"
	"path/filepath"
	"sqlit/generator"
	"sqlit/io"
	"sqlit/parser"
	"sqlit/tokenizer"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	removeContents("tmp")

	if len(os.Args) > 1 {
		// SQLFilename := os.Args[2]
		// fileContents, err := ioutil.ReadFile("tmp/" + SQLFilename)
		// check(err)

		// runSQLScript()
		launchConsole()
	} else {
		// launchConsole()
	}

	// fileContents, err := ioutil.ReadFile("tmp/PA1_test.sql")

	// filepath := "tmp/PA2_test.sql"
	// somestring := tokenizer.GetString()

	// Slurp file's entire contents to memory
	fileContents, err := ioutil.ReadFile("PA1_test.sql")
	check(err)
	// fmt.Print(string(fileContents))

	// More control over which parts of file are read
	// f, err := os.Open(filepath)
	// check(err)

	// Read some bytes from the beginning of the file. Allow up to 5 to be read but also note how many actually were read.
	// b1 := make([]byte, 5)
	// n1, err := f.Read(b1)
	// check(err)
	// fmt.Printf("%d bytes: %s\n", n1, string(b1))

	// You can also Seek to a known location in the file and Read from there.
	// o2, err := f.Seek(6, 0)
	// check(err)
	// b2 := make([]byte, 2)
	// n2, err := f.Read(b2)
	// check(err)
	// fmt.Printf("%d bytes @ %d: %s\n", n2, o2, string(b2))

	// The io package provides some functions that may be helpful for file reading. For example, reads like the ones above can be more robustly implemented with ReadAtLeast.
	// o3, err := f.Seek(6, 0)
	// check(err)
	// b3 := make([]byte, 2)
	// n3, err := io.ReadAtLeast(f, b3, 2)
	// check(err)
	// fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	// There is no built-in rewind, but Seek(0, 0) accomplishes this.
	// _, err = f.Seek(0, 0)
	// check(err)

	// The bufio package implements a buffered reader that may be useful both for its efficiency with many small reads and because of the additional reading methods it provides.
	// r4 := bufio.NewReader(f)
	// b4, err := r4.Peek(5)
	// check(err)
	// fmt.Printf("5 bytes: %s\n", string(b4))

	// Close the file when you’re done (usually this would be scheduled immediately after Opening with defer).

	// f.Close()

	// Open the test file for reading
	// sqlFile, err := os.Open(filepath)
	// check(err)

	// sqlFileReader := bufio.NewReader(sqlFile)
	// bufStr := sqlFileReader.ReadString(' ')

	// Break file into strings
	// fmt.Print("\n\n\n\n\n\nHMM\n\n\n", string(fileContents), "\n\n\n\n ARRAY BELOW \n\n\n\n")
	// fileContentsString := string(fileContents)

	// Make string array from fileData by line
	linesCommentsRemoved := make([]string, 0, 1000)
	lines := strings.Split(string(fileContents), "\n")
	for _, line := range lines {

		// Strip comment lines
		if strings.HasPrefix(line, COMMENT) == false {

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
			if strings.EqualFold(line, EXIT) == false {
				a := []string{line, linesCarridgeReturnsRemoved[i+1]}
				linesCarridgeReturnsRemoved[i+1] = strings.Join(a, "")
			}
		}
		i++
	}

	// now that statements are sanitized, generate Tokens
	// tokenSequences := make([]string, 0, 1000)

	fmt.Println("")
	for _, line := range linesBreaksRemoved {
		onLine(line)
	}

	// launchConsole()

	// print
	// for _, tokenSeq := range tokenSequences {
	// fmt.Println(tokenSeq)
	// }

	// Bind sqlfiledata to string array
	// testArray := strings.Fields(fileContentsString)

	// Iterate through string array and print
	// for _, v := range testArray {
	// fmt.Println(v)
	// }

	fmt.Println("...done.")
	os.Exit(0)
}

func onLine(line string) {
	statements := make([]tokenizer.Statement, 1000)

	statement := tokenizer.TokenizeStatement(line)
	statements = append(statements, statement)
	statement = parser.ParseStatement(statement)
	function := generator.Generate(statement)
	io.Execute()
	printFunction(function)
	// printStatement(statement)
}

func printStatement(statement tokenizer.Statement) {
	fmt.Print("type	", statement.Type, "\n")
	for _, token := range statement.Tokens {
		tokenizer.PrintToken(token)
		fmt.Print("\n")
	}

	fmt.Print("\n")
}

func printFunction(function generator.Function) {
	// fmt.Print("\n")
}

func launchConsole() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("🔥 ")
		line, _ := reader.ReadString('\n')
		// convert CRLF to LF
		line = strings.Replace(line, "\n", "", -1)

		if len(line) == 0 {
			// cmd := exec.Command("go run main.go", "1")
			// log.Printf("Running command and waiting for it to finish...")
			// err := cmd.Run()
			// log.Printf("Command finished with error: %v", err)
			os.Exit(1)
		}

		onLine(line)

		if strings.Compare("hi", line) == 0 {
			fmt.Println("hello, Yourself")
		}

	}
}

func runSQLScript() {

}

// debugging helper to clear a dir - https://stackoverflow.com/a/33451503/7977208
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

// SEMICOLON delimiter
const SEMICOLON = ";"

// COMMENT delimiter
const COMMENT = "--"

// EXIT delimiter
const EXIT = ".exit"
