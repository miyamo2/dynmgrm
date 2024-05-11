package dynmgrm

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm/schema"
	"reflect"
)

// ErrIncompatibleNestedStruct occurs when an incompatible with nested-struct.
var ErrIncompatibleNestedStruct = errors.New("incompatible nested struct")

var (
	_ schema.SerializerInterface       = (*nestedStructSerializer)(nil)
	_ schema.SerializerValuerInterface = (*nestedStructSerializer)(nil)
)

// nestedStructSerializer is a serializer for nested struct.
type nestedStructSerializer struct{}

func (n nestedStructSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		return nil
	}
	fieldValue := reflect.New(field.FieldType)
	switch dbValue := dbValue.(type) {
	case map[string]interface{}:
		err := assignMapValueToReflectValue(field.FieldType, fieldValue, dbValue)
		if err != nil {
			return err
		}
	default:
		return ErrIncompatibleNestedStruct
	}
	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return nil
}

func (n nestedStructSerializer) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	av, err := toDocumentAttributeValue[*types.AttributeValueMemberM](fieldValue)
	if av == (*types.AttributeValueMemberM)(nil) {
		return nil, err
	}
	return *av, err
}

func init() {
	schema.RegisterSerializer("nested", nestedStructSerializer{})
}
