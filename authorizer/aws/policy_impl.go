package aws

import "github.com/aws/aws-lambda-go/events"

type PolicyImpl struct{}

// GeneratePolicy implements Policy.
func (p *PolicyImpl) GeneratePolicy(principalID string, effect string, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}
	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	authResponse.Context = map[string]interface{}{
		"user_id": principalID,
	}
	return authResponse
}

func NewPolicyImpl() Policy {
	return &PolicyImpl{}
}
