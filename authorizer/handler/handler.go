package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type AuthorizerHandler interface {
	HandleAuthorizer(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)
}
