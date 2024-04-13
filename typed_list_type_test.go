package dynmgrm

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
)

type Something struct {
	A string
	B int
	C bool
	D float64
	E *string
	F *int
	G *bool
	H *float64
	I Set[string]
}

func TestTypedList_Scan(t *testing.T) {
	type testCase struct {
		sut           TypedList[Something]
		args          interface{}
		expectedState TypedList[Something]
		want          error
	}
	tests := map[string]testCase{
		"happy-path/empty-list": {
			sut:           TypedList[Something]{},
			args:          []interface{}{},
			expectedState: TypedList[Something]{},
		},
		"happy-path/single-value": {
			sut: TypedList[Something]{},
			args: []interface{}{map[string]interface{}{
				"A": "foo",
				"B": 1,
				"C": true,
				"D": 1.1,
				"E": "foo",
				"F": 1,
				"G": true,
				"H": 1.1,
				"I": []string{"foo", "bar", "baz"},
			}},
			expectedState: TypedList[Something]{
				Something{
					A: "foo",
					B: 1,
					C: true,
					D: 1.1,
					E: aws.String("foo"),
					F: aws.Int(1),
					G: aws.Bool(true),
					H: aws.Float64(1.1),
					I: Set[string]{"foo", "bar", "baz"},
				},
			},
		},
		"happy-path/with-nil-attribute": {
			sut: TypedList[Something]{},
			args: []interface{}{map[string]interface{}{
				"A": "foo",
				"B": 1,
				"C": true,
				"D": 1.1,
				"E": nil,
				"F": nil,
				"G": nil,
				"H": nil,
				"I": []string{"foo", "bar", "baz"},
			}},
			expectedState: TypedList[Something]{
				Something{
					A: "foo",
					B: 1,
					C: true,
					D: 1.1,
					E: nil,
					F: nil,
					G: nil,
					H: nil,
					I: Set[string]{"foo", "bar", "baz"},
				},
			},
		},
		"happy-path/multiple-values": {
			sut: TypedList[Something]{},
			args: []interface{}{
				map[string]interface{}{
					"A": "foo",
					"B": 1,
					"C": true,
					"D": 1.1,
					"E": "foo",
					"F": 1,
					"G": true,
					"H": 1.1,
					"I": []string{"foo", "bar", "baz"},
				},
				map[string]interface{}{
					"A": "bar",
					"B": 2,
					"C": false,
					"D": 2.2,
					"E": "bar",
					"F": 2,
					"G": false,
					"H": 2.2,
					"I": []string{"foo", "bar", "baz"},
				},
			},
			expectedState: TypedList[Something]{
				Something{
					A: "foo",
					B: 1,
					C: true,
					D: 1.1,
					E: aws.String("foo"),
					F: aws.Int(1),
					G: aws.Bool(true),
					H: aws.Float64(1.1),
					I: Set[string]{"foo", "bar", "baz"},
				},
				Something{
					A: "bar",
					B: 2,
					C: false,
					D: 2.2,
					E: aws.String("bar"),
					F: aws.Int(2),
					G: aws.Bool(false),
					H: aws.Float64(2.2),
					I: Set[string]{"foo", "bar", "baz"},
				},
			},
		},
		"unhappy-path/non-slice-value": {
			sut:           TypedList[Something]{},
			args:          "non-slice",
			expectedState: TypedList[Something]{},
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

func TestTypedList_GormValue(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           TypedList[Something]
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: TypedList[Something]{
				{
					A: "foo",
					B: 1,
					C: true,
					D: 1.1,
					E: aws.String("foo"),
					F: aws.Int(1),
					G: aws.Bool(true),
					H: aws.Float64(1.1),
					I: Set[string]{"foo", "bar", "baz"},
				},
			},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberL{
						Value: []types.AttributeValue{
							&types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"A": &types.AttributeValueMemberS{Value: "foo"},
									"B": &types.AttributeValueMemberN{Value: "1"},
									"C": &types.AttributeValueMemberBOOL{Value: true},
									"D": &types.AttributeValueMemberN{Value: "1.1"},
									"E": &types.AttributeValueMemberS{Value: "foo"},
									"F": &types.AttributeValueMemberN{Value: "1"},
									"G": &types.AttributeValueMemberBOOL{Value: true},
									"H": &types.AttributeValueMemberN{Value: "1.1"},
									"I": &types.AttributeValueMemberSS{Value: []string{"foo", "bar", "baz"}},
								},
							},
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
		cmp.AllowUnexported(types.AttributeValueMemberBOOL{}),
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
