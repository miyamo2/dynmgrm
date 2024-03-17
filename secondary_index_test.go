package dynamgorm

import (
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
)

func TestSecondaryIndexExpression_Build(t *testing.T) {
	type test struct {
		sut                secondaryIndexExpression
		builder            clause.Builder
		exprectedTable     string
		exprectedTableExpr clause.Expr
	}
	tests := map[string]test{
		"happy-path/with-table-name": {
			sut: SecondaryIndex("testTable", "testIndex"),
			builder: &gorm.Statement{
				DB: &gorm.DB{
					Config: &gorm.Config{
						Dialector: &mockDialector{},
					},
				},
			},
			exprectedTableExpr: clause.Expr{
				SQL: `"testTable"."testIndex"`,
			},
			exprectedTable: "testTable.testIndex",
		},
		"happy-path/with-table-clause": {
			sut: SecondaryIndex(
				clause.Table{
					Name: "testTable",
				}, "testIndex"),
			builder: &gorm.Statement{
				DB: &gorm.DB{
					Config: &gorm.Config{
						Dialector: &mockDialector{},
					},
				},
			},
			exprectedTableExpr: clause.Expr{
				SQL: `"testTable"."testIndex"`,
			},
			exprectedTable: "testTable.testIndex",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.sut.Build(tt.builder)
			if statement, ok := tt.builder.(*gorm.Statement); ok {
				if diff := cmp.Diff(statement.Table, tt.exprectedTable); diff != "" {
					t.Errorf("Statement.Table mismatch (-want +got):\n%s", diff)
					return
				}
				if diff := cmp.Diff(statement.TableExpr, &tt.exprectedTableExpr); diff != "" {
					t.Errorf("Statement.TableExpr mismatch (-want +got):\n%s", diff)
					return
				}
			}
		})
	}
}
