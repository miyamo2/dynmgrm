package dynmgrm

import (
	"context"
	"errors"
	"gorm.io/gorm/schema"
	"reflect"
)

var ErrIncompatibleNestedStruct = errors.New("incompatible nested struct")

var (
	_ schema.SerializerInterface       = (*nestedStructSerializer)(nil)
	_ schema.SerializerValuerInterface = (*nestedStructSerializer)(nil)
)

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

func (n nestedStructSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	schema.RegisterSerializer("nested", nestedStructSerializer{})
}
