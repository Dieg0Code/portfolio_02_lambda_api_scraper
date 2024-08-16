package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dieg0code/authorizer/auth"
	"github.com/dieg0code/authorizer/aws"
	"github.com/dieg0code/authorizer/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting authorizer...")
	jwtValidator := auth.NewJWTValidator()
	policy := aws.NewPolicyImpl()
	handler := handler.NewAuthorizerHandler(policy, jwtValidator)

	lambda.Start(handler.HandleAuthorizer)

	logrus.Info("Authorizer started")
}
