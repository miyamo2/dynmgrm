package integrationtest

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/btnguyen2k/godynamo"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

type CreateGSITestTable struct {
	PK    string `dynmgrm:"pk"`
	SK    string `dynmgrm:"sk"`
	GsiPK string `dynmgrm:"gsi-pk:gsi_pk-gsi_sk-index"`
	GsiSK string `dynmgrm:"gsi-sk:gsi_pk-gsi_sk-index"`
}

func Test_CreateGSI(t *testing.T) {
	tableName := "create_gsi_test_tables"
	gsiName := "gsi_pk-gsi_sk-index"
	db := getGormDB(t)

	createTable(t,
		tableName,
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("pk"),
			AttributeType: aws.String("S"),
		},
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("sk"),
			AttributeType: aws.String("S"),
		},
		1,
		3)
	err := db.Migrator().CreateIndex(&CreateGSITestTable{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer deleteTable(t, tableName, 0)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	ctx, cancelF := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelF()
	if err := godynamo.WaitForGSIStatus(ctx, sqlDB, tableName, gsiName, []string{"ACTIVE"}, 100*time.Millisecond); err != nil {
		return
	}
	table := getTableWithRetry(t, tableName, 0)
	if table == nil {
		t.Fatalf("table not found")
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
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_pk"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_sk"),
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
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndexDescription{
				{
					IndexName: aws.String(gsiName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("gsi_pk"),
							KeyType:       aws.String("HASH"),
						},
						{
							AttributeName: aws.String("gsi_sk"),
							KeyType:       aws.String("RANGE"),
						},
					},
					Projection: &dynamodb.Projection{
						ProjectionType: aws.String("ALL"),
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
			"Table.OnDemandThroughput",
		),
		cmpopts.IgnoreFields(*table.Table.GlobalSecondaryIndexes[0], "IndexSizeBytes", "IndexStatus", "ItemCount", "ProvisionedThroughput", "IndexArn", "OnDemandThroughput"),
	}
	if diff := cmp.Diff(except, *table, cmpopt...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}

type CreateGSIWithCapacityTestTable struct {
	PK    string `dynmgrm:"pk"`
	SK    string `dynmgrm:"sk"`
	GsiPK string `dynmgrm:"gsi-pk:gsi_pk-gsi_sk-index"`
	GsiSK string `dynmgrm:"gsi-sk:gsi_pk-gsi_sk-index"`
}

func (t CreateGSIWithCapacityTestTable) WCU() int {
	return 1
}

func (t CreateGSIWithCapacityTestTable) RCU() int {
	return 1
}

func Test_CreateGSI_With_Capacity(t *testing.T) {
	tableName := "create_gsi_with_capacity_test_tables"
	gsiName := "gsi_pk-gsi_sk-index"
	db := getGormDB(t)

	createTable(t,
		tableName,
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("pk"),
			AttributeType: aws.String("S"),
		},
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("sk"),
			AttributeType: aws.String("S"),
		},
		0,
		3)
	err := db.Migrator().CreateIndex(&CreateGSIWithCapacityTestTable{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer deleteTable(t, tableName, 0)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	ctx, cancelF := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelF()
	if err := godynamo.WaitForGSIStatus(ctx, sqlDB, tableName, gsiName, []string{"ACTIVE"}, 100*time.Millisecond); err != nil {
		return
	}
	table := getTableWithRetry(t, tableName, 0)
	if table == nil {
		t.Fatalf("table not found")
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
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_pk"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_sk"),
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
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndexDescription{
				{
					IndexName: aws.String(gsiName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("gsi_pk"),
							KeyType:       aws.String("HASH"),
						},
						{
							AttributeName: aws.String("gsi_sk"),
							KeyType:       aws.String("RANGE"),
						},
					},
					Projection: &dynamodb.Projection{
						ProjectionType: aws.String("ALL"),
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughputDescription{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
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
			"Table.OnDemandThroughput",
		),
		cmpopts.IgnoreFields(*table.Table.GlobalSecondaryIndexes[0], "IndexSizeBytes", "IndexStatus", "ItemCount", "IndexArn", "OnDemandThroughput"),
	}
	if diff := cmp.Diff(except, *table, cmpopt...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}

type CreateGSIWithNonProjectiveTestTable struct {
	PK                string `dynmgrm:"pk"`
	SK                int    `dynmgrm:"sk"`
	GsiPK             string `dynmgrm:"gsi-pk:gsi_pk-gsi_sk-index;non-projective:[lsi_sk-index]"`
	GsiSK             string `dynmgrm:"gsi-sk:gsi_pk-gsi_sk-index;non-projective:[lsi_sk-index]"`
	ProjectiveAttr    string
	NonProjectiveAttr string `dynmgrm:"non-projective:[gsi_pk-gsi_sk-index]"`
}

func Test_CreateGSI_With_Non_Projective_Attrs(t *testing.T) {
	tableName := "create_gsi_with_non_projective_test_tables"
	gsiName := "gsi_pk-gsi_sk-index"
	db := getGormDB(t)

	createTable(t,
		tableName,
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("pk"),
			AttributeType: aws.String("S"),
		},
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("sk"),
			AttributeType: aws.String("S"),
		},
		1,
		3)
	err := db.Migrator().CreateIndex(&CreateGSIWithNonProjectiveTestTable{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer deleteTable(t, tableName, 0)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	ctx, cancelF := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelF()
	if err := godynamo.WaitForGSIStatus(ctx, sqlDB, tableName, gsiName, []string{"ACTIVE"}, 100*time.Millisecond); err != nil {
		return
	}
	table := getTableWithRetry(t, tableName, 0)
	if table == nil {
		t.Fatalf("table not found")
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
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_pk"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_sk"),
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
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndexDescription{
				{
					IndexName: aws.String(gsiName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("gsi_pk"),
							KeyType:       aws.String("HASH"),
						},
						{
							AttributeName: aws.String("gsi_sk"),
							KeyType:       aws.String("RANGE"),
						},
					},
					Projection: &dynamodb.Projection{
						NonKeyAttributes: []*string{aws.String("pk"), aws.String("sk"), aws.String("projective_attr")},
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
		cmpopts.SortSlices(func(a, b *string) bool {
			return *a < *b
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
			"Table.OnDemandThroughput",
		),
		cmpopts.IgnoreFields(*table.Table.GlobalSecondaryIndexes[0], "IndexSizeBytes", "IndexStatus", "ItemCount", "ProvisionedThroughput", "IndexArn", "OnDemandThroughput"),
	}
	if diff := cmp.Diff(except, *table, cmpopt...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}

type CreateGSIWithKeysOnlyTestTable struct {
	PK    string `dynmgrm:"pk;non-projective:[gsi_pk-gsi_sk-index]"`
	SK    int    `dynmgrm:"sk;non-projective:[gsi_pk-gsi_sk-index]"`
	GsiPK string `dynmgrm:"gsi-pk:gsi_pk-gsi_sk-index;"`
	GsiSK string `dynmgrm:"gsi-sk:gsi_pk-gsi_sk-index;"`
}

func Test_CreateGSI_With_Keys_Only(t *testing.T) {
	tableName := "create_gsi_with_keys_only_test_tables"
	gsiName := "gsi_pk-gsi_sk-index"
	db := getGormDB(t)

	createTable(t,
		tableName,
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("pk"),
			AttributeType: aws.String("S"),
		},
		dynamodb.AttributeDefinition{
			AttributeName: aws.String("sk"),
			AttributeType: aws.String("S"),
		},
		1,
		3)
	err := db.Migrator().CreateIndex(&CreateGSIWithKeysOnlyTestTable{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer deleteTable(t, tableName, 0)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	ctx, cancelF := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelF()
	if err := godynamo.WaitForGSIStatus(ctx, sqlDB, tableName, gsiName, []string{"ACTIVE"}, 100*time.Millisecond); err != nil {
		return
	}
	table := getTableWithRetry(t, tableName, 0)
	if table == nil {
		t.Fatalf("table not found")
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
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_pk"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("gsi_sk"),
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
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndexDescription{
				{
					IndexName: aws.String(gsiName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("gsi_pk"),
							KeyType:       aws.String("HASH"),
						},
						{
							AttributeName: aws.String("gsi_sk"),
							KeyType:       aws.String("RANGE"),
						},
					},
					Projection: &dynamodb.Projection{
						ProjectionType: aws.String("KEYS_ONLY"),
					},
				},
			},
		},
	}
	cmpopt := []cmp.Option{
		cmpopts.SortSlices(func(a, b *dynamodb.AttributeDefinition) bool {
			return *a.AttributeName < *b.AttributeName
		}),
		cmpopts.SortSlices(func(a, b *string) bool {
			return *a < *b
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
			"Table.OnDemandThroughput",
		),
		cmpopts.IgnoreFields(*table.Table.GlobalSecondaryIndexes[0], "IndexSizeBytes", "IndexStatus", "ItemCount", "ProvisionedThroughput", "IndexArn", "OnDemandThroughput"),
	}
	if diff := cmp.Diff(except, *table, cmpopt...); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}
}
