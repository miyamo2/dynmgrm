package dynmgrm

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// errDocumentAttributeValueIsIncompatible occurs when an incompatible conversion to following:
//   - *types.AttributeValueMemberL
//   - *types.AttributeValueMemberM
//   - *types.AttributeValueMemberSS
//   - *types.AttributeValueMemberNS
//   - *types.AttributeValueMemberBS
var errDocumentAttributeValueIsIncompatible = fmt.Errorf("document-attribute-value is incompatible")

func toAttibuteValue(value interface{}) (types.AttributeValue, error) {
	switch value := value.(type) {
	case List:
		avs := make([]types.AttributeValue, 0, len(value))
		for _, v := range value {
			av, err := toAttibuteValue(v)
			if err != nil {
				return nil, err
			}
			avs = append(avs, av)
		}
		return &types.AttributeValueMemberL{Value: avs}, nil
	case Map:
		avm := make(map[string]types.AttributeValue)
		for k, v := range value {
			av, err := toAttibuteValue(v)
			if err != nil {
				return nil, err
			}
			avm[k] = av
		}
		return &types.AttributeValueMemberM{Value: avm}, nil
	case Set[string]:
		return &types.AttributeValueMemberSS{Value: value}, nil
	case Set[int]:
		ss := make([]string, 0, len(value))
		for _, v := range value {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
		return &types.AttributeValueMemberNS{Value: ss}, nil
	case Set[float64]:
		ss := make([]string, 0, len(value))
		for _, v := range value {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
		return &types.AttributeValueMemberNS{Value: ss}, nil
	case Set[[]byte]:
		return &types.AttributeValueMemberBS{Value: value}, nil
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Struct:
			avm := make(map[string]types.AttributeValue)
			for i := 0; i < rv.NumField(); i++ {
				fv := rv.Field(i)
				ft := rv.Type().Field(i)
				if fv.CanInterface() {
					av, err := toAttibuteValue(fv.Interface())
					if err != nil {
						return nil, err
					}
					avm[getDBNameFromStructField(ft)] = av
				} else if fv.CanAddr() {
					av, err := toAttibuteValue(fv.Addr().Interface())
					if err != nil {
						return nil, err
					}
					avm[getDBNameFromStructField(ft)] = av
				}
			}
			return &types.AttributeValueMemberM{Value: avm}, nil
		case reflect.Ptr:
			if rv.IsNil() {
				return &types.AttributeValueMemberNULL{}, nil
			}
			if !rv.CanAddr() {
				return attributevalue.Marshal(value)
			}
			return toAttibuteValue(rv.Addr().Interface())
		}
		return attributevalue.Marshal(value)
	}
}

type documentAttributeMember interface {
	*types.AttributeValueMemberL | *types.AttributeValueMemberM | *types.AttributeValueMemberSS | *types.AttributeValueMemberNS | *types.AttributeValueMemberBS
}

func toDocumentAttributeValue[T documentAttributeMember](value interface{}) (T, error) {
	v, err := toAttibuteValue(value)
	if err != nil {
		return nil, err
	}
	if v, ok := v.(T); ok {
		return v, nil
	}
	return nil, errDocumentAttributeValueIsIncompatible
}
