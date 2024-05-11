package dynmgrm

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
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
	prfl := stmt.Schema.PrimaryFieldDBNames
	for i, column := range columns {
		v := items[i]
		if isZeroValue(v) && !slices.Contains[[]string](prfl, column.Name) {
			continue
		}
		if i > 0 {
			stmt.WriteString(", ")
		}
		stmt.WriteString(fmt.Sprintf(`'%s'`, column.Name))
		stmt.WriteString(" : ")

		// NOTE: this is a temporary hack, if `btnguyen2k/godynamo` will support `driver.Valuer`, remove the entire switch block.
		switch dv := v.(type) {
		case driver.Valuer:
			dvv, err := dv.Value()
			if err != nil {
				stmt.AddError(err)
				continue
			}
			v = dvv
		}
		stmt.AddVar(stmt, v)
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
		stmt.WriteQuoted(asgcol)
		stmt.WriteByte('=')
		asgv := assignment.Value
		switch asgv := asgv.(type) {
		case functionForPartiQLUpdates:
			stmt.WriteString(asgv.expression(stmt.DB, asgcol))
			stmt.AddVar(stmt, asgv.bindVariable().GormValue(stmt.Context, stmt.DB))
			stmt.WriteByte(')')
			continue
		}
		stmt.AddVar(stmt, asgv)
	}
}

func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsZero()
}
