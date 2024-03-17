package dynamgorm

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynamgorm/internal/mocks"
	"go.uber.org/mock/gomock"
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

func (s mockDialector) Initialize(_ *gorm.DB) error {
	return nil
}

func (s mockDialector) Migrator(_ *gorm.DB) gorm.Migrator {
	return nil
}

func (s mockDialector) DataTypeOf(_ *schema.Field) string {
	return ""
}

func (s mockDialector) DefaultValueOf(_ *schema.Field) clause.Expression {
	return nil
}

func (s mockDialector) BindVarTo(writer clause.Writer, _ *gorm.Statement, _ interface{}) {
	_ = writer.WriteByte('?')
}

func (s mockDialector) QuoteTo(writer clause.Writer, s2 string) {
	_, _ = writer.WriteString(fmt.Sprintf(`"%s"`, s2))
}

func (s mockDialector) Explain(_ string, _ ...interface{}) string {
	return ""
}

func TestValuesClause(t *testing.T) {
	type test struct {
		args         clause.Values
		expectedSQL  string
		expectedVars []interface{}
	}
	tests := map[string]test{
		"happy-path/single-column": {
			args: clause.Values{
				Columns: []clause.Column{{
					Name: "column1",
				}},
				Values: [][]interface{}{{"value1"}},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{"value1"},
		},
		"happy-path/multiple-columns": {
			args: clause.Values{
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
			expectedSQL:  "VALUE {'column1' : ?, 'column2' : ?}",
			expectedVars: []interface{}{"value1", "value2"},
		},
		"happy-path/with-sets": {
			args: clause.Values{
				Columns: []clause.Column{
					{
						Name: "column1",
					},
				},
				Values: [][]interface{}{
					{Sets[string]{"value1", "value2"}},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{Sets[string]{"value1", "value2"}},
		},
		"happy-path/with-map": {
			args: clause.Values{
				Columns: []clause.Column{
					{
						Name: "column1",
					},
				},
				Values: [][]interface{}{
					{Map{"key1": "value1"}},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{Map{"key1": "value1"}},
		},
		"happy-path/with-list": {
			args: clause.Values{
				Columns: []clause.Column{
					{
						Name: "column1",
					},
				},
				Values: [][]interface{}{
					{List{"value1", "value2"}},
				},
			},
			expectedSQL:  "VALUE {'column1' : ?}",
			expectedVars: []interface{}{List{"value1", "value2"}},
		},
		"unhappy-path/empty-columns": {
			args: clause.Values{
				Columns: []clause.Column{},
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
		stmt  *gorm.Statement
		value interface{}
	}
	type test struct {
		args args
		want bool
	}
	tests := map[string]test{
		"happy-path/not-collection-type": {
			args: args{
				stmt: &gorm.Statement{
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
			if gotBound := bindVarIfCollectionType(tt.args.stmt, tt.args.value); gotBound != tt.want {
				t.Errorf("bindVarIfCollectionType() = %v, want %v", gotBound, tt.want)
			}
		})
	}
}

type mockExpressionProp struct {
	callCount int
}

type setupMockExpressionOptions func(*mockExpressionProp)

func setupMockExpressionWithCallCount(count int) setupMockExpressionOptions {
	return func(prop *mockExpressionProp) {
		prop.callCount = count
	}
}

func setupMockExpression(t *testing.T, xp *mocks.MockExpression, options ...setupMockExpressionOptions) *mocks.MockExpression {
	t.Helper()
	prop := mockExpressionProp{}
	for _, opt := range options {
		opt(&prop)
	}
	xp.EXPECT().Build(gomock.Any()).Times(prop.callCount)
	return xp
}

func Test_toClauseBuilder(t *testing.T) {
	type test struct {
		builder       clause.Builder
		xpOverride    clause.Expression
		mockXpOptions []setupMockExpressionOptions
	}
	xpBuilder := func(c *mocks.MockExpression, stmt *gorm.Statement) {
		c.Build(stmt)
	}
	tests := map[string]test{
		"happy-path": {
			builder: &gorm.Statement{},
			mockXpOptions: []setupMockExpressionOptions{
				setupMockExpressionWithCallCount(1),
			},
		},
		"unhappy-path/not-statement": {
			builder: &mocks.MockBuilder{},
		},
		"unhappy-path/expression-not-matching": {
			builder:    &gorm.Statement{},
			xpOverride: clause.Values{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			xp := mocks.NewMockExpression(ctrl)
			setupMockExpression(t, xp, tt.mockXpOptions...)
			c := clause.Clause{
				Expression: xp,
			}
			if tt.xpOverride != nil {
				c.Expression = tt.xpOverride
			}
			toClauseBuilder(xpBuilder)(c, tt.builder)
		})
	}
}
