package tests

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
	"testing"
)

func Test_Update_With_Save(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

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
		"some_string_sets": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_sets": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_sets": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_sets": {
			BS: [][]byte{[]byte("XYZ"), []byte("ABC")},
		},
		"any": {
			S: aws.String("UPDATED"),
		},
	}

	db.Save(TestTable{
		PK:             "Partition1",
		SK:             1,
		SomeString:     "UPDATED",
		SomeInt:        5,
		SomeFloat:      5.5,
		SomeBool:       false,
		SomeBinary:     []byte("XYZ"),
		SomeList:       dynmgrm.List{"UPDATED", 5, 5.5, false, []byte("XYZ")},
		SomeMap:        dynmgrm.Map{"some_string": "UPDATED", "some_number": 5.5, "some_bool": false, "some_binary": []byte("XYZ")},
		SomeStringSets: dynmgrm.Sets[string]{"UPDATED"},
		SomeIntSets:    dynmgrm.Sets[int]{5, 10},
		SomeFloatSets:  dynmgrm.Sets[float64]{5.5, 11.0},
		SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("XYZ"), []byte("ABC")},
		Any:            "UPDATED",
	})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setsCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Updates(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

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
		"some_string_sets": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_sets": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_sets": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_sets": {
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
			"some_string_sets": dynmgrm.Sets[string]{"UPDATED"},
			"some_int_sets":    dynmgrm.Sets[int]{5, 10},
			"some_float_sets":  dynmgrm.Sets[float64]{5.5, 11.0},
			"some_binary_sets": dynmgrm.Sets[[]byte]{[]byte("XYZ"), []byte("ABC")},
			"any":              "UPDATED",
		})

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setsCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Update_Clause(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

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
		"some_string_sets": {
			SS: []*string{aws.String("Hello"), aws.String("World")},
		},
		"some_int_sets": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_sets": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_sets": {
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

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setsCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Tx_Commit(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

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
		"some_string_sets": {
			SS: []*string{aws.String("UPDATED")},
		},
		"some_int_sets": {
			NS: []*string{aws.String("5"), aws.String("10")},
		},
		"some_float_sets": {
			NS: []*string{aws.String("5.5"), aws.String("11")},
		},
		"some_binary_sets": {
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
			"some_string_sets": dynmgrm.Sets[string]{"UPDATED"},
			"some_int_sets":    dynmgrm.Sets[int]{5, 10},
			"some_float_sets":  dynmgrm.Sets[float64]{5.5, 11.0},
			"some_binary_sets": dynmgrm.Sets[[]byte]{[]byte("XYZ"), []byte("ABC")},
			"any":              "UPDATED",
		})
	tx.Commit()

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setsCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_Tx_Rollback(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expect := testDataForSelect[0]

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
			"some_string_sets": dynmgrm.Sets[string]{"UPDATED"},
			"some_int_sets":    dynmgrm.Sets[int]{5, 10},
			"some_float_sets":  dynmgrm.Sets[float64]{5.5, 11.0},
			"some_binary_sets": dynmgrm.Sets[[]byte]{[]byte("XYZ"), []byte("ABC")},
			"any":              "UPDATED",
		})
	tx.Rollback()

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setsCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func Test_Update_With_SetAdd(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

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
		"some_string_sets": {
			SS: []*string{aws.String("Hello"), aws.String("World"), aws.String("こんにちは"), aws.String("世界")},
		},
		"some_int_sets": {
			NS: []*string{aws.String("1"), aws.String("2")},
		},
		"some_float_sets": {
			NS: []*string{aws.String("1.1"), aws.String("2.2")},
		},
		"some_binary_sets": {
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
		}).Update("some_string_sets", gorm.Expr("set_add(some_string_sets, ?)", dynmgrm.Sets[string]{"こんにちは", "世界"}))

	result := getData(t, testTableName, "Partition1", 1)

	if diff := cmp.Diff(expect, result, append(avCmpOpts, setsCmpOpts...)...); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}
