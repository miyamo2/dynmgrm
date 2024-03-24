package tests

import (
	"fmt"
	"testing"
)

func Test_Delete(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect[1:], testTableName)

	data := TestTable{
		PK: "Partition1",
		SK: 1,
	}
	db.Delete(TestTable{}, `pk = ? AND sk = ?`, data.PK, data.SK)

	actual := scanData(t, testTableName)
	for _, item := range actual {
		if *item["pk"].S == data.PK && *item["sk"].N == fmt.Sprint(data.SK) {
			t.Errorf("data has not been deleted: %v", data)
		}
	}
}

func Test_Delete_With_Where_Clause(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect[1:], testTableName)

	data := TestTable{
		PK: "Partition1",
		SK: 1,
	}
	db.Where(`pk = ? AND sk = ?`, data.PK, data.SK).Delete(TestTable{})

	actual := scanData(t, testTableName)
	for _, item := range actual {
		if *item["pk"].S == data.PK && *item["sk"].N == fmt.Sprint(data.SK) {
			t.Errorf("data has not been deleted: %v", data)
		}
	}
}

func Test_Delete_With_Tx_Commit(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect[1:], testTableName)

	data := TestTable{
		PK: "Partition1",
		SK: 1,
	}
	tx := db.Begin()
	tx.Delete(TestTable{}, `pk = ? AND sk = ?`, data.PK, data.SK)
	tx.Commit()

	actual := scanData(t, testTableName)
	for _, item := range actual {
		if *item["pk"].S == data.PK && *item["sk"].N == fmt.Sprint(data.SK) {
			t.Errorf("data has not been deleted: %v", data)
		}
	}
}

func Test_Delete_With_Tx_Rollback(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	data := TestTable{
		PK: "Partition1",
		SK: 1,
	}
	tx := db.Begin()
	tx.Delete(TestTable{}, `pk = ? AND sk = ?`, data.PK, data.SK)
	tx.Rollback()

	actual := scanData(t, testTableName)
	exist := false
	for _, item := range actual {
		if *item["pk"].S == data.PK && *item["sk"].N == fmt.Sprint(data.SK) {
			exist = true
		}
	}
	if !exist {
		t.Errorf("data has been deleted: %v", data)
	}
}
