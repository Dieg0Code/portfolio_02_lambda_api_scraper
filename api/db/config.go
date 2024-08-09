package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewDynamoDB(region string) *dynamodb.DynamoDB {
	var dynamoEndpoint string
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		dynamoEndpoint = "http://localhost:8000"
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(dynamoEndpoint),
	}))

	return dynamodb.New(sess)
}
