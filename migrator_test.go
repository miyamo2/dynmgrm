package dynmgrm

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynmgrm/internal/mocks"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
)

func TestTableClass_String(t *testing.T) {
	type test struct {
		sut  TableClass
		want string
	}
	tests := map[string]test{
		"happy_path/standard": {
			sut:  TableClassStandard,
			want: "STANDARD",
		},
		"happy_path/standard_ia": {
			sut:  TableClassStandardIA,
			want: "STANDARD_IA",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.sut.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigrator_FullDataTypeOf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bm := mocks.NewMockBaseMigrator(ctrl)
	want := clause.Expr{
		SQL: "test",
	}
	bm.EXPECT().FullDataTypeOf(gomock.Any()).Return(want).Times(1)
	sut := Migrator{
		db:   nil,
		base: bm,
	}
	got := sut.FullDataTypeOf(nil)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FullDataTypeOf() mismatch (-want +got):\n%s", diff)
	}
}

func TestMigrator_currentTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bm := mocks.NewMockBaseMigrator(ctrl)
	want := "test"
	bm.EXPECT().CurrentTable(gomock.Any()).Return(
		clause.Table{
			Name: want,
		}).Times(1)
	sut := Migrator{
		db:   nil,
		base: bm,
	}
	got := sut.currentTable(nil)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("currentTable() mismatch (-want +got):\n%s", diff)
	}
}

type mockDBExecArgs struct {
	sql    string
	values []interface{}
}

type mockDBExecProp struct {
	args   *mockDBExecArgs
	result *gorm.DB
	times  int
}

type mockBaseMigratorCurrentTableProp struct {
	result clause.Table
	times  int
}

var (
	_ CapacityUnitsSpecifier = (*CreateTableATable)(nil)
	_ TableClassSpecifier    = (*CreateTableATable)(nil)
)

type CreateTableATable struct {
	PK              string `dynmgrm:"pk"`
	SK              int    `dynmgrm:"sk"`
	LsiSK           []byte `dynmgrm:"lsi-sk:lsi_sk-index"`
	GsiPK           string `dynmgrm:"gsi-pk:gsi_pk-gsi_sk-index;non-projective:[lsi_sk-index]"`
	GsiSK           string `dynmgrm:"gsi-sk:gsi_pk-gsi_sk-index;non-projective:[lsi_sk-index]"`
	ProjectiveAttr1 string
	ProjectiveAttr2 string
	NonProjective   string `dynmgrm:"non-projective:[lsi_sk-index,gsi_pk-gsi_sk-index]"`
}

func (t CreateTableATable) TableClass() TableClass {
	return TableClassStandard
}

func (t CreateTableATable) WCU() int {
	return 10
}

func (t CreateTableATable) RCU() int {
	return 10
}

