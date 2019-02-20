package generator

import (
	"fmt"
	"sqlit/io"
	"sqlit/tokenizer"
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

// Function ...
type Function struct {
	name string
}

// Generate ...
func Generate(statement tokenizer.Statement) Function {
	function := Function{}

	switch statement.Type {
	case statementTypes["CREATE_DATABASE"]:
		function = generateCreateDatabase(statement)
	case statementTypes["DROP_DATABASE"]:
		function = generateDropDatabase(statement)

		/*
			case statementTypes["USE_DATABASE"]:
				function = generateUseDatabase(statement)
			case statementTypes["CREATE_TABLE"]:
				function = generateCreateTable(statement)
			case statementTypes["ALTER_TABLE"]:
				function = generateAlterTable(statement)
			case statementTypes["DROP_TABLE"]:
				function = generateDropTable(statement)
			case statementTypes["INSERT"]:
				function = generateInsert(statement)
			case statementTypes["SELECT"]:
				function = generateSelect(statement)
		*/
	}
	return function
}

func generateCreateDatabase(statement tokenizer.Statement) Function {

	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	if io.CheckIfDatabaseExists(name) == true {
		fmt.Println("!Failed to create database db_1 because it already exists.")
		// TODO: handle error
	} else {
		io.CreateDatabase(name)
		fmt.Println("Database " + name + " created.")
	}

	return Function{}
}

func generateDropDatabase(statement tokenizer.Statement) Function {

	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	if io.CheckIfDatabaseExists(name) == false {
		fmt.Println("error - database with name " + name + " does not exist")
		// TODO: handle error
	} else {
		io.DeleteDatabase(name)
	}

	return Function{}
}

func getFirstTokenOfName(statement tokenizer.Statement, name string) string {
	for _, token := range statement.Tokens {
		if token.Name == name {
			return token.Special
		}
	}
	return "err"
}
