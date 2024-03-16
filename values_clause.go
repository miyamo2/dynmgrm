package dynamgorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// buildValuesClause builds VALUES clause
func buildValuesClause(c clause.Clause, builder clause.Builder) {
	values, ok := c.Expression.(clause.Values)
	if !ok {
		return
	}
	columns := values.Columns
	if len(columns) <= 0 {
		return
	}
	// PartiQL for DynamoDB does not support multiple rows in VALUES clause
	items := values.Values[0]

	builder.WriteString("VALUE ")
	builder.WriteByte('{')
	for i, column := range columns {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf(`'%s'`, column.Name))
		builder.WriteString(" : ")
		v := items[i]
		if bindVarIfCollectionType(builder, v) {
			continue
		}
		builder.AddVar(builder, items[i])
	}
	builder.WriteByte('}')
}

func bindVarIfCollectionType(builder clause.Builder, value interface{}) (bound bool) {
	stmt, ok := builder.(*gorm.Statement)
	if !ok {
		return
	}
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
