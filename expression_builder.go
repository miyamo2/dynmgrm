package dynamgorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		v := items[i]
		if bindVarIfCollectionType(stmt, v) {
			continue
		}
		stmt.AddVar(stmt, items[i])
	}
	stmt.WriteByte('}')
}

// bindVarIfCollectionType binds a variable if it is a collection type
func bindVarIfCollectionType(stmt *gorm.Statement, value interface{}) (bound bool) {
	switch value := (value).(type) {
	case Sets[string],
		Sets[[]byte],
		Sets[int],
		Sets[float64]:
		stmt.Vars = append(stmt.Vars, value)
		stmt.DB.Dialector.BindVarTo(stmt, stmt, value)
		bound = true
	case Map:
		if err := resolveCollectionsNestedInMap(&value); err != nil {
			break
		}
		stmt.Vars = append(stmt.Vars, value)
		stmt.DB.Dialector.BindVarTo(stmt, stmt, value)
		bound = true
	case List:
		if err := resolveCollectionsNestedInList(&value); err != nil {
			break
		}
		stmt.Vars = append(stmt.Vars, value)
		stmt.DB.Dialector.BindVarTo(stmt, stmt, value)
		bound = true
	}
	return
}

// buildSetClause builds SET clause
func buildSetClause(set clause.Set, stmt *gorm.Statement) {
	if len(set) <= 0 {
		return
	}
	for idx, assignment := range set {
		if idx > 0 {
			stmt.WriteByte(' ')
		}
		stmt.WriteString("SET ")
		stmt.WriteString(assignment.Column.Name)
		stmt.WriteByte('=')
		asgv := assignment.Value
		if bindVarIfCollectionType(stmt, asgv) {
			continue
		}
		stmt.AddVar(stmt, asgv)
	}
}
