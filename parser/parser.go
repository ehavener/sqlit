// Package parser ... 
package parser

import (
	"sqlit/tokenizer"
)

var names = map[string]string{
	"CREATE":   "CREATE",
	"DROP":     "DROP",
	"USE":      "USE",
	"DATABASE": "DATABASE",
	"TABLE":    "TABLE",
	"INSERT":   "INSERT",
	"ALTER":    "ALTER",
	"SELECT":   "SELECT",
	"DELETE":   "DELETE",
	"LITERAL":  "LITERAL",
	"special":  "special",
}

var specialNames = map[string]string{
	"DATABASE_NAME": "DATABASE_NAME",
	"TABLE_NAME":    "TABLE_NAME",
	"COL_NAME":      "COL_NAME",
	"COL_TYPE":      "COL_TYPE",
	"ADD_COL":       "ADD_COL",
	"ALL":           "ALL",
	"FROM":          "FROM",
}

// Types are general classes for statements
var Types = map[string]string{
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

// ParseStatement ....
func ParseStatement(statement tokenizer.Statement) tokenizer.Statement {
	// TODO: statement = constructParseTree(statement)
	statement = inferStatementType(statement)
	statement = inferTokenspecialNames(statement)
	return statement
}

// TODO: a parse tree might be worthwhile for complex queries instead of all these conditionals
// func constructParseTree() {}

// inferStatementType infers a statement's general type based on how it begins
func inferStatementType(statement tokenizer.Statement) tokenizer.Statement {
	if len(statement.Tokens) < 2 {
		return statement
	}

	if statement.Tokens[0].Name == names["CREATE"] && statement.Tokens[1].Name == names["DATABASE"] {
		statement.Type = Types["CREATE_DATABASE"]
		return statement
	}

	if statement.Tokens[0].Name == names["DROP"] && statement.Tokens[1].Name == names["DATABASE"] {
		statement.Type = Types["DROP_DATABASE"]
		return statement
	}

	if statement.Tokens[0].Name == names["USE"] {
		statement.Type = Types["USE_DATABASE"]
		return statement
	}

	if statement.Tokens[0].Name == names["CREATE"] && statement.Tokens[1].Name == names["TABLE"] {
		statement.Type = Types["CREATE_TABLE"]
		return statement
	}

	if statement.Tokens[0].Name == names["DROP"] && statement.Tokens[1].Name == names["TABLE"] {
		statement.Type = Types["DROP_TABLE"]
		return statement
	}

	if statement.Tokens[0].Name == names["ALTER"] && statement.Tokens[1].Name == names["TABLE"] {
		statement.Type = Types["ALTER_TABLE"]
		return statement
	}

	if statement.Tokens[0].Name == names["INSERT"] {
		statement.Type = Types["INSERT"]
		return statement
	}

	if statement.Tokens[0].Name == names["SELECT"] {
		statement.Type = Types["SELECT"]
		return statement
	}

	if statement.Tokens[0].Name == names["DELETE"] {
		statement.Type = Types["DELETE"]
		return statement
	}

	return statement
}

// inferTokenspecialNames infers the meaning of remaining tokens in a statement
func inferTokenspecialNames(statement tokenizer.Statement) tokenizer.Statement {
	switch statement.Type {
	case Types["CREATE_DATABASE"]:
		statement = parseCreateDatabase(statement)
	case Types["DROP_DATABASE"]:
		statement = parseDropDatabase(statement)
	case Types["USE_DATABASE"]:
		statement = parseUseDatabase(statement)
	case Types["CREATE_TABLE"]:
		statement = parseCreateTable(statement)
	case Types["ALTER_TABLE"]:
		statement = parseAlterTable(statement)
	case Types["DROP_TABLE"]:
		statement = parseDropTable(statement)
	case Types["INSERT"]:
		statement = parseInsert(statement)
	case Types["SELECT"]:
		statement = parseSelect(statement)
	}

	return statement
}

func parseCreateDatabase(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialNames["DATABASE_NAME"])
	return statement
}

func parseDropDatabase(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialNames["DATABASE_NAME"])
	return statement
}

func parseUseDatabase(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 1, specialNames["DATABASE_NAME"])
	return statement
}

func parseDropTable(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialNames["TABLE_NAME"])
	return statement
}

func parseCreateTable(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialNames["TABLE_NAME"])

	i := 3
	for i < len(statement.Tokens) {
		if i%2 == 0 {
			setSpecialNameIfTokenExists(statement, i, specialNames["COL_TYPE"])
		} else {
			setSpecialNameIfTokenExists(statement, i, specialNames["COL_NAME"])
		}
		i++
	}

	return statement
}

func parseAlterTable(statement tokenizer.Statement) tokenizer.Statement {
	setSpecialNameIfTokenExists(statement, 2, specialNames["TABLE_NAME"])
	if statement.Tokens[3].Special == "ADD" {
		setSpecialNameIfTokenExists(statement, 3, specialNames["ADD_COL"])
		setSpecialNameIfTokenExists(statement, 4, specialNames["COL_TYPE"])
		setSpecialNameIfTokenExists(statement, 5, specialNames["COL_NAME"])
	}
	return statement
}

// https://www.sqlite.org/draft/syntaxdiagrams.html#select-stmt
func parseSelect(statement tokenizer.Statement) tokenizer.Statement {
	if statement.Tokens[1].Special == "*" {
		setSpecialNameIfTokenExists(statement, 1, specialNames["ALL"])
	}

	if statement.Tokens[2].Special == "FROM" {
		setSpecialNameIfTokenExists(statement, 2, specialNames["FROM"])
		setSpecialNameIfTokenExists(statement, 3, specialNames["TABLE_NAME"])
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
