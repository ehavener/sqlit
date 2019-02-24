/* UNR CS 457 | SPRING 2019 | emerson@nevada.unr.edu */

// Package tokenizer transforms a string of SQL into a new Statement structure
package tokenizer

import (
	"fmt"
	"strings"
)

// A Token represents a SQL word
type Token struct {
	Name    string
	Special string
}

// Names are general classes for tokens
var Names = [11]string{
	"CREATE",
	"DROP",
	"USE",
	"DATABASE",
	"TABLE",
	"INSERT",
	"ALTER",
	"DELETE",
	"SELECT",
	"LITERAL",
	"special",
}

// A Statement is an array of tokens and a general class of meaning
type Statement struct {
	Tokens []Token
	Type   string
}

// TokenizeStatement breaks a string of SQL into a statement
func TokenizeStatement(rawStatement string) Statement {
	rawWords := strings.Fields(rawStatement)

	tokens := make([]Token, 0, len(rawWords))

	for _, rawWord := range rawWords {
		tokens = append(tokens, TokenizeWord(rawWord))
	}

	statement := Statement{Tokens: tokens}

	return statement
}

// TokenizeWord maps a string to a token and classifies it
func TokenizeWord(word string) Token {
	for _, name := range Names {
		if strings.EqualFold(word, name) {
			return Token{Name: name}
		}
	}

	return Token{Name: "special", Special: word}
}

//
//			Helper functions
//

// PrintStatement is used to debug statement properties
func PrintStatement(statement Statement) {
	fmt.Print("type	", statement.Type, "\n")
	for _, token := range statement.Tokens {
		PrintToken(token)
		fmt.Print("\n")
	}

	fmt.Print("\n")
}

// PrintToken is used to debug token properties
func PrintToken(token Token) {
	fmt.Print("	", token.Name, " ")
	fmt.Print("	", token.Special, " ")
}
