package dynmgrm

import (
	"database/sql"
)

type dbOpener struct {
	dsn        string
	driverName string
}

func (o dbOpener) DSN() string {
	return o.dsn
}

func (o dbOpener) DriverName() string {
	return o.driverName
}

func (o dbOpener) Apply() (*sql.DB, error) {
	return sql.Open(o.DriverName(), o.DSN())
}
