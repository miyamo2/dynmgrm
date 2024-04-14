package dynmgrm

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func Test_newGormTag(t *testing.T) {
	type A struct {
		A string `gorm:"column:a"`
		B string `gorm:"column"`
		C string `gorm:""`
		D string
	}

	rt := reflect.TypeOf(A{})
	type args struct {
		tag reflect.StructTag
	}

	type test struct {
		args args
		want gormTag
	}

	tests := map[string]test{
		"happy_path": {
			args: args{
				tag: rt.Field(0).Tag,
			},
			want: gormTag{
				Column: "a",
			},
		},
		"unhappy_path/column_name_is_empty": {
			args: args{
				tag: rt.Field(1).Tag,
			},
			want: gormTag{
				Column: "",
			},
		},
		"unhappy_path/empty_tag": {
			args: args{
				tag: rt.Field(2).Tag,
			},
			want: gormTag{
				Column: "",
			},
		},
		"unhappy_path/no_tag": {
			args: args{
				tag: rt.Field(3).Tag,
			},
			want: gormTag{
				Column: "",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := newGormTag(tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("newGormTag() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getNameFromStructField(t *testing.T) {
	type A struct {
		A string `gorm:"column:aaa"`
		B string `gorm:"column"`
		C string `gorm:""`
		D string
	}

	rt := reflect.TypeOf(A{})
	type args struct {
		tag reflect.StructField
	}

	type test struct {
		args args
		want string
	}

	tests := map[string]test{
		"happy_path/with_tag": {
			args: args{
				tag: rt.Field(0),
			},
			want: "aaa",
		},
		"happy_path/column_name_is_empty": {
			args: args{
				tag: rt.Field(1),
			},
			want: "b",
		},
		"happy_path/empty_tag": {
			args: args{
				tag: rt.Field(2),
			},
			want: "c",
		},
		"happy_path/no_tag": {
			args: args{
				tag: rt.Field(3),
			},
			want: "d",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := getDBNameFromStructField(tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("newGormTag() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
