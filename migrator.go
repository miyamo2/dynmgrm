package dynmgrm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
	"reflect"
	"slices"
	"strings"
)

// CapacityUnitsSpecifier could specify WCUs and RCU
type CapacityUnitsSpecifier interface {
	WCU() int
	RCU() int
}

// TableClass is the type of table class
type TableClass int

func (t TableClass) String() string {
	switch t {
	case TableClassStandard:
		return "STANDARD"
	case TableClassStandardIA:
		return "STANDARD_IA"
	}
	return ""
}

// TableClassStandard and TableClassStandardIA are the supported table classes
const (
	TableClassStandard TableClass = iota
	TableClassStandardIA
)

// TableClassSpecifier could specify table class.
type TableClassSpecifier interface {
	TableClass() TableClass
}

type dbForMigrator interface {
	AddError(err error) error
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
}

// compatibility check
var _ gorm.Migrator = (*Migrator)(nil)

// Migrator is gorm.Migrator implementation for dynamodb
type Migrator struct {
	db   dbForMigrator
	base migrator.Migrator
}

func (m Migrator) CurrentDatabase() string {
	return ""
}

func (m Migrator) AutoMigrate(dst ...interface{}) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) FullDataTypeOf(field *schema.Field) clause.Expr {
	return m.base.FullDataTypeOf(field)
}

func (m Migrator) GetTypeAliases(databaseTypeName string) []string {
	return []string{}
}

func (m Migrator) CreateTable(models ...interface{}) error {
	for _, model := range models {
		err := m.base.RunWithValue(model, func(stmt *gorm.Statement) (err error) {
			var (
				wcu, rcu   int
				tableClass string
			)

			if ws, ok := model.(CapacityUnitsSpecifier); ok {
				wcu = ws.WCU()
				rcu = ws.RCU()
			}
			if tcs, ok := model.(TableClassSpecifier); ok {
				tableClass = tcs.TableClass().String()
			}

			rv := reflect.ValueOf(model)
			var rt reflect.Type
			switch rv.Kind() {
			case reflect.Ptr:
				rt = reflect.TypeOf(reflect.ValueOf(model).Elem().Interface())
			default:
				rt = reflect.TypeOf(model)
			}

			ddlBuilder := strings.Builder{}
			ddlBuilder.WriteString(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s`, m.currentTable(stmt)))

			td := newDynmgrmTableDefine(rt)
			// `CREATE TABLE` are proprietary PartiQL syntax by btnguyen2k/godynamo
			// This is why place holder/bind variables are not used.
			ddlBuilder.WriteString(fmt.Sprintf(` WITH PK=%s:%s`, td.PK.Name, td.PK.DataType))
			if skn := td.SK.Name; skn != "" {
				ddlBuilder.WriteString(fmt.Sprintf(`, WITH SK=%s:%s`, skn, td.SK.DataType))
			}
			opts := make([]string, 0, 7)
			if wcu > 0 {
				opts = append(opts, fmt.Sprintf(`WITH wcu=%d`, wcu))
			}
			if rcu > 0 {
				opts = append(opts, fmt.Sprintf(`WITH rcu=%d`, rcu))
			}
			if tableClass != "" {
				opts = append(opts, fmt.Sprintf(`WITH table-class=%s`, tableClass))
			}
			for k, v := range td.LSI {
				lsi := fmt.Sprintf(`WITH LSI=%s:%s:%s`, k, v.SK.Name, v.SK.DataType)
				projective := slices.DeleteFunc(td.NonKeyAttr, func(s string) bool {
					if s == v.SK.Name {
						return true
					}
					return slices.Contains(v.NonProjectiveAttrs, s)
				})
				if len(projective) > 0 {
					lsi += fmt.Sprintf(`:%s`, strings.Join(projective, ","))
				}
				opts = append(opts, lsi)
			}

			for _, o := range opts {
				ddlBuilder.WriteString(", ")
				ddlBuilder.WriteString(o)
			}
			if err := m.db.Exec(ddlBuilder.String()).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) currentTable(stmt *gorm.Statement) string {
	txp := m.base.CurrentTable(stmt)
	if txp, ok := txp.(clause.Table); ok {
		return txp.Name
	}
	return ""
}

func (m Migrator) DropTable(dst ...interface{}) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) HasTable(dst interface{}) bool {
	return false
}

func (m Migrator) RenameTable(oldName, newName interface{}) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) GetTables() (tableList []string, err error) {
	return
}

func (m Migrator) TableType(dst interface{}) (gorm.TableType, error) {
	return nil, ErrDynmgrmAreNotSupported
}

func (m Migrator) AddColumn(dst interface{}, field string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) DropColumn(dst interface{}, field string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) AlterColumn(dst interface{}, field string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) MigrateColumn(dst interface{}, field *schema.Field, columnType gorm.ColumnType) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) MigrateColumnUnique(dst interface{}, field *schema.Field, columnType gorm.ColumnType) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) HasColumn(dst interface{}, field string) bool {
	return false
}

func (m Migrator) RenameColumn(dst interface{}, oldName, field string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) ColumnTypes(dst interface{}) ([]gorm.ColumnType, error) {
	return []gorm.ColumnType{}, ErrDynmgrmAreNotSupported
}

func (m Migrator) CreateView(name string, option gorm.ViewOption) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) DropView(name string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) CreateConstraint(dst interface{}, name string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) DropConstraint(dst interface{}, name string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) HasConstraint(dst interface{}, name string) bool {
	return false
}

func (m Migrator) CreateIndex(dst interface{}, name string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) DropIndex(dst interface{}, name string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) HasIndex(dst interface{}, name string) bool {
	return false
}

func (m Migrator) RenameIndex(dst interface{}, oldName, newName string) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) GetIndexes(dst interface{}) ([]gorm.Index, error) {
	return []gorm.Index{}, ErrDynmgrmAreNotSupported
}
