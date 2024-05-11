package dynmgrm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"slices"
)

// compatibility check
var (
	_ gorm.Valuer = (*TypedList[interface{}])(nil)
	_ sql.Scanner = (*TypedList[interface{}])(nil)
)

// TypedList is a DynamoDB list type with type specification.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type TypedList[T any] []T

// GormDataType returns the data type for Gorm.
func (l *TypedList[T]) GormDataType() string {
	return "dgtypedlist"
}

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (l *TypedList[T]) Scan(value interface{}) error {
	if len(*l) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	sv, ok := value.([]interface{})
	if !ok {
		return errors.Join(ErrFailedToCast, fmt.Errorf("incompatible %T and %T", l, value))
	}
	*l = slices.Grow(*l, len(sv))
	for _, v := range sv {
		mv, ok := v.(map[string]interface{})
		if !ok {
			var t T
			return errors.Join(ErrFailedToCast, fmt.Errorf("incompatible %T and %T", t, v))
		}
		dest := new(T)
		rv := reflect.ValueOf(dest)
		rt := reflect.TypeOf(*dest)
		err := assignMapValueToReflectValue(rt, rv, mv)
		if err != nil {
			return err
		}
		*l = append(*l, *dest)
	}
	return nil
}

// GormValue implements the [gorm.Valuer] interface.
//
// [gorm.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (l TypedList[T]) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	avl := types.AttributeValueMemberL{Value: make([]types.AttributeValue, 0, len(l))}
	for _, v := range l {
		av, err := toDocumentAttributeValue[*types.AttributeValueMemberM](v)
		if err != nil {
			_ = db.AddError(err)
			return clause.Expr{}
		}
		avl.Value = append(avl.Value, av)
	}
	return clause.Expr{SQL: "?", Vars: []interface{}{avl}}
}
