package tests

import (
	"github.com/google/go-cmp/cmp"
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm/clause"
	"testing"
)

func Test_Select_All(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition1",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
		{
			PK:         "Partition2",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition2",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
		{
			PK:         "Partition3",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
		},
		{
			PK:         "Partition3",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
		},
	}

	var result []TestTable
	err := db.Select("*").Table("test_tables").Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_Columns_Specify(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK: "Partition1",
			SK: 1,
		},
		{
			PK: "Partition1",
			SK: 2,
		},
		{
			PK: "Partition2",
			SK: 1,
		},
		{
			PK: "Partition2",
			SK: 2,
		},
		{
			PK: "Partition3",
			SK: 1,
		},
		{
			PK: "Partition3",
			SK: 2,
		},
	}

	var result []TestTable
	err := db.Select("pk", "sk").Table("test_tables").Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_PK(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition1",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`pk = ?`, "Partition1").Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_PK_And_SK(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`pk = ? AND sk = ?`, "Partition1", 1).Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_Secondary_Index(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
	}

	var withTableClause []TestTable
	err := db.Table("test_tables").Clauses(
		dynmgrm.SecondaryIndex("pk-some_string-index"),
	).Where(`pk = ? AND some_string = ?`, "Partition1", "こんにちは").Scan(&withTableClause).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, withTableClause, setsCmpOpts...); diff != "" {
		t.Errorf("withTableClause mismatch (-want +got):\n%s", diff)
	}

	var withSecondaryIdxOfString []TestTable
	err = db.Clauses(
		dynmgrm.SecondaryIndex("pk-some_string-index",
			dynmgrm.SecondaryIndexOf("test_tables"),
		),
	).Where(`pk = ? AND some_string = ?`, "Partition1", "こんにちは").Scan(&withSecondaryIdxOfString).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, withSecondaryIdxOfString, setsCmpOpts...); diff != "" {
		t.Errorf("with_secondary_idx_of_string mismatch (-want +got):\n%s", diff)
	}

	var withSecondaryIdxOfTableClause []TestTable
	err = db.Clauses(
		dynmgrm.SecondaryIndex("pk-some_string-index",
			dynmgrm.SecondaryIndexOf(clause.Table{Name: "test_tables"}),
		),
	).Where(`pk = ? AND some_string = ?`, "Partition1", "こんにちは").
		Scan(&withSecondaryIdxOfTableClause).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, withSecondaryIdxOfTableClause, setsCmpOpts...); diff != "" {
		t.Errorf("withSecondaryIdxOfTableClause mismatch (-want +got):\n%s", diff)
	}

	var withTableNameDotIndexName []TestTable
	err = db.Clauses(
		dynmgrm.SecondaryIndex("test_tables.pk-some_string-index"),
	).Where(`pk = ? AND some_string = ?`, "Partition1", "こんにちは").Scan(&withTableNameDotIndexName).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, withTableNameDotIndexName, setsCmpOpts...); diff != "" {
		t.Errorf("withTableNameDotIndexName mismatch (-want +got):\n%s", diff)
	}

	// Not supported yet
	//withModelExpect := TestTablePKSomeStringIndex{
	//	PK:         "Partition1",
	//	SomeString: "こんにちは",
	//	SK:         2,
	//	SomeInt:    2,
	//	SomeFloat:  2.2,
	//	SomeBool:   false,
	//	SomeBinary: []byte("GHI"),
	//	SomeList: dynmgrm.List{
	//		"こんにちは",
	//		float64(2),
	//		2.2,
	//		false,
	//		[]byte("GHI"),
	//	},
	//	SomeMap: dynmgrm.Map{
	//		"some_string": "こんにちは",
	//		"some_number": 2.2,
	//		"some_bool":   false,
	//		"some_binary": []byte("GHI"),
	//	},
	//	SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
	//	SomeIntSets:    dynmgrm.Sets[int]{2, 4},
	//	SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
	//	SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
	//}
	//var withModel = TestTablePKSomeStringIndex{
	//	PK:         "Partition1",
	//	SomeString: "こんにちは",
	//}
	//err = db.Clauses(
	//	dynmgrm.SecondaryIndex(
	//		"test_tables.pk-some_string-index"),
	//).First(&withModel).Error
	//if err != nil {
	//	t.Errorf("unexpected error: %v", err)
	//	err = nil
	//}
	//if diff := cmp.Diff(withModelExpect, withModel, setsCmpOpts...); diff != "" {
	//	t.Errorf("withModelAs mismatch (-want +got):\n%s", diff)
	//}
}

func Test_Select_With_BeginWith(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition2",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition3",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`begins_with("some_string", ?)`, "H").Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_IsMissing(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition3",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
		},
		{
			PK:         "Partition3",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`some_list IS MISSING`).Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_IsNotMissing(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition1",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
		{
			PK:         "Partition2",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition2",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`some_list IS NOT MISSING`).Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_Contains(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
		{
			PK:         "Partition2",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`CONTAINS("some_string_sets", ?)`, "こんにちは").Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_Size(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition1",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
		{
			PK:         "Partition2",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition2",
			SK:         2,
			SomeString: "こんにちは",
			SomeInt:    2,
			SomeFloat:  2.2,
			SomeBool:   false,
			SomeBinary: []byte("GHI"),
			SomeList: dynmgrm.List{
				"こんにちは",
				float64(2),
				2.2,
				false,
				[]byte("GHI"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "こんにちは",
				"some_number": 2.2,
				"some_bool":   false,
				"some_binary": []byte("GHI"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"こんにちは", "世界"},
			SomeIntSets:    dynmgrm.Sets[int]{2, 4},
			SomeFloatSets:  dynmgrm.Sets[float64]{2.2, 4.4},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("GHI"), []byte("JKL")},
			Any:            "0",
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`SIZE("some_string_sets") = ?`, 2).Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Select_With_AttributeType(t *testing.T) {
	db := getGormDB(t)
	dataPreparation(t, testDataForSelect, testTableName)
	defer dataCleanup(t, testDataForSelect, testTableName)

	expected := []TestTable{
		{
			PK:         "Partition1",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
		{
			PK:         "Partition2",
			SK:         1,
			SomeString: "Hello",
			SomeInt:    1,
			SomeFloat:  1.1,
			SomeBool:   true,
			SomeBinary: []byte("ABC"),
			SomeList: dynmgrm.List{
				"Hello",
				float64(1),
				1.1,
				true,
				[]byte("ABC"),
			},
			SomeMap: dynmgrm.Map{
				"some_string": "Hello",
				"some_number": 1.1,
				"some_bool":   true,
				"some_binary": []byte("ABC"),
			},
			SomeStringSets: dynmgrm.Sets[string]{"Hello", "World"},
			SomeIntSets:    dynmgrm.Sets[int]{1, 2},
			SomeFloatSets:  dynmgrm.Sets[float64]{1.1, 2.2},
			SomeBinarySets: dynmgrm.Sets[[]byte]{[]byte("ABC"), []byte("DEF")},
			Any:            "any",
		},
	}

	var result []TestTable
	err := db.Table("test_tables").Where(`ATTRIBUTE_TYPE("any", ?)`, "S").Scan(&result).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		err = nil
	}
	if diff := cmp.Diff(expected, result, setsCmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
