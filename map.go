package dynmgrm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ gorm.Valuer = (*Map)(nil)

// Map is a DynamoDB map type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Map map[string]interface{}

// GormDataType returns the data type for Gorm.
func (m Map) GormDataType() string {
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
	return resolveCollectionsNestedInMap(m)
}

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
		if isCompatibleWithSets[int](v) {
			s := newSets[int]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSets[float64](v) {
			s := newSets[float64]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSets[string](v) {
			s := newSets[string]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSets[[]byte](v) {
			s := newSets[[]byte]()
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
