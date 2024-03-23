package dynmgrm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var errValueInCompatible = fmt.Errorf("value is incompatible")

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
	case Sets[string]:
		return &types.AttributeValueMemberSS{Value: value}, nil
	case Sets[int]:
		ss := make([]string, 0, len(value))
		for _, v := range value {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
		return &types.AttributeValueMemberNS{Value: ss}, nil
	case Sets[float64]:
		ss := make([]string, 0, len(value))
		for _, v := range value {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
		return &types.AttributeValueMemberNS{Value: ss}, nil
	case Sets[[]byte]:
		return &types.AttributeValueMemberBS{Value: value}, nil
	default:
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
	return nil, errValueInCompatible
}
