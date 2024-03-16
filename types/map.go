package types

// Map is a DynamoDB map type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Map map[string]interface{}

// GormDataType returns the data type for Gorm.
func (m *Map) GormDataType() string {
	return "dgmap"
}

// Scan implements the sql.Scanner#Scan
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
	return m.ResolveNestedDocument()
}

// ResolveNestedDocument resolves nested document type attribute.
func (m *Map) ResolveNestedDocument() error {
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
		if s := newSets[int](); s.IsCompatible(v) {
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if s := newSets[float64](); s.IsCompatible(v) {
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if s := newSets[string](); s.IsCompatible(v) {
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if s := newSets[[]byte](); s.IsCompatible(v) {
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
