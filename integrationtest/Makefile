export
AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=dummy
AWS_SECRET_ACCESS_KEY=dummy
DYNAMODB_ENDPOINT=http://localhost:8000

test:
	@env
	docker-compose up -d
	# Somehow the yes command does not work...
	# Therefore, the following tentative hacks are being used
	# yes q | aws dynamodb create-table --cli-input-json file://table-def.json --endpoint-url http://localhost:8000
	aws dynamodb create-table --cli-input-json file://table-def.json --endpoint-url http://localhost:8000 > /dev/null
	go mod tidy
	go test -v ./... | tee /dev/null
	docker-compose down