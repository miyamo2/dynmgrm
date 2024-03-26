package dynmgrm

import (
	"context"
	"errors"
	"math"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrValueIsIncompatibleOfStringSlice  = errors.New("value is incompatible of string slice")
	ErrValueIsIncompatibleOfIntSlice     = errors.New("value is incompatible of int slice")
	ErrValueIsIncompatibleOfFloat64Slice = errors.New("value is incompatible of float64 slice")
	ErrValueIsIncompatibleOfBinarySlice  = errors.New("value is incompatible of []byte slice")
)

// SetSupportable are the types that support the Set
type SetSupportable interface {
	string | []byte | int | float64
}

var _ gorm.Valuer = (*Set[string])(nil)

// Set is a DynamoDB set type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Set[T SetSupportable] []T

// GormDataType returns the data type for Gorm.
func (s *Set[T]) GormDataType() string {
	return "dgsets"
}

// Scan implements the sql.Scanner#Scan
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

// GormValue implements the gorm.Valuer interface.
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
