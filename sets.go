package dynmgrm

import (
	"context"
	"errors"
	"fmt"
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

// SetsSupportable are the types that support the Set
type SetsSupportable interface {
	string | []byte | int | float64
}

var _ gorm.Valuer = (*Sets[string])(nil)

// Sets is a set of values.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Sets[T SetsSupportable] []T

// GormDataType returns the data type for Gorm.
func (s *Sets[T]) GormDataType() string {
	return "dgsets"
}

// Scan implements the sql.Scanner#Scan
func (s *Sets[T]) Scan(value interface{}) error {
	if len(*s) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	if value == nil {
		*s = nil
		return nil
	}
	var t T
	switch (interface{})(t).(type) {
	case int:
		return scanAsIntSets((interface{})(s).(*Sets[int]), value)
	case float64:
		return scanAsFloat64Sets((interface{})(s).(*Sets[float64]), value)
	case string:
		return scanAsStringSets((interface{})(s).(*Sets[string]), value)
	case []byte:
		return scanAsBinarySets((interface{})(s).(*Sets[[]byte]), value)
	default:
		// never happens (now).
		return fmt.Errorf(
			"unsupported type %T. Sets supports only the following types: string, []byte, int, float32, float64", t)
	}
}

func (s Sets[T]) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	switch s := (interface{})(s).(type) {
	case Sets[int]:
		av, err := numericSetsToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	case Sets[float64]:
		av, err := numericSetsToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	case Sets[string]:
		av, err := stringSetsToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	case Sets[[]byte]:
		av, err := binarySetsToAttributeValue(s)
		if err != nil {
			_ = db.AddError(err)
			break
		}
		return clause.Expr{SQL: "?", Vars: []interface{}{*av}}
	}
	return clause.Expr{}
}

func numericSetsToAttributeValue[T Sets[int] | Sets[float64]](s T) (*types.AttributeValueMemberNS, error) {
	return toDocumentAttributeValue[*types.AttributeValueMemberNS](s)
}

func stringSetsToAttributeValue(s Sets[string]) (*types.AttributeValueMemberSS, error) {
	return toDocumentAttributeValue[*types.AttributeValueMemberSS](s)
}

func binarySetsToAttributeValue(s Sets[[]byte]) (*types.AttributeValueMemberBS, error) {
	return toDocumentAttributeValue[*types.AttributeValueMemberBS](s)
}

// scanAsIntSets scans the value as Sets[int]
func scanAsIntSets(s *Sets[int], value interface{}) error {
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

// scanAsFloat64Sets scans the value as Sets[float64]
func scanAsFloat64Sets(s *Sets[float64], value interface{}) error {
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

// scanAsStringSets scans the value as Sets[string]
func scanAsStringSets(s *Sets[string], value interface{}) error {
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

// scanAsBinarySets scans the value as Sets[[]byte]
func scanAsBinarySets(s *Sets[[]byte], value interface{}) error {
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

func isCompatibleWithSets[T SetsSupportable](value interface{}) (compatible bool) {
	var t T
	switch (interface{})(t).(type) {
	case string:
		compatible = isStringSetsCompatible(value)
	case int:
		compatible = isIntSetsCompatible(value)
	case float64:
		compatible = isFloat64SetsCompatible(value)
	case []byte:
		compatible = isBinarySetsCompatible(value)
	}
	return
}

func isIntSetsCompatible(value interface{}) (compatible bool) {
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

func isStringSetsCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]string); ok {
		compatible = true
	}
	return
}

func isFloat64SetsCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]float64); ok {
		compatible = true
	}
	return
}

func isBinarySetsCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([][]byte); ok {
		compatible = true
	}
	return
}

func newSets[T SetsSupportable]() Sets[T] {
	return Sets[T]{}
}
