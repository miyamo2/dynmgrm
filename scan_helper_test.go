package dynmgrm

import (
	"database/sql"
	"errors"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

type A struct {
	Str string
}

type B struct {
	Str string
}

var _ sql.Scanner = (*B)(nil)

var errScanToB = errors.New("scan to B")

func (b *B) Scan(value interface{}) error {
	m, ok := value.(map[string]interface{})
	if !ok {
		return errScanToB
	}
	b.Str = m["str"].(string)
	return nil
}

func Test_assignInterfaceValueToReflectValue(t *testing.T) {
	type args struct {
		rt    reflect.Type
		rv    reflect.Value
		value interface{}
	}
	type test struct {
		args     args
		want     error
		expected interface{}
	}
	tests := map[string]test{
		"happy_path/string": {
			args: args{
				rt:    reflect.TypeOf(""),
				rv:    reflect.New(reflect.TypeOf("")).Elem(),
				value: "a",
			},
			want:     nil,
			expected: "a",
		},
		"happy_path/int": {
			args: args{
				rt:    reflect.TypeOf(int(0)),
				rv:    reflect.New(reflect.TypeOf(int(0))).Elem(),
				value: float64(1.0),
			},
			want:     nil,
			expected: 1,
		},
		"happy_path/bool": {
			args: args{
				rt:    reflect.TypeOf(false),
				rv:    reflect.New(reflect.TypeOf(false)).Elem(),
				value: true,
			},
			want:     nil,
			expected: true,
		},
		"happy_path/float64": {
			args: args{
				rt:    reflect.TypeOf(0.0),
				rv:    reflect.New(reflect.TypeOf(0.0)).Elem(),
				value: 1.0,
			},
			want:     nil,
			expected: 1.0,
		},
		"happu_path/binary": {
			args: args{
				rt:    reflect.TypeOf([]byte{}),
				rv:    reflect.New(reflect.TypeOf([]byte{})).Elem(),
				value: []byte("a"),
			},
			want:     nil,
			expected: []byte("a"),
		},
		"happy_path/struct": {
			args: args{
				rt:    reflect.TypeOf(A{}),
				rv:    reflect.ValueOf(&A{}),
				value: map[string]interface{}{"str": "a"},
			},
			want:     nil,
			expected: &A{Str: "a"},
		},
		"happy_path/pointer": {
			args: args{
				rt:    reflect.TypeOf(&A{}),
				rv:    reflect.New(reflect.TypeOf(&A{})).Elem(),
				value: map[string]interface{}{"str": "a"},
			},
			want:     nil,
			expected: &A{Str: "a"},
		},
		"happy_path/scanner": {
			args: args{
				rt:    reflect.TypeOf(B{}),
				rv:    reflect.ValueOf(&B{}),
				value: map[string]interface{}{"str": "a"},
			},
			want:     nil,
			expected: &B{Str: "a"},
		},
		"unhappy_path/incompatible_with_string": {
			args: args{
				rt:    reflect.TypeOf(""),
				rv:    reflect.New(reflect.TypeOf("")).Elem(),
				value: struct{}{},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: "",
		},
		"unhappy_path/incompatible_with_int": {
			args: args{
				rt:    reflect.TypeOf(int(0)),
				rv:    reflect.New(reflect.TypeOf(int(0))).Elem(),
				value: struct{}{},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: 0,
		},
		"unhappy_path/incompatible_with_bool": {
			args: args{
				rt:    reflect.TypeOf(false),
				rv:    reflect.New(reflect.TypeOf(false)).Elem(),
				value: struct{}{},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: false,
		},
		"unhappy_path/incompatible_with_float64": {
			args: args{
				rt:    reflect.TypeOf(0.0),
				rv:    reflect.New(reflect.TypeOf(0.0)).Elem(),
				value: struct{}{},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: 0.0,
		},
		"unhappy_path/incompatible_with_binary": {
			args: args{
				rt:    reflect.TypeOf([]byte{}),
				rv:    reflect.New(reflect.TypeOf([]byte{})).Elem(),
				value: struct{}{},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: ([]byte)(nil),
		},
		"unhappy_path/non_byte_slice": {
			args: args{
				rt:    reflect.TypeOf([]string{}),
				rv:    reflect.New(reflect.TypeOf([]string{})).Elem(),
				value: []string{"a"},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: ([]string)(nil),
		},
		"unhappy_path/incompatible_with_struct": {
			args: args{
				rt:    reflect.TypeOf(A{}),
				rv:    reflect.ValueOf(&A{}),
				value: "a",
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: &A{},
		},
		"unhappy_path/incompatible_with_struct_attribute": {
			args: args{
				rt:    reflect.TypeOf(A{}),
				rv:    reflect.ValueOf(&A{}),
				value: map[string]interface{}{"str": struct{}{}},
			},
			want:     ErrNestedStructHasIncompatibleAttributes,
			expected: &A{},
		},
		"happy_path/nil_pointer": {
			args: args{
				rt:    reflect.TypeOf(&A{}),
				rv:    reflect.New(reflect.TypeOf(&A{})).Elem(),
				value: nil,
			},
			want:     nil,
			expected: (*A)(nil),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := assignInterfaceValueToReflectValue(tt.args.rt, tt.args.rv, tt.args.value)
			if !errors.Is(err, tt.want) {
				t.Errorf("assignInterfaceValueToReflectValue() error = %v, want %v", err, tt.want)
			}
			if diff := cmp.Diff(tt.expected, tt.args.rv.Interface()); diff != "" {
				t.Errorf("assignInterfaceValueToReflectValue() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
