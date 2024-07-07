export
AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=dummy
AWS_SECRET_ACCESS_KEY=dummy
DYNAMODB_ENDPOINT=http://localhost:8000

.PHONY: up create-table test
up:
	@env
	docker compose up -d

create-table:up
	@env
	# Somehow the yes command does not work...
	# Therefore, the following tentative hacks are being used
	aws dynamodb create-table --cli-input-json file://table-def.json --endpoint-url http://localhost:8000 | tee /dev/null

test:create-table
	@env
	docker compose up -d
	go mod tidy
	go test -v ./... | tee /dev/null
	go clean -testcache
	docker compose stop