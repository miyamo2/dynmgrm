package tests

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

type CreateTableTestTable struct {
	PK string `dynmgrm:"pk""`
	SK int    `dynmgrm:"sk""`
}

func Test_CreateTable(t *testing.T) {
	tableName := "create_table_test_tables"
	db := getGormDB(t)

	err := db.Migrator().CreateTable(&CreateTableTestTable{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer deleteTable(t, tableName, 0)

	table := getTable(t, tableName)
	if table == nil {
		t.Errorf("table not created")
	}
	except := dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableName: &tableName,
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("pk"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("sk"),
					AttributeType: aws.String("N"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("pk"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("sk"),
					KeyType:       aws.String("RANGE"),
				},
			},
		},
	}
	cmpopt := []cmp.Option{
		cmpopts.IgnoreUnexported(dynamodb.DescribeTableOutput{}),
		cmpopts.IgnoreFields(*table,
			"Table.ArchivalSummary",
			"Table.BillingModeSummary",
			"Table.CreationDateTime",
			"Table.DeletionProtectionEnabled",
			"Table.GlobalTableVersion",
			"Table.ItemCount",
			"Table.LatestStreamArn",
			"Table.LatestStreamLabel",
			"Table.ProvisionedThroughput",
			"Table.Replicas",
			"Table.RestoreSummary",
			"Table.SSEDescription",
			"Table.StreamSpecification",
			"Table.TableArn",
			"Table.TableClassSummary",
			"Table.TableId",
			"Table.TableSizeBytes",
			"Table.TableStatus",
		),
	}
	if diff := cmp.Diff(except, *table, cmpopt...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}

type CreateTableTestWithLSITable struct {
	PK            string `dynmgrm:"pk""`
	SK            int    `dynmgrm:"sk""`
	Name          string `dynmgrm:"lsi:pk-name-index"`
	Projective    string
	NonProjective string `dynmgrm:"non-projective:[pk-name-index]"`
}

func Test_CreateTable_With_LSI(t *testing.T) {
	tableName := "create_table_test_with_lsi_tables"
	db := getGormDB(t)

	err := db.Migrator().CreateTable(&CreateTableTestWithLSITable{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer deleteTable(t, tableName, 0)

	table := getTable(t, tableName)
	if table == nil {
		t.Errorf("table not created")
	}
	except := dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableName: &tableName,
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("pk"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("sk"),
					AttributeType: aws.String("N"),
				},
				{
					AttributeName: aws.String("name"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("pk"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("sk"),
					KeyType:       aws.String("RANGE"),
				},
			},
			LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndexDescription{
				{
					IndexName: aws.String("pk-name-index"),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("pk"),
							KeyType:       aws.String("HASH"),
						},
						{
							AttributeName: aws.String("name"),
							KeyType:       aws.String("RANGE"),
						},
					},
					Projection: &dynamodb.Projection{
						NonKeyAttributes: []*string{aws.String("projective")},
						ProjectionType:   aws.String("INCLUDE"),
					},
				},
			},
		},
	}
	cmpopt := []cmp.Option{
		cmpopts.SortSlices(func(a, b *dynamodb.AttributeDefinition) bool {
			return *a.AttributeName < *b.AttributeName
		}),
		cmpopts.IgnoreUnexported(dynamodb.DescribeTableOutput{}),
		cmpopts.IgnoreFields(*table,
			"Table.ArchivalSummary",
			"Table.BillingModeSummary",
			"Table.CreationDateTime",
			"Table.DeletionProtectionEnabled",
			"Table.GlobalTableVersion",
			"Table.ItemCount",
			"Table.LatestStreamArn",
			"Table.LatestStreamLabel",
			"Table.ProvisionedThroughput",
			"Table.Replicas",
			"Table.RestoreSummary",
			"Table.SSEDescription",
			"Table.StreamSpecification",
			"Table.TableArn",
			"Table.TableClassSummary",
			"Table.TableId",
			"Table.TableSizeBytes",
			"Table.TableStatus",
		),
		cmpopts.IgnoreFields(dynamodb.LocalSecondaryIndexDescription{},
			"IndexArn",
			"IndexSizeBytes",
			"ItemCount",
			"IndexArn",
		),
	}
	if diff := cmp.Diff(except, *table, cmpopt...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}
