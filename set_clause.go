package dynamgorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func buildSetClause(set clause.Set, statement *gorm.Statement) {
	if len(set) <= 0 {
		return
	}
	for idx, assignment := range set {
		if idx > 0 {
			statement.WriteByte(',')
		}
		statement.WriteQuoted(assignment.Column)
		statement.WriteByte('=')
		statement.AddVar(statement, assignment.Value)
	}
}
