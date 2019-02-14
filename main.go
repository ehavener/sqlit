package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sqlit/tokenizer"
	"strings"
	// "bufio"
	// "io"
	// "go/scanner"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	filepath := "tmp/PA1_test.sql"
	// somestring := tokenizer.GetString()

	// Slurp file's entire contents to memory
	fileContents, err := ioutil.ReadFile("tmp/PA2_test.sql")
	check(err)
	// fmt.Print(string(fileContents))

	// More control over which parts of file are read
	f, err := os.Open(filepath)
	check(err)

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

	f.Close()

	// Open the test file for reading
	// sqlFile, err := os.Open(filepath)
	check(err)

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
		if strings.HasSuffix(line, ";") == true {
			linesBreaksRemoved = append(linesBreaksRemoved, line)
		} else if strings.HasSuffix(line, ";") == false {
			if line != EXIT {
				a := []string{line, linesCarridgeReturnsRemoved[i+1]}
				linesCarridgeReturnsRemoved[i+1] = strings.Join(a, "")
			}
		}
		i++
	}

	// now that statements are sanitized, generate Tokens
	tokenSequences := make([]string, 0, 1000)
	for _, line := range linesBreaksRemoved {
		tokenSequences = append(linesBreaksRemoved, tokenizer.GetTokenSequence(line))
	}

	// print
	for _, tokenSeq := range tokenSequences {
		fmt.Println(tokenSeq)
	}

	// Bind sqlfiledata to string array
	// testArray := strings.Fields(fileContentsString)

	// Iterate through string array and print
	// for _, v := range testArray {
	// fmt.Println(v)
	// }
}

// SEMICOLON delimiter
const SEMICOLON = ";"

// COMMENT delimiter
const COMMENT = "--"

// EXIT delimiter
const EXIT = ".exit"
