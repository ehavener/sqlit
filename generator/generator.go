package generator

import (
	"fmt"
	"sqlit/io"
	"sqlit/tokenizer"
	"strings"
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
	case statementTypes["USE_DATABASE"]:
		function = generateUseDatabase(statement)
	case statementTypes["CREATE_TABLE"]:
		function = generateCreateTable(statement)
	case statementTypes["ALTER_TABLE"]:
		function = generateAlterTable(statement)
	case statementTypes["DROP_TABLE"]:
		function = generateDropTable(statement)
	case statementTypes["SELECT"]:
		function = generateSelect(statement)
		//case statementTypes["INSERT"]:
		//	function = generateInsert(statement)
	}

	return function
}

func generateCreateDatabase(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	if io.CheckIfDatabaseExists(name) == true {
		fmt.Println("!Failed to create database " + name + " because it already exists.")
		// TODO: handle error
		return Function{}
	}
	io.CreateDatabase(name)
	fmt.Println("Database " + name + " created.")

	return Function{}
}

func generateDropDatabase(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	if io.CheckIfDatabaseExists(name) == false {
		fmt.Println("!Failed to delete " + name + " because it does not exist.")
		// TODO: handle error
		return Function{}
	}

	io.DeleteDatabase(name)
	fmt.Println("Database " + name + " deleted.")

	return Function{}
}

func generateUseDatabase(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	if io.CheckIfDatabaseExists(name) == false {
		fmt.Println("!Failed to use database " + name + " because it does not exist.")
		// TODO: handle error
		return Function{}
	}

	io.UseDatabase(name)
	fmt.Println("Using database " + name)

	return Function{}
}

func generateCreateTable(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "TABLE_NAME")
	columns := getAllTokensOfName(statement, "COL_NAME")
	constraints := getAllTokensOfName(statement, "COL_TYPE")

	// TODO: move this to tokenizer?
	columns = sanitizeColumns(columns)
	constraints = sanitizeConstraints(constraints)

	if io.CheckIfAnyDatabaseIsInUse() == false {
		fmt.Println("!Failed to create table " + name + " because no database is in use.")
		// TODO: handle error
		return Function{}
	}

	if io.CheckIfTableExists(name) == true {
		fmt.Println("!Failed to create table " + name + " because it already exists.")
		// TODO: handle error
		return Function{}
	}

	io.CreateTable(name, columns, constraints)
	fmt.Println("Table " + name + " created.")

	return Function{}
}

func generateDropTable(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "TABLE_NAME")

	if io.CheckIfAnyDatabaseIsInUse() == false {
		fmt.Println("!Failed to delete table " + name + " because no database is in use.")
		// TODO: handle error
		return Function{}
	}

	if io.CheckIfTableExists(name) == false {
		fmt.Println("!Failed to delete table " + name + " because it does not exist.")
		// TODO: handle error
		return Function{}
	}

	io.DropTable(name)
	fmt.Println("Table " + name + " deleted.")

	return Function{}
}

func generateAlterTable(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "TABLE_NAME")
	method := getFirstTokenOfName(statement, "ADD_COL")
	column := getFirstTokenOfName(statement, "COL_NAME")
	constraint := getFirstTokenOfName(statement, "COL_TYPE")

	if io.CheckIfAnyDatabaseIsInUse() == false {
		fmt.Println("!Failed to alter table " + name + " because no database is in use.")
		// TODO: handle error
		return Function{}
	}

	if io.CheckIfTableExists(name) == false {
		fmt.Println("!Failed to alter table " + name + " because it does not exist.")
		// TODO: handle error
		return Function{}
	}

	io.AlterTable(name, method, column, constraint)
	fmt.Println("Table " + name + " modified.")

	return Function{}
}

func generateSelect(statement tokenizer.Statement) Function {
	name := getFirstTokenOfName(statement, "TABLE_NAME")

	if io.CheckIfAnyDatabaseIsInUse() == false {
		fmt.Println("!Failed to query table " + name + " because no database is in use.")
		// TODO: handle error
		return Function{}
	}

	if io.CheckIfTableExists(name) == false {
		fmt.Println("!Failed to query table " + name + " because it does not exist.")
		// TODO: handle error
		return Function{}
	}

	result := io.SelectAll(name)
	fmt.Println(result)

	return Function{}
}

func getAllTokensOfName(statement tokenizer.Statement, name string) []string {
	var specials []string
	for _, token := range statement.Tokens {
		if token.Name == name {
			specials = append(specials, token.Special)
		}
	}
	return specials
}

func getFirstTokenOfName(statement tokenizer.Statement, name string) string {
	for _, token := range statement.Tokens {
		if token.Name == name {
			return token.Special
		}
	}
	return "err"
}

func sanitizeColumns(columns []string) []string {
	for i := 0; i < len(columns); i++ {
		columns[i] = strings.Replace(columns[i], "(", "", 1)
	}
	return columns
}

func sanitizeConstraints(constraints []string) []string {
	for i := 0; i < len(constraints); i++ {
		constraints[i] = strings.Replace(constraints[i], ")", "", 1)
		constraints[i] = strings.Replace(constraints[i], ",", "", -1)
	}
	return constraints
}
