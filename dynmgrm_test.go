package dynmgrm

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func TestNew(t *testing.T) {
	type args struct {
		option []DialectorOption
	}

	type want struct {
		dsn string
	}

	type test struct {
		args args
		want want
		env  map[string]string
	}

	tests := map[string]test{
		"happy_path/without_option": {
			args: args{
				option: []DialectorOption{},
			},
			want: want{
				dsn: "",
			},
		},
		"happy_path/with_region": {
			args: args{
				option: []DialectorOption{
					WithRegion("ap-north-east1"),
				},
			},
			want: want{
				dsn: "region=ap-north-east1",
			},
		},
		"happy_path/with_ak_id": {
			args: args{
				option: []DialectorOption{
					WithAccessKeyID("ACCESS_KEY_ID"),
				},
			},
			want: want{
				dsn: "akId=ACCESS_KEY_ID",
			},
		},
		"happy_path/with_secret": {
			args: args{
				option: []DialectorOption{
					WithSecretKey("SECRET"),
				}},
			want: want{
				dsn: "secretKey=SECRET",
			},
		},
		"happy_path/with_endpoint": {
			args: args{
				option: []DialectorOption{
					WithEndpoint("http://localhost:8000"),
				},
			},
			want: want{
				dsn: "endpoint=http://localhost:8000",
			},
		},
		"happy_path/with_timeout": {
			args: args{
				option: []DialectorOption{
					WithTimeout(1000),
				},
			},
			want: want{
				dsn: "timeout=1000",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got := New(tt.args.option...)
			d, ok := got.(*Dialector)
			if !ok {
				t.Fatal("dialector is not *Dialector")
			}
			if diff := cmp.Diff(tt.want.dsn, d.dbOpener.DSN()); diff != "" {
				t.Errorf("DSN() mismatch (-want +got): \n%v", diff)
			}
		})

	}
}

func TestOpen(t *testing.T) {
	type args struct {
		dsn string
	}

	type want struct {
		dsn string
	}

	type test struct {
		args args
		want want
	}

	tests := map[string]test{
		"happy_path": {
			args: args{
				dsn: "region=ap-north-east1;akId=ACCESS_KEY_ID;secret",
			},
			want: want{
				dsn: "region=ap-north-east1;akId=ACCESS_KEY_ID;secret",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Open(tt.args.dsn)
			d, ok := got.(*Dialector)
			if !ok {
				t.Fatal("dialector is not *Dialector")
			}
			if diff := cmp.Diff(tt.want.dsn, d.dbOpener.DSN()); diff != "" {
				t.Errorf("DSN() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}

type mockClauseWriter struct {
	v byte
}

func (m *mockClauseWriter) WriteByte(b byte) error {
	m.v = b
	return nil
}

func (m *mockClauseWriter) WriteString(s string) (int, error) {
	return 0, nil
}

func TestDialector_BindVarTo(t *testing.T) {
	type args struct {
		writer clause.Writer
		stmt   *gorm.Statement
		v      interface{}
	}
	type test struct {
		args   args
		expect byte
	}
	tests := map[string]test{
		"happy_path": {
			args: args{
				writer: &mockClauseWriter{},
				stmt:   nil,
				v:      nil,
			},
			expect: '?',
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dialector := Dialector{}
			dialector.BindVarTo(tt.args.writer, tt.args.stmt, tt.args.v)
			actual := tt.args.writer.(*mockClauseWriter).v
			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("BindVarTo() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}

func TestDialector_Name(t *testing.T) {
	type test struct {
		want string
	}

	tests := map[string]test{
		"happy_path": {
			want: DBName,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dialector := Dialector{}
			got := dialector.Name()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Name() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}

func TestDialector_DefaultValueOf(t *testing.T) {
	type args struct {
		field *schema.Field
	}

	type test struct {
		args args
		want clause.Expression
	}

	tests := map[string]test{
		"happy_path": {
			args: args{
				field: &schema.Field{},
			},
			want: clause.Expr{SQL: ""},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dialector := Dialector{}
			got := dialector.DefaultValueOf(tt.args.field)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DefaultValueOf() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}

func Test_writeConnectionParameter(t *testing.T) {
	type args struct {
		dsnbuf *strings.Builder
		key    string
		value  string
	}
	type test struct {
		args         args
		preAddition  map[string]string
		expectString string
	}

	tests := map[string]test{
		"happy_path": {
			args: args{
				dsnbuf: &strings.Builder{},
				key:    "key",
				value:  "value",
			},
			expectString: "key=value",
		},
		"happy_path/with_pre_addition": {
			args: args{
				dsnbuf: &strings.Builder{},
				key:    "key2",
				value:  "value2",
			},
			preAddition: map[string]string{
				"key1": "value1",
			},
			expectString: "key1=value1;key2=value2",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range tt.preAddition {
				writeConnectionParameter(tt.args.dsnbuf, k, v)
			}
			writeConnectionParameter(tt.args.dsnbuf, tt.args.key, tt.args.value)
			actual := tt.args.dsnbuf.String()
			if diff := cmp.Diff(tt.expectString, actual); diff != "" {
				t.Errorf("writeConnectionParameter() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}

func TestDialector_QuoteTo(t *testing.T) {
	type args struct {
		writer clause.Writer
		str    string
	}
	type test struct {
		args     args
		expected string
	}
	tests := map[string]test{
		"happy_path": {
			args: args{
				writer: &gorm.Statement{},
				str:    "test",
			},
			expected: `"test"`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dialector := Dialector{}
			dialector.QuoteTo(tt.args.writer, tt.args.str)
			actual := tt.args.writer.(*gorm.Statement).SQL.String()
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("QuoteTo() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}

func TestDialector_DataTypeOf(t *testing.T) {
	type args struct {
		field *schema.Field
	}
	type test struct {
		args args
		want string
	}
	tests := map[string]test{
		"happy_path/bool": {
			args: args{
				field: &schema.Field{
					DataType: schema.Bool,
				},
			},
			want: DataTypeString.String(),
		},
		"happy_path/uint": {
			args: args{
				field: &schema.Field{
					DataType: schema.Uint,
				},
			},
			want: DataTypeNumber.String(),
		},
		"happy_path/int": {
			args: args{
				field: &schema.Field{
					DataType: schema.Int,
				},
			},
			want: DataTypeNumber.String(),
		},
		"happy_path/float": {
			args: args{
				field: &schema.Field{
					DataType: schema.Float,
				},
			},
			want: DataTypeNumber.String(),
		},
		"happy_path/string": {
			args: args{
				field: &schema.Field{
					DataType: schema.String,
				},
			},
			want: DataTypeString.String(),
		},
		"happy_path/bytes": {
			args: args{
				field: &schema.Field{
					DataType: schema.Bytes,
				},
			},
			want: DataTypeBinary.String(),
		},
		"happy_path/time": {
			args: args{
				field: &schema.Field{
					DataType: schema.Time,
				},
			},
			want: DataTypeString.String(),
		},
		"happy_path/default": {
			args: args{
				field: &schema.Field{
					DataType: schema.DataType("jsonb"),
				},
			},
			want: DataTypeString.String(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dialector := Dialector{}
			got := dialector.DataTypeOf(tt.args.field)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DataTypeOf() mismatch (-want +got): \n%v", diff)
			}
		})
	}
}
