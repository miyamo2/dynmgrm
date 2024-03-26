package dynmgrm

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_toAttibuteValue(t *testing.T) {
	type args struct {
		value interface{}
	}

	type want struct {
		av  types.AttributeValue
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"happy-path/string-sets": {
			args: args{
				value: Set[string]{"a", "b", "c"},
			},
			want: want{
				av:  &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}},
				err: nil,
			},
		},
		"happy-path/int-sets": {
			args: args{
				value: Set[int]{1, 2, 3},
			},
			want: want{
				av:  &types.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
				err: nil,
			},
		},
		"happy-path/float-sets": {
			args: args{
				value: Set[float64]{1.1, 2.2, 3.3},
			},
			want: want{
				av:  &types.AttributeValueMemberNS{Value: []string{"1.1", "2.2", "3.3"}},
				err: nil,
			},
		},
		"happy-path/byte-sets": {
			args: args{
				value: Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
			},
			want: want{
				av:  &types.AttributeValueMemberBS{Value: [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
				err: nil,
			},
		},
		"happy-path/list": {
			args: args{
				value: List{1, 2, 3},
			},
			want: want{
				av: &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberN{Value: "1"},
						&types.AttributeValueMemberN{Value: "2"},
						&types.AttributeValueMemberN{Value: "3"},
					},
				},
				err: nil,
			},
		},
		"happy-path/map": {
			args: args{
				value: Map{"a": 1, "b": 2, "c": 3},
			},
			want: want{
				av: &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"a": &types.AttributeValueMemberN{Value: "1"},
						"b": &types.AttributeValueMemberN{Value: "2"},
						"c": &types.AttributeValueMemberN{Value: "3"},
					},
				},
			},
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
			got, err := toAttibuteValue(tt.args.value)
			if !errors.Is(tt.want.err, err) {
				t.Errorf("toAttibuteValue() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.want.av, got, opts...); diff != "" {
				t.Errorf("toAttibuteValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDocumentAttributeValue_StringSet(t *testing.T) {
	type args struct {
		value interface{}
	}

	type want struct {
		av  types.AttributeValue
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"happy-path/string-sets": {
			args: args{
				value: Set[string]{"a", "b", "c"},
			},
			want: want{
				av:  &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}},
				err: nil,
			},
		},
		"unhappy-path/int-sets": {
			args: args{
				value: Set[int]{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberSS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/float-sets": {
			args: args{
				value: Set[float64]{1.1, 2.2, 3.3},
			},
			want: want{
				av:  (*types.AttributeValueMemberSS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/byte-sets": {
			args: args{
				value: Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
			},
			want: want{
				av:  (*types.AttributeValueMemberSS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/list": {
			args: args{
				value: List{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberSS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/map": {
			args: args{
				value: Map{"a": 1, "b": 2, "c": 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberSS)(nil),
				err: errValueInCompatible,
			},
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
			got, err := toDocumentAttributeValue[*types.AttributeValueMemberSS](tt.args.value)
			if !errors.Is(tt.want.err, err) {
				t.Errorf("toDocumentAttributeValue() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.want.av, got, opts...); diff != "" {
				t.Errorf("toDocumentAttributeValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDocumentAttributeValue_NumberSet(t *testing.T) {
	type args struct {
		value interface{}
	}

	type want struct {
		av  types.AttributeValue
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"unhappy-path/int-sets": {
			args: args{
				value: Set[int]{1, 2, 3},
			},
			want: want{
				av:  &types.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
				err: nil,
			},
		},
		"happy-path/float-sets": {
			args: args{
				value: Set[float64]{1.1, 2.2, 3.3},
			},
			want: want{
				av:  &types.AttributeValueMemberNS{Value: []string{"1.1", "2.2", "3.3"}},
				err: nil,
			},
		},
		"happy-path/string-sets": {
			args: args{
				value: Set[string]{"a", "b", "c"},
			},
			want: want{
				av:  (*types.AttributeValueMemberNS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/byte-sets": {
			args: args{
				value: Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
			},
			want: want{
				av:  (*types.AttributeValueMemberNS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/list": {
			args: args{
				value: List{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberNS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/map": {
			args: args{
				value: Map{"a": 1, "b": 2, "c": 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberNS)(nil),
				err: errValueInCompatible,
			},
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
			got, err := toDocumentAttributeValue[*types.AttributeValueMemberNS](tt.args.value)
			if !errors.Is(tt.want.err, err) {
				t.Errorf("toDocumentAttributeValue() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.want.av, got, opts...); diff != "" {
				t.Errorf("toDocumentAttributeValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDocumentAttributeValue_BinarySet(t *testing.T) {
	type args struct {
		value interface{}
	}

	type want struct {
		av  types.AttributeValue
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"happy-path/binary-sets": {
			args: args{
				value: Set[[]byte]{[]byte("1"), []byte("2"), []byte("3")},
			},
			want: want{
				av:  &types.AttributeValueMemberBS{Value: [][]byte{[]byte("1"), []byte("2"), []byte("3")}},
				err: nil,
			},
		},
		"unhappy-path/int-sets": {
			args: args{
				value: Set[int]{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberBS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/float-sets": {
			args: args{
				value: Set[float64]{1.0, 2.0, 3.0},
			},
			want: want{
				av:  (*types.AttributeValueMemberBS)(nil),
				err: errValueInCompatible,
			},
		},
		"happy-path/string-sets": {
			args: args{
				value: Set[string]{"a", "b", "c"},
			},
			want: want{
				av:  (*types.AttributeValueMemberBS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/list": {
			args: args{
				value: List{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberBS)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/map": {
			args: args{
				value: Map{"a": 1, "b": 2, "c": 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberBS)(nil),
				err: errValueInCompatible,
			},
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
			got, err := toDocumentAttributeValue[*types.AttributeValueMemberBS](tt.args.value)
			if !errors.Is(tt.want.err, err) {
				t.Errorf("toDocumentAttributeValue() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.want.av, got, opts...); diff != "" {
				t.Errorf("toDocumentAttributeValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDocumentAttributeValue_List(t *testing.T) {
	type args struct {
		value interface{}
	}

	type want struct {
		av  types.AttributeValue
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"happy-path/list": {
			args: args{
				value: List{1, 2, 3},
			},
			want: want{
				av: &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberN{Value: "1"},
						&types.AttributeValueMemberN{Value: "2"},
						&types.AttributeValueMemberN{Value: "3"},
					},
				},
				err: nil,
			},
		},
		"happy-path/binary-sets": {
			args: args{
				value: Set[[]byte]{[]byte("1"), []byte("2"), []byte("3")},
			},
			want: want{
				av:  (*types.AttributeValueMemberL)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/int-sets": {
			args: args{
				value: Set[int]{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberL)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/float-sets": {
			args: args{
				value: Set[float64]{1.0, 2.0, 3.0},
			},
			want: want{
				av:  (*types.AttributeValueMemberL)(nil),
				err: errValueInCompatible,
			},
		},
		"happy-path/string-sets": {
			args: args{
				value: Set[string]{"a", "b", "c"},
			},
			want: want{
				av:  (*types.AttributeValueMemberL)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/map": {
			args: args{
				value: Map{"a": 1, "b": 2, "c": 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberL)(nil),
				err: errValueInCompatible,
			},
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
			got, err := toDocumentAttributeValue[*types.AttributeValueMemberL](tt.args.value)
			if !errors.Is(tt.want.err, err) {
				t.Errorf("toDocumentAttributeValue() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.want.av, got, opts...); diff != "" {
				t.Errorf("toDocumentAttributeValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDocumentAttributeValue_Map(t *testing.T) {
	type args struct {
		value interface{}
	}

	type want struct {
		av  types.AttributeValue
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"happy-path/map": {
			args: args{
				value: Map{"a": 1, "b": 2, "c": 3},
			},
			want: want{
				av: &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"a": &types.AttributeValueMemberN{Value: "1"},
						"b": &types.AttributeValueMemberN{Value: "2"},
						"c": &types.AttributeValueMemberN{Value: "3"},
					},
				},
				err: nil,
			},
		},
		"unhappy-path/list": {
			args: args{
				value: List{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberM)(nil),
				err: errValueInCompatible,
			},
		},
		"happy-path/binary-sets": {
			args: args{
				value: Set[[]byte]{[]byte("1"), []byte("2"), []byte("3")},
			},
			want: want{
				av:  (*types.AttributeValueMemberM)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/int-sets": {
			args: args{
				value: Set[int]{1, 2, 3},
			},
			want: want{
				av:  (*types.AttributeValueMemberM)(nil),
				err: errValueInCompatible,
			},
		},
		"unhappy-path/float-sets": {
			args: args{
				value: Set[float64]{1.0, 2.0, 3.0},
			},
			want: want{
				av:  (*types.AttributeValueMemberM)(nil),
				err: errValueInCompatible,
			},
		},
		"happy-path/string-sets": {
			args: args{
				value: Set[string]{"a", "b", "c"},
			},
			want: want{
				av:  (*types.AttributeValueMemberM)(nil),
				err: errValueInCompatible,
			},
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
			got, err := toDocumentAttributeValue[*types.AttributeValueMemberM](tt.args.value)
			if !errors.Is(tt.want.err, err) {
				t.Errorf("toDocumentAttributeValue() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.want.av, got, opts...); diff != "" {
				t.Errorf("toDocumentAttributeValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
