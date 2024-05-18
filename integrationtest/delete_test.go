package integrationtest

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"testing"
)

func Test_Delete(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData[1:], testTableName)

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
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData[1:], testTableName)

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
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData[1:], testTableName)

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
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

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

func Test_Delete_With_Transaction_Success(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData[1:], testTableName)

	data := TestTable{
		PK: "Partition1",
		SK: 1,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(TestTable{}, `pk = ? AND sk = ?`, data.PK, data.SK).Error
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	actual := scanData(t, testTableName)
	for _, item := range actual {
		if *item["pk"].S == data.PK && *item["sk"].N == fmt.Sprint(data.SK) {
			t.Errorf("data has not been deleted: %v", data)
		}
	}
}

func Test_Delete_With_Transaction_Fail(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	data := TestTable{
		PK: "Partition1",
		SK: 1,
	}

	errFoo := errors.New("foo")
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Delete(TestTable{}, `pk = ? AND sk = ?`, data.PK, data.SK)
		return errFoo
	})

	if !errors.Is(err, errFoo) {
		t.Errorf("unexpected error: %v", err)
	}

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
