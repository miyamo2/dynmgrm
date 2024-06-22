package dynmgrm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"reflect"
	"slices"
)

var (
	ErrValueIsIncompatibleOfStringSlice  = errors.New("value is incompatible of string slice")
	ErrValueIsIncompatibleOfIntSlice     = errors.New("value is incompatible of int slice")
	ErrValueIsIncompatibleOfFloat64Slice = errors.New("value is incompatible of float64 slice")
	ErrValueIsIncompatibleOfBinarySlice  = errors.New("value is incompatible of []byte slice")
	ErrCollectionAlreadyContainsItem     = errors.New("collection already contains item")
	ErrFailedToCast                      = errors.New("failed to cast")
)

// SetSupportable are the types that support the Set
type SetSupportable interface {
	string | []byte | int | float64
}

// compatibility check
var (
	_ gorm.Valuer = (*Set[string])(nil)
	_ sql.Scanner = (*Set[string])(nil)
)

// Set is a DynamoDB set type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Set[T SetSupportable] []T

// GormDataType returns the data type for Gorm.
func (s *Set[T]) GormDataType() string {
	return "dgsets"
}

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (s *Set[T]) Scan(value interface{}) error {
	if len(*s) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	if value == nil {
		*s = nil
		return nil
	}
	switch (interface{})(s).(type) {
	case *Set[int]:
		return scanAsIntSet((interface{})(s).(*Set[int]), value)
	case *Set[float64]:
		return scanAsFloat64Set((interface{})(s).(*Set[float64]), value)
	case *Set[string]:
		return scanAsStringSet((interface{})(s).(*Set[string]), value)
	case *Set[[]byte]:
		return scanAsBinarySet((interface{})(s).(*Set[[]byte]), value)
	}
	return nil
}

// GormValue implements the [gorm.Valuer] interface.
//
// [gorm.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (s Set[T]) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	switch s := (interface{})(s).(type) {
	case Set[int]:
		av, err := numericSetToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	case Set[float64]:
		av, err := numericSetToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	case Set[string]:
		av, err := stringSetToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	case Set[[]byte]:
		av, err := binarySetToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	}
	return clause.Expr{}
}

func numericSetToAttributeValue[T Set[int] | Set[float64]](s T) (*types.AttributeValueMemberNS, error) {
	return toDocumentAttributeValue[*types.AttributeValueMemberNS](s)
}

func stringSetToAttributeValue(s Set[string]) (*types.AttributeValueMemberSS, error) {
	return toDocumentAttributeValue[*types.AttributeValueMemberSS](s)
}

func binarySetToAttributeValue(s Set[[]byte]) (*types.AttributeValueMemberBS, error) {
	return toDocumentAttributeValue[*types.AttributeValueMemberBS](s)
}

// scanAsIntSet scans the value as Set[int]
func scanAsIntSet(s *Set[int], value interface{}) error {
	sv, ok := value.([]float64)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfIntSlice
	}
	for _, v := range sv {
		if math.Floor(v) != v {
			*s = nil
			return ErrValueIsIncompatibleOfIntSlice
		}
		*s = append(*s, int(v))
	}
	return nil
}

// scanAsFloat64Set scans the value as Set[float64]
func scanAsFloat64Set(s *Set[float64], value interface{}) error {
	sv, ok := value.([]float64)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfFloat64Slice
	}
	for _, v := range sv {
		*s = append(*s, v)
	}
	return nil
}

// scanAsStringSet scans the value as Set[string]
func scanAsStringSet(s *Set[string], value interface{}) error {
	sv, ok := value.([]string)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfStringSlice
	}
	for _, v := range sv {
		*s = append(*s, v)
	}
	return nil
}

// scanAsBinarySet scans the value as Set[[]byte]
func scanAsBinarySet(s *Set[[]byte], value interface{}) error {
	sv, ok := value.([][]byte)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfBinarySlice
	}
	for _, v := range sv {
		*s = append(*s, v)
	}
	return nil
}

func isCompatibleWithSet[T SetSupportable](value interface{}) (compatible bool) {
	var t T
	switch (interface{})(t).(type) {
	case string:
		compatible = isStringSetCompatible(value)
	case int:
		compatible = isIntSetCompatible(value)
	case float64:
		compatible = isFloat64SetCompatible(value)
	case []byte:
		compatible = isBinarySetCompatible(value)
	}
	return
}

func isIntSetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]int); ok {
		compatible = true
		return
	}
	if value, ok := value.([]float64); ok {
		compatible = true
		for _, v := range value {
			if math.Floor(v) == v {
				compatible = true
				continue
			}
			compatible = false
			return
		}
	}
	return
}

func isStringSetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]string); ok {
		compatible = true
	}
	return
}

func isFloat64SetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]float64); ok {
		compatible = true
	}
	return
}

func isBinarySetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([][]byte); ok {
		compatible = true
	}
	return
}

func newSet[T SetSupportable]() Set[T] {
	return Set[T]{}
}

// compatibility check
var (
	_ gorm.Valuer = (*List)(nil)
	_ sql.Scanner = (*List)(nil)
)

// List is a DynamoDB list type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type List []interface{}

// GormDataType returns the data type for Gorm.
func (l *List) GormDataType() string {
	return "dglist"
}

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (l *List) Scan(value interface{}) error {
	if len(*l) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	sv, ok := value.([]interface{})
	if !ok {
		return errors.Join(ErrFailedToCast, fmt.Errorf("incompatible %T and %T", l, value))
	}
	*l = sv
	return resolveCollectionsNestedInList(l)
}

// GormValue implements the [gorm.Valuer] interface.
//
// [gorm.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (l List) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	if err := resolveCollectionsNestedInList(&l); err != nil {
		_ = db.AddError(err)
		return clause.Expr{}
	}
	av, err := toDocumentAttributeValue[*types.AttributeValueMemberL](l)
	if err != nil {
		_ = db.AddError(err)
		return clause.Expr{}
	}
	return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
}

// resolveCollectionsNestedInList resolves nested collection type attribute.
func resolveCollectionsNestedInList(l *List) error {
	for i, v := range *l {
		if v, ok := v.(map[string]interface{}); ok {
			m := Map{}
			err := m.Scan(v)
			if err != nil {
				*l = nil
				return err
			}
			(*l)[i] = m
			continue
		}
		if isCompatibleWithSet[int](v) {
			s := newSet[int]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSet[float64](v) {
			s := newSet[float64]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSet[string](v) {
			s := newSet[string]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSet[[]byte](v) {
			s := newSet[[]byte]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if v, ok := v.([]interface{}); ok {
			il := List{}
			err := il.Scan(v)
			if err != nil {
				*l = nil
				return err
			}
			(*l)[i] = il
		}
	}
	return nil
}

// compatibility check
var (
	_ gorm.Valuer = (*Map)(nil)
	_ sql.Scanner = (*Map)(nil)
)

// Map is a DynamoDB map type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Map map[string]interface{}

// GormDataType returns the data type for Gorm.
func (m Map) GormDataType() string {
	return "dgmap"
}

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (m *Map) Scan(value interface{}) error {
	if len(*m) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	mv, ok := value.(map[string]interface{})
	if !ok {
		*m = nil
		return ErrFailedToCast
	}
	*m = mv
	return resolveCollectionsNestedInMap(m)
}

// GormValue implements the [gorm.Valuer] interface.
//
// [gorm.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (m Map) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	if err := resolveCollectionsNestedInMap(&m); err != nil {
		_ = db.AddError(err)
		return clause.Expr{}
	}
	av, err := toDocumentAttributeValue[*types.AttributeValueMemberM](m)
	if err != nil {
		_ = db.AddError(err)
		return clause.Expr{}
	}
	return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
}

// resolveCollectionsNestedInMap resolves nested document type attribute.
func resolveCollectionsNestedInMap(m *Map) error {
	for k, v := range *m {
		if v, ok := v.(map[string]interface{}); ok {
			im := Map{}
			err := im.Scan(v)
			if err != nil {
				*m = nil
				return err
			}
			(*m)[k] = im
			continue
		}
		if isCompatibleWithSet[int](v) {
			s := newSet[int]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSet[float64](v) {
			s := newSet[float64]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSet[string](v) {
			s := newSet[string]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSet[[]byte](v) {
			s := newSet[[]byte]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if v, ok := v.([]interface{}); ok {
			l := List{}
			err := l.Scan(v)
			if err != nil {
				*m = nil
				return err
			}
			(*m)[k] = l
		}
	}
	return nil
}

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
