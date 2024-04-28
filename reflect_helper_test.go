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
		D string `gorm:"type:string"`
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
				DBType: "string",
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
				t.Errorf("getDBNameFromStructField() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getDBTypeFromStructField(t *testing.T) {
	type A struct {
		A string `gorm:"type:string"`
		B string `gorm:"type"`
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
			want: "string",
		},
		"happy_path/type_name_is_empty": {
			args: args{
				tag: rt.Field(1),
			},
			want: "",
		},
		"happy_path/empty_tag": {
			args: args{
				tag: rt.Field(2),
			},
			want: "",
		},
		"happy_path/no_tag": {
			args: args{
				tag: rt.Field(3),
			},
			want: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := getDBTypeFromStructField(tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getDBTypeFromStructField() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_extractDBTypeFromStructField(t *testing.T) {
	type A struct {
		A string `gorm:"type:string"`
		B string
		C int
		D float64
		E []byte
		F rune
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
			want: "string",
		},
		"happy_path/string": {
			args: args{
				tag: rt.Field(1),
			},
			want: "string",
		},
		"happy_path/int": {
			args: args{
				tag: rt.Field(2),
			},
			want: "number",
		},
		"happy_path/float64": {
			args: args{
				tag: rt.Field(3),
			},
			want: "number",
		},
		"happy_path/byte_slice": {
			args: args{
				tag: rt.Field(4),
			},
			want: "binary",
		},
		"happy_path/other_type": {
			args: args{
				tag: rt.Field(5),
			},
			want: "string",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := extractDBTypeFromStructField(tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("extractDBTypeFromStructField() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newDynmgrmTag(t *testing.T) {
	type A struct {
		A string `dynmgrm:"pk"`
		B string `dynmgrm:"sk"`
		C string `dynmgrm:"gsi-pk:c_d-index"`
		D string `dynmgrm:"gsi-sk:c_d-index"`
		E string `dynmgrm:"lsi:e-index"`
		F string `dynmgrm:"non-projective:[c_d-index,e-index]"`
		G string `dynmgrm:"pk;sk"`
		H string `dynmgrm:"sk;pk"`
		I string `dynmgrm:""`
	}
	rt := reflect.TypeOf(A{})
	type args struct {
		tag reflect.StructTag
	}

	type test struct {
		args args
		want dynmgrmTag
	}

	tests := map[string]test{
		"happy_path/pk": {
			args: args{
				tag: rt.Field(0).Tag,
			},
			want: dynmgrmTag{
				PK: true,
			},
		},
		"happy_path/sk": {
			args: args{
				tag: rt.Field(1).Tag,
			},
			want: dynmgrmTag{
				SK: true,
			},
		},
		"happy_path/gsi-pk": {
			args: args{
				tag: rt.Field(2).Tag,
			},
			want: dynmgrmTag{
				IndexProperty: []secondaryIndexProperty{
					{
						PK:   true,
						Name: "c_d-index",
						Kind: secondaryIndexKindGSI,
					},
				},
			},
		},
		"happy_path/gsi-sk": {
			args: args{
				tag: rt.Field(3).Tag,
			},
			want: dynmgrmTag{
				IndexProperty: []secondaryIndexProperty{
					{
						SK:   true,
						Name: "c_d-index",
						Kind: secondaryIndexKindGSI,
					},
				},
			},
		},
		"happy_path/lsi": {
			args: args{
				tag: rt.Field(4).Tag,
			},
			want: dynmgrmTag{
				IndexProperty: []secondaryIndexProperty{
					{
						Name: "e-index",
						Kind: secondaryIndexKindLSI,
						SK:   true,
					},
				},
			},
		},
		"happy_path/non-projective": {
			args: args{
				tag: rt.Field(5).Tag,
			},
			want: dynmgrmTag{
				NonProjective: []string{"c_d-index", "e-index"},
			},
		},
		"unhappy_path/pk_and_sk": {
			args: args{
				tag: rt.Field(6).Tag,
			},
			want: dynmgrmTag{
				PK: true,
			},
		},
		"unhappy_path/sk_and_pk": {
			args: args{
				tag: rt.Field(7).Tag,
			},
			want: dynmgrmTag{
				SK: true,
			},
		},
		"unhappy_path/empty_value": {
			args: args{
				tag: rt.Field(8).Tag,
			},
			want: dynmgrmTag{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := newDynmgrmTag(tt.args.tag)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("newDynmgrmTag() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_newDynmgrmTableDefine(t *testing.T) {
	type A struct {
		A string `dynmgrm:"pk"`
		B string `dynmgrm:"sk"`
		C string `dynmgrm:"gsi-pk:c_d-index"`
		D string `dynmgrm:"gsi-sk:c_d-index"`
		E string `dynmgrm:"lsi:e-index"`
		F string `dynmgrm:"non-projective:[c_d-index,e-index]"`
	}
	type test struct {
		args reflect.Type
		want dynmgrmTableDefine
	}
	tests := map[string]test{
		"happy_path": {
			args: reflect.TypeOf(A{}),
			want: dynmgrmTableDefine{
				PK: dynmgrmKeyDefine{
					Name:     "a",
					DataType: "string",
				},
				SK: dynmgrmKeyDefine{
					Name:     "b",
					DataType: "string",
				},
				NonKeyAttr: []string{"c", "d", "e", "f"},
				GSI: map[string]*dynmgrmSecondaryIndexDefine{
					"c_d-index": {
						PK: dynmgrmKeyDefine{
							Name:     "c",
							DataType: "string",
						},
						SK: dynmgrmKeyDefine{
							Name:     "d",
							DataType: "string",
						},
						NonProjectiveAttrs: []string{"f"},
					},
				},
				LSI: map[string]*dynmgrmSecondaryIndexDefine{
					"e-index": {
						SK: dynmgrmKeyDefine{
							Name:     "e",
							DataType: "string",
						},
						NonProjectiveAttrs: []string{"f"},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := newDynmgrmTableDefine(tt.args)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("newDynmgrmTableDefine() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
