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

func TestSet_Scan(t *testing.T) {
	type testCase struct {
		sut           Set[string]
		args          interface{}
		want          error
		expectedState Set[string]
	}
	tests := map[string]testCase{
		"unhappy-path/already-contains-item": {
			sut:           Set[string]{"value"},
			args:          []interface{}{"test"},
			expectedState: Set[string]{"value"},
			want:          ErrCollectionAlreadyContainsItem,
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

func TestSet_Scan_String(t *testing.T) {
	type testCase struct {
		sut           Set[string]
		args          interface{}
		want          error
		expectedState Set[string]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           newSet[string](),
			args:          []string{"test"},
			expectedState: Set[string]{"test"},
		},
		"happy-path/multiple-values": {
			sut:           newSet[string](),
			args:          []string{"test1", "test2"},
			expectedState: Set[string]{"test1", "test2"},
		},
		"happy-path/numeric-value": {
			sut:           newSet[string](),
			args:          []string{"1.1", "1.2"},
			expectedState: Set[string]{"1.1", "1.2"},
		},
		"happy-path/null": {
			sut:           newSet[string](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           newSet[string](),
			args:          "test",
			want:          ErrValueIsIncompatibleOfStringSlice,
			expectedState: nil,
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

func TestSet_Scan_Int(t *testing.T) {
	type testCase struct {
		sut           Set[int]
		args          interface{}
		want          error
		expectedState Set[int]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           newSet[int](),
			args:          []float64{1},
			expectedState: Set[int]{1},
		},
		"happy-path/multiple-values": {
			sut:           newSet[int](),
			args:          []float64{1, 2},
			expectedState: Set[int]{1, 2},
		},
		"happy-path/float-value": {
			sut:           newSet[int](),
			args:          []float64{1.0, 2.0},
			expectedState: Set[int]{1, 2},
		},
		"happy-path/null": {
			sut:           newSet[int](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           newSet[int](),
			args:          1,
			want:          ErrValueIsIncompatibleOfIntSlice,
			expectedState: nil,
		},
		"unhappy-path/string-slice": {
			sut:           newSet[int](),
			args:          []string{"A"},
			want:          ErrValueIsIncompatibleOfIntSlice,
			expectedState: nil,
		},
		"unhappy-path/float-value": {
			sut:           newSet[int](),
			args:          []float64{1.1, 2.1},
			want:          ErrValueIsIncompatibleOfIntSlice,
			expectedState: nil,
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

func TestSet_Scan_Float64(t *testing.T) {
	type testCase struct {
		sut           Set[float64]
		args          interface{}
		want          error
		expectedState Set[float64]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           newSet[float64](),
			args:          []float64{1.1},
			expectedState: Set[float64]{1.1},
		},
		"happy-path/multiple-values": {
			sut:           newSet[float64](),
			args:          []float64{1.1, 1.2},
			expectedState: Set[float64]{1.1, 1.2},
		},
		"happy-path/int-value": {
			sut:           newSet[float64](),
			args:          []float64{1, 2},
			expectedState: Set[float64]{1, 2},
		},
		"happy-path/int-and-float64": {
			sut:           newSet[float64](),
			args:          []float64{1, 1.2},
			expectedState: Set[float64]{1, 1.2},
		},
		"happy-path/null": {
			sut:           newSet[float64](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:  newSet[float64](),
			args: 1,
			want: ErrValueIsIncompatibleOfFloat64Slice,
		},
		"unhappy-path/not-compatible-slice": {
			sut:  newSet[float64](),
			args: []string{"A"},
			want: ErrValueIsIncompatibleOfFloat64Slice,
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

func TestSet_Scan_Binary(t *testing.T) {
	type testCase struct {
		sut           Set[[]byte]
		args          interface{}
		want          error
		expectedState Set[[]byte]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           newSet[[]byte](),
			args:          [][]byte{{116, 101, 115, 116, 49}},
			expectedState: Set[[]byte]{[]byte{116, 101, 115, 116, 49}},
		},
		"happy-path/multiple-value": {
			sut: newSet[[]byte](),
			args: [][]byte{
				{116, 101, 115, 116, 49},
				{116, 101, 115, 116, 50},
			},
			expectedState: Set[[]byte]{
				[]byte{116, 101, 115, 116, 49},
				[]byte{116, 101, 115, 116, 50},
			},
		},
		"happy-path/null": {
			sut:           newSet[[]byte](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:  newSet[[]byte](),
			args: 1,
			want: ErrValueIsIncompatibleOfBinarySlice,
		},
		"unhappy-path/not-compatible-slice": {
			sut:  newSet[[]byte](),
			args: []interface{}{1},
			want: ErrValueIsIncompatibleOfBinarySlice,
		},
		"unhappy-path/string-value": {
			sut:           newSet[[]byte](),
			args:          []interface{}{"test"},
			expectedState: nil,
			want:          ErrValueIsIncompatibleOfBinarySlice,
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

func TestSet_GormDataType(t *testing.T) {
	type testCase struct {
		sut  Set[string]
		want string
	}
	tests := map[string]testCase{
		"happy-path": {
			sut:  newSet[string](),
			want: "dgsets",
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

func TestSet_IsCompatibleWithSet_String(t *testing.T) {
	type testCase struct {
		sut  Set[string]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: []string{"test"},
			want: true,
		},
		"happy-path/int-value": {
			args: []int{1},
			want: false,
		},
		"happy-path/float64-value": {
			args: []float64{1.1},
			want: false,
		},
		"happy-path/float32-value": {
			args: []float32{1.1},
			want: false,
		},
		"happy-path/binary-value": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSet[string](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsCompatibleWithSet_Int(t *testing.T) {
	type testCase struct {
		sut  Set[int]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: []int{1},
			want: true,
		},
		"happy-path/string-value": {
			args: []string{"test"},
			want: false,
		},
		"happy-path/float64-value": {
			args: []float64{1.1},
			want: false,
		},
		"happy-path/float64-typed-int-value": {
			args: []float64{1},
			want: true,
		},
		"happy-path/binary-value": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSet[int](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsCompatibleWithSet_Float64(t *testing.T) {
	type testCase struct {
		sut  Set[float64]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: []float64{1.1},
			want: true,
		},
		"unhappy-path/string-value": {
			args: []string{"test"},
			want: false,
		},
		"unhappy-path/int-value": {
			args: []int{1},
			want: false,
		},
		"unhappy-path/float32-value": {
			args: []float32{1.1},
			want: false,
		},
		"unhappy-path/binary-value": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSet[float64](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsCompatibleWithSet_Binary(t *testing.T) {
	type testCase struct {
		sut  Set[[]byte]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: true,
		},
		"unhappy-path/string-value": {
			args: []string{"test"},
			want: false,
		},
		"happy-path/int-value": {
			args: []int{1},
			want: false,
		},
		"happy-path/float64-value": {
			args: []float64{1.1},
			want: false,
		},
		"happy-path/float32-value": {
			args: []float32{1.1},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSet[[]byte](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_GormValue_String(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           Set[string]
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: Set[string]{"1"},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberSS{Value: []string{"1"}},
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

func TestSet_GormValue_Int(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           Set[int]
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: Set[int]{1},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberNS{Value: []string{"1"}},
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

func TestSet_GormValue_Float64(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           Set[float64]
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: Set[float64]{1.1},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberNS{Value: []string{"1.1"}},
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

func TestSet_GormValue_Binary(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	type test struct {
		sut           Set[[]byte]
		args          args
		want          clause.Expr
		expectDBError error
	}
	tests := map[string]test{
		"happy-path": {
			sut: Set[[]byte]{[]byte("A")},
			args: args{
				ctx: context.Background(),
				db:  &gorm.DB{},
			},
			want: clause.Expr{
				SQL: "?",
				Vars: []interface{}{
					types.AttributeValueMemberBS{Value: [][]byte{[]byte("A")}},
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
