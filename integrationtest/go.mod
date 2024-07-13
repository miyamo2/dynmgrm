module github.com/miyamo2/dynmgrm/integrationtest

go 1.21

replace github.com/miyamo2/dynmgrm => ../

require (
	github.com/aws/aws-sdk-go v1.54.15
	github.com/miyamo2/godynamo v1.4.0
	github.com/google/go-cmp v0.6.0
	github.com/joho/godotenv v1.5.1
	github.com/miyamo2/dynmgrm v0.0.0-00010101000000-000000000000
	github.com/miyamo2/sqldav v0.2.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/aws/aws-sdk-go-v2 v1.30.3 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.24 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.14.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.16 // indirect
	github.com/aws/smithy-go v1.20.3 // indirect
	github.com/btnguyen2k/consu/g18 v0.1.0 // indirect
	github.com/btnguyen2k/consu/reddo v0.1.9 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
