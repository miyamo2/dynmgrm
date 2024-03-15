package types

import (
	"errors"
	"fmt"
)

var (
	ErrValueIsIncompatibleOfInterfaceSlice = errors.New("value is incompatible of interface slice")
	ErrValueIsIncompatibleOfIntSlice       = errors.New("value is incompatible of int slice")
	ErrValueIsIncompatibleOfFloat64Slice   = errors.New("value is incompatible of float64 slice")
	ErrValueIsIncompatibleOfBinarySlice    = errors.New("value is incompatible of []byte slice")
)

// SetsSupportable are the types that support the Set
type SetsSupportable interface {
	string | []byte | int | float64
}

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

// scanAsIntSets scans the value as Sets[int]
func scanAsIntSets(s *Sets[int], value interface{}) error {
	sv, ok := value.([]interface{})
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfInterfaceSlice
	}
	for _, v := range sv {
		cv, ok := v.(int)
		if !ok {
			switch v := v.(type) {
			case float32:
				cv = int(v)
				ok = true
			case float64:
				cv = int(v)
				ok = true
			}
		}
		if !ok {
			*s = nil
			return ErrValueIsIncompatibleOfIntSlice
		}
		*s = append(*s, cv)
	}
	return nil
}

// scanAsFloat64Sets scans the value as Sets[float64]
func scanAsFloat64Sets(s *Sets[float64], value interface{}) error {
	sv, ok := value.([]interface{})
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfInterfaceSlice
	}
	for _, v := range sv {
		cv, ok := v.(float64)
		if !ok {
			switch v := v.(type) {
			case int:
				cv = float64(v)
				ok = true
			}
		}
		if !ok {
			*s = nil
			return ErrValueIsIncompatibleOfFloat64Slice
		}
		*s = append(*s, cv)
	}
	return nil
}

// scanAsStringSets scans the value as Sets[string]
func scanAsStringSets(s *Sets[string], value interface{}) error {
	sv, ok := value.([]interface{})
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfInterfaceSlice
	}
	for _, v := range sv {
		cv, ok := v.(string)
		if !ok {
			cv = fmt.Sprintf("%v", v)
		}
		*s = append(*s, cv)
	}
	return nil
}

// scanAsBinarySets scans the value as Sets[[]byte]
func scanAsBinarySets(s *Sets[[]byte], value interface{}) error {
	sv, ok := value.([]interface{})
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfInterfaceSlice
	}
	for _, v := range sv {
		cv, ok := v.([]byte)
		if !ok {
			switch v := v.(type) {
			case string:
				cv = []byte(v)
				ok = true
			default:
				*s = nil
				return ErrValueIsIncompatibleOfBinarySlice
			}
		}
		*s = append(*s, cv)
	}
	return nil
}

func (s Sets[T]) IsCompatible(value interface{}) (compatible bool) {
	sValue, ok := value.([]interface{})
	if !ok {
		return
	}
	var t T
	switch (interface{})(t).(type) {
	case string:
		compatible = true
	case int:
		compatible = isIntSetsCompatible(sValue)
	case float64:
		compatible = isFloat63SetsCompatible(sValue)
	case []byte:
		compatible = isBinarySetsCompatible(sValue)
	}
	return
}

func isIntSetsCompatible(value []interface{}) (compatible bool) {
	compatible = true
	for _, v := range value {
		if _, ok := v.(int); ok {
			continue
		}
		switch v.(type) {
		case float32, float64:
			continue
		default:
			compatible = false
			return
		}
	}
	return true
}

func isFloat63SetsCompatible(value []interface{}) (compatible bool) {
	compatible = true
	for _, v := range value {
		if _, ok := v.(float64); ok {
			continue
		}
		switch v.(type) {
		case float32, float64, int:
			continue
		default:
			compatible = false
			return
		}
	}
	return true
}

func isBinarySetsCompatible(value []interface{}) (compatible bool) {
	compatible = true
	for _, v := range value {
		if _, ok := v.([]byte); ok {
			continue
		}
		compatible = false
		return
	}
	return
}
