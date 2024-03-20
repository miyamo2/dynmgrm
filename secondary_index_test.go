package dynmgrm

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
		expectedDist       interface{}
	}

	type TestTable struct{}

	tests := map[string]test{
		"happy-path": {
			sut: SecondaryIndex("testTable.testIndex"),
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
		"happy-path/with-secondary-index-of": {
			sut: SecondaryIndex("testIndex", SecondaryIndexOf("testTable")),
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
		"happy-path/with-model-as": {
			sut: SecondaryIndex("testTable.testIndex", ModelAs(TestTable{})),
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
			expectedDist:   TestTable{},
		},
		"happy-path/with-statement-table": {
			sut: SecondaryIndex("testIndex"),
			builder: &gorm.Statement{
				Table: "testTable",
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
				if diff := cmp.Diff(statement.Dest, tt.expectedDist); diff != "" {
					t.Errorf("Statement.Dest mismatch (-want +got):\n%s", diff)
					return
				}
			}
		})
	}
}
