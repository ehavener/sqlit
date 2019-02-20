package parser

import (
	// "fmt"
	"sqlit/tokenizer"
	// "strings"
)

// consts ...
const (
	CREATE   = "CREATE"
	DROP     = "DROP"
	USE      = "USE"
	DATABASE = "DATABASE"
	TABLE    = "TABLE"
	INSERT   = "INSERT"
	ALTER    = "ALTER"
	DELETE   = "DELETE"
	SELECT   = "SELECT"
	LITERAL  = "LITERAL"
	special  = "special"
)

var statementTypes = map[string]string{
	"CREATE_DATABASE": "CREATE_DATABASE",
	"DROP_DATABASE":   "DROP_DATABASE",
	"USE_DATABASE":    "USE_DATABASE",
	"CREATE_TABLE":    "CREATE_TABLE",
	"DROP_TABLE":      "DROP_TABLE",
	"ALTER_TABLE":     "ALTER_TABLE",
	"INSERT":          "INSERT",
	"SELECT":          "SELECT",
	"DELETE":          "DELETE",
}

var specialTypes = map[string]string{
	"DATABASE_NAME": "DATABASE_NAME",
	"TABLE_NAME":    "TABLE_NAME",
	"COL_NAME":      "COL_NAME",
	"COL_TYPE":      "COL_TYPE",
	"ADD_COL":       "ADD_COL",
	"ALL":           "ALL",
	"FROM":          "FROM",
}

// ParseStatement ....
func ParseStatement(statement tokenizer.Statement) tokenizer.Statement {
	statement = inferStatementType(statement)
	// TODO // statement = constructParseTree(statement)
	statement = inferTokenSpecialTypes(statement)
	return statement
}

// inferStatementType infers a class of meaning based on the beginning of a statement
func inferStatementType(statement tokenizer.Statement) tokenizer.Statement {

	// DATABASE
	if statement.Tokens[0].Name == CREATE && statement.Tokens[1].Name == DATABASE {
		statement.Type = statementTypes["CREATE_DATABASE"]
		return statement
	}

	if statement.Tokens[0].Name == DROP && statement.Tokens[1].Name == DATABASE {
		statement.Type = statementTypes["DROP_DATABASE"]
		return statement
	}

	if statement.Tokens[0].Name == USE {
		statement.Type = statementTypes["USE_DATABASE"]
		return statement
	}

	if statement.Tokens[0].Name == CREATE && statement.Tokens[1].Name == TABLE {
		statement.Type = statementTypes["CREATE_TABLE"]
		return statement
	}

	if statement.Tokens[0].Name == DROP && statement.Tokens[1].Name == TABLE {
		statement.Type = statementTypes["DROP_TABLE"]
		return statement
	}

	if statement.Tokens[0].Name == ALTER && statement.Tokens[1].Name == TABLE {
		statement.Type = statementTypes["ALTER_TABLE"]
		return statement
	}

	if statement.Tokens[0].Name == INSERT {
		statement.Type = statementTypes["INSERT"]
		return statement
	}

	if statement.Tokens[0].Name == SELECT {
		statement.Type = statementTypes["SELECT"]
		return statement
	}

	if statement.Tokens[0].Name == DELETE {
		statement.Type = statementTypes["DELETE"]
		return statement
	}

	return statement
}

// TODO a parse tree might be worthwhile for complex queries
// func constructParseTree()

// inferTokenSpecialTypes infers contextual token types
func inferTokenSpecialTypes(statement tokenizer.Statement) tokenizer.Statement {
	switch statement.Type {
	case statementTypes["CREATE_DATABASE"]:
		statement = parseCreateDatabase(statement)
	case statementTypes["DROP_DATABASE"]:
		statement = parseDropDatabase(statement)
	case statementTypes["USE_DATABASE"]:
		statement = parseUseDatabase(statement)
	case statementTypes["CREATE_TABLE"]:
		statement = parseCreateTable(statement)
	case statementTypes["ALTER_TABLE"]:
		statement = parseAlterTable(statement)
	case statementTypes["DROP_TABLE"]:
		statement = parseDropTable(statement)
	case statementTypes["INSERT"]:
		statement = parseInsert(statement)
	case statementTypes["SELECT"]:
		statement = parseSelect(statement)
	}

	return statement
}

func parseCreateDatabase(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialTypes["DATABASE_NAME"])
	return statement
}

func parseDropDatabase(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialTypes["DATABASE_NAME"])
	return statement
}

func parseUseDatabase(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 1, specialTypes["DATABASE_NAME"])
	return statement
}

func parseDropTable(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialTypes["TABLE_NAME"])
	return statement
}

func parseCreateTable(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialTypes["TABLE_NAME"])

	i := 3
	for i < len(statement.Tokens) {
		if i%2 == 0 {
			setSpecialNameIfTokenExists(statement, i, specialTypes["COL_TYPE"])
		} else {
			setSpecialNameIfTokenExists(statement, i, specialTypes["COL_NAME"])
		}
		i++
	}

	return statement
}

func parseAlterTable(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialTypes["TABLE_NAME"])
	if statement.Tokens[3].Special == "ADD" {
		setSpecialNameIfTokenExists(statement, 3, specialTypes["ADD_COL"])
		setSpecialNameIfTokenExists(statement, 4, specialTypes["COL_TYPE"])
		setSpecialNameIfTokenExists(statement, 5, specialTypes["COL_NAME"])
	}
	return statement
}

// https://www.sqlite.org/draft/syntaxdiagrams.html#select-stmt
func parseSelect(statement tokenizer.Statement) tokenizer.Statement {
	if statement.Tokens[1].Special == "*" {
		setSpecialNameIfTokenExists(statement, 1, specialTypes["ALL"])
	}

	if statement.Tokens[2].Special == "FROM" {
		setSpecialNameIfTokenExists(statement, 2, specialTypes["FROM"])
		setSpecialNameIfTokenExists(statement, 3, specialTypes["TABLE_NAME"])
	}
	return statement
}

func parseInsert(statement tokenizer.Statement) tokenizer.Statement {
	return statement
}

func parseDelete(statement tokenizer.Statement) tokenizer.Statement {
	return statement
}

func setSpecialNameIfTokenExists(statement tokenizer.Statement, i int, name string) {
	if len(statement.Tokens) >= i {
		statement.Tokens[i].Name = name
	}
}
