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
		for i := 0; i < rt.NumField(); i++ {
			tf := rt.Field(i)
			vf := rv.Elem().Field(i)
			name := getDBNameFromStructField(tf)
			a, ok := mv[name]
			if !ok {
				continue
			}
			switch vf.Interface().(type) {
			case string:
				str, ok := a.(string)
				if !ok {
					return fmt.Errorf("incompatible %T and %T", vf.Interface(), a)
				}
				vf.SetString(str)
				continue
			case int:
				f64, ok := a.(float64)
				if !ok {
					return fmt.Errorf("incompatible %T and %T", vf.Interface(), a)
				}
				vf.Set(reflect.ValueOf(int(f64)))
				continue
			case bool:
				b, ok := a.(bool)
				if !ok {
					return fmt.Errorf("incompatible %T and %T", vf.Interface(), a)
				}
				vf.SetBool(b)
				continue
			case float64:
				f64, ok := a.(float64)
				if !ok {
					return fmt.Errorf("incompatible %T and %T", vf.Interface(), a)
				}
				vf.SetFloat(f64)
				continue
			case []byte:
				b, ok := a.([]byte)
				if !ok {
					return fmt.Errorf("incompatible %T and %T", vf.Interface(), a)
				}
				vf.SetBytes(b)
			case *string:
				str, ok := a.(string)
				if !ok {
					break
				}
				vf.Set(reflect.ValueOf(&str))
				continue
			case *int:
				f64, ok := a.(float64)
				if !ok {
					break
				}
				i := int(f64)
				vf.Set(reflect.ValueOf(&i))
				continue
			case *bool:
				b, ok := a.(bool)
				if !ok {
					break
				}
				vf.Set(reflect.ValueOf(&b))
				continue
			case *float64:
				f64, ok := a.(float64)
				if !ok {
					break
				}
				vf.Set(reflect.ValueOf(&f64))
				continue
			case *[]byte:
				b, ok := a.([]byte)
				if !ok {
					break
				}
				vf.Set(reflect.ValueOf(&b))
			}
			if !vf.CanAddr() {
				continue
			}
			switch ptr := vf.Addr().Interface().(type) {
			case sql.Scanner:
				if err := ptr.Scan(a); err != nil {
					return err
				}
			}
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
