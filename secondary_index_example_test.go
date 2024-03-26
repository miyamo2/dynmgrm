package dynmgrm_test

import (
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TestTable struct {
	PK     string `gorm:"primaryKey"`
	SK     int    `gorm:"primaryKey"`
	GSIKey string
}

func ExampleSecondaryIndex() {
	db, err := gorm.Open(dynmgrm.New(), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	result := TestTable{}
	db.Table("something").Clauses(
		dynmgrm.SecondaryIndex("gsi_key-sk-index")).
		Where(`gsi_key = ?`, "1").
		Scan(&result)
}

func ExampleSecondaryIndex_withSecondaryIndexOf_withString() {
	db, err := gorm.Open(dynmgrm.New(), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	result := TestTable{}
	db.Clauses(
		dynmgrm.SecondaryIndex("gsi_key-sk-index",
			dynmgrm.SecondaryIndexOf("something"))).
		Where(`gsi_key = ?`, "1").
		Scan(&result)
}

func ExampleSecondaryIndex_withSecondaryIndexOf_withTableClause() {
	db, err := gorm.Open(dynmgrm.New(), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	result := TestTable{}
	db.Clauses(
		dynmgrm.SecondaryIndex("gsi_key-sk-index",
			dynmgrm.SecondaryIndexOf(
				clause.Table{
					Name: "something",
				}))).
		Where(`gsi_key = ?`, "1").
		Scan(&result)
}
