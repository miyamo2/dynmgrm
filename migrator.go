package dynmgrm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// compatibility check
var _ gorm.Migrator = (*Migrator)(nil)

// Migrator is gorm.Migrator implementation for dynamodb
//
// Deprecated: Migrator is not implemented.
type Migrator struct{}

func (m Migrator) AutoMigrate(dst ...interface{}) error {
	return ErrDynmgrmAreNotSupported
}

func (m Migrator) CurrentDatabase() string {
	return ""
}

func (m Migrator) FullDataTypeOf(field *schema.Field) clause.Expr {
	return clause.Expr{}
}

func (m Migrator) GetTypeAliases(databaseTypeName string) []string {
	return []string{}
}

func (m Migrator) CreateTable(dst ...interface{}) error {
	return ErrDynmgrmAreNotSupported
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
