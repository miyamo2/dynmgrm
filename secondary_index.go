package dynamgorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// compatibility
var _ clause.Expression = (*secondaryIndexExpression)(nil)

// SecondaryIndexOf is the table with the index to be used.
type SecondaryIndexOf interface {
	string | clause.Table
}

// secondaryIndexExpression is a clause.Expression that represents a secondary index
type secondaryIndexExpression struct {
	table     clause.Table
	tableName string
	indexName string
}

// Build builds the secondaryIndexExpression
func (s secondaryIndexExpression) Build(builder clause.Builder) {
	stmt, ok := builder.(*gorm.Statement)
	if !ok {
		return
	}
	tn := s.tableName
	if tn == "" {
		tn = s.table.Name
	}
	stmt.Table = fmt.Sprintf(`%s.%s`, tn, s.indexName)

	qtn := &strings.Builder{}
	stmt.Dialector.QuoteTo(qtn, tn)
	qin := &strings.Builder{}
	stmt.Dialector.QuoteTo(qin, s.indexName)

	stmt.TableExpr = &clause.Expr{
		SQL:                fmt.Sprintf(`%s.%s`, qtn.String(), qin.String()),
		Vars:               nil,
		WithoutParentheses: false,
	}
}

// SecondaryIndex enables queries using a secondary index
func SecondaryIndex[T SecondaryIndexOf](table T, indexName string) secondaryIndexExpression {
	var xp secondaryIndexExpression
	switch table := (interface{})(table).(type) {
	case string:
		xp = secondaryIndexExpression{
			tableName: table,
			indexName: indexName,
		}
	case clause.Table:
		xp = secondaryIndexExpression{
			table:     table,
			indexName: indexName,
		}
	}
	return xp
}
