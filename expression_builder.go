package dynmgrm

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
			v = toPhysicalDocumentAttributeValue(dvv)
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
			dvv, err := asgv.bindVariable().Value()
			if err != nil {
				stmt.AddError(err)
				continue
			}
			stmt.WriteString(asgv.expression(stmt.DB, asgcol))
			stmt.AddVar(stmt, toPhysicalDocumentAttributeValue(dvv))
			stmt.WriteByte(')')
			continue
		}
		// NOTE: this is a temporary hack, if `btnguyen2k/godynamo` will support `driver.Valuer`, remove the entire switch block.
		switch dv := asgv.(type) {
		case driver.Valuer:
			dvv, err := dv.Value()
			if err != nil {
				stmt.AddError(err)
				continue
			}
			asgv = toPhysicalDocumentAttributeValue(dvv)
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

func toPhysicalDocumentAttributeValue(v interface{}) driver.Value {
	switch v := v.(type) {
	case *types.AttributeValueMemberM:
		return *v
	case *types.AttributeValueMemberL:
		return *v
	case *types.AttributeValueMemberSS:
		return *v
	case *types.AttributeValueMemberNS:
		return *v
	case *types.AttributeValueMemberBS:
		return *v
	}
	return v
}
