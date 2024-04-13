package dynmgrm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_ gorm.Valuer = (*TypedList[interface{}])(nil)
	_ sql.Scanner = (*TypedList[interface{}])(nil)
)

// TypedList is a DynamoDB list type with type specification.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type TypedList[T any] []T

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
		switch v := v.(type) {
		case map[string]interface{}:
			dest := new(T)
			rv := reflect.ValueOf(dest)
			for k, a := range v {
				f := rv.Elem().FieldByName(k)
				switch f.Interface().(type) {
				case string:
					str, ok := a.(string)
					if !ok {
						return fmt.Errorf("incompatible %T and %T", f.Interface(), a)
					}
					f.SetString(str)
					continue
				case int:
					i, ok := a.(int)
					if !ok {
						return fmt.Errorf("incompatible %T and %T", f.Interface(), a)
					}
					f.SetInt(int64(i))
					continue
				case bool:
					b, ok := a.(bool)
					if !ok {
						return fmt.Errorf("incompatible %T and %T", f.Interface(), a)
					}
					f.SetBool(b)
					continue
				case float64:
					f64, ok := a.(float64)
					if !ok {
						return fmt.Errorf("incompatible %T and %T", f.Interface(), a)
					}
					f.SetFloat(f64)
					continue
				case *string:
					str, ok := a.(string)
					if !ok {
						break
					}
					f.Set(reflect.ValueOf(&str))
					continue
				case *int:
					i, ok := a.(int)
					if !ok {
						break
					}
					f.Set(reflect.ValueOf(&i))
					continue
				case *bool:
					b, ok := a.(bool)
					if !ok {
						break
					}
					f.Set(reflect.ValueOf(&b))
					continue
				case *float64:
					f64, ok := a.(float64)
					if !ok {
						break
					}
					f.Set(reflect.ValueOf(&f64))
					continue
				}
				if !f.CanAddr() {
					continue
				}
				switch ptr := f.Addr().Interface().(type) {
				case sql.Scanner:
					if err := ptr.Scan(a); err != nil {
						return err
					}
					continue
				}
			}
			*l = append(*l, *dest)
		}
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
