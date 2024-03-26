package dynmgrm

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestList_GormDataType(t *testing.T) {
	l := &List{}
	if got := l.GormDataType(); got != "dglist" {
		t.Errorf("GormDataType() = %v, want %v", got, "dglist")
	}
}

func TestList_Scan(t *testing.T) {
	type testCase struct {
		sut           List
		args          interface{}
		expectedState List
		want          error
	}
	tests := map[string]testCase{
		"happy-path/empty-list": {
			sut:           List{},
			args:          []interface{}{},
			expectedState: List{},
		},
		"happy-path/single-value": {
			sut:           List{},
			args:          []interface{}{1},
			expectedState: List{1},
		},
		"happy-path/multiple-values": {
			sut:           List{},
			args:          []interface{}{1, "2"},
			expectedState: List{1, "2"},
		},
		"happy-path/single-nested-map": {
			sut:           List{},
			args:          []interface{}{map[string]interface{}{"a": 1}},
			expectedState: List{Map{"a": 1}},
		},
		"unhappy-path/sut-is-not-empty": {
			sut:           List{1, 2, 3},
			args:          []interface{}{4, 5, 6},
			expectedState: List{1, 2, 3},
			want:          ErrCollectionAlreadyContainsItem,
		},
		"unhappy-path/non-slice-value": {
			sut:           List{},
			args:          "non-slice",
			expectedState: List{},
			want:          ErrFailedToCast,
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

func TestList_ResolveNestedCollections(t *testing.T) {
	type testCase struct {
		sut           List
		expectedState List
		want          error
	}
	tests := map[string]testCase{
		"happy-path/empty-list": {
			sut:           List{},
			expectedState: List{},
		},
		"happy-path/single-nested-map": {
			sut:           List{map[string]interface{}{"a": 1}},
			expectedState: List{Map{"a": 1}},
		},
		"happy-path/multiple-nested-maps": {
			sut:           List{map[string]interface{}{"a": 1}, map[string]interface{}{"b": 2}},
			expectedState: List{Map{"a": 1}, Map{"b": 2}},
		},
		"happy-path/nested-list": {
			sut:           List{[]interface{}{1, "b"}},
			expectedState: List{List{1, "b"}},
		},
		"happy-path/nested-int-sets": {
			sut:           List{[]float64{1, 2, 3}},
			expectedState: List{Set[int]{1, 2, 3}},
		},
		"happy-path/nested-float-sets": {
			sut:           List{[]float64{1.1, 2.1, 3.1}},
			expectedState: List{Set[float64]{1.1, 2.1, 3.1}},
		},
		"happy-path/nested-string-sets": {
			sut:           List{[]string{"1", "2", "3"}},
			expectedState: List{Set[string]{"1", "2", "3"}},
		},
		"happy-path/nested-binary-sets": {
			sut:           List{[][]byte{[]byte("1"), []byte("2"), []byte("3")}},
			expectedState: List{Set[[]byte]{[]byte("1"), []byte("2"), []byte("3")}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := resolveCollectionsNestedInList(&tt.sut)
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

func TestList_GormValue(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           List
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: List{1},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberL{
						Value: []types.AttributeValue{
							&types.AttributeValueMemberN{Value: "1"},
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
