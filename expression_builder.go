package dynamgorm

import (
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
