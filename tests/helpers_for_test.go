package tests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/joho/godotenv"
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
	"sort"
	"testing"
)

type TestTable struct {
	PK             string `gorm:"primaryKey"`
	SK             int    `gorm:"primaryKey"`
	SomeString     string
	SomeInt        int
	SomeFloat      float64
	SomeBool       bool
	SomeBinary     []byte
	SomeList       dynmgrm.List
	SomeMap        dynmgrm.Map
	SomeStringSets dynmgrm.Sets[string]
	SomeIntSets    dynmgrm.Sets[int]
	SomeFloatSets  dynmgrm.Sets[float64]
	SomeBinarySets dynmgrm.Sets[[]byte]
	Any            string
}

type TestTablePKSomeStringIndex struct {
	PK             string `gorm:"primaryKey"`
	SomeString     string `gorm:"primaryKey"`
	SK             int
	SomeInt        int
	SomeFloat      float64
	SomeBool       bool
	SomeBinary     []byte
	SomeList       dynmgrm.List
	SomeMap        dynmgrm.Map
	SomeStringSets dynmgrm.Sets[string]
	SomeIntSets    dynmgrm.Sets[int]
	SomeFloatSets  dynmgrm.Sets[float64]
	SomeBinarySets dynmgrm.Sets[[]byte]
	Any            string
}

var dynamoDBClient dynamodbiface.DynamoDBAPI

var setsCmpOpts = []cmp.Option{
	cmpopts.SortSlices(func(i, j int) bool {
		return i < j
	}),
	cmpopts.SortSlices(func(i, j float64) bool {
		return i < j
	}),
	cmpopts.SortSlices(func(i, j string) bool {
		ss := []string{i, j}
		sort.Strings(ss)
		return ss[0] == j
	}),
}

var avCmpOpts = []cmp.Option{
	cmp.AllowUnexported(dynamodb.AttributeValue{}),
}

func init() {
	// for local testing
	_ = godotenv.Load("./.env")

	dynamoDBClient = dynamodb.New(session.New())
}

func getGormDB(t *testing.T) *gorm.DB {
	t.Helper()
	d := dynmgrm.New()
	db, err := gorm.Open(
		d,
		&gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		t.Fatalf("failed to open database, got error %v", err)
	}
	return db
}

func dataPreparation(t *testing.T, testData []map[string]*dynamodb.AttributeValue, tableName string) {
	t.Helper()
	for _, av := range testData {
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		_, err := dynamoDBClient.PutItem(input)
		if err != nil {
			t.Fatalf("failed to put item: %s", err)
		}
	}
}

func dataCleanup(t *testing.T, testData []map[string]*dynamodb.AttributeValue, tableName string) {
	t.Helper()
	for _, av := range testData {
		input := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"pk": av["pk"],
				"sk": av["sk"],
			},
			TableName: aws.String(tableName),
		}
		_, err := dynamoDBClient.DeleteItem(input)
		if err != nil {
			t.Fatalf("failed to delete item: %s", err)
		}
	}
}

func getData(t *testing.T, tableName string, pk string, sk int) map[string]*dynamodb.AttributeValue {
	t.Helper()
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {S: aws.String(pk)},
			"sk": {N: aws.String(fmt.Sprint(sk))},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynamoDBClient.GetItem(input)
	if err != nil {
		t.Fatalf("failed to get item: %s", err)
	}
	return result.Item
}

func scanData(t *testing.T, tableName string) []map[string]*dynamodb.AttributeValue {
	t.Helper()
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	items := make([]map[string]*dynamodb.AttributeValue, 0)
	dynamoDBClient.ScanPages(input, func(output *dynamodb.ScanOutput, lastPage bool) bool {
		items = append(items, output.Items...)
		input.ExclusiveStartKey = output.LastEvaluatedKey
		return !lastPage
	})
	result, err := dynamoDBClient.Scan(input)
	if err != nil {
		t.Fatalf("failed to scan: %s", err)
	}
	return result.Items
}
