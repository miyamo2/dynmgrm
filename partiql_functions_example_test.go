package dynmgrm_test

import (
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
	"log"
)

func ExampleListAppend() {
	db, err := gorm.Open(
		dynmgrm.New(),
		&gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		log.Fatalf("failed to open database, got error %v", err)
	}
	db.Model(&TestTable{PK: "Partition1", SK: 1}).
		Update("list_type_attr",
			dynmgrm.ListAppend(
				dynmgrm.List{dynmgrm.Map{"Foo": "Bar"}}))
}
