package response

type LogInUserResponse struct {
	Token string `json:"token" validate:"required"`
}
