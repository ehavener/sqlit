package tokenizer

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

type token struct {
	name string
}

// GetTokenSequence ...
func GetTokenSequence(line string) string {
	return "we tokenizing"
}

// GetToken ...
func GetToken(chunk string) string {
	return chunk
}
