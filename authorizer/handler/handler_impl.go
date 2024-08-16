package handler

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dieg0code/authorizer/auth"
	"github.com/dieg0code/authorizer/aws"
	"github.com/sirupsen/logrus"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthorizerHandlerImpl struct {
	Policy       aws.Policy
	JWTValidator auth.JWTValidator
}

// HandleAuthorizer implements AuthorizerHandler.
func (a *AuthorizerHandlerImpl) HandleAuthorizer(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := strings.TrimPrefix(event.AuthorizationToken, "Bearer ")

	claims, err := a.JWTValidator.ValidateToken(token, jwtSecret)
	if err != nil {
		logrus.WithError(err).Error("error validating token")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("invalid token")
	}

	// Check if the token has expired
	exp, ok := claims["exp"].(float64)
	if !ok {
		logrus.Error("invalid exp in token")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("invalid exp in token")
	}

	if time.Now().Unix() > int64(exp) {
		logrus.Error("token has expired")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("token has expired")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		logrus.Error("invalid user_id in token")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("invalid user_id in token")
	}

	return a.Policy.GeneratePolicy(userID, "Allow", event.MethodArn), nil
}

func NewAuthorizerHandler(policy aws.Policy, jwtValidator auth.JWTValidator) AuthorizerHandler {
	return &AuthorizerHandlerImpl{
		Policy:       policy,
		JWTValidator: jwtValidator,
	}
}
