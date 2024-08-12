compile_lambda:
	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o bootstrap main.go

zip_lambda:
	zip api_scraper_lambda.zip bootstrap


start_db:
	docker run -d --name dynamodb -p 8000:8000 amazon/dynamodb-local


create_table:
	aws dynamodb create-table \
		--table-name products \
		--attribute-definitions \
			AttributeName=ProductID,AttributeType=S \
		--key-schema \
			AttributeName=ProductID,KeyType=HASH \
		--provisioned-throughput \
			ReadCapacityUnits=5,WriteCapacityUnits=5 \
		--endpoint-url http://localhost:8000
stop_db:
	docker stop dynamodb
	docker rm dynamodb