package dynamgorm

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"testing"
)

var _ gorm.Dialector = (*mockDialector)(nil)

type mockDialector struct{}

func (s mockDialector) Name() string {
	return ""
}

func (s mockDialector) Initialize(db *gorm.DB) error {
	return nil
}

func (s mockDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return nil
}

func (s mockDialector) DataTypeOf(field *schema.Field) string {
	return ""
}

func (s mockDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return nil
}

func (s mockDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}

func (s mockDialector) QuoteTo(writer clause.Writer, s2 string) {
	writer.WriteString(fmt.Sprintf(`"%s"`, s2))
}

func (s mockDialector) Explain(sql string, vars ...interface{}) string {
	return ""
}

var _ clause.Builder = (*mockBuilder)(nil)

type mockBuilder struct{}

func (m mockBuilder) WriteByte(b byte) error {
	return nil
}

func (m mockBuilder) WriteString(s string) (int, error) {
	return 0, nil
}

func (m mockBuilder) WriteQuoted(field interface{}) {
	return
}

func (m mockBuilder) AddVar(writer clause.Writer, i ...interface{}) {
	return
}

func (m mockBuilder) AddError(err error) error {
	return nil
}

func TestValuesClause(t *testing.T) {
	type test struct {
		args         clause.Clause
		expectedSQL  string
		expectedVars []interface{}
	}
	tests := map[string]test{
		"happy-path/single-column": {
			args: clause.Clause{
				Expression: clause.Values{
					Columns: []clause.Column{{
						Name: "column1",
					}},
					Values: [][]interface{}{{"value1"}},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{"value1"},
		},
		"happy-path/multiple-columns": {
			args: clause.Clause{
				Expression: clause.Values{
					Columns: []clause.Column{
						{
							Name: "column1",
						},
						{
							Name: "column2",
						},
					},
					Values: [][]interface{}{
						{"value1", "value2"},
					},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?, 'column2' : ?}",
			expectedVars: []interface{}{"value1", "value2"},
		},
		"happy-path/with-sets": {
			args: clause.Clause{
				Expression: clause.Values{
					Columns: []clause.Column{
						{
							Name: "column1",
						},
					},
					Values: [][]interface{}{
						{Sets[string]{"value1", "value2"}},
					},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{Sets[string]{"value1", "value2"}},
		},
		"happy-path/with-map": {
			args: clause.Clause{
				Expression: clause.Values{
					Columns: []clause.Column{
						{
							Name: "column1",
						},
					},
					Values: [][]interface{}{
						{Map{"key1": "value1"}},
					},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{Map{"key1": "value1"}},
		},
		"happy-path/with-list": {
			args: clause.Clause{
				Expression: clause.Values{
					Columns: []clause.Column{
						{
							Name: "column1",
						},
					},
					Values: [][]interface{}{
						{List{"value1", "value2"}},
					},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{List{"value1", "value2"}},
		},
		"unhappy-path/empty-columns": {
			args: clause.Clause{
				Expression: clause.Values{
					Columns: []clause.Column{},
				},
			},
			expectedSQL:  "",
			expectedVars: nil,
		},
		"unhappy-path/expression-is-not-values-type": {
			args: clause.Clause{
				Expression: clause.Clause{},
			},
			expectedSQL:  "",
			expectedVars: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sut := &gorm.Statement{
				DB: &gorm.DB{
					Config: &gorm.Config{
						Dialector: &mockDialector{},
					},
				},
			}
			// Call the function we are testing
			buildValuesClause(tt.args, sut)

			acutalSQL := sut.SQL.String()
			if diff := cmp.Diff(acutalSQL, tt.expectedSQL); diff != "" {
				t.Errorf("SQL mismatch (-want +got):\n%s", diff)
			}
			acutalVars := sut.Vars
			if diff := cmp.Diff(acutalVars, tt.expectedVars); diff != "" {
				t.Errorf("Vars mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_bindVarIfCollectionType(t *testing.T) {
	type args struct {
		builder clause.Builder
		value   interface{}
	}
	type test struct {
		args args
		want bool
	}
	tests := map[string]test{
		"unhappy-path/not-gorm-statement": {
			args: args{
				builder: &mockBuilder{},
				value:   Map{},
			},
			want: false,
		},
		"happy-path/not-collection-type": {
			args: args{
				builder: &gorm.Statement{
					DB: &gorm.DB{
						Config: &gorm.Config{
							Dialector: &mockDialector{},
						},
					},
				},
				value: "not-collection-type",
			},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if gotBound := bindVarIfCollectionType(tt.args.builder, tt.args.value); gotBound != tt.want {
				t.Errorf("bindVarIfCollectionType() = %v, want %v", gotBound, tt.want)
			}
		})
	}
}
