package types

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
	return l.ResolveNestedDocument()
}

// ResolveNestedDocument resolves nested document type attribute.
func (l *List) ResolveNestedDocument() error {
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
		if s := newSets[int](); s.IsCompatible(v) {
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if s := newSets[float64](); s.IsCompatible(v) {
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if s := newSets[string](); s.IsCompatible(v) {
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if s := newSets[[]byte](); s.IsCompatible(v) {
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
