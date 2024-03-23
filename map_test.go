package dynmgrm

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
)

func TestMap_GormDataType(t *testing.T) {
	type testCase struct {
		sut  Map
		want string
	}
	tests := map[string]testCase{
		"happy-path": {
			sut:  Map{},
			want: "dgmap",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.sut.GormDataType(); got != tt.want {
				t.Errorf("GormDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_Scan(t *testing.T) {
	type testCase struct {
		sut           Map
		args          interface{}
		want          error
		expectedState Map
	}
	tests := map[string]testCase{
		"happy-path/empty-map": {
			sut:           Map{},
			args:          map[string]interface{}{},
			expectedState: Map{},
		},
		"happy-path/single-value": {
			sut:           Map{},
			args:          map[string]interface{}{"a": 1},
			expectedState: Map{"a": 1},
		},
		"happy-path/multiple-values": {
			sut:           Map{},
			args:          map[string]interface{}{"a": 1, "b": "2"},
			expectedState: Map{"a": 1, "b": "2"},
		},
		"unhappy-path/non-map-value": {
			sut:  Map{},
			args: "non-map",
			want: ErrFailedToCast,
		},
		"unhappy-path/sut-is-not-empty": {
			sut:           Map{"existing": 1},
			args:          map[string]interface{}{"a": 1},
			want:          ErrCollectionAlreadyContainsItem,
			expectedState: Map{"existing": 1},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestMap_ResolveNestedCollections(t *testing.T) {
	type testCase struct {
		sut           Map
		expectedState Map
		want          error
	}
	tests := map[string]testCase{
		"happy-path/empty-map": {
			sut:           Map{},
			expectedState: Map{},
		},
		"happy-path/single-nested-map": {
			sut:           Map{"a": map[string]interface{}{"b": 1}},
			expectedState: Map{"a": Map{"b": 1}},
		},
		"happy-path/multiple-nested-maps": {
			sut:           Map{"a": map[string]interface{}{"b": 1}, "c": map[string]interface{}{"d": 2}},
			expectedState: Map{"a": Map{"b": 1}, "c": Map{"d": 2}},
		},
		"happy-path/nested-list": {
			sut:           Map{"a": []interface{}{1, "b"}},
			expectedState: Map{"a": List{1, "b"}},
		},
		"happy-path/nested-int-sets": {
			sut:           Map{"a": []float64{1, 2, 3}},
			expectedState: Map{"a": Sets[int]{1, 2, 3}},
		},
		"happy-path/nested-float-sets": {
			sut:           Map{"a": []float64{1.1, 2.1, 3.1}},
			expectedState: Map{"a": Sets[float64]{1.1, 2.1, 3.1}},
		},
		"happy-path/nested-string-sets": {
			sut:           Map{"a": []string{"1", "2", "3"}},
			expectedState: Map{"a": Sets[string]{"1", "2", "3"}},
		},
		"happy-path/nested-binary-sets": {
			sut:           Map{"a": [][]byte{[]byte("1"), []byte("2"), []byte("3")}},
			expectedState: Map{"a": Sets[[]byte]{[]byte("1"), []byte("2"), []byte("3")}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := resolveCollectionsNestedInMap(&tt.sut)
			if !errors.Is(err, tt.want) {
				t.Errorf("ResolveNestedDocument() error = %v, want %v", err, tt.want)
				return

			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("ResolveNestedDocument() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestMap_GormValue(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           Map
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: Map{"a": 1},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
						"a": &types.AttributeValueMemberN{Value: "1"},
					}},
				}},
		},
	}
	opts := []cmp.Option{
		cmp.AllowUnexported(types.AttributeValueMemberS{}),
		cmp.AllowUnexported(types.AttributeValueMemberSS{}),
		cmp.AllowUnexported(types.AttributeValueMemberN{}),
		cmp.AllowUnexported(types.AttributeValueMemberNS{}),
		cmp.AllowUnexported(types.AttributeValueMemberB{}),
		cmp.AllowUnexported(types.AttributeValueMemberBS{}),
		cmp.AllowUnexported(types.AttributeValueMemberL{}),
		cmp.AllowUnexported(types.AttributeValueMemberM{}),
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.sut.GormValue(tt.args.ctx, tt.args.db)
			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("GormValue() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
