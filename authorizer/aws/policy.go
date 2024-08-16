package aws

import "github.com/aws/aws-lambda-go/events"

type Policy interface {
	GeneratePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse
}
