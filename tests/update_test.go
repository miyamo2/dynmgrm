package tests

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
	"testing"
)

func Test_Update_With_Save(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("UPDATED"),
		},
		"some_int": {
			N: aws.String("5"),
		},
		"some_float": {
			N: aws.String("5.5"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("XYZ"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("UPDATED"),
				},
				{
					N: aws.String("5"),
				},
				{
					N: aws.String("5.5"),
				},
				{
					BOOL: aws.Bool(false),
				},
				{
					B: []byte("XYZ"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("UPDATED"),
				},
				"some_number": {
					N: aws.String("5.5"),
				},
				"some_bool": {
					BOOL: aws.Bool(false),
				},
				"some_binary": {
					B: []byte("XYZ"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_set": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_set": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("XYZ"), []byte("ABC")},
		},
		"any": {
			S: aws.String("UPDATED"),
		},
	}

	db.Save(TestTable{
		PK:            "Partition1",
		SK:            1,
		SomeString:    "UPDATED",
		SomeInt:       5,
		SomeFloat:     5.5,
		SomeBool:      false,
		SomeBinary:    []byte("XYZ"),
		SomeList:      dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
		SomeMap:       dynmgrm.Map{"some_string": "UPDATED", "some_number": 5.5, "some_bool": false, "some_binary": []byte("XYZ")},
		SomeStringSet: dynmgrm.Set[string]{"UPDATED"},
		SomeIntSet:    dynmgrm.Set[int]{5, 10},
		SomeFloatSet:  dynmgrm.Set[float64]{5.5, 11.0},
		SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("XYZ"), []byte("ABC")},
		Any:           "UPDATED",
	})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Updates(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("UPDATED"),
		},
		"some_int": {
			N: aws.String("5"),
		},
		"some_float": {
			N: aws.String("5.5"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("XYZ"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("UPDATED"),
				},
				{
					N: aws.String("5"),
				},
				{
					N: aws.String("5.5"),
				},
				{
					BOOL: aws.Bool(false),
				},
				{
					B: []byte("XYZ"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("UPDATED"),
				},
				"some_number": {
					N: aws.String("5.5"),
				},
				"some_bool": {
					BOOL: aws.Bool(false),
				},
				"some_binary": {
					B: []byte("XYZ"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_set": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_set": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("XYZ"), []byte("ABC")},
		},
		"any": {
			S: aws.String("UPDATED"),
		},
	}

	// zero-value is omitted if the argument is a struct.
	// This is the behavior of Updates.
	// https://gorm.io/docs/update.html#Update-Changed-Fields
	db.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Updates(
		map[string]interface{}{
			"some_string": "UPDATED",
			"some_int":    5,
			"some_float":  5.5,
			"some_bool":   false,
			"some_binary": []byte("XYZ"),
			"some_list":   dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
			"some_map": dynmgrm.Map{
				"some_string": "UPDATED",
				"some_number": 5.5,
				"some_bool":   false,
				"some_binary": []byte("XYZ")},
			"some_string_set": dynmgrm.Set[string]{"UPDATED"},
			"some_int_set":    dynmgrm.Set[int]{5, 10},
			"some_float_set":  dynmgrm.Set[float64]{5.5, 11.0},
			"some_binary_set": dynmgrm.Set[[]byte]{[]byte("XYZ"), []byte("ABC")},
			"any":             "UPDATED",
		})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Update_Clause(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("UPDATED"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	}

	db.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Update("some_string", "UPDATED")

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Tx_Commit(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("UPDATED"),
		},
		"some_int": {
			N: aws.String("5"),
		},
		"some_float": {
			N: aws.String("5.5"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("XYZ"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("UPDATED"),
				},
				{
					N: aws.String("5"),
				},
				{
					N: aws.String("5.5"),
				},
				{
					BOOL: aws.Bool(false),
				},
				{
					B: []byte("XYZ"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("UPDATED"),
				},
				"some_number": {
					N: aws.String("5.5"),
				},
				"some_bool": {
					BOOL: aws.Bool(false),
				},
				"some_binary": {
					B: []byte("XYZ"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_set": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_set": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("XYZ"), []byte("ABC")},
		},
		"any": {
			S: aws.String("UPDATED"),
		},
	}

	tx := db.Begin()
	tx.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Updates(
		map[string]interface{}{
			"some_string": "UPDATED",
			"some_int":    5,
			"some_float":  5.5,
			"some_bool":   false,
			"some_binary": []byte("XYZ"),
			"some_list":   dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
			"some_map": dynmgrm.Map{
				"some_string": "UPDATED",
				"some_number": 5.5,
				"some_bool":   false,
				"some_binary": []byte("XYZ")},
			"some_string_set": dynmgrm.Set[string]{"UPDATED"},
			"some_int_set":    dynmgrm.Set[int]{5, 10},
			"some_float_set":  dynmgrm.Set[float64]{5.5, 11.0},
			"some_binary_set": dynmgrm.Set[[]byte]{[]byte("XYZ"), []byte("ABC")},
			"any":             "UPDATED",
		})
	tx.Commit()

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Tx_Rollback(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := testData[0]

	tx := db.Begin()
	tx.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Updates(
		map[string]interface{}{
			"some_string": "UPDATED",
			"some_int":    5,
			"some_float":  5.5,
			"some_bool":   false,
			"some_binary": []byte("XYZ"),
			"some_list":   dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
			"some_map": dynmgrm.Map{
				"some_string": "UPDATED",
				"some_number": 5.5,
				"some_bool":   false,
				"some_binary": []byte("XYZ")},
			"some_string_set": dynmgrm.Set[string]{"UPDATED"},
			"some_int_set":    dynmgrm.Set[int]{5, 10},
			"some_float_set":  dynmgrm.Set[float64]{5.5, 11.0},
			"some_binary_set": dynmgrm.Set[[]byte]{[]byte("XYZ"), []byte("ABC")},
			"any":             "UPDATED",
		})
	tx.Rollback()

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_SetAdd(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World"), aws.String("Bye")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	}

	db.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Update("some_string_set", gorm.Expr("set_add(some_string_set, ?)", dynmgrm.Set[string]{"Bye"}))

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_SetDelete(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	}

	db.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Update("some_string_set", gorm.Expr("set_delete(some_string_set, ?)", dynmgrm.Set[string]{"Hello"}))

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_ListAppend(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
				{
					M: map[string]*dynamodb.AttributeValue{"append_item": {S: aws.String("Foo")}},
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	}

	db.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Update("some_list", gorm.Expr("list_append(some_list, ?)", dynmgrm.List{dynmgrm.Map{"append_item": "Foo"}}))

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Transaction_Success(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("UPDATED"),
		},
		"some_int": {
			N: aws.String("5"),
		},
		"some_float": {
			N: aws.String("5.5"),
		},
		"some_bool": {
			BOOL: aws.Bool(false),
		},
		"some_binary": {
			B: []byte("XYZ"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("UPDATED"),
				},
				{
					N: aws.String("5"),
				},
				{
					N: aws.String("5.5"),
				},
				{
					BOOL: aws.Bool(false),
				},
				{
					B: []byte("XYZ"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("UPDATED"),
				},
				"some_number": {
					N: aws.String("5.5"),
				},
				"some_bool": {
					BOOL: aws.Bool(false),
				},
				"some_binary": {
					B: []byte("XYZ"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_set": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_set": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("XYZ"), []byte("ABC")},
		},
		"any": {
			S: aws.String("UPDATED"),
		},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(
			&TestTable{
				PK: "Partition1",
				SK: 1,
			}).Updates(
			map[string]interface{}{
				"some_string": "UPDATED",
				"some_int":    5,
				"some_float":  5.5,
				"some_bool":   false,
				"some_binary": []byte("XYZ"),
				"some_list":   dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
				"some_map": dynmgrm.Map{
					"some_string": "UPDATED",
					"some_number": 5.5,
					"some_bool":   false,
					"some_binary": []byte("XYZ")},
				"some_string_set": dynmgrm.Set[string]{"UPDATED"},
				"some_int_set":    dynmgrm.Set[int]{5, 10},
				"some_float_set":  dynmgrm.Set[float64]{5.5, 11.0},
				"some_binary_set": dynmgrm.Set[[]byte]{[]byte("XYZ"), []byte("ABC")},
				"any":             "UPDATED",
			}).Error
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Transaction_Fail(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	errFoo := errors.New("foo")

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Model(
			&TestTable{
				PK: "Partition1",
				SK: 1,
			}).Updates(
			map[string]interface{}{
				"some_string": "UPDATED",
				"some_int":    5,
				"some_float":  5.5,
				"some_bool":   false,
				"some_binary": []byte("XYZ"),
				"some_list":   dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
				"some_map": dynmgrm.Map{
					"some_string": "UPDATED",
					"some_number": 5.5,
					"some_bool":   false,
					"some_binary": []byte("XYZ")},
				"some_string_set": dynmgrm.Set[string]{"UPDATED"},
				"some_int_set":    dynmgrm.Set[int]{5, 10},
				"some_float_set":  dynmgrm.Set[float64]{5.5, 11.0},
				"some_binary_set": dynmgrm.Set[[]byte]{[]byte("XYZ"), []byte("ABC")},
				"any":             "UPDATED",
			})
		return errFoo
	})
	if !errors.Is(err, errFoo) {
		t.Errorf("unexpected error: %v", err)
	}

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Save_Has_TypedList_Column(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForTypedList, testTableName)
	defer dataCleanup(t, testDataForTypedList, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"typed_list": {
			L: []*dynamodb.AttributeValue{
				{
					M: map[string]*dynamodb.AttributeValue{
						"some_string": {S: aws.String("World")},
						"some_int":    {N: aws.String("2")},
						"some_float":  {N: aws.String("2.2")},
						"some_bool":   {BOOL: aws.Bool(false)},
						"some_binary": {B: []byte("DEF")},
						"some_map": {
							M: map[string]*dynamodb.AttributeValue{
								"k": {S: aws.String("v")},
							}},
						"some_list": {
							L: []*dynamodb.AttributeValue{
								{S: aws.String("b")},
								{M: map[string]*dynamodb.AttributeValue{"k": {S: aws.String("v")}}},
								{SS: aws.StringSlice([]string{"d", "e", "f"})},
							},
						},
						"some_string_set": {SS: aws.StringSlice([]string{"d", "e", "f"})},
						"some_int_set":    {NS: aws.StringSlice([]string{"4", "5", "6"})},
						"some_float_set":  {NS: aws.StringSlice([]string{"4.4", "5.5", "6.6"})},
						"some_binary_set": {BS: [][]byte{[]byte("d"), []byte("e"), []byte("f")}},
					},
				},
			},
		},
	}

	db.Save(TestTableWithTypedList{
		PK:         "Partition1",
		SK:         1,
		SomeString: "Hello",
		TypedList: dynmgrm.TypedList[TypedListValue]{
			{
				SomeString:    "World",
				SomeInt:       2,
				SomeFloat:     2.2,
				SomeBool:      false,
				SomeBinary:    []byte("DEF"),
				SomeMap:       dynmgrm.Map{"k": "v"},
				SomeList:      dynmgrm.List{"b", dynmgrm.Map{"k": "v"}, dynmgrm.Set[string]{"d", "e", "f"}},
				SomeStringSet: dynmgrm.Set[string]{"d", "e", "f"},
				SomeIntSet:    dynmgrm.Set[int]{4, 5, 6},
				SomeFloatSet:  dynmgrm.Set[float64]{4.4, 5.5, 6.6},
				SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("d"), []byte("e"), []byte("f")},
			},
		},
	})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Updates_Has_TypedList_Column(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForTypedList, testTableName)
	defer dataCleanup(t, testDataForTypedList, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk":          {S: aws.String("Partition1")},
		"sk":          {N: aws.String("1")},
		"some_string": {S: aws.String("Hello")},
		"typed_list": {
			L: []*dynamodb.AttributeValue{
				{
					M: map[string]*dynamodb.AttributeValue{
						"some_string": {S: aws.String("UPDATED")},
						"some_int":    {N: aws.String("1")},
						"some_float":  {N: aws.String("1.1")},
						"some_bool":   {BOOL: aws.Bool(true)},
						"some_binary": {B: []byte("ABC")},
						"some_map": {
							M: map[string]*dynamodb.AttributeValue{
								"UPDATED": {S: aws.String("UPDATED")},
							}},
						"some_list": {
							L: []*dynamodb.AttributeValue{
								{S: aws.String("UPDATED")},
								{M: map[string]*dynamodb.AttributeValue{"UPDATED": {S: aws.String("UPDATED")}}},
								{SS: aws.StringSlice([]string{"UPDATED"})},
							},
						},
						"some_string_set": {SS: aws.StringSlice([]string{"UPDATED"})},
						"some_int_set":    {NS: aws.StringSlice([]string{"4", "5", "6"})},
						"some_float_set":  {NS: aws.StringSlice([]string{"2.2", "4.4", "6.6"})},
						"some_binary_set": {BS: [][]byte{[]byte("d"), []byte("e"), []byte("f")}},
					},
				},
			},
		},
	}

	// zero-value is omitted if the argument is a struct.
	// This is the behavior of Updates.
	// https://gorm.io/docs/update.html#Update-Changed-Fields
	db.Model(
		&TestTableWithTypedList{
			PK: "Partition1",
			SK: 1,
		}).Updates(
		map[string]interface{}{
			"typed_list": dynmgrm.TypedList[TypedListValue]{
				{
					SomeString: "UPDATED",
					SomeInt:    1,
					SomeFloat:  1.1,
					SomeBool:   true,
					SomeBinary: []byte("ABC"),
					SomeMap: dynmgrm.Map{
						"UPDATED": "UPDATED",
					},
					SomeList: dynmgrm.List{
						"UPDATED",
						dynmgrm.Map{
							"UPDATED": "UPDATED",
						},
						dynmgrm.Set[string]{"UPDATED"},
					},
					SomeStringSet: dynmgrm.Set[string]{"UPDATED"},
					SomeIntSet:    dynmgrm.Set[int]{4, 5, 6},
					SomeFloatSet:  dynmgrm.Set[float64]{2.2, 4.4, 6.6},
					SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("d"), []byte("e"), []byte("f")},
				},
			},
		})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Update_Clause_TypedList_Column(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForTypedList, testTableName)
	defer dataCleanup(t, testDataForTypedList, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk":          {S: aws.String("Partition1")},
		"sk":          {N: aws.String("1")},
		"some_string": {S: aws.String("Hello")},
		"typed_list": {
			L: []*dynamodb.AttributeValue{
				{
					M: map[string]*dynamodb.AttributeValue{
						"some_string": {S: aws.String("UPDATED")},
						"some_int":    {N: aws.String("1")},
						"some_float":  {N: aws.String("1.1")},
						"some_bool":   {BOOL: aws.Bool(true)},
						"some_binary": {B: []byte("ABC")},
						"some_map": {
							M: map[string]*dynamodb.AttributeValue{
								"UPDATED": {S: aws.String("UPDATED")},
							}},
						"some_list": {
							L: []*dynamodb.AttributeValue{
								{S: aws.String("UPDATED")},
								{M: map[string]*dynamodb.AttributeValue{"UPDATED": {S: aws.String("UPDATED")}}},
								{SS: aws.StringSlice([]string{"UPDATED"})},
							},
						},
						"some_string_set": {SS: aws.StringSlice([]string{"UPDATED"})},
						"some_int_set":    {NS: aws.StringSlice([]string{"4", "5", "6"})},
						"some_float_set":  {NS: aws.StringSlice([]string{"2.2", "4.4", "6.6"})},
						"some_binary_set": {BS: [][]byte{[]byte("d"), []byte("e"), []byte("f")}},
					},
				},
			},
		},
	}

	db.Model(
		&TestTableWithTypedList{
			PK: "Partition1",
			SK: 1,
		}).Update("typed_list", dynmgrm.TypedList[TypedListValue]{
		{
			SomeString: "UPDATED",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeMap: dynmgrm.Map{
				"UPDATED": "UPDATED",
			},
			SomeList: dynmgrm.List{
				"UPDATED",
				dynmgrm.Map{
					"UPDATED": "UPDATED",
				},
				dynmgrm.Set[string]{"UPDATED"},
			},
			SomeStringSet: dynmgrm.Set[string]{"UPDATED"},
			SomeIntSet:    dynmgrm.Set[int]{4, 5, 6},
			SomeFloatSet:  dynmgrm.Set[float64]{2.2, 4.4, 6.6},
			SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("d"), []byte("e"), []byte("f")},
		},
	})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_ListAppend_helper(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testData, testTableName)
	defer dataCleanup(t, testData, testTableName)

	expect := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("Partition1"),
		},
		"sk": {
			N: aws.String("1"),
		},
		"some_string": {
			S: aws.String("Hello"),
		},
		"some_int": {
			N: aws.String("1"),
		},
		"some_float": {
			N: aws.String("1.1"),
		},
		"some_bool": {
			BOOL: aws.Bool(true),
		},
		"some_binary": {
			B: []byte("ABC"),
		},
		"some_list": {
			L: []*dynamodb.AttributeValue{
				{
					S: aws.String("Hello"),
				},
				{
					N: aws.String("1"),
				},
				{
					N: aws.String("1.1"),
				},
				{
					BOOL: aws.Bool(true),
				},
				{
					B: []byte("ABC"),
				},
				{
					M: map[string]*dynamodb.AttributeValue{"append_item": {S: aws.String("Foo")}},
				},
			},
		},
		"some_map": {
			M: map[string]*dynamodb.AttributeValue{
				"some_string": {
					S: aws.String("Hello"),
				},
				"some_number": {
					N: aws.String("1.1"),
				},
				"some_bool": {
					BOOL: aws.Bool(true),
				},
				"some_binary": {
					B: []byte("ABC"),
				},
			},
		},
		"some_string_set": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_set": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_set": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_set": {
			BS: [][]byte{[]byte("ABC"), []byte("DEF")},
		},
		"any": {
			S: aws.String("any"),
		},
	}

	db.Model(
		&TestTable{
			PK: "Partition1",
			SK: 1,
		}).Update("some_list", dynmgrm.ListAppend(dynmgrm.Map{"append_item": "Foo"}))

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}
