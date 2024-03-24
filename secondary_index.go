package dynmgrm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// compatibility
var _ clause.Expression = (*secondaryIndexExpression)(nil)
var _ gorm.StatementModifier = (*secondaryIndexExpression)(nil)

// SecondaryIndexOption is a functional option for secondaryIndexExpression
type SecondaryIndexOption func(*secondaryIndexExpression)

// SecondaryIndexOf is the table with the index to be used.
func SecondaryIndexOf[T string | clause.Table](table T) SecondaryIndexOption {
	return func(s *secondaryIndexExpression) {
		switch table := (interface{})(table).(type) {
		case string:
			s.tableName = table
		case clause.Table:
			s.table = table
		}
	}
}

// secondaryIndexExpression is a clause.Expression that represents a secondary index
type secondaryIndexExpression struct {
	table     clause.Table
	tableName string
	indexName string
}

// ModifyStatement modifies the gorm.Statement to use the secondary index
func (s secondaryIndexExpression) ModifyStatement(stmt *gorm.Statement) {
	tn := s.tableName
	if tn == "" {
		tn = s.table.Name
	}
	if tn == "" {
		// if specified indexName is in the format "table_name.index_name", then split it into table_name and index_name
		// e.g.
		//	dynmgrm.SecondaryIndex("table_name.index_name")
		stn := strings.Split(s.indexName, ".")
		if len(stn) == 2 {
			tn = stn[0]
			s.indexName = stn[1]
		}
	}
	if tn == "" {
		tn = stmt.Table
	}

	stmt.Table = fmt.Sprintf(`%s.%s`, tn, s.indexName)
	qtn := &strings.Builder{}
	stmt.Dialector.QuoteTo(qtn, tn)
	qin := &strings.Builder{}
	stmt.Dialector.QuoteTo(qin, s.indexName)
	stmt.TableExpr = &clause.Expr{
		SQL: fmt.Sprintf(`%s.%s`, qtn.String(), qin.String()),
	}
}

// Build builds the secondaryIndexExpression
func (s secondaryIndexExpression) Build(builder clause.Builder) {
	stmt, ok := builder.(*gorm.Statement)
	if !ok {
		return
	}
	s.ModifyStatement(stmt)
}

// SecondaryIndex enables queries using a secondary index
func SecondaryIndex(indexName string, options ...SecondaryIndexOption) secondaryIndexExpression {
	var xp secondaryIndexExpression
	xp.indexName = indexName
	for _, opt := range options {
		opt(&xp)
	}
	return xp
}
