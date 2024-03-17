package dynamgorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type expressionBuilder[T clause.Expression] func(expression T, statement *gorm.Statement)

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
