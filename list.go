package dynmgrm

import (
	"errors"
	"fmt"
)

// List is a DynamoDB list type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type List []interface{}

// GormDataType returns the data type for Gorm.
func (l *List) GormDataType() string {
	return "dglist"
}

// Scan implements the sql.Scanner#Scan
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
		if isCompatibleWithSets[int](v) {
			s := newSets[int]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSets[float64](v) {
			s := newSets[float64]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSets[string](v) {
			s := newSets[string]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSets[[]byte](v) {
			s := newSets[[]byte]()
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
