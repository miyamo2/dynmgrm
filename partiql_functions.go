package dynmgrm

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"regexp"
)

var (
	reInvalidColumnName  = regexp.MustCompile(`^[\w._]+$`)
	ErrInvalidColumnName = errors.New("column name contains invalid characters")
)

var (
	_ functionForPartiQLUpdates = (*listAppend)(nil)
)

// functionForPartiQLUpdates is an interface for PartiQL functions that can be used in updates.
type functionForPartiQLUpdates interface {
	// expression returns a string that represents the function in the SQL query.
	expression(db *gorm.DB, column string) string
	// bindVariable returns a gorm.Valuer that represents the bind variable for the function.
	bindVariable() gorm.Valuer
}

// listAppend is a struct that implements functionForPartiQLUpdates interface for `list_append` function.
type listAppend struct {
	value List
}

func (la *listAppend) expression(db *gorm.DB, column string) string {
	if !reInvalidColumnName.MatchString(column) {
		db.AddError(ErrInvalidColumnName)
		return ""
	}
	return fmt.Sprintf("list_append(%s, ", column)
}

func (la *listAppend) bindVariable() gorm.Valuer {
	return la.value
}

// ListAppend returns a functionForPartiQLUpdates implementation for `list_append` function.
func ListAppend(item ...interface{}) *listAppend {
	return &listAppend{value: item}
}
