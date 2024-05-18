package integrationtest

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
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

func Test_Insert_With_Transaction_Success(t *testing.T) {
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

	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(data).Error
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

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

func Test_Insert_With_Transaction_Fail(t *testing.T) {
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

	errFoo := errors.New("foo")
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Create(data)
		return errFoo
	})

	if !errors.Is(err, errFoo) {
		t.Errorf("unexpected error: %v", err)
	}

	expect := make([]map[string]*dynamodb.AttributeValue, 0)
	actual := scanData(t, testTableName)
	if diff := cmp.Diff(expect, actual, setCmpOpts...); diff != "" {
		t.Errorf("mismatch (-got +want)\n%s", diff)
	}
}

func Test_Insert_With_Create_Has_TypedList_Column(t *testing.T) {
	db := getGormDB(t)

	data := TestTableWithTypedList{
		PK:         "Partition4",
		SK:         1,
		SomeString: "Hello",
		TypedList: dynmgrm.TypedList[TypedListValue]{
			{
				SomeString:    "Hello",
				SomeInt:       1,
				SomeFloat:     1.1,
				SomeBool:      true,
				SomeBinary:    []byte("ABC"),
				SomeMap:       dynmgrm.Map{"key": "value"},
				SomeList:      dynmgrm.List{"a", dynmgrm.Map{"key": "value"}, dynmgrm.Set[string]{"a", "b", "c"}},
				SomeStringSet: dynmgrm.Set[string]{"a", "b", "c"},
				SomeIntSet:    dynmgrm.Set[int]{1, 2, 3},
				SomeFloatSet:  dynmgrm.Set[float64]{1.1, 2.2, 3.3},
				SomeBinarySet: dynmgrm.Set[[]byte]{[]byte("a"), []byte("b"), []byte("c")},
			},
		},
	}
	defer deleteData(t, testTableName, data.PK, data.SK)

	db.Create(data)

	expect := []map[string]*dynamodb.AttributeValue{
		{
			"pk":          &dynamodb.AttributeValue{S: aws.String(data.PK)},
			"sk":          &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(data.SK))},
			"some_string": &dynamodb.AttributeValue{S: aws.String("Hello")},
			"typed_list": &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{
					{
						M: map[string]*dynamodb.AttributeValue{
							"some_string": {S: aws.String("Hello")},
							"some_int":    {N: aws.String("1")},
							"some_float":  {N: aws.String("1.1")},
							"some_bool":   {BOOL: aws.Bool(true)},
							"some_binary": {B: []byte("ABC")},
							"some_map": {
								M: map[string]*dynamodb.AttributeValue{
									"key": {S: aws.String("value")},
								}},
							"some_list": {
								L: []*dynamodb.AttributeValue{
									{S: aws.String("a")},
									{M: map[string]*dynamodb.AttributeValue{"key": {S: aws.String("value")}}},
									{SS: aws.StringSlice([]string{"a", "b", "c"})},
								},
							},
							"some_string_set": {SS: aws.StringSlice([]string{"a", "b", "c"})},
							"some_int_set":    {NS: aws.StringSlice([]string{"1", "2", "3"})},
							"some_float_set":  {NS: aws.StringSlice([]string{"1.1", "2.2", "3.3"})},
							"some_binary_set": {BS: [][]byte{[]byte("a"), []byte("b"), []byte("c")}},
						},
					},
				},
			},
		},
	}
	actual := scanData(t, testTableName)
	if diff := cmp.Diff(expect, actual, setCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}

func Test_Insert_With_Create_Has_NetedStruct_Column(t *testing.T) {
	db := getGormDB(t)

	data := TestTableWithNested{
		PK: "Partition4",
		SK: 1,
		SomeMap: NestedAttribute{
			SomeString: "Hello",
			SomeNumber: 1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
		},
	}
	defer deleteData(t, testTableName, data.PK, data.SK)

	db.Create(data)

	expect := []map[string]*dynamodb.AttributeValue{
		{
			"pk": {
				S: aws.String("Partition4"),
			},
			"sk": {
				N: aws.String("1"),
			},
			"some_map": &dynamodb.AttributeValue{
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
				}},
		},
	}
	actual := scanData(t, testTableName)
	if diff := cmp.Diff(expect, actual, setCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}
