package dynmgrm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slices"
)

// expressionBuilder is a function that builds a clause.Expression
type expressionBuilder[T clause.Expression] func(expression T, statement *gorm.Statement)

// toClauseBuilder converts expressionBuilder to clause.ClauseBuilder
func toClauseBuilder[T clause.Expression](xpBuilder expressionBuilder[T]) clause.ClauseBuilder {
	return func(c clause.Clause, builder clause.Builder) {
		xp, ok := c.Expression.(T)
		if !ok {
			return
		}
		statement, ok := builder.(*gorm.Statement)
		if !ok {
			return
		}
		xpBuilder(xp, statement)
	}
}

// buildValuesClause builds VALUES clause
func buildValuesClause(values clause.Values, stmt *gorm.Statement) {
	columns := values.Columns
	if len(columns) <= 0 {
		return
	}
	// PartiQL for DynamoDB does not support multiple rows in VALUES clause
	items := values.Values[0]

	stmt.WriteString("VALUE ")
	stmt.WriteByte('{')
	for i, column := range columns {
		if i > 0 {
			stmt.WriteString(", ")
		}
		stmt.WriteString(fmt.Sprintf(`'%s'`, column.Name))
		stmt.WriteString(" : ")
		stmt.AddVar(stmt, items[i])
	}
	stmt.WriteByte('}')
}

// buildSetClause builds SET clause
func buildSetClause(set clause.Set, stmt *gorm.Statement) {
	if len(set) <= 0 {
		return
	}
	prfl := stmt.Schema.PrimaryFieldDBNames
	for idx, assignment := range set {
		asgcol := assignment.Column.Name
		if slices.Contains[[]string](prfl, asgcol) {
			continue
		}
		if idx > 0 {
			stmt.WriteByte(' ')
		}
		stmt.WriteString("SET ")
		stmt.WriteString(assignment.Column.Name)
		stmt.WriteByte('=')
		asgv := assignment.Value
		stmt.AddVar(stmt, asgv)
	}
}
