package tests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynmgrm"
	"testing"
)

func Test_Insert_With_Create(t *testing.T) {
	db := getGormDB(t)

	data := TestTable{
		PK:            "Partition4",
		SK:            1,
		SomeMap:       dynmgrm.Map{"key": "value"},
		SomeList:      dynmgrm.List{"a", dynmgrm.Map{"key": "value"}, dynmgrm.Set[string]{"a", "b", "c"}},
		SomeStringSet: dynmgrm.Set[string]{"a", "b", "c"},
		SomeIntSet:    dynmgrm.Set[int]{1, 2, 3},
		SomeFloatSet:  dynmgrm.Set[float64]{1.1, 2.2, 3.3},
		SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
	}
	defer deleteData(t, testTableName, data.PK, data.SK)

	db.Create(data)

	expect := []map[string]*dynamodb.AttributeValue{
		{
			"pk": &dynamodb.AttributeValue{S: aws.String(data.PK)},
			"sk": &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(data.SK))},
			"some_map": &dynamodb.AttributeValue{
				M: map[string]*dynamodb.AttributeValue{
					"key": {S: aws.String("value")},
				}},
			"some_list": &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{
					{S: aws.String("a")},
					{M: map[string]*dynamodb.AttributeValue{"key": {S: aws.String("value")}}},
					{SS: aws.StringSlice([]string{"a", "b", "c"})},
				},
			},
			"some_string_set": &dynamodb.AttributeValue{SS: aws.StringSlice([]string{"a", "b", "c"})},
			"some_int_set":    &dynamodb.AttributeValue{NS: aws.StringSlice([]string{"1", "2", "3"})},
			"some_float_set":  &dynamodb.AttributeValue{NS: aws.StringSlice([]string{"1.1", "2.2", "3.3"})},
			"some_binary_set": &dynamodb.AttributeValue{BS: [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
		},
	}
	actual := scanData(t, testTableName)
	if diff := cmp.Diff(expect, actual, setCmpOpts...); diff != "" {
		t.Errorf("mismatch (-got +want)\n%s", diff)
	}
}

func Test_Insert_With_Tx_Commit(t *testing.T) {
	db := getGormDB(t)

	data := TestTable{
		PK:            "Partition4",
		SK:            1,
		SomeMap:       dynmgrm.Map{"key": "value"},
		SomeList:      dynmgrm.List{"a", dynmgrm.Map{"key": "value"}, dynmgrm.Set[string]{"a", "b", "c"}},
		SomeStringSet: dynmgrm.Set[string]{"a", "b", "c"},
		SomeIntSet:    dynmgrm.Set[int]{1, 2, 3},
		SomeFloatSet:  dynmgrm.Set[float64]{1.1, 2.2, 3.3},
		SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
	}
	defer deleteData(t, testTableName, data.PK, data.SK)

	tx := db.Begin()
	tx.Create(data)
	tx.Commit()

	expect := []map[string]*dynamodb.AttributeValue{
		{
			"pk": &dynamodb.AttributeValue{S: aws.String(data.PK)},
			"sk": &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(data.SK))},
			"some_map": &dynamodb.AttributeValue{
				M: map[string]*dynamodb.AttributeValue{
					"key": {S: aws.String("value")},
				}},
			"some_list": &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{
					{S: aws.String("a")},
					{M: map[string]*dynamodb.AttributeValue{"key": {S: aws.String("value")}}},
					{SS: aws.StringSlice([]string{"a", "b", "c"})},
				},
			},
			"some_string_set": &dynamodb.AttributeValue{SS: aws.StringSlice([]string{"a", "b", "c"})},
			"some_int_set":    &dynamodb.AttributeValue{NS: aws.StringSlice([]string{"1", "2", "3"})},
			"some_float_set":  &dynamodb.AttributeValue{NS: aws.StringSlice([]string{"1.1", "2.2", "3.3"})},
			"some_binary_set": &dynamodb.AttributeValue{BS: [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
		},
	}
	actual := scanData(t, testTableName)
	if diff := cmp.Diff(expect, actual, setCmpOpts...); diff != "" {
		t.Errorf("mismatch (-got +want)\n%s", diff)
	}
}

func Test_Insert_With_Tx_Rollback(t *testing.T) {
	db := getGormDB(t)

	data := TestTable{
		PK:            "Partition4",
		SK:            1,
		SomeMap:       dynmgrm.Map{"key": "value"},
		SomeList:      dynmgrm.List{"a", dynmgrm.Map{"key": "value"}, dynmgrm.Set[string]{"a", "b", "c"}},
		SomeStringSet: dynmgrm.Set[string]{"a", "b", "c"},
		SomeIntSet:    dynmgrm.Set[int]{1, 2, 3},
		SomeFloatSet:  dynmgrm.Set[float64]{1.1, 2.2, 3.3},
		SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
	}
	defer deleteData(t, testTableName, data.PK, data.SK)

	tx := db.Begin()
	tx.Create(data)
	tx.Rollback()

	expect := make([]map[string]*dynamodb.AttributeValue, 0)
	actual := scanData(t, testTableName)
	if diff := cmp.Diff(expect, actual, setCmpOpts...); diff != "" {
		t.Errorf("mismatch (-got +want)\n%s", diff)
	}
}