func TestMigrator_CreateTable(t *testing.T) {
	type test struct {
		args                                []interface{}
		mockDBExecOptions                   []func(*mockDBExecProp)
		mockBaseMigratorCurrentTableOptions []func(*mockBaseMigratorCurrentTableProp)
		want                                error
	}
	errDBExec := errors.New("db exec error")
	tests := map[string]test{
		"happy_path/pointer": {
			args: []interface{}{&CreateTableATable{}},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE TABLE IF NOT EXISTS create_table_a_tables WITH PK=pk:string, WITH SK=sk:number, WITH LSI=lsi_sk-index:lsi_sk:binary:projective_attr_1,projective_attr_2, WITH wcu=10, WITH rcu=10, WITH table-class=STANDARD`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
		},
		"happy_path/physical": {
			args: []interface{}{CreateTableATable{}},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE TABLE IF NOT EXISTS create_table_a_tables WITH PK=pk:string, WITH SK=sk:number, WITH LSI=lsi_sk-index:lsi_sk:binary:projective_attr_1,projective_attr_2, WITH wcu=10, WITH rcu=10, WITH table-class=STANDARD`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
		},
		"unhappy_path/db_exec_returns_error": {
			args: []interface{}{&CreateTableATable{}, CreateTableATable{}},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE TABLE IF NOT EXISTS create_table_a_tables WITH PK=pk:string, WITH SK=sk:number, WITH LSI=lsi_sk-index:lsi_sk:binary:projective_attr_1,projective_attr_2, WITH wcu=10, WITH rcu=10, WITH table-class=STANDARD`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{Error: errDBExec}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
			want: errDBExec,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mdb := mocks.NewMockDBForMigrator(ctrl)
			mbm := mocks.NewMockBaseMigrator(ctrl)
			mbm.EXPECT().RunWithValue(gomock.Any(), gomock.Any()).DoAndReturn(
				func(model interface{}, f func(stmt *gorm.Statement) error) error {
					return f(nil)
				},
			).Times(1)

			ctp := mockBaseMigratorCurrentTableProp{}
			for _, o := range tt.mockBaseMigratorCurrentTableOptions {
				o(&ctp)
			}
			mbm.EXPECT().CurrentTable(gomock.Any()).Return(ctp.result).Times(ctp.times)

			ep := mockDBExecProp{}
			for _, o := range tt.mockDBExecOptions {
				o(&ep)
			}
			var dbExecCall *gomock.Call
			if ep.args != nil {
				dbExecCall = mdb.EXPECT().Exec(ep.args.sql, ep.args.values...)
			} else {
				dbExecCall = mdb.EXPECT().Exec(gomock.Any(), gomock.Any())
			}
			dbExecCall.Return(ep.result).Times(ep.times)
			sut := Migrator{
				db:   mdb,
				base: mbm,
			}
			err := sut.CreateTable(tt.args...)
			if !errors.Is(err, tt.want) {
				t.Errorf("CreateTable() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func mockDBForMigratorExecWithArgs(t *testing.T, args mockDBExecArgs) func(*mockDBExecProp) {
	t.Helper()
	return func(prop *mockDBExecProp) {
		prop.args = &args
	}
}

func mockDBForMigratorExecWithResult(t *testing.T, tx *gorm.DB) func(*mockDBExecProp) {
	t.Helper()
	return func(prop *mockDBExecProp) {
		prop.result = tx
	}
}

func mockDBForMigratorExecWithTimes(t *testing.T, times int) func(*mockDBExecProp) {
	t.Helper()
	return func(prop *mockDBExecProp) {
		prop.times = times
	}
}

func mockBaseMigratorCurrentTableWithResult(t *testing.T, result clause.Table) func(*mockBaseMigratorCurrentTableProp) {
	t.Helper()
	return func(prop *mockBaseMigratorCurrentTableProp) {
		prop.result = result
	}
}

func mockBaseMigratorCurrentTableWithTimes(t *testing.T, times int) func(*mockBaseMigratorCurrentTableProp) {
	t.Helper()
	return func(prop *mockBaseMigratorCurrentTableProp) {
		prop.times = times
	}
}

func TestMigrator_CreateIndex(t *testing.T) {
	type args struct {
		dest interface{}
		name string
	}
	type test struct {
		args                                args
		mockDBExecOptions                   []func(*mockDBExecProp)
		mockBaseMigratorCurrentTableOptions []func(*mockBaseMigratorCurrentTableProp)
		want                                error
	}
	errDBExec := errors.New("db exec error")
	tests := map[string]test{
		"happy_path/pointer": {
			args: args{&CreateTableATable{}, ""},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE GSI IF NOT EXISTS gsi_pk-gsi_sk-index ON create_table_a_tables WITH PK=gsi_pk:string, WITH SK=gsi_sk:string, WITH wcu=10, WITH rcu=10, WITH projection=lsi_sk,projective_attr_1,projective_attr_2,pk,sk`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
		},
		"happy_path/physical": {
			args: args{&CreateTableATable{}, ""},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE GSI IF NOT EXISTS gsi_pk-gsi_sk-index ON create_table_a_tables WITH PK=gsi_pk:string, WITH SK=gsi_sk:string, WITH wcu=10, WITH rcu=10, WITH projection=lsi_sk,projective_attr_1,projective_attr_2,pk,sk`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
		},
		"happy_path/with_name": {
			args: args{&CreateTableATable{}, "gsi_pk-gsi_sk-index"},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE GSI IF NOT EXISTS gsi_pk-gsi_sk-index ON create_table_a_tables WITH PK=gsi_pk:string, WITH SK=gsi_sk:string, WITH wcu=10, WITH rcu=10, WITH projection=lsi_sk,projective_attr_1,projective_attr_2,pk,sk`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
		},
		"unhappy_path/db_exec_returns_error": {
			args: args{&CreateTableATable{}, ""},
			mockDBExecOptions: []func(*mockDBExecProp){
				mockDBForMigratorExecWithArgs(t,
					mockDBExecArgs{
						`CREATE GSI IF NOT EXISTS gsi_pk-gsi_sk-index ON create_table_a_tables WITH PK=gsi_pk:string, WITH SK=gsi_sk:string, WITH wcu=10, WITH rcu=10, WITH projection=lsi_sk,projective_attr_1,projective_attr_2,pk,sk`,
						nil}),
				mockDBForMigratorExecWithTimes(t, 1),
				mockDBForMigratorExecWithResult(t, &gorm.DB{Error: errDBExec}),
			},
			mockBaseMigratorCurrentTableOptions: []func(*mockBaseMigratorCurrentTableProp){
				mockBaseMigratorCurrentTableWithResult(t, clause.Table{Name: "create_table_a_tables"}),
				mockBaseMigratorCurrentTableWithTimes(t, 1),
			},
			want: errDBExec,
		},
		"unhappy_path/with-undefined-gsi-name": {
			args: args{&CreateTableATable{}, "test"},
			want: fmt.Errorf("gsi '%s' is not defined in %T", `test`, &CreateTableATable{}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mdb := mocks.NewMockDBForMigrator(ctrl)
			mbm := mocks.NewMockBaseMigrator(ctrl)
			mbm.EXPECT().RunWithValue(gomock.Any(), gomock.Any()).DoAndReturn(
				func(model interface{}, f func(stmt *gorm.Statement) error) error {
					return f(nil)
				},
			).Times(1)

			ctp := mockBaseMigratorCurrentTableProp{}
			for _, o := range tt.mockBaseMigratorCurrentTableOptions {
				o(&ctp)
			}
			mbmxp := mbm.EXPECT().CurrentTable(gomock.Any())
			if cmp.Diff(ctp.result, clause.Table{}) != "" {
				mbmxp = mbmxp.Return(ctp.result)
			}
			mbmxp.Times(ctp.times)

			ep := mockDBExecProp{}
			for _, o := range tt.mockDBExecOptions {
				o(&ep)
			}

			var dbExecCall *gomock.Call
			if ep.args != nil {
				dbExecCall = mdb.EXPECT().Exec(ep.args.sql, ep.args.values...)
			} else {
				dbExecCall = mdb.EXPECT().Exec(gomock.Any(), gomock.Any())
			}
			dbExecCall.Return(ep.result).Times(ep.times)
			sut := Migrator{
				db:   mdb,
				base: mbm,
			}
			err := sut.CreateIndex(tt.args.dest, tt.args.name)
			if !errors.Is(err, tt.want) && err.Error() != tt.want.Error() {
				t.Errorf("CreateIndex() error = %v, want %v", err, tt.want)
			}
		})
	}
}
