package sqlops

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)

// For tokenizer - https://github.com/xwb1989/sqlparser
func (x *SQLOperator) GetTablesInQuery(query string) ([]string, error) {
	log.Println("query to be parsed", query)
	stmt, err := x.parser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %v", err)
	}

	var tableNames []string

	// Handle different SQL statement types
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		tableNames = append(tableNames, x.extractTablesFromSelect(stmt)...)
	case *sqlparser.Insert:
		// tableNames = append(tableNames, stmt.Table.As.String())
		return nil, fmt.Errorf("DDL statements are not supported")
	case *sqlparser.Update:
		// tableNames = append(tableNames, x.extractTablesFromTableExprs(stmt.TableExprs)...)
		return nil, fmt.Errorf("DDL statements are not supported")
	case *sqlparser.Delete:
		// tableNames = append(tableNames, x.extractTablesFromTableExprs(stmt.TableExprs)...)
		return nil, fmt.Errorf("DDL statements are not supported")
	// case *sqlparser.Union:
	// 	for _, sel := range stmt.Left. {
	// 		sel.
	// 		tableNames = append(tableNames, x.extractTablesFromSelect(sel)...)
	// 	}
	case sqlparser.DDLStatement:
		return nil, fmt.Errorf("DDL statements are not supported")
	}

	return unique(tableNames), nil
}

func (x *SQLOperator) extractTablesFromSelect(sel *sqlparser.Select) []string {
	var tables []string
	for _, tableExpr := range sel.From {
		log.Println("tableExpr ===============", tableExpr)
		switch tableExpr := tableExpr.(type) {
		case *sqlparser.AliasedTableExpr:
			switch table := tableExpr.Expr.(type) {
			case sqlparser.TableName:
				tables = append(tables, table.Name.String())
			}
		case *sqlparser.JoinTableExpr:
			tables = append(tables, x.extractTablesFromJoin(tableExpr)...)
		}
	}
	log.Println("tables ===============", tables)
	log.Println("sel ===============", sel.From)
	if sel.Where != nil {
		tables = append(tables, x.extractTablesFromExpr(sel.Where.Expr)...)
	}
	for _, expr := range sel.SelectExprs {
		if aliasedExpr, ok := expr.(*sqlparser.AliasedExpr); ok {
			tables = append(tables, x.extractTablesFromExpr(aliasedExpr.Expr)...)
		}
	}

	return tables
}

func (x *SQLOperator) extractTablesFromJoin(joinExpr *sqlparser.JoinTableExpr) []string {
	var tables []string
	if leftTable, ok := joinExpr.LeftExpr.(*sqlparser.AliasedTableExpr); ok {
		if table, ok := leftTable.Expr.(sqlparser.TableName); ok {
			tables = append(tables, table.Name.String())
		}
	}
	if rightTable, ok := joinExpr.RightExpr.(*sqlparser.AliasedTableExpr); ok {
		if table, ok := rightTable.Expr.(sqlparser.TableName); ok {
			tables = append(tables, table.Name.String())
		}
	}
	return tables
}

func (x *SQLOperator) extractTablesFromExpr(expr sqlparser.Expr) []string {
	var tables []string
	switch expr := expr.(type) {
	case *sqlparser.Subquery:
		// Handle subqueries recursively
		if sel, ok := expr.Select.(*sqlparser.Select); ok {
			tables = append(tables, x.extractTablesFromSelect(sel)...)
		}
	case *sqlparser.ColName:
		// Do nothing for column names
	default:
		// Handle other types of expressions (optional, can be expanded based on needs)
	}
	return tables
}

// extractTablesFromExpr extracts tables from expressions (used for subqueries).
func (x *SQLOperator) extractTablesFromTableExprs(expr []sqlparser.TableExpr) []string {
	var tables []string
	for _, table := range expr {
		tables = append(tables, table.(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String())
	}
	return tables
}

// unique removes duplicate values from a slice
func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func (x *SQLOperator) IsDMLQuery(query string) bool {
	return sqlparser.IsDML(query)
}

func (x *SQLOperator) IsDDLQuery(query string) bool {
	log.Println("query to be validated for DDL -- ", query)
	stmt, err := x.parser.Parse(query)
	if err != nil {
		x.logger.Error().Err(err).Msg("failed to parse query")
		return true
	}
	log.Println("stmt ===============", stmt)
	_, ok := stmt.(sqlparser.DDLStatement)
	return ok
}

func (x *SQLOperator) IsNotROQuery(query string) bool {
	return sqlparser.IsDML(query) || x.IsDDLQuery(query)
}

func isSQL(input string) bool {
	// List of SQL keywords (you can expand this list)
	sqlKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "FROM", "WHERE", "JOIN", "ORDER BY", "GROUP BY", "HAVING", "LIMIT", "DISTINCT", "VALUES", "ALTER", "DROP", "CREATE", "INTO",
	}

	// Convert input to uppercase for case-insensitive comparison
	input = strings.ToUpper(input)

	// Check if the input contains any SQL keywords
	for _, keyword := range sqlKeywords {
		if strings.Contains(input, keyword) {
			return true
		}
	}

	// Check if the input contains SQL operators
	sqlOperators := []string{
		"=", ">", "<", "<>", "AND", "OR", "LIKE", "BETWEEN", "IS NULL", "IN",
	}
	for _, operator := range sqlOperators {
		if strings.Contains(input, operator) {
			return true
		}
	}

	// You can extend this with regex for more complex checks, e.g., checking for SQL clauses

	// Basic check for SQL-like structure using a simple regex (e.g., SELECT query structure)
	sqlRegex := `(?i)(SELECT|INSERT|UPDATE|DELETE|FROM|WHERE)`
	re := regexp.MustCompile(sqlRegex)
	if re.MatchString(input) {
		return true
	}

	// If no SQL-like patterns were found, assume it's normal text
	return false
}
