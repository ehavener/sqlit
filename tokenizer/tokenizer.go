package tokenizer

import (
	"fmt"
	// "io/ioutil"
	// "os"
	// "sqlit/tokenizer"
	"strings"
	// "bufio"
	// "io"
	// "go/scanner"
)

/*
var keywords = map[string]int{
	"add":      ADD,
	"against":  AGAINST,
	"alter":    ALTER,
	"create":   CREATE,
	"database": DATABASE,
	"drop":     DROP,
	"float":    FLOAT,
	"from":     FROM,
	"insert":   INSERT,
	"int":      INT,
	"into":     INTO,
	"select":   SELECT,
	"set":      SET,
	"table":    TABLE,
	"update":   UPDATE,
	"use":      USE,
	"values":   VALUES,
	"varchar":  VARCHAR,
	"where":    WHERE,
}
*/

// Constants

/*
const (
	CREATE   = 0
	DROP     = 1
	USE      = 2
	DATABASE = 3
	TABLE    = 4
	special  = 5
)
*/
/*
var kw = map[string]int{
	"CREATE":   CREATE,
	"DROP":     DROP,
	"USE":      USE,
	"DATABASE": DATABASE,
	"TABLE":    TABLE,
	"special":  special,
}
*/

// langWords are general classes for tokens
var langWords = [11]string{
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

// A Token represents a SQL word
type Token struct {
	Name    string
	Special string
}

// A Statement is an array of tokens and a general class of meaning
type Statement struct {
	Tokens []Token
	Type   string
}

// TokenizeStatement breaks a string of SQL into a statement
func TokenizeStatement(rawStatement string) Statement {

	rawWords := strings.Fields(rawStatement)

	tokens := make([]Token, 0, 100000)

	for _, rawWord := range rawWords {
		tokens = append(tokens, TokenizeWord(rawWord))
	}

	statement := Statement{Tokens: tokens}

	return statement
}

// TokenizeWord maps a string to a token and classifies it
func TokenizeWord(word string) Token {
	for _, langWord := range langWords {
		if strings.EqualFold(word, langWord) {
			return Token{Name: langWord}
		}
	}

	return Token{Name: "special", Special: word}
}

// PrintToken is used to debug token properties
func PrintToken(token Token) {
	fmt.Print("	", token.Name, " ")
	fmt.Print("	", token.Special, " ")
}
