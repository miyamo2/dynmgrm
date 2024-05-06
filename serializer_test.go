package dynmgrm

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm/schema"
	"reflect"
	"testing"
)

type Inner struct {
	A string
	B string
}
type Outer struct {
	Inner Inner
}

func Test_nestedStructSerializer_Scan(t *testing.T) {
	type args struct {
		ctx     context.Context
		field   *schema.Field
		dst     Outer
		dbValue interface{}
	}
	type want struct {
		err error
	}
	type expected struct {
		dst Outer
	}
	type test struct {
		args args
		want want
		expected
	}
	tests := map[string]test{
		"happy_path": {
			args: args{
				ctx: context.Background(),
				field: &schema.Field{
					FieldType: reflect.TypeOf(Outer{}),
					ReflectValueOf: func(ctx context.Context, value reflect.Value) reflect.Value {
						return value.Elem()
					},
				},
				dst: Outer{},
				dbValue: map[string]interface{}{
					"inner": map[string]interface{}{
						"a": "foo",
						"b": "bar",
					},
				},
			},
			want: want{
				err: nil,
			},
			expected: expected{
				Outer{
					Inner: Inner{
						A: "foo",
						B: "bar",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			n := nestedStructSerializer{}
			err := n.Scan(tt.args.ctx, tt.args.field, reflect.ValueOf(&tt.args.dst), tt.args.dbValue)
			if !errors.Is(err, tt.want.err) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want.err)
			}
			if diff := cmp.Diff(tt.expected.dst, tt.args.dst); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
